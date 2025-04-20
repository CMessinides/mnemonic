package config

import (
	"encoding/json"
	"fmt"
	"path/filepath"
)

type ServerConfig struct {
	Host string `json:"host"`
	Port uint   `json:"port"`
}

type UserConfig struct {
	Server *ServerConfig `json:"server"`
}

type UserDirs struct {
	Home       string
	ConfigHome string
	DataHome   string
}

type Config struct {
	Dirs *UserDirs
	UserConfig
}

func (c *Config) DBConnString() string {
	return filepath.Join(c.Dirs.DataHome, "mnemonic.sqlite") + "?_time_format=sqlite"
}

type (
	LookupEnv  func(name string) (value string, exists bool)
	FileExists func(path string) bool
	ReadFile   func(path string) ([]byte, error)
)

func getConfigDir(lookup LookupEnv, fallback string) string {
	if c, exists := lookup("MNEMONIC_CONFIG_HOME"); exists {
		return c
	}

	return filepath.Join(fallback, "mnemonic")
}

func getDataDir(lookup LookupEnv, fallback string) string {
	if d, exists := lookup("MNEMONIC_DATA_HOME"); exists {
		return d
	}

	return filepath.Join(fallback, "mnemonic")
}

var defaultConfig = &UserConfig{
	Server: &ServerConfig{
		Host: "127.0.0.1",
		Port: 9753,
	},
}

func ReadConfig(
	userHome string,
	configHome string,
	dataHome string,
	lookup LookupEnv,
	fileExists FileExists,
	readFile ReadFile,
) (*Config, error) {
	userDirs := &UserDirs{
		Home:       userHome,
		ConfigHome: getConfigDir(lookup, configHome),
		DataHome:   getDataDir(lookup, dataHome),
	}

	var userConfig UserConfig
	configFile := filepath.Join(userDirs.ConfigHome, "config.json")
	if !fileExists(configFile) {
		userConfig = *defaultConfig
	} else {
		data, err := readFile(configFile)
		if err != nil {
			return nil, fmt.Errorf("could not read config file at %s: %w", configFile, err)
		}

		err = json.Unmarshal(data, &userConfig)
		if err != nil {
			return nil, fmt.Errorf("could not parse config file at %s: %w", configFile, err)
		}
	}

	return &Config{
		Dirs:       userDirs,
		UserConfig: userConfig,
	}, nil
}
