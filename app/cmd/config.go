package cmd

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

func findOrCreateConfigDir(target string) (bool, error) {
	return doesConfigExistOrCreate(target, false, []string{})
}

// Check whether or nor a config file exists and if it does not we are going to attempt to create one
func doesConfigExistOrCreate(target string, isSingleFilePreview bool, excludePaths []string) (bool, error) {
	// Configs can be `yaml` or `yml`
	configYamlPath := ""
	if strings.HasSuffix(target, "/") {
		configYamlPath = target + "config.yaml"
	} else {
		configYamlPath = target + "/config.yaml"
	}

	configYmlPath := ""
	if strings.HasSuffix(target, "/") {
		configYmlPath = target + "config.yml"
	} else {
		configYmlPath = target + "/config.yml"
	}

	createdConfig := false
	_, yamlExists := os.Stat(configYamlPath)

	if yamlExists == nil { // Yaml exists
		if isSingleFilePreview == false {
			fmt.Printf("INFO: Using existing config.yaml. ")
		}

		return createdConfig, nil
	} else if os.IsNotExist(yamlExists) {
		_, ymlExists := os.Stat(configYmlPath)

		if ymlExists == nil { // Yml exists
			if isSingleFilePreview == false {
				fmt.Printf("INFO: Using existing config.yaml. ")
			}

			return createdConfig, nil
		} else if os.IsNotExist(ymlExists) {
			if isSingleFilePreview == false {
				// Neither exists so we are going to create one
				fmt.Printf("INFO: No configuration found, generating autoconfig.yaml ")
			}
			if target == tmpSingleFileDir {
				err := createAutoConfig(target, ".", excludePaths)
				if err != nil {
					return false, err
				}
			} else {
				// UnitsDirectory supplied from a string flag
				err := createAutoConfig(target, UnitsDirectory, excludePaths)
				if err != nil {
					return false, err
				}
			}
			createdConfig = true
		}
	}
	return createdConfig, nil
}

// Creates a config file based on three things:
// 1. Did you give us a units directory?
// 2. Do you have a units directory?
// Units must exist in units dir or one provided!
func createAutoConfig(target, requestedUnitsDir string, excludePaths []string) error {
	blockRoot := ""

	// Make sure we have an ending slash on the root dir
	if strings.HasSuffix(target, "/") {
		blockRoot = target
	} else {
		blockRoot = target + "/"
	}

	// The config file location that we will be creating
	autoConfigYamlPath := blockRoot + "autoconfig.yaml"

	// Remove the existing one if its around
	_, err := os.Stat(autoConfigYamlPath)
	if err == nil {
		os.Remove(autoConfigYamlPath)
	}

	// Create tmpSingleFileDir if it does not exist
	if _, err := os.Stat(tmpSingleFileDir); os.IsNotExist(err) {
		os.Mkdir(tmpSingleFileDir, os.FileMode(0777))
	}

	// Create the config file
	configFile, err := os.Create(autoConfigYamlPath)
	if err != nil {
		return err
	}
	defer configFile.Sync()
	defer configFile.Close()

	// If no unitsDir was passed in, create a Units directory string
	unitsDir := ""
	unitsDirName := ""
	unitsRootDirName := "units"

	if requestedUnitsDir == "" {
		unitsDir = blockRoot + unitsRootDirName
		unitsDirName = "Unit 1"
	} else {
		unitsDir = blockRoot + requestedUnitsDir
		unitsDirName = requestedUnitsDir
		unitsRootDirName = requestedUnitsDir
	}

	unitToContentFileMap := map[string][]string{}

	// Check to see if units directory exists
	_, err = os.Stat(unitsDir)

	whereToLookForUnits := blockRoot

	if err == nil {
		whereToLookForUnits = unitsDir

		allItems, err := ioutil.ReadDir(whereToLookForUnits)
		if err != nil {
			return err
		}

		for _, info := range allItems {
			if info.Mode().IsRegular() && strings.HasSuffix(info.Name(), ".md") {
				unitToContentFileMap[unitsDirName] = append(unitToContentFileMap[unitsDirName], unitsRootDirName+"/"+info.Name())
			}
		}
	}

	// Find all the directories in the block
	directories := []string{}

	allDirs, err := ioutil.ReadDir(whereToLookForUnits)
	if err != nil {
		return err
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

					if len(blockRoot) > 0 && len(path) > len(blockRoot) && strings.HasSuffix(path, ".md") {
						localPath := path
						if blockRoot != "./" {
							localPath = path[len(blockRoot):]
						}
						if strings.Contains(localPath, "\\") {
							localPath = strings.Replace(localPath, "\\", "/", -1)
						}
						unitToContentFileMap[dirName] = append(unitToContentFileMap[dirName], localPath)
					}

					return nil
				})
				if err != nil {
					return err
				}
			}
		}
	}

	configFile.WriteString("# This file is auto-generated and orders your content based on the file structure of your repo.\n")
	configFile.WriteString("# Do not edit this file; it will be replaced the next time you run the preview command.\n")
	configFile.WriteString("\n")
	configFile.WriteString("# To manually order the contents of this curriculum rather than using the auto-generated file,\n")
	configFile.WriteString("# include a config.yaml in your repo following the same conventions as this auto-generated file.\n")
	configFile.WriteString("# A user-created config.yaml will have priority over the auto-generated one.\n")
	configFile.WriteString("\n")
	configFile.WriteString("---\n")
	configFile.WriteString("Standards:\n")

	if len(unitToContentFileMap) == 0 {
		return fmt.Errorf("No content found at '%s'. Preview of an individual unit is not supported, make sure '%s' is the root of a repo or a single lesson.", target, target)
	}

	// sort unit keys in lexigraphical order
	unitKeys := make([]string, 0, len(unitToContentFileMap))
	for unit := range unitToContentFileMap {
		unitKeys = append(unitKeys, unit)
	}
	sort.Strings(unitKeys)

	formattedTargetName := formattedName(target)
	for _, unit := range unitKeys {
		parts := strings.Split(unit, "/")
		if strings.HasPrefix(parts[0], "__") {
			continue
		}

		// skip the unit when all content files are excluded
		allFilesExcluded := true
		for _, path := range unitToContentFileMap[unit] {
			if !anyMatchingPrefix("/"+path, excludePaths) {
				allFilesExcluded = false
				break
			}
		}
		if allFilesExcluded {
			continue
		}

		configFile.WriteString("  -\n")

		formattedUnitName := formattedName(unit)
		if formattedUnitName != "" {
			configFile.WriteString("    Title: " + formattedUnitName + "\n")
		} else {
			configFile.WriteString("    Title: " + formattedTargetName + "\n")
		}

		var unitUID = []byte(formattedUnitName)
		var md5unitUID = md5.Sum(unitUID)

		if formattedUnitName != "" {
			configFile.WriteString("    Description: " + formattedUnitName + "\n")
		} else {
			configFile.WriteString("    Description: " + formattedTargetName + "\n")
		}

		configFile.WriteString("    UID: " + hex.EncodeToString(md5unitUID[:]) + "\n")
		configFile.WriteString("    SuccessCriteria:\n")
		configFile.WriteString("      - success criteria\n")
		configFile.WriteString("    ContentFiles:\n")

		for _, path := range unitToContentFileMap[unit] {
			if anyMatchingPrefix("/"+path, excludePaths) {
				continue
			}
			parts := strings.Split(path, "/")
			if strings.HasPrefix(parts[len(parts)-1], "__") {
				continue
			}
			if path != "README.md" {
				configFile.WriteString("      -\n")

				contentFileType := detectContentType(path)
				configFile.WriteString("        Type: " + contentFileType + "\n")

				if strings.Contains(strings.ToLower(path), "hidden") {
					configFile.WriteString("        DefaultVisibility: hidden\n")
				}

				if strings.Contains(strings.ToLower(path), "..") {
					path = strings.Replace(path, "..", ".", 1)
				}

				var cfUID = []byte(formattedUnitName + path)
				var md5cfUID = md5.Sum(cfUID)

				configFile.WriteString("        UID: " + hex.EncodeToString(md5cfUID[:]) + "\n")

				if strings.HasPrefix(path, "./") {
					configFile.WriteString("        Path: " + path[1:] + "\n")
				} else {
					configFile.WriteString("        Path: /" + path + "\n")
				}
			}
		}
	}
	if err != nil {
		return err
	}
	return nil
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

// anyMatchingPrefix reports if any of the given prefixes haveb been found to be a prefix of the target
func anyMatchingPrefix(target string, prefixes []string) bool {
	for _, prefix := range prefixes {
		if strings.HasPrefix(target, prefix) {
			return true
		}
	}
	return false
}
