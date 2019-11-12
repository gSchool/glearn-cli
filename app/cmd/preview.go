package cmd

import (
	"archive/zip"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const tmpFile string = "preview-curriculum.zip"

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
		err = uploadToS3(f, checksum)
		if err != nil {
			fmt.Printf("Failed to upload zip file to s3. Err: %v", err)
			os.Exit(1)
			return
		}

		fmt.Println("Sucessfully uploaded your curriculum preview!")
	},
}

func uploadToS3(file *os.File, checksum string) error {
	// Coerce AWS credentials to strings
	accessKeyID, ok := viper.Get("aws_access_key_id").(string)
	if !ok {
		return errors.New("Your aws_access_key_id must be a string")
	}
	secretAccessKey, ok := viper.Get("aws_secret_access_key").(string)
	if !ok {
		return errors.New("Your aws_secret_access_key must be a string")
	}
	bucketName, ok := viper.Get("aws_s3_bucket").(string)
	if !ok {
		return errors.New("Your aws_s3_bucket must be a string")
	}
	keyPrefix, ok := viper.Get("aws_s3_key_prefix").(string)
	if !ok {
		return errors.New("Your aws_s3_key_prefix must be a string")
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
	uploader := s3manager.NewUploader(sess)

	// Upload compressed zip file to s3
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fmt.Sprintf("%s/%s-%s", keyPrefix, checksum, tmpFile)),
		Body:   file,
	})
	if err != nil {
		return fmt.Errorf("Error uploading assets to s3: %v", err)
	}

	return nil
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
