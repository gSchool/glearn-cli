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
	"path/filepath"
	"strings"
	"time"

	proxyReader "github.com/Galvanize-IT/glearn-cli/app/proxy_reader"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/cheggaaa/pb/v3"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const tmpFile string = "preview-curriculum.zip"

type LearnPreviewResponse struct {
	Url string `json:"url"`
}

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

		res, err := notifyLearn(bucketKey)
		if err != nil {
			fmt.Printf("Failed to notify learn of new preview content. Err: %v", err)
			os.Exit(1)
			return
		}

		// POST to /api/v1/releases create method just need s3 key
		// Need to add functionality to endpoing on learn to take some sort of param saying I am a single content file or a full release
		// Would be great to ping an endpoint seeing when the build is complete and return the url to go to

		// Maybes?
		// Potentially make stage bucket a flag or something instead of changing env vars every time?
		// Possibly a loader bar for the zip upload? Maybe something like heroku's little animation showing things are processing

		fmt.Printf("Sucessfully uploaded your preview! You can find your content at: %s", res.Url)
	},
}

func notifyLearn(bucketKey string) (*LearnPreviewResponse, error) {
	apiToken, ok := viper.Get("api_token").(string)
	if !ok {
		return nil, errors.New("Please set your api_token in ~/.glearn-config.yaml")
	}

	payload := map[string]string{
		"bucket_key_name": bucketKey,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	client := &http.Client{Timeout: time.Second * 10}

	req, err := http.NewRequest("POST", "https://httpbin.org/post", bytes.NewBuffer(payloadBytes))
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

	// Set up an AWS session and create an s3 manager uploader
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
		u.PartSize = 5 * 1024 * 1024
		u.LeavePartsOnError = false
	})

	// Generate the bucket key using the key prefix, checksum, and tmpFile name
	bucketKey := fmt.Sprintf("%s/%s-%s", keyPrefix, checksum, tmpFile)

	// Obtain FileInfo so we can look at length in bytes
	fileStats, err := file.Stat()
	if err != nil {
		return "", fmt.Errorf("Could not obtain file stats for %s", file.Name())
	}

	// Create and start a new progress bar
	bar := pb.Full.Start64(fileStats.Size())

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

func cleanUpFiles() {
	err := os.Remove(tmpFile)
	if err != nil {
		fmt.Println("Sorry, we had trouble cleaning up the zip file created for curriculum preview")
	}
}

func compressDirectory(source, target string) error {
	zipfile, err := os.Create(target)
	if err != nil {
		return err
	}
	defer zipfile.Close()

	archive := zip.NewWriter(zipfile)
	defer archive.Close()

	info, err := os.Stat(source)
	if err != nil {
		return nil
	}

	var baseDir string
	if info.IsDir() {
		baseDir = filepath.Base(source)
	}

	filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		if baseDir != "" {
			header.Name = filepath.Join(baseDir, strings.TrimPrefix(path, source))
		}

		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}

		writer, err := archive.CreateHeader(header)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

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
