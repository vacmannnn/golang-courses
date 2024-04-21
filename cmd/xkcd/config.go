package main

import (
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	SourceUrl string `yaml:"source_url"`
	DBFile    string `yaml:"db_file"`
}

func newConfig(configPath string) (*Config, error) {
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
