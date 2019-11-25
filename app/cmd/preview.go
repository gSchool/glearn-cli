package cmd

import (
	"archive/zip"
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	proxyReader "github.com/Galvanize-IT/glearn-cli/app/proxy_reader"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	pb "github.com/cheggaaa/pb/v3"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// tmpFile is used throughout as the temporary zip file target location.
const tmpFile string = "preview-curriculum.zip"

// LearnPreviewResponse is a simple struct defining the shape of data we care about
// that comes back from notifying Learn for decoding into.
type LearnPreviewResponse struct {
	ReleaseID  int    `json:"release_id"`
	PreviewURL string `json:"preview_url"`
	Errors     string `json:"errors"`
	Status     string `json:"status"`
}

// previewCmd is executed when the `glearn preview` command is used. Preview's concerns:
// 1. Compress directory/file into target location.
// 2. Defer cleaning up the file after command is finished.
// 3. Create a checksum for the zip file.
// 4. Upload the zip file to s3.
// 5. Notify learn that new content is available for building.
// 6. Handle progress bar for s3 upload.
var previewCmd = &cobra.Command{
	Use:   "preview [file_path]",
	Short: "Preview your content",
	Long:  `Long description for previewing`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Takes one argument which is the filepath to the directory you want zipped/previewed
		if len(args) != 1 {
			fmt.Println("Usage: `learn preview` takes one argument")
			os.Exit(1)
			return
		}

		// Compress directory, output -> tmpFile
		err := compressDirectory(args[0], tmpFile)
		if err != nil {
			fmt.Printf("Error compressing directory %s: %v", args[0], err)
			os.Exit(1)
			return
		}

		// Removes artifacts on user's machine
		defer cleanUpFiles()

		// Open file so we can get a checksum as well as send to s3
		f, err := os.Open(tmpFile)
		if err != nil {
			fmt.Printf("Failed to open file %q, %v", tmpFile, err)
			return
		}
		defer f.Close()

		// Create checksum of files in directory
		checksum, err := createChecksumFromZip(f)
		if err != nil {
			fmt.Printf("Failed to create checksum for compressed file. Err: %v", err)
			os.Exit(1)
			return
		}

		// Send compressed zip file to s3
		bucketKey, err := uploadToS3(f, checksum)
		if err != nil {
			fmt.Printf("Failed to upload zip file to s3. Err: %v", err)
			os.Exit(1)
			return
		}

		fileInfo, err := os.Stat(args[0])
		if err != nil {
			fmt.Printf("Failed to get stats on file. Err: %v", err)
			os.Exit(1)
			return
		}

		isDirectory := fileInfo.IsDir()
		res, err := notifyLearn(bucketKey, isDirectory)
		if err != nil {
			fmt.Printf("Failed to notify learn of new preview content. Err: %v", err)
			os.Exit(1)
			return
		}

		// Should the above call to notify learn just return an identifier to track/poll? Instead of
		// passing bucket key again?
		var attempts uint8 = 20
		res, err = pollForBuildResponse(res.ReleaseID, &attempts)
		if err != nil {
			fmt.Printf("Failed to poll Learn for your new preview build. Err: %v", err)
			os.Exit(1)
			return
		}

		fmt.Printf("Sucessfully uploaded your preview! You can find your content at: %s", res.PreviewURL)
		exec.Command("bash", "-c", fmt.Sprintf("open %s", res.PreviewURL)).Output()
	},
}

func pollForBuildResponse(releaseID int, attempts *uint8) (*LearnPreviewResponse, error) {
	apiToken, ok := viper.Get("api_token").(string)
	if !ok {
		return nil, errors.New("Please set your api_token in ~/.glearn-config.yaml")
	}

	client := &http.Client{Timeout: time.Second * 30}

	req, err := http.NewRequest("GET", fmt.Sprintf("http://localhost:3003/api/v1/releases/%d/release_polling", releaseID), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiToken))

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Error: response status: %d", res.StatusCode)
	}

	var l LearnPreviewResponse
	err = json.NewDecoder(res.Body).Decode(&l)
	if err != nil {
		return nil, err
	}

	if l.Status == "processing" || l.Status == "pending" {
		*attempts--
		time.Sleep(2 * time.Second)

		if *attempts == uint8(0) {
			return nil, errors.New(
				"Sorry, we are having trouble requesting your preview build from Learn. Please try again",
			)
		}

		return pollForBuildResponse(releaseID, attempts)
	}

	return &l, nil
}

// notifyLearn takes an s3 bucket key name as an argument is used to tell Learn there is new preview
// content on s3 and where to find it so it can build/preview.
func notifyLearn(bucketKey string, isDirectory bool) (*LearnPreviewResponse, error) {
	apiToken, ok := viper.Get("api_token").(string)
	if !ok {
		return nil, errors.New("Please set your api_token in ~/.glearn-config.yaml")
	}

	payload := map[string]string{
		"s3_key": bucketKey,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	var endpoint string
	if isDirectory {
		endpoint = "/api/v1/releases"
	} else {
		endpoint = "/api/v1/content_file"
	}

	client := &http.Client{Timeout: time.Second * 30}

	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("http://localhost:3003%s", endpoint),
		bytes.NewBuffer(payloadBytes),
	)
	if err != nil {
		return nil, err
	}
	defer req.Body.Close()

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiToken))

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Error: response status: %d", res.StatusCode)
	}

	l := &LearnPreviewResponse{}
	json.NewDecoder(res.Body).Decode(l)

	return l, nil
}

// uploadToS3 takes a file and it's checksum and uploads it to s3 in the appropriate bucket/key
func uploadToS3(file *os.File, checksum string) (string, error) {
	// Coerce AWS credentials to strings
	accessKeyID, ok := viper.Get("aws_access_key_id").(string)
	if !ok {
		return "", errors.New("Your aws_access_key_id must be a string")
	}
	secretAccessKey, ok := viper.Get("aws_secret_access_key").(string)
	if !ok {
		return "", errors.New("Your aws_secret_access_key must be a string")
	}
	bucketName, ok := viper.Get("aws_s3_bucket").(string)
	if !ok {
		return "", errors.New("Your aws_s3_bucket must be a string")
	}
	keyPrefix, ok := viper.Get("aws_s3_key_prefix").(string)
	if !ok {
		return "", errors.New("Your aws_s3_key_prefix must be a string")
	}

	// Set up an AWS session with the user's credentials
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2"),
		Credentials: credentials.NewStaticCredentials(
			accessKeyID,
			secretAccessKey,
			"",
		),
	})

	// Create new uploader and specify buffer size (in bytes) to use when buffering
	// data into chunks and sending them as parts to S3 and clean up on error
	uploader := s3manager.NewUploader(sess, func(u *s3manager.Uploader) {
		u.PartSize = 5 * 1024 * 1024 // 5,242,880 bytes or 5.24288 Mb which is the default minimum here
		u.LeavePartsOnError = false  // If an error occurs during upload to s3, clean up & don't leave partial upload there
	})

	// Generate the bucket key using the key prefix, checksum, and tmpFile name
	bucketKey := fmt.Sprintf("%s/%s-%s", keyPrefix, checksum, tmpFile)

	// Obtain FileInfo so we can look at length in bytes
	fileStats, err := file.Stat()
	if err != nil {
		return "", fmt.Errorf("Could not obtain file stats for %s", file.Name())
	}

	// Create and start a new progress bar with a fixed width
	bar := pb.Full.Start64(fileStats.Size()).SetWidth(100)

	// Create a ProxyReader and attach the file and progress bar
	pr := proxyReader.New(file, bar)

	// Upload compressed zip file to s3
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(bucketKey),
		Body:   pr, // As our file is read and uploaded, our proxy reader will update/render the progress bar
	})
	if err != nil {
		return "", fmt.Errorf("Error uploading assets to s3: %v", err)
	}

	return bucketKey, nil
}

// createChecksumFromZip takes a pointer to a file and creates a sha256 checksum
// of the content. We use this for naming the s3 bucket key so that we don't write
// duplicates to s3. The call to io.Copy actually consumes the read position of
// the file to EOF so we call file.Seek and set the read position back to the
// beginning of the file
func createChecksumFromZip(file *os.File) (string, error) {
	// Create a sha256 hash of the curriculum directory
	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	// Make the hash URL safe with base64
	checksum := base64.URLEncoding.EncodeToString(hash.Sum(nil))

	// The io.Copy call for producing the hash consumed the read position of the
	// file (file now at EOF). Need to reset to beginning for sending to s3
	_, err := file.Seek(0, io.SeekStart)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
		return "", err
	}

	return checksum, nil
}

// cleanUpFiles removes the tmp zipfile that was created for uploading to s3. We
// wouldn't want to leave artifacts on user's machines
func cleanUpFiles() {
	err := os.Remove(tmpFile)
	if err != nil {
		fmt.Println("Sorry, we had trouble cleaning up the zip file created for curriculum preview")
	}
}

// compressDirectory takes a source file path (where the content you want zipped lives)
// and a target file path (where to put the zip file) and recursively compresses the source.
// Source can either be a directory or a single file
func compressDirectory(source, target string) error {
	// Create file with target name and defer its closing
	zipfile, err := os.Create(target)
	if err != nil {
		return err
	}
	defer zipfile.Close()

	// Create a new zip writer and pass our zipfile in
	archive := zip.NewWriter(zipfile)
	defer archive.Close()

	// Get os.FileInfo about our source
	info, err := os.Stat(source)
	if err != nil {
		return nil
	}

	// Check to see if the provided source file is a directory and set baseDir if so
	var baseDir string
	if info.IsDir() {
		baseDir = filepath.Base(source)
	}

	// Walk the whole filepath
	filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Creates a partially-populated FileHeader from an os.FileInfo
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		// Check if baseDir has been set (from the IsDir check) and if it has not been
		// set, update the header.Name to reflect the correct path
		if baseDir != "" {
			header.Name = filepath.Join(baseDir, strings.TrimPrefix(path, source))
		}

		// Check if the file we are iterating is a directory and update the header.Name
		// or the header.Method appropriately
		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}

		//  Add a file to the zip archive using the provided FileHeader for the file metadata
		writer, err := archive.CreateHeader(header)
		if err != nil {
			return err
		}

		// Return nil if at this point if info is a directory
		if info.IsDir() {
			return nil
		}

		// If it was not a directory, we open the file and copy it into the archive writer
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		_, err = io.Copy(writer, file)

		return err
	})

	return err
}
