package main

import (
	"courses/core"
	"flag"
	"gopkg.in/yaml.v3"
	"log/slog"
	"os"
)

func getConfig(configPath string) (core.Config, error) {
	config := core.Config{}

	file, err := os.ReadFile(configPath)
	if err != nil {
		return config, err
	}

	if err = yaml.Unmarshal(file, &config); err != nil {
		return core.Config{}, err
	}

	return config, nil
}

func getFlags() (string, int, slog.Level) {
	var configPath string
	flag.StringVar(&configPath, "c", "config.yaml", "path to config file")
	var port int
	flag.IntVar(&port, "p", -1, "server port")
	var showDebugMsg bool
	flag.BoolVar(&showDebugMsg, "d", false, "show debug messages in log")
	flag.Parse()

	level := slog.LevelInfo
	if showDebugMsg {
		level = slog.LevelDebug
	}
	return configPath, port, level
}
