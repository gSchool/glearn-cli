package cmd

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

const autoComment = `# This file is auto-generated and orders your content based on the file structure of your repo.
# Do not edit this file; it will be replaced the next time you run the preview command.

# To manually order the contents of this curriculum rather than using the auto-generated file,
# include a config.yaml in your repo following the same conventions as this auto-generated file.
# A user-created config.yaml will have priority over the auto-generated one.

`

var validContentFileAttrs = []string{"Type", "UID", "DefaultVisibility", "MaxCheckpointSubmissions", "EmailOnCompletion", "TimeLimit", "Autoscore"}

type ConfigBuilder struct {
	ConfigYaml          ConfigYaml
	target              string
	isSingleFilePreview bool
	publishContext      bool
	excludePaths        []string
	blockRoot           string
	unitsDir            string
	unitsDirName        string
	unitsRootDirName    string
}

// Note: struct fields must be public in order for unmarshal to
// correctly populate the data.
type ConfigYaml struct {
	Standards []Standard `yaml:"Standards"`
}

type Standard struct {
	Title           string             `yaml:"Title"`
	UID             string             `yaml:"UID"`
	Description     string             `yaml:"Description"`
	SuccessCriteria []string           `yaml:"SuccessCriteria,omitempty"`
	ContentFiles    []ContentFileAttrs `yaml:"ContentFiles"`
}

type ContentFileAttrs struct {
	Type                     string `yaml:"Type"`
	Path                     string `yaml:"Path"`
	UID                      string `yaml:"UID"`
	DefaultVisibility        string `yaml:"DefaultVisibility,omitempty"`
	MaxCheckpointSubmissions int    `yaml:"MaxCheckpointSubmissions,omitempty"`
	EmailOnCompletion        bool   `yaml:"EmailOnCompletion,omitempty"`
	TimeLimit                int    `yaml:"TimeLimit,omitempty"`
	Autoscore                bool   `yaml:"Autoscore,omitempty"`
	fromHeader               bool   // fromHeader is set true when the attrs were parsed from the header
}

var gitTopLevelCmd = "git rev-parse --show-toplevel"

// only used from publish, just going to send
func publishFindOrCreateConfig(target string) (bool, error) {
	cb := NewConfigBuilder(target, false, true, []string{})
	return cb.findOrCreateConfig()
}

// previewFindOrCreateConfig ensures a config for the previewed curriculum exists to be read by Learn. Because
// docker directory paths can contain content that looks like lesson curriculum, they are passed as arguments to
// prevent those directories from being included as lesson content in a generated config file.
func previewFindOrCreateConfig(target string, isSingleFilePreview bool, excludePaths []string) (bool, error) {
	cb := NewConfigBuilder(target, isSingleFilePreview, false, excludePaths)
	return cb.findOrCreateConfig()
}

func NewConfigBuilder(target string, isSingleFilePreview, publishContext bool, excludePaths []string) *ConfigBuilder {
	// Make sure we have an ending slash on the root dir
	blockRoot := ""
	if publishContext {
		target, _ = GitTopLevelDir()
		blockRoot = target + "/"
	} else {
		if strings.HasSuffix(target, "/") {
			blockRoot = target
		} else {
			blockRoot = target + "/"
		}
	}

	// UnitsDirectory supplied from a string flag
	unitsDirectory := UnitsDirectory
	if target == tmpSingleFileDir {
		unitsDirectory = "."
	}

	// If no unitsDir was passed in, create a Units directory string
	// UnitsDir is always blockroot + unitsDirName
	unitsDir := ""
	unitsDirName := ""
	unitsRootDirName := "units"

	if unitsDirectory == "" {
		unitsDir = blockRoot + unitsRootDirName
		unitsDirName = "Unit 1"
	} else {
		unitsDir = blockRoot + unitsDirectory
		unitsDirName = unitsDirectory
		unitsRootDirName = unitsDirectory
	}

	return &ConfigBuilder{
		target:              target,
		isSingleFilePreview: isSingleFilePreview,
		publishContext:      publishContext,
		excludePaths:        excludePaths,
		blockRoot:           blockRoot,
		unitsDir:            unitsDir,
		unitsDirName:        unitsDirName,
		unitsRootDirName:    unitsRootDirName,
	}
}

func (cb *ConfigBuilder) ConfigExists() bool {
	// Configs can be `yaml` or `yml`
	configYamlPath := ""
	configYmlPath := ""
	if strings.HasSuffix(cb.target, "/") {
		configYamlPath = cb.target + "config.yaml"
		configYmlPath = cb.target + "config.yml"
	} else {
		configYamlPath = cb.target + "/config.yaml"
		configYmlPath = cb.target + "/config.yml"
	}

	_, yamlErr := os.Stat(configYamlPath)
	_, ymlErr := os.Stat(configYmlPath)

	yamlPresent := !os.IsNotExist(yamlErr)
	ymlPresent := !os.IsNotExist(ymlErr)

	if !yamlPresent && !ymlPresent {
		return false
	}

	if !cb.isSingleFilePreview {
		if yamlPresent {
			fmt.Printf("INFO: Using existing config.yaml. \n")
		} else if ymlPresent {
			fmt.Printf("INFO: Using existing config.yml. \n")
		}
	}

	return yamlPresent || ymlPresent
}

// findOrCreateConfig returns true when a config.yaml was created
func (cb *ConfigBuilder) findOrCreateConfig() (bool, error) {
	if cb.ConfigExists() {
		return false, nil
	}

	if !cb.isSingleFilePreview {
		fmt.Printf("INFO: Using existing config.yaml. \n")
	}

	err := cb.createYamlConfig()
	if err != nil {
		return false, err
	}

	return true, nil
}

// createYamlConfig determines the context for creating a new config yaml and handles file creation and encoding
func (cb *ConfigBuilder) createYamlConfig() error {
	// The config file location that we will be creating
	autoConfigYamlPath := cb.blockRoot + "autoconfig.yaml"

	// Remove the existing one if its around
	_, err := os.Stat(autoConfigYamlPath)
	if err == nil {
		os.Remove(autoConfigYamlPath)
	}

	// Create tmpSingleFileDir if it does not exist
	if _, err := os.Stat(tmpSingleFileDir); os.IsNotExist(err) {
		os.Mkdir(tmpSingleFileDir, os.FileMode(0777))
	}

	configFile, err := os.Create(autoConfigYamlPath)
	if err != nil {
		return err
	}
	configFile.WriteString(autoComment)
	defer configFile.Sync()
	defer configFile.Close()

	encoder := yaml.NewEncoder(configFile)
	defer encoder.Close()

	autoConfig, err := cb.newConfigYaml()
	if err != nil {
		return err
	}
	encoder.Encode(autoConfig)
	return nil
}

// newConfigYaml creates a ConfigYaml struct given certain conditions
// 1. Did you give us a units directory?
// 2. Do you have a units directory?
// Units must exist in units dir or one provided!
func (cb *ConfigBuilder) newConfigYaml() (ConfigYaml, error) {
	config := ConfigYaml{Standards: []Standard{}}

	unitToContentFileMap, err := cb.buildUnitToContentFileMap()
	if err != nil {
		return config, err
	}
	if len(unitToContentFileMap) == 0 {
		return config, fmt.Errorf("No content found at '%s'. Preview of an individual unit is not supported, make sure '%s' is the root of a repo or a single lesson.", cb.target, cb.target)
	}

	// sort unit keys in lexicographical order
	unitKeys := make([]string, 0, len(unitToContentFileMap))
	for unit := range unitToContentFileMap {
		unitKeys = append(unitKeys, unit)
	}
	sort.Strings(unitKeys)

	formattedTargetName := formattedName(cb.target)
	for _, unit := range unitKeys {
		parts := strings.Split(unit, "/")
		if strings.HasPrefix(parts[0], "__") {
			continue
		}

		// skip the unit when all content files are excluded
		allFilesExcluded := true
		for _, contentFile := range unitToContentFileMap[unit] {
			if !anyMatchingPrefix("/"+contentFile.Path, cb.excludePaths) {
				allFilesExcluded = false
				break
			}
		}
		if allFilesExcluded {
			continue
		}

		unitsDirectoryFile, err := os.Stat(cb.unitsDir)

		whereToLookForUnits := cb.blockRoot
		if err == nil && unitsDirectoryFile.IsDir() {
			whereToLookForUnits = fmt.Sprintf("%s%s", cb.blockRoot, cb.unitsRootDirName)
		}
		standard := newStandard(whereToLookForUnits, unit)
		if standard.Title == "" {
			standard.Title = formattedTargetName
		}
		if standard.Description == "" {
			standard.Description = formattedTargetName
		}

		for _, contentFile := range unitToContentFileMap[unit] {
			if anyMatchingPrefix("/"+contentFile.Path, cb.excludePaths) {
				continue
			}
			parts := strings.Split(contentFile.Path, "/")
			if strings.HasPrefix(parts[len(parts)-1], "__") {
				continue
			}
			if contentFile.Path != "README.md" {
				if strings.Contains(strings.ToLower(contentFile.Path), "..") {
					contentFile.Path = strings.Replace(contentFile.Path, "..", ".", 1)
				}
				if strings.HasPrefix(contentFile.Path, "./") {
					contentFile.Path = contentFile.Path[1:]
				} else {
					contentFile.Path = "/" + contentFile.Path
				}
				if contentFile.fromHeader {
					// when it came from the header but Type is not set, fall back to detecting from path
					if contentFile.Type == "" {
						contentFile.Type = detectContentType(contentFile.Path)
					}
					// when it came from the header but UID is not set, fall back to detecting from path
					if contentFile.UID == "" {
						cfUID := []byte(standard.Title + contentFile.Path)
						md5cfUID := md5.Sum(cfUID)
						contentFile.UID = hex.EncodeToString(md5cfUID[:])
					}
					// when it came from the header but DefaultVisibility is not set, fall back to detecting from path
					if contentFile.DefaultVisibility == "" && strings.Contains(strings.ToLower(contentFile.Path), "hidden") {
						contentFile.DefaultVisibility = "hidden"
					}
					standard.ContentFiles = append(standard.ContentFiles, contentFile)
				} else {
					cfUID := []byte(standard.Title + contentFile.Path)
					md5cfUID := md5.Sum(cfUID)

					contentFile.Type = detectContentType(contentFile.Path)
					contentFile.UID = hex.EncodeToString(md5cfUID[:])
					if strings.Contains(strings.ToLower(contentFile.Path), "hidden") {
						contentFile.DefaultVisibility = "hidden"
					}
					standard.ContentFiles = append(standard.ContentFiles, contentFile)
				}
			}
		}
		config.Standards = append(config.Standards, standard)
	}

	return config, nil
}

func detectContentType(p string) string {
	fullpath := strings.ToLower(p)
	parts := strings.Split(fullpath, "/")
	path := parts[len(parts)-1]
	instructorMatch, _ := regexp.MatchString("^instructor[.-]|[.-]instructor[.-]", path)
	checkpointMatch, _ := regexp.MatchString("^checkpoint[.-]|[.-]checkpoint[.-]", path)
	resourceMatch, _ := regexp.MatchString("^resource[.-]|[.-]resource[.-]", path)
	surveyMatch, _ := regexp.MatchString("^survey[.-]|[.-]survey[.-]", path)
	if instructorMatch {
		return "Instructor"
	} else if checkpointMatch {
		return "Checkpoint"
	} else if resourceMatch {
		return "Resource"
	} else if surveyMatch {
		return "Survey"
	}
	return "Lesson"
}

func formattedName(name string) string {
	parts := strings.Split(name, "/")

	a := regexp.MustCompile(`\-`)
	parts = a.Split(parts[0], -1)

	if len(parts) == 1 {
		a = regexp.MustCompile(`\.`)
		parts = a.Split(parts[0], -1)
	}

	a = regexp.MustCompile(`\_`)
	parts = a.Split(strings.Join(parts, " "), -1)

	formattedName := ""
	for _, piece := range parts {
		formattedName = formattedName + " " + strings.Title(piece)
	}
	// remove leading numbers like '01'
	a = regexp.MustCompile(`^([0-9]{1,3} :?)`)
	parts = a.Split(strings.TrimSpace(formattedName), -1)
	return strings.TrimSpace(strings.Join(parts, ""))
}

// anyMatchingPrefix reports if any of the given prefixes have been found to be a prefix of the target
func anyMatchingPrefix(target string, prefixes []string) bool {
	for _, prefix := range prefixes {
		if strings.HasPrefix(target, prefix) {
			return true
		}
	}
	return false
}

// tries to find the config yaml or autoconfig yaml
func findConfig(target string) (string, error) {
	configPath := ""
	if strings.HasSuffix(target, "/") {
		configPath = target + "config.yaml"
	} else {
		configPath = target + "/config.yaml"
	}
	_, yamlExists := os.Stat(configPath)
	if yamlExists != nil {
		if strings.HasSuffix(target, "/") {
			configPath = target + "config.yml"
		} else {
			configPath = target + "/config.yml"
		}
		_, yamlExists = os.Stat(configPath)
		if yamlExists != nil {
			if strings.HasSuffix(target, "/") {
				configPath = target + "autoconfig.yaml"
			} else {
				configPath = target + "/autoconfig.yaml"
			}

			_, yamlExists = os.Stat(configPath)

			if yamlExists != nil {
				return "", fmt.Errorf("Could not find config or autoconfig yaml")
			}
		}
	}

	return configPath, nil
}

// readContentFileAttrs takes a file path and a readPath and returns contentFileAttrs if they were present in the header yaml
func readContentFileAttrs(path, readPath string) (contentFile ContentFileAttrs, err error) {
	file, err := os.Open(readPath)
	if err != nil {
		return contentFile, err
	}
	defer file.Close()
	bufferSize := 1024 * 1024
	buf := make([]byte, bufferSize)
	scanner := bufio.NewScanner(file)
	scanner.Buffer(buf, bufferSize)
	scanner.Split(split)
	// read to the first yaml delimiter
	scanner.Scan()
	yamlText := scanner.Text() // extract yaml
	if err = scanner.Err(); err != nil {
		return contentFile, err
	}

	if strings.TrimSpace(yamlText) != "" {
		err = printExtras(yamlText, path)
		if err != nil {
			return contentFile, err
		}
		err = yaml.Unmarshal([]byte(yamlText), &contentFile)
		if err != nil {
			return contentFile, fmt.Errorf("Error parsing yaml header for '%s': %s\n'", path, err)
		}
		contentFile.Path = path
		contentFile.fromHeader = true
		return contentFile, err
	}
	contentFile.Path = path
	return contentFile, nil
}

// printExtras prints unknown content file header keys
func printExtras(yamlText, path string) error {
	attributes := map[string]interface{}{}
	err := yaml.Unmarshal([]byte(yamlText), &attributes)
	if err != nil {
		return fmt.Errorf("yaml header for '%s' is not valid:\n%s\n", path, err)
	}
	for key := range attributes {
		acceptableKey := false
		for _, validKey := range validContentFileAttrs {
			if key == validKey {
				acceptableKey = true
			}
		}
		if !acceptableKey {
			fmt.Printf("Found unknown content file header key '%s' in file %s\n", key, path)
		}
	}
	return nil
}

// buildUnitToContentFileMap reads contents from the unit directory and includes md files. It returns attributes from the header for each file
// TODO refactor inputs, should be simplified like unitsDir is just the first and last inputs put together; example from test
// bockRoot ../../fixtures/test-block-no-config/
// unitsDir ../../fixtures/test-block-no-config/units
// unitsDirName Unit 1
// unitsRootDirName units
func (cb *ConfigBuilder) buildUnitToContentFileMap() (map[string][]ContentFileAttrs, error) {
	unitToContentFileMap := map[string][]ContentFileAttrs{}

	// Check to see if units directory exists
	_, err := os.Stat(cb.unitsDir)

	whereToLookForUnits := cb.blockRoot

	if err == nil {
		whereToLookForUnits = cb.unitsDir

		allItems, err := ioutil.ReadDir(whereToLookForUnits)
		if err != nil {
			return unitToContentFileMap, err
		}

		for _, info := range allItems {
			if info.Mode().IsRegular() && strings.HasSuffix(info.Name(), ".md") {
				readPath := cb.blockRoot + cb.unitsRootDirName + "/" + info.Name()
				path := cb.unitsRootDirName + "/" + info.Name()
				contentFile, err := readContentFileAttrs(path, readPath)
				if err != nil {
					return unitToContentFileMap, err
				}
				unitToContentFileMap[cb.unitsDirName] = append(unitToContentFileMap[cb.unitsDirName], contentFile)
			}
		}
	}

	// Find all the directories in the block
	directories := []string{}
	allDirs, err := ioutil.ReadDir(whereToLookForUnits)
	if err != nil {
		return unitToContentFileMap, err
	}

	for _, info := range allDirs {
		if info.IsDir() {
			directories = append(directories, info.Name())
		}
	}

	if len(directories) > 0 {
		for _, dirName := range directories {
			nestedFolder := ""

			if dirName != ".git" {
				if strings.HasSuffix(whereToLookForUnits, "/") {
					nestedFolder = whereToLookForUnits + dirName
				} else {
					nestedFolder = whereToLookForUnits + "/" + dirName
				}

				err = filepath.Walk(nestedFolder, func(path string, info os.FileInfo, err error) error {
					if err != nil {
						return err
					}

					if len(cb.blockRoot) > 0 && len(path) > len(cb.blockRoot) && strings.HasSuffix(path, ".md") {
						localPath := path
						if cb.blockRoot != "./" {
							localPath = path[len(cb.blockRoot):]
						}

						readPath := cb.blockRoot + "/" + localPath
						contentFile, err := readContentFileAttrs(localPath, readPath)
						if err != nil {
							return err
						}
						unitToContentFileMap[dirName] = append(unitToContentFileMap[dirName], contentFile)
					}

					return nil
				})
				if err != nil {
					return unitToContentFileMap, err
				}
			}
		}
	}
	return unitToContentFileMap, nil
}

// get the root dir of the git project
func GitTopLevelDir() (string, error) {
	out, err := exec.Command("bash", "-c", gitTopLevelCmd).CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("%s", out)
	}

	return strings.TrimSpace(string(out)), err
}

// newStandard returns a standard from tne unitDir and unit name combination
// unitDir is the location of the individual unit, with unit the directory beneath it
// Either a description yaml file is read from, or the unit name is used to build a standard
// func standardAttributes(unitDir, unit string) (title, UID, description string, successCriteria []string) {
func newStandard(unitDir, unit string) Standard {
	yamlLocation := fmt.Sprintf("%s/%s/%s", unitDir, unit, "description.yaml")
	yamlBytes, err := os.ReadFile(yamlLocation)
	if err != nil {
		// no description yaml found,
		return standardFromUnit(unit)
	} else {
		// read yaml contents of file
		//return standardFromUnit(unit)
		standard := Standard{}
		if err = yaml.NewDecoder(bytes.NewReader(yamlBytes)).Decode(&standard); err != nil {
			return standardFromUnit(unit)
		}
		return standard
	}
}

func standardFromUnit(unit string) Standard {
	title := formattedName(unit)
	unitUID := []byte(title)
	md5unitUID := md5.Sum(unitUID)
	UID := hex.EncodeToString(md5unitUID[:])
	description := title
	successCriteria := []string{"success criteria"}
	return Standard{
		Title:           title,
		UID:             UID,
		Description:     description,
		SuccessCriteria: successCriteria,
	}
}

// split is the bufio Scanner Split interface implementation for fetching content file header attributes
func split(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if len(data) < 4 {
		return 0, nil, nil // the content simply does not exist
	}
	// if the first three characters aren't '---' don't continue
	if string(data[:3]) != "---" {
		return 0, nil, nil
	}

	index := bytes.Index(data, []byte("---"))
	if index >= 0 {
		// data[index+3:] reads past the incidence of '---' up to the index of the next '---'
		next := bytes.Index(data[index+3:], []byte("---"))
		if next > 0 {
			// the next '---' is found, advance just past it with next + 3, the token here is the terminating '---'
			return next + 3, bytes.TrimSpace(data[:next+3]), nil
		}
		// when no second instance, advance to the end of the file
		return len(data), bytes.TrimSpace(data[index+3:]), nil
	}
	// when atEOF return the final index
	if atEOF {
		return len(data), data, nil
	}
	return 0, nil, nil
}
