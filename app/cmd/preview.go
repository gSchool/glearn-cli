package cmd

import (
	"archive/zip"
	"bufio"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
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
	di "github.com/gSchool/glearn-cli/ignorematcher"
	"github.com/gSchool/glearn-cli/mdresourceparser"
	proxyReader "github.com/gSchool/glearn-cli/proxy_reader"
)

// tmpSingleFileDir is used throughout as the temporary single file directory location. This
// is the name of the tmp dir we build when needing to attach relative links.
const tmpSingleFileDir string = "single-file-upload"

type previewBuilder struct {
	target              string
	fileInfo            os.FileInfo
	singleFileLinkPaths []string
	dataPaths           []string
	dockerPaths         []string
	configYamlPaths     []string
	startOfCmd          time.Time
	bench               *learn.CLIBenchmark
}

func NewPreviewBuilder(args []string) (*previewBuilder, error) {
	setupLearnAPI()

	if viper.Get("api_token") == "" || viper.Get("api_token") == nil {
		return &previewBuilder{}, fmt.Errorf(setAPITokenMessage)
	}

	fileInfo, err := os.Stat(args[0])
	if err != nil {
		return &previewBuilder{}, fmt.Errorf("Failed to get stats on file. Err: %v", err)
	}
	p := &previewBuilder{
		target:              args[0],
		fileInfo:            fileInfo,
		singleFileLinkPaths: []string{},
		dataPaths:           []string{},
		dockerPaths:         []string{},
		configYamlPaths:     []string{},
		startOfCmd:          time.Now(),
	}

	if p.invalidPreviewTarget() {
		return &previewBuilder{}, fmt.Errorf("The preview file that you chose is not able to be rendered as a single file preview in learn")
	}

	return p, nil
}

// collectPaths reads from a file and collects required docker paths, file link paths, and other resources needed for preview
func (p *previewBuilder) collectPaths() error {
	// collect nothing if we do not include links
	if !p.includesLinks() || filepath.Ext(p.target) != ".md" {
		return nil
	}

	// TODO collectDataPaths and collectResourcePaths parse the file twice, should be able to do this in one pass
	dataPaths, err := collectDataPaths(p.target)
	if err != nil {
		return fmt.Errorf("Failed to attach local images for single file preview for: (%s). Err: %v", p.target, err)
	}
	singleFileLinkPaths, dockerPaths, err := collectResourcePaths(p.target)
	if err != nil {
		return fmt.Errorf("Failed to attach local images for single file preview for: (%s). Err: %v", p.target, err)
	}

	p.dataPaths = dataPaths
	p.singleFileLinkPaths = singleFileLinkPaths
	p.dockerPaths = dockerPaths

	return nil
}

func (p *previewBuilder) buildAlternateTarget() error {
	paths := append(p.singleFileLinkPaths, p.dataPaths...)
	alternateTarget, err := createNewTarget(p.target, paths, p.dockerPaths)
	if err != nil {
		return err
	}

	// reset target if there has been a new one created
	if alternateTarget != "" {
		p.target = alternateTarget
	}
	return nil
}

// setConfigYaml finds or creates a config.yaml file for the preview environment. The paths on the file are read and set.
func (p *previewBuilder) setConfigYaml() error {
	_, err := previewFindOrCreateConfig(p.target, p.isSingleFilePreview(), p.dockerPaths)
	if err != nil {
		return fmt.Errorf("Failed to find or create a config file for: (%s).\nErr: %v", p.target, err)
	}

	configYamlPaths, err := parseConfigAndGatherLinkedPaths(p.target)
	if err != nil {
		return fmt.Errorf("Failed to parse config/autoconfig yaml for: (%s).\nErr: %v", p.target, err)
	}
	p.configYamlPaths = configYamlPaths

	return nil
}

// compressDirectory takes a source file path (where the content you want zipped lives)
// and a target file path (where to put the zip file) and recursively compresses the source.
// Source can either be a directory or a single file. When singleFile is true, all files in
// the zip are added. resourcePaths specify non-config paths which should be included in the zip.
func (p *previewBuilder) compressDirectory(zipTarget string) error {
	// Start a processing spinner that runs until a user's content is compressed
	fmt.Println("Compressing your content...")
	zipSpinner := spinner.New(spinner.CharSets[26], 100*time.Millisecond)
	zipSpinner.Color("blue")
	zipSpinner.Start()

	// Start benchmark for compressDirectory
	startOfCompression := time.Now()

	resourcePaths := append(p.dockerPaths, p.dataPaths...)
	// Create file with zipTarget name and defer its closing
	zipfile, err := os.Create(zipTarget)
	if err != nil {
		return err
	}
	defer zipfile.Close()

	// Create a new zip writer and pass our zipfile in
	archive := zip.NewWriter(zipfile)
	defer archive.Close()

	// Get os.FileInfo about our source
	info, err := os.Stat(p.target)
	if err != nil {
		return nil
	}

	// Check to see if the provided source file is a directory and set baseDir if so
	var baseDir string
	if info.IsDir() {
		baseDir = filepath.Base(p.target)
	}

	// Walk the whole filepath
	filepath.Walk(p.target, func(path string, info os.FileInfo, err error) error {
		path = filepath.ToSlash(path)

		fileIsIncluded := false
		for _, p := range p.configYamlPaths {
			var configPathSplits = strings.Split(p, string(os.PathSeparator))
			var fileName = configPathSplits[len(configPathSplits)-1]
			if strings.Contains(path, fileName) {
				fileIsIncluded = true
			}
		}
		for _, d := range resourcePaths {
			if strings.Contains(path, d) {
				fileIsIncluded = true
			}
		}
		if len(p.configYamlPaths) == 0 && !info.IsDir() {
			// This accounts for the single file preview which won't have yaml files and won't be a directory
			fileIsIncluded = true
		}

		var isConfigFile = strings.Contains(path, "config.yml") || strings.Contains(path, "config.yaml") || strings.Contains(path, "autoconfig.yaml")
		ext := filepath.Ext(path)
		// Ignoring all files over 1mb for preivew and warning users if the file is over 20mb that it will be ignored in publish action as well.
		if !info.IsDir() && info.Size() > 1000000 && !strings.Contains(path, ".git/") {
			if path == "preview-curriculum.zip" { // don't warn on preview-curriculum, it gets read here but still cleaned up
				return nil
			}
			if info.Size() > 20000000 {
				fmt.Printf("\nWARNING: Ignoring File For Preview: File chosen/linked is too large to preview and too large to publish: %s\n", path)
			} else {
				fmt.Printf("\nWARNING: Ignoring File For Preview: File chosen/linked is too large to preview, but will successfully publish: %s\n", path)
			}
			return nil
		}

		if isConfigFile || fileIsIncluded || (info.IsDir() && !strings.Contains(path, ".git/") && (ext != ".git" && path != "node_modules")) {
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
				header.Name = filepath.Join(baseDir, strings.TrimPrefix(path, p.target))
			}

			// Check if the file we are iterating is a directory and update the header.Name
			// or the header.Method appropriately
			if info.IsDir() {
				header.Name += string(os.PathSeparator)
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
			// ignore zip files
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			_, err = io.Copy(writer, file)
			if err != nil {
				return err
			}
		}

		return err
	})

	p.bench = &learn.CLIBenchmark{
		Compression: time.Since(startOfCompression).Milliseconds(),
		CmdName:     "preview",
	}

	zipSpinner.Stop()
	printlnGreen("√")

	return err
}

// uploadZip is responsible for taking a compressed preview directory and uploading it to be built by Learn.
func (p *previewBuilder) uploadZip(tmpZipFile string) (bucketKey string, err error) {
	f, err := os.Open(tmpZipFile)
	if err != nil {
		return bucketKey, fmt.Errorf("Failed opening file (%q). Err: %v", tmpZipFile, err)
	}
	defer f.Close()

	// Create checksum of files in directory
	checksum, err := createChecksumFromZip(f)
	if err != nil {
		return bucketKey, fmt.Errorf("Failed to create a checksum for compressed file. Err: %v", err)
	}

	// Start benchmark for uploadToS3
	startOfUploadToS3 := time.Now()

	// Send compressed zip file to s3
	bucketKey, err = uploadToS3(tmpZipFile, f, checksum, learn.API.Credentials)
	if err != nil {
		return bucketKey, fmt.Errorf("Failed to upload zip file to s3. Err: %v", err)
	}

	// Add benchmark in milliseconds for uploadToS3
	p.bench.UploadToS3 = time.Since(startOfUploadToS3).Milliseconds()
	return bucketKey, nil
}

// buildLearnPreview triggers the Learn preview building process and montiors its completion via polling
func (p *previewBuilder) buildLearnPreview(bucketKey string) error {
	fmt.Println("\nBuilding preview...")

	// Start a processing spinner that runs until Learn is finished building the preview
	s := spinner.New(spinner.CharSets[32], 100*time.Millisecond)
	s.Color("blue")
	s.Start()

	// Start benchmark for BuildReleaseFromS3 & PollForBuildResponse (Learn build stage)
	startBuildAndPollRelease := time.Now()

	// Let Learn know there is new preview content on s3, where it is, and to build it
	res, err := learn.API.BuildReleaseFromS3(bucketKey, (p.isDirectory() || p.fileContainsSQLPaths() || p.fileContainsDocker()))
	if err != nil {
		return fmt.Errorf("Failed to build new preview content in learn. Err: %v", err)
	}

	// If content is a directory, rewrite the res from polling for build response. Directories
	// can take much longer to build, however single files build instantly so we do not need to
	// poll for them because the call to BuildReleaseFromS3 will get a preview_url right away
	if p.isDirectory() || p.fileContainsSQLPaths() || p.fileContainsDocker() {
		var attempts uint8 = 30
		res, err = learn.API.PollForBuildResponse(res.ReleaseID, p.fileInfo.IsDir(), p.fileInfo.Name(), &attempts)
		if err != nil {
			return fmt.Errorf("Failed to poll Learn for your new preview build. Err: %v", err)
		}
	}

	// Add benchmark in milliseconds for the Learn build stage and total time in preview cmd
	p.bench.LearnBuild = time.Since(startBuildAndPollRelease).Milliseconds()
	p.bench.TotalCmdTime = time.Since(p.startOfCmd).Milliseconds()

	// Set final message for display
	s.FinalMSG = fmt.Sprintf("Successfully uploaded your preview! You can find your content at: %s\n", res.PreviewURL)

	// Stop the processing spinner
	s.Stop()
	printlnGreen("√")

	if OpenPreview {
		openURL(res.PreviewURL)
	}

	return nil
}

// invalidForPreview requires that non-dierctory targets be of markdown or ipynb type
func (p *previewBuilder) invalidPreviewTarget() bool {
	return !p.fileInfo.IsDir() && (!strings.HasSuffix(p.target, ".md") && !strings.HasSuffix(p.target, ".ipynb"))
}

// includesLinks reports if the target file info is not a directory and not overidden by the FileOnly flag
func (p *previewBuilder) includesLinks() bool {
	return !p.fileInfo.IsDir() && !FileOnly
}

func (p *previewBuilder) fileContainsLinks() bool {
	return len(p.singleFileLinkPaths) > 0
}

func (p *previewBuilder) fileContainsSQLPaths() bool {
	return len(p.dataPaths) > 0
}

func (p *previewBuilder) fileContainsDocker() bool {
	return len(p.dockerPaths) > 0
}

func (p *previewBuilder) containsAnyResources() bool {
	return p.fileContainsLinks() || p.fileContainsSQLPaths() || p.fileContainsDocker()
}

func (p *previewBuilder) isSingleFilePreview() bool {
	return !p.isDirectory() && p.containsAnyResources()
}

// isDirectory reports if either the target is a directory, or if the target contains resources and a directory is required for content
func (p *previewBuilder) isDirectory() bool {
	return p.fileInfo.IsDir() || (!p.fileInfo.IsDir() && (p.fileContainsLinks() || p.fileContainsDocker()))
}

// previewCmd is executed when the `learn preview` command is used. Preview's concerns:
// 1. Compress directory/file into target location.
// 2. Defer cleaning up the file after command is finished.
// 3. Create a checksum for the zip file.
// 4. Upload the zip file to s3.
// 5. Notify learn that new content is available for building.
// 6. Handle progress bar for s3 upload.
var previewCmd = &cobra.Command{
	Use:   "preview [options] <directory|file_path>",
	Short: "Uploads content and builds a preview.",
	Long: `
The preview command takes a path to either a directory or a single file and
uploads the content to Learn through the Learn API. Learn will build the
preview and return/open the preview URL when it is complete.
	`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		tmpZipFile := "preview-curriculum.zip"
		previewer, err := NewPreviewBuilder(args)
		if err != nil {
			previewCmdError(fmt.Sprintf("%v", err), tmpZipFile)
			return
		}

		err = previewer.collectPaths()
		if err != nil {
			previewCmdError(fmt.Sprintf("%v", err), tmpZipFile)
			return
		}

		if previewer.containsAnyResources() {
			err = previewer.buildAlternateTarget()
			if err != nil {
				previewCmdError(fmt.Sprintf("%v", err), tmpZipFile)
				return
			}
		}

		if previewer.containsAnyResources() || previewer.isDirectory() {
			err = previewer.setConfigYaml()
			if err != nil {
				previewCmdError(fmt.Sprintf("%v", err), tmpZipFile)
				return
			}
		}

		err = previewer.compressDirectory(tmpZipFile)
		if err != nil {
			previewCmdError(fmt.Sprintf("Failed to compress provided directory (%s). Err: %v", previewer.target, err), tmpZipFile)
			return
		}

		// Removes artifacts on user's machine
		defer removeArtifacts(tmpZipFile)

		bucketKey, err := previewer.uploadZip(tmpZipFile)
		if err != nil {
			previewCmdError(fmt.Sprintf("%v", err), tmpZipFile)
			return
		}

		err = previewer.buildLearnPreview(bucketKey)
		if err != nil {
			previewCmdError(fmt.Sprintf("%v", err), tmpZipFile)
			return
		}

		err = learn.API.SendMetadataToLearn(&learn.CLIBenchmarkPayload{
			CLIBenchmark: previewer.bench,
		})

		if err != nil {
			removeArtifacts(tmpZipFile)
			learn.API.NotifySlack(err)
			os.Exit(1)
		}
	},
}

// createNewTarget will set up and create everything needed for single file previews if they are needed.
// Returns a string representing the source name which if not single file tmp dir is needed, will return the
// original
func createNewTarget(target string, singleFilePaths, dockerPaths []string) (string, error) {
	// Tmp dir so we can build out a new dir with the correct links in their correct
	// paths based on relative link paths supplied in the single markdown file
	newSrcPath := tmpSingleFileDir

	// Get the name of the single target file
	srcArray := strings.Split(target, "/")
	srcMDFile := srcArray[len(srcArray)-1]
	substringPaths := []string{}

	for _, filePath := range singleFilePaths {
		if !strings.HasPrefix(filePath, "/") {
			filePath = fmt.Sprintf("/%s", filePath)
		}

		// Ex. images/something-else/my_neat_image.png -> ["images", "something-else", "my_neat_image.png"]
		pathArray := []string{}
		var containsPeriodPeriod bool
		for _, dir := range strings.Split(filePath, "/") {
			if dir != ".." { // sanitize any .. so we don't have to worry about nested things
				containsPeriodPeriod = true
				pathArray = append(pathArray, dir)
			}
		}
		// We need to modify the actual markdown file so it no longer has `..` in links, since we're putting
		// everything in the newSrcPath
		if containsPeriodPeriod {
			substringPaths = append(substringPaths, filePath)
		}

		imageName := pathArray[len(pathArray)-1] // -> "my_neat_image.png"

		// create an linkDirs var and depending on how long the image file path is, update it to include
		// everything up to the image itself
		var linkDirs string

		if len(pathArray) == 1 {
			linkDirs = ""
		} else if len(pathArray) == 2 {
			linkDirs = pathArray[0]
		} else {
			// Collect everything up until the image name (last item) and join it back together
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
		sourceLinkPath := trimFirstRune(filePath)
		if len(targetArray[:len(targetArray)-1]) != 0 && !strings.HasSuffix(sourceLinkPath, ".sql") {
			oneDirBackFromTarget := strings.Join(targetArray[:len(targetArray)-1], "/")
			sourceLinkPath = oneDirBackFromTarget + filePath
		}

		if _, err := os.Stat(sourceLinkPath); os.IsNotExist(err) {
			if strings.HasSuffix(sourceLinkPath, ".sql") {
				useThisPath := ""
				parent := "../" + sourceLinkPath
				for i := 1; i <= 5; i++ {
					_, parentExists := os.Stat(parent)
					if parentExists == nil {
						useThisPath = parent
						break
					} else {
						parent = "../" + parent
					}
				}

				if useThisPath != "" {
					err = Copy(useThisPath, newSrcPath+linkDirs+"/"+imageName)
					if err != nil {
						return "", err
					}
				}
			} else {
				log.Printf("Link not found with path '%s'\n", sourceLinkPath)
			}
		} else {
			// Copy the actual image into our new temp directory in it's appropriate spot
			err = Copy(sourceLinkPath, newSrcPath+linkDirs+"/"+imageName)
			if err != nil {
				return "", err
			}
		}
	} // End of loop over files

	// iterate over docker directories as their contents must be recursively copied
	for _, dirPath := range dockerPaths {
		fmt.Printf("Including docker_directory_path: %s\n", dirPath)
		dirPath = trimFirstRune(dirPath)
		fileDir, err := os.Stat(dirPath)

		// when the directory does not exist, keep moving back in the directory structure until it is found
		if os.IsNotExist(err) {
			newDirPath := ""
			parent := "../" + dirPath
			for i := 1; i <= 5; i++ {
				f, parentExists := os.Stat(parent)
				if parentExists == nil && f.IsDir() {
					newDirPath = parent
					break
				} else if parentExists == nil && !f.IsDir() {
					return "", fmt.Errorf("docker_directory_path %s is not a directory", dirPath)
				} else {
					parent = "../" + parent
				}
			}

			if newDirPath != "" {
				// the directory was found after checkpoint parents, copy contents
				ignorePatterns, err := DockerIgnorePatterns(newDirPath)
				if err != nil {
					fmt.Print(err.Error())
				}

				err = CopyDirectoryContents(newDirPath, newSrcPath+"/"+dirPath, ignorePatterns)
				if err != nil {
					return "", err
				}
			}

		} else if !fileDir.IsDir() {
			return "", fmt.Errorf("docker_directory_path %s is not a directory", dirPath)
		} else {
			ignorePatterns, err := DockerIgnorePatterns(dirPath)
			if err != nil {
				fmt.Print(err.Error())
			}

			err = CopyDirectoryContents(dirPath, newSrcPath+"/"+dirPath, ignorePatterns)
			if err != nil {
				return "", err
			}
		}
	} // End docker path loop

	// Copy original single markdown file into the base of our new tmp dir
	newTarget := newSrcPath + "/" + srcMDFile
	err := Copy(target, newTarget)
	if err != nil {
		return "", err
	}

	if len(singleFilePaths) > 0 || len(dockerPaths) > 0 {
		if len(substringPaths) > 0 {
			// open contents of new target
			b, err := ioutil.ReadFile(newTarget)
			if err != nil {
				return "", fmt.Errorf("Could not read copied target file: %s", err)
			}
			contents := string(b)
			for _, pathToSub := range substringPaths {
				pathToSub = strings.TrimPrefix(pathToSub, "/")
				imgPathWithoutPeriodPeriod := strings.Replace(pathToSub, "../", "", -1)
				contents = strings.ReplaceAll(contents, pathToSub, imgPathWithoutPeriodPeriod)
			}
			// overwrite target file with the contents
			err = ioutil.WriteFile(newTarget, []byte(contents), 0777)
			if err != nil {
				return "", fmt.Errorf("Could not write copied target file with cleaned up links: %s", err)
			}
		}
		return newSrcPath, nil
	}

	return target, nil
}

// previewCmdError is a small wrapper for all errors within the preview command. It ensures
// artifacts are cleaned up with a call to removeArtifacts
func previewCmdError(msg, tmpZipFile string) {
	fmt.Println(msg)
	removeArtifacts(tmpZipFile)
	learn.API.NotifySlack(errors.New(msg))
	os.Exit(1)
}

// printlnGreen simply prints a green string
func printlnGreen(text string) {
	fmt.Printf("\033[32m%s\033[0m\n", text)
}

// collectResourcePaths takes a target, reads it, and passes it's contents (slice of bytes)
// to our MDResourceParser as a string. All relative/local markdown flavored images are parsed
// into an array of strings and returned
func collectResourcePaths(target string) ([]string, []string, error) {
	contents, err := ioutil.ReadFile(target)
	if err != nil {
		return []string{}, []string{}, fmt.Errorf("Failure to read file '%s'. Err: %s", string(contents), err)
	}

	m := mdresourceparser.New([]rune(string(contents)))
	m.ParseResources()

	uniqueMap := make(map[string]struct{})
	for _, v := range m.DockerDirectoryPaths {
		uniqueMap[v] = struct{}{}
	}
	var uniqueDockerPaths []string
	for v := range uniqueMap {
		uniqueDockerPaths = append(uniqueDockerPaths, v)
	}

	return m.Links, uniqueDockerPaths, nil
}

// collectDataPaths takes a target, reads it, and scans the file for data paths and collects them
func collectDataPaths(target string) ([]string, error) {
	contents, err := ioutil.ReadFile(target)
	if err != nil {
		return []string{}, fmt.Errorf("Failure to read file '%s'. Err: %s", string(contents), err)
	}

	if strings.Contains(string(contents), "* data_path: ") {
		file, err := os.Open(target)
		if err != nil {
			return []string{}, fmt.Errorf("Failure to read file '%s'. Err: %s", string(contents), err)
		}
		defer file.Close()

		dataPathsMap := map[string]string{}
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			if strings.Contains(scanner.Text(), "* data_path: ") {
				s := strings.Replace(scanner.Text(), "* data_path: ", "", 1)
				dataPathsMap[s] = s
			}
		}
		ret := []string{}
		for k := range dataPathsMap {
			ret = append(ret, k)
		}

		return ret, nil
	}

	return []string{}, nil
}

// uploadToS3 takes a file and it's checksum and uploads it to s3 in the appropriate bucket/key
func uploadToS3(tmpZipFile string, file *os.File, checksum string, creds *learn.Credentials) (string, error) {
	// Set up an AWS session with the user's credentials
	region := "us-west-2"
	alternateRegion := os.Getenv("S3_REGION")
	if alternateRegion != "" {
		region = alternateRegion
	}
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String(region),
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
func removeArtifacts(tmpZipFile string) {
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

func DockerIgnorePatterns(src string) ([]string, error) {
	dockerIgnore := src + ".dockerignore"
	if !strings.HasSuffix(src, "/") {
		dockerIgnore = src + "/.dockerignore"
	}
	_, err := os.Stat(dockerIgnore)
	if err != nil && os.IsNotExist(err) {
		return []string{}, nil
	}

	ignoreFile, err := ioutil.ReadFile(dockerIgnore)
	if err != nil {
		return []string{}, fmt.Errorf("Could not parse dockerignore file: %s\n", err)
	}

	return strings.Split(string(ignoreFile), "\n"), nil
}

func CopyDirectoryContents(src, dst string, ignorePatterns []string) error {
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}
	if !srcInfo.IsDir() {
		return fmt.Errorf("path specified is not a directory: %s\n", src)
	}

	err = os.MkdirAll(dst, srcInfo.Mode())
	if err != nil {
		return err
	}

	files, err := ioutil.ReadDir(src)
	if err != nil {
		return err
	}
	var alwaysAllowSlice = []string{"test.sh", "Dockerfile", "docker-compose.yaml", "docker-compose.yml"}
	for _, file := range files {
		source := filepath.Join(src, file.Name())
		destination := filepath.Join(dst, file.Name())
		ignore := false
		for _, pattern := range ignorePatterns {
			var alwaysAllowBool = false
			for _, aa := range alwaysAllowSlice {
				if aa == pattern {
					alwaysAllowBool = true
				}
			}
			if alwaysAllowBool {
				continue
			}
			localizedPattern := src + "/" + pattern
			matched, err := di.IgnoreMatches(localizedPattern, source)
			if err != nil {
				fmt.Printf("error while parsing at: %s", err)
			}
			if matched {
				ignore = matched
				break
			}
		}
		if ignore {
			continue
		}

		if file.IsDir() {
			err = CopyDirectoryContents(source, destination, ignorePatterns)
			if err != nil {
				return err
			}
		} else {
			err = Copy(source, destination)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func trimFirstRune(s string) string {
	for i := range s {
		if i > 0 {
			// The value i is the index in s of the second
			// rune.  Slice to remove the first rune.
			return s[i:]
		}
	}
	// There are 0 or 1 runes in the string.
	return ""
}
