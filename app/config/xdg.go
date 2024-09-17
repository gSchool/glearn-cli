package config

import (
	"fmt"
	"os"
	"path"

	"github.com/adrg/xdg"
	"github.com/fatih/color"
	"github.com/spf13/viper"
)

var warningShown bool = false

var defaultFileContent = "api_token: "

var xdgDirPath = path.Join(xdg.ConfigHome, "galvanize")
var xdgFileName = "config"
var xdgConfig = viper.New()

var homeDirPath = xdg.Home
var homeFileName = ".glearn-config"
var homeConfig = viper.New()
var configInUse *viper.Viper

func ConfigPath() string {
	setUp()
	return configInUse.ConfigFileUsed()
}

func Upgrade() (bool, error) {
	setUpWithWarning(false)
	if configInUse == homeConfig {
		os.MkdirAll(xdgDirPath, 0755)
		homeConfig.SetConfigFile(path.Join(xdgDirPath, xdgFileName+".yaml"))
		if err := homeConfig.WriteConfig(); err != nil {
			return false, err
		}
		if err := os.Remove(path.Join(homeDirPath, homeFileName+".yaml")); err != nil {
			return false, err
		}
		if err := xdgConfig.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
				return false, fmt.Errorf("unable to load new configuration: %v", err)
			}
		}
		configInUse = xdgConfig
		return true, nil
	}
	return false, nil
}

func GetString(key string) (string, error) {
	if err := setUp(); err != nil {
		return "", err
	}
	return configInUse.GetString(key), nil
}

func Set(key string, value any) error {
	if err := setUp(); err != nil {
		return err
	}
	configInUse.Set(key, value)
	return nil
}

func Write() error {
	if err := setUp(); err != nil {
		return err
	}
	return configInUse.WriteConfig()
}

func setUp() error {
	return setUpWithWarning(true)
}

func setUpWithWarning(showWarning bool) error {
	xdgConfig.AddConfigPath(xdgDirPath)
	xdgConfig.SetConfigName(xdgFileName)
	xdgErr := xdgConfig.ReadInConfig()
	_, xdgMissing := xdgErr.(viper.ConfigFileNotFoundError)

	homeConfig.AddConfigPath(homeDirPath)
	homeConfig.SetConfigName(homeFileName)
	homeErr := homeConfig.ReadInConfig()
	_, homeMissing := homeErr.(viper.ConfigFileNotFoundError)

	if (xdgErr != nil && !xdgMissing) || (homeErr != nil && !homeMissing) {
		// Something bad happened that we can't recover from
		// Return an error
		xdgPath := path.Join(xdgDirPath, xdgFileName+".yaml")
		homePath := path.Join(homeDirPath, homeFileName+".yaml")
		return fmt.Errorf("failed to read %s or %s", xdgPath, homePath)
	}

	if xdgMissing && homeMissing {
		// When both are missing, create the directory and file in xdg.ConfigHome
		// Then, run setUp again
		os.MkdirAll(xdgDirPath, 0755)
		os.WriteFile(path.Join(xdgDirPath, xdgFileName+".yaml"), []byte(defaultFileContent), 0700)
		return setUpWithWarning(false)
	} else if xdgMissing && !homeMissing {
		// When only the config file in $HOME exists, print a warning
		// Then use it
		if !warningShown && showWarning {
			warningShown = true
			color.Set(color.FgYellow)
			fmt.Println("Config file ~/.glearn-config is deprecated.")
			fmt.Println("Please run 'learn set --upgrade' to move the config file.")
			fmt.Println()
			color.Unset()
		}
		configInUse = homeConfig
	} else if !xdgMissing && homeMissing {
		// When config files exist only in XDG config directory, just use it
		configInUse = xdgConfig
	} else if !xdgMissing && !homeMissing {
		// When config files exist in both places, print a warning
		// Then, use the one in the XDG config directory
		if !warningShown && showWarning {
			warningShown = true
			color.Set(color.FgYellow)
			fmt.Println("Found both ~/.config/galvanize/config.yaml and ~/.glearn-config.yaml.")
			fmt.Println("Using ~/.config/galvanize/config.yaml")
			fmt.Println()
			color.Unset()
		}
		configInUse = xdgConfig
	}
	return nil
}
