package main

import (
	"flag"
	"gopkg.in/yaml.v3"
	"log/slog"
	"os"
)

type Config struct {
	SourceUrl string `yaml:"source_url"`
	DBFile    string `yaml:"db_file"`
}

func getConfig(configPath string) (*Config, error) {
	config := &Config{}

	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	d := yaml.NewDecoder(file)
	if err = d.Decode(&config); err != nil {
		return nil, err
	}

	return config, nil
}

func getFlags() (string, string, bool, slog.Level) {
	// parse flags
	var configPath string
	flag.StringVar(&configPath, "c", "config.yaml", "path to config.yml file")
	var inputString string
	flag.StringVar(&inputString, "s", "", "string to find")
	var byIndex bool
	flag.BoolVar(&byIndex, "i", false, "find comics by index")
	var showDebugMsg bool
	flag.BoolVar(&showDebugMsg, "d", false, "show debug messages in log")
	flag.Parse()

	level := slog.LevelInfo
	if showDebugMsg {
		level = slog.LevelDebug
	}
	return configPath, inputString, byIndex, level
}

func getGoroutinesNum() (int, error) {
	defaultValue := 500
	obj := make(map[string]int)

	yamlFile, err := os.ReadFile("parallel")
	if err != nil {
		return defaultValue, err
	}
	err = yaml.Unmarshal(yamlFile, obj)
	if err != nil {
		return defaultValue, err
	}

	if obj["goroutines"] == 0 {
		obj["goroutines"] = defaultValue
	}
	return obj["goroutines"], nil
}
