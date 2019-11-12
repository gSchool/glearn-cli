package cmd

import (
	"archive/zip"
	"crypto/sha256"
	"encoding/base64"
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
	Use:   "preview [not_sure_yet]",
	Short: "Preview your content",
	Long:  `Long description for previewing`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Takes one argument which is the filepath to the directory you want zipped/previewed
		if len(args) != 1 {
			fmt.Println("Usage: `learn preview` takes one argument")
			os.Exit(1)
		}

		// Compress directory into tmpFile name
		compressDirectory(args[0], tmpFile)
		// Removes artifacts on user's machine
		defer cleanUpFiles()

		// Open file so we can get a checksum as well as send to s3
		f, err := os.Open(tmpFile)
		if err != nil {
			fmt.Printf("Failed to open file %q, %v", tmpFile, err)
			return
		}
		defer f.Close()

		// Create a sha256 hash of the curriculum directory
		hash := sha256.New()
		if _, err := io.Copy(hash, f); err != nil {
			log.Fatal(err)
			return
		}
		// Make the hash URL safe with base64
		checksum := base64.URLEncoding.EncodeToString(hash.Sum(nil))

		// The io.Copy call for producing the hash consumed the read position of the
		// file (file now at EOF). Need to reset to beginning for sending to s3
		_, err = f.Seek(0, os.SEEK_SET)
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
			return
		}

		// TODO: FIRST CHECK IF THE CREDENTIALS CAN BE COERCED TO A STRING BELOW

		// Set up an AWS session and create an s3 manager uploader
		sess, err := session.NewSession(&aws.Config{
			Region: aws.String("us-west-2"),
			Credentials: credentials.NewStaticCredentials(
				viper.Get("aws_access_key_id").(string),
				viper.Get("aws_secret_access_key").(string),
				"",
			),
		})
		uploader := s3manager.NewUploader(sess)

		// Upload compressed zip file to s3
		_, err = uploader.Upload(&s3manager.UploadInput{
			Bucket: aws.String("BUCKET_NAME_ENV"),
			Key:    aws.String(fmt.Sprintf("KEY_NAME_ENV")),
			Body:   f,
		})
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
			return
		}

		fmt.Println("Sucessfully uploaded your curriculum preview!")
	},
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
