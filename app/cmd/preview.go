package cmd

import (
	"archive/zip"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/briandowns/spinner"
	pb "github.com/cheggaaa/pb/v3"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/gSchool/glearn-cli/api/learn"
	"github.com/gSchool/glearn-cli/mdlinkparser"
	proxyReader "github.com/gSchool/glearn-cli/proxy_reader"
)

// tmpZipFile is used throughout as the temporary zip file target location.
const tmpZipFile string = "preview-curriculum.zip"

// tmpSingleFileDir is used throughout as the temporary single file directory location. This
// is the name of the tmp dir we build when needing to attach relative links.
const tmpSingleFileDir string = "single-file-upload"

// previewCmd is executed when the `learn preview` command is used. Preview's concerns:
// 1. Compress directory/file into target location.
// 2. Defer cleaning up the file after command is finished.
// 3. Create a checksum for the zip file.
// 4. Upload the zip file to s3.
// 5. Notify learn that new content is available for building.
// 6. Handle progress bar for s3 upload.
var previewCmd = &cobra.Command{
	Use:   "preview [file_path]",
	Short: "Uploads content and builds a preview.",
	Long: `
		The preview command takes a path to either a directory or a single file and
		uploads the content to Learn through the Learn API. Learn will build the preview
		and return/open the preview URL when it is complete.
	`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Start benchmarking the total time spent in preview cmd
		startOfCmd := time.Now()

		setupLearnAPI()

		if viper.Get("api_token") == "" || viper.Get("api_token") == nil {
			previewCmdError("Please set your API token first with `learn set --api_token=your_token`")
		}

		// Takes one argument which is the filepath to the directory you want zipped/previewed
		if len(args) != 1 {
			previewCmdError("Usage: `learn preview` takes just one argument")
			return
		}

		// Set target file path
		target := args[0]

		// Get os.FileInfo from call to os.Stat so we can see if it is a single file or directory
		fileInfo, err := os.Stat(target)
		if err != nil {
			previewCmdError(fmt.Sprintf("Failed to get stats on file. Err: %v", err))
			return
		}
		isDirectory := fileInfo.IsDir()

		// If it is a single file preview we need to parse the target for any md link tags
		// linking to local files. If there are any, add them to the target
		var singleFileLinkPaths []string
		if !isDirectory {
			if filepath.Ext(target) == ".md" {
				singleFileLinkPaths, err = collectLinkPaths(target)
				if err != nil {
					previewCmdError(fmt.Sprintf("Failed to attach local images for single file preview for: (%s). Err: %v", target, err))
					return
				}
			} else {
				previewCmdError("Sorry we only support markdown files for single file previews")
				return
			}
		}

		// variable holding whether or not source is a dir OR when it is a single file preview
		// AND singleFileLinkPaths is > 0 that means it is now a dir again (tmp one we created)
		isDirectory = isDirectory || (!isDirectory && len(singleFileLinkPaths) > 0)

		if len(singleFileLinkPaths) > 0 {
			target, err = createNewTarget(target, singleFileLinkPaths)
			if err != nil {
				previewCmdError(fmt.Sprintf("Failed build tmp files around single file preview for: (%s). Err: %v", target, err))
				return
			}
		}

		// Detect config file
		_, err = doesConfigExistOrCreate(target, UnitsDirectory)
		if err != nil {
			previewCmdError(fmt.Sprintf("Failed to find or create a config file for: (%s). Err: %v", target, err))
			return
		}

		// Start a processing spinner that runs until a user's content is compressed
		fmt.Println("Compressing your content...")
		s := spinner.New(spinner.CharSets[26], 100*time.Millisecond)
		s.Color("blue")
		s.Start()

		// Start benchmark for compressDirectory
		startOfCompression := time.Now()

		// Compress directory, output -> tmpZipFile
		err = compressDirectory(target, tmpZipFile)
		if err != nil {
			previewCmdError(fmt.Sprintf("Failed to compress provided directory (%s). Err: %v", target, err))
			return
		}

		// Add benchmark in milliseconds for compressDirectory
		bench := &learn.CLIBenchmark{
			Compression: time.Since(startOfCompression).Milliseconds(),
			CmdName:     "preview",
		}

		// Stop the processing spinner
		s.Stop()
		printlnGreen("√")

		// Removes artifacts on user's machine
		defer removeArtifacts()

		// Open file so we can get a checksum as well as send to s3
		f, err := os.Open(tmpZipFile)
		if err != nil {
			previewCmdError(fmt.Sprintf("Failed opening file (%q). Err: %v", tmpZipFile, err))
			return
		}
		defer f.Close()

		// Create checksum of files in directory
		checksum, err := createChecksumFromZip(f)
		if err != nil {
			previewCmdError(fmt.Sprintf("Failed to create a checksum for compressed file. Err: %v", err))
			return
		}

		// Start benchmark for uploadToS3
		startOfUploadToS3 := time.Now()

		// Send compressed zip file to s3
		bucketKey, err := uploadToS3(f, checksum, learn.API.Credentials)
		if err != nil {
			previewCmdError(fmt.Sprintf("Failed to upload zip file to s3. Err: %v", err))
			return
		}

		// Add benchmark in milliseconds for uploadToS3
		bench.UploadToS3 = time.Since(startOfUploadToS3).Milliseconds()

		fmt.Println("\nPlease wait while Learn builds your preview...")

		// Start a processing spinner that runs until Learn is finsihed building the preview
		s = spinner.New(spinner.CharSets[32], 100*time.Millisecond)
		s.Color("blue")
		s.Start()

		// Start benchmark for BuildReleaseFromS3 & PollForBuildResponse (Learn build stage)
		startBuildAndPollRelease := time.Now()

		// Let Learn know there is new preview content on s3, where it is, and to build it
		res, err := learn.API.BuildReleaseFromS3(bucketKey, isDirectory)
		if err != nil {
			previewCmdError(fmt.Sprintf("Failed to notify learn of new preview content. Err: %v", err))
			return
		}

		// If content is a directory, rewrite the res from polling for build response. Directories
		// can take much longer to build, however single files build instantly so we do not need to
		// poll for them because the call to BuildReleaseFromS3 will get a preview_url right away
		if isDirectory {
			var attempts uint8 = 30
			res, err = learn.API.PollForBuildResponse(res.ReleaseID, &attempts)
			if err != nil {
				previewCmdError(fmt.Sprintf("Failed to poll Learn for your new preview build. Err: %v", err))
				return
			}
		}

		// Add benchmark in milliseconds for the Learn build stage and total time in preview cmd
		bench.LearnBuild = time.Since(startBuildAndPollRelease).Milliseconds()
		bench.TotalCmdTime = time.Since(startOfCmd).Milliseconds()

		// Set final message for dislpay
		s.FinalMSG = fmt.Sprintf("Sucessfully uploaded your preview! You can find your content at: %s\n", res.PreviewURL)

		// Stop the processing spinner
		s.Stop()
		printlnGreen("√")

		exec.Command("bash", "-c", fmt.Sprintf("open %s", res.PreviewURL)).Output()

		err = learn.API.SendMetadataToLearn(&learn.CLIBenchmarkPayload{
			CLIBenchmark: bench,
		})
		if err != nil {
			removeArtifacts()
			learn.API.NotifySlack(err)
			os.Exit(1)
		}
	},
}

// createNewTarget will set up and create everything needed for single file previews if they are needed.
// Returns a string representing the source name which if not single file tmp dir is needed, will return the
// original
func createNewTarget(target string, singleFileLinkPaths []string) (string, error) {
	// Tmp dir so we can build out a new dir with the correct links in their correct
	// paths based on relative link paths supplied in the single markdown file
	newSrcPath := tmpSingleFileDir

	// Get the name of the single target file
	srcArray := strings.Split(target, "/")
	srcMDFile := srcArray[len(srcArray)-1]

	for _, imgPath := range singleFileLinkPaths {
		if !strings.HasPrefix(imgPath, "/") {
			imgPath = fmt.Sprintf("/%s", imgPath)
		}

		// Ex. images/something-else/my_neat_image.png -> ["images", "something-else", "my_neat_image.png"]
		pathArray := strings.Split(imgPath, "/")
		imageName := pathArray[len(pathArray)-1] // -> "my_neat_image.png"

		// create an linkDirs var and depending on how long the image file path is, update it to include
		// everything up to the image itself
		var linkDirs string

		if len(pathArray) == 1 {
			linkDirs = ""
		} else if len(pathArray) == 2 {
			linkDirs = pathArray[0]
		} else {
			// Collect verything up until the image name (last item) and join it back together
			// This gives us the name of the directory(ies) to make to put the image in
			linkDirs = strings.Join(pathArray[:len(pathArray)-1], "/")
		}

		// Create appropriate directory for each link using the linkDirs
		err := os.MkdirAll(newSrcPath+linkDirs, os.FileMode(0777))
		if err != nil {
			return "", err
		}

		// Get "oneDirBackFromTarget" because target will be an .md file with relative
		// links so we need to go one back from "target" so things aren't trying
		// to be nested in the .md file itself
		targetArray := strings.Split(target, "/")
		oneDirBackFromTarget := strings.Join(targetArray[:len(targetArray)-1], "/")

		sourceLinkPath := oneDirBackFromTarget + imgPath
		if _, err := os.Stat(sourceLinkPath); os.IsNotExist(err) {
			fmt.Printf("Link not found with path '%s'\n", sourceLinkPath)
		} else {
			// Copy the actual image into our new temp directory in it's appropriate spot
			err = Copy(sourceLinkPath, newSrcPath+linkDirs+"/"+imageName)
			if err != nil {
				return "", err
			}
		}
	}

	// Copy original single markdown file into the base of our new tmp dir
	err := Copy(target, newSrcPath+"/"+srcMDFile)
	if err != nil {
		return "", err
	}

	if len(singleFileLinkPaths) > 0 {
		return newSrcPath, nil
	}

	return target, nil
}

// previewCmdError is a small wrapper for all errors within the preview command. It ensures
// artifacts are cleaned up with a call to removeArtifacts
func previewCmdError(msg string) {
	fmt.Println(msg)
	removeArtifacts()
	learn.API.NotifySlack(errors.New(msg))
	os.Exit(1)
}

// printlnGreen simply prints a green string
func printlnGreen(text string) {
	fmt.Printf("\033[32m%s\033[0m\n", text)
}

// collectLinkPaths takes a target, reads it, and passes it's contents (slice of bytes)
// to our MDLinkParser as a string. All relative/local markdown flavored images are parsed
// into an array of strings and returned
func collectLinkPaths(target string) ([]string, error) {
	contents, err := ioutil.ReadFile(target)
	if err != nil {
		return []string{}, fmt.Errorf("Failure to read file '%s'. Err: %s", string(contents), err)
	}

	m := mdlinkparser.New(string(contents))
	m.ParseLinks()

	return m.Links, nil
}

// uploadToS3 takes a file and it's checksum and uploads it to s3 in the appropriate bucket/key
func uploadToS3(file *os.File, checksum string, creds *learn.Credentials) (string, error) {
	// Set up an AWS session with the user's credentials
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2"),
		Credentials: credentials.NewStaticCredentials(
			creds.AccessKeyID,
			creds.SecretAccessKey,
			"",
		),
	})

	// Create new uploader and specify buffer size (in bytes) to use when buffering
	// data into chunks and sending them as parts to S3 and clean up on error
	uploader := s3manager.NewUploader(sess, func(u *s3manager.Uploader) {
		u.PartSize = 5 * 1024 * 1024 // 5,242,880 bytes or 5.24288 Mb which is the default minimum here
		u.LeavePartsOnError = false  // If an error occurs during upload to s3, clean up & don't leave partial upload there
	})

	// Generate the bucket key using the key prefix, checksum, and tmpZipFile name
	bucketKey := fmt.Sprintf("%s/%s-%s", creds.KeyPrefix, checksum, tmpZipFile)

	// Obtain FileInfo so we can look at length in bytes
	fileStats, err := file.Stat()
	if err != nil {
		return "", fmt.Errorf("Could not obtain file stats for %s", file.Name())
	}

	// Create and start a new progress bar with a fixed width
	bar := pb.Full.Start64(fileStats.Size()).SetWidth(100)

	// Create a ProxyReader and attach the file and progress bar
	pr := proxyReader.New(file, bar)

	fmt.Println("Uploading assets to Learn...")

	// Upload compressed zip file to s3
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(creds.BucketName),
		Key:    aws.String(bucketKey),
		Body:   pr, // As our file is read and uploaded, our proxy reader will update/render the progress bar
	})
	if err != nil {
		return "", fmt.Errorf("Error uploading assets to s3: %v", err)
	}

	bar.Finish()
	printlnGreen("√")

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

// removeArtifacts removes the tmp zipfile and the tmp single file preview directory (if there
// was one) that were created for uploading to s3 and including local images. We wouldn't
// want to leave artifacts on user's machines
func removeArtifacts() {
	err := os.Remove(tmpZipFile)
	if err != nil {
		fmt.Println("Sorry, we had trouble cleaning up the zip file created for curriculum preview")
	}

	// Remove tmpSingleFileDir if it exists at this point
	if _, err := os.Stat(tmpSingleFileDir); !os.IsNotExist(err) {
		err = os.RemoveAll(tmpSingleFileDir)
		if err != nil {
			fmt.Println("Sorry, we had trouble cleaning up the tmp single file preview directory")
		}
	}
}

// Copy the src file to target dest. Any existing file will be overwritten and will not copy file attributes.
func Copy(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
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
		ext := filepath.Ext(path)
		_, ok := fileExtWhitelist[ext]

		if ok || (info.IsDir() && (ext != ".git" && path != "node_modules")) {
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
			// ingore zip files
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()
			_, err = io.Copy(writer, file)
		}

		return err
	})

	return err
}
