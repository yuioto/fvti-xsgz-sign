package config

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
)

func GetConfigFilePath(appname string) string {
	var configDir string
	const configFileName = "Config"
	const configFile = configFileName + ".toml"

	switch runtime.GOOS {
	case "windows":
		configDir = filepath.Join(os.Getenv("APPDATA"), appname)
	case "linux":
		configDir = filepath.Join(os.Getenv("HOME"), ".config", appname)
	case "darwin":
		configDir = filepath.Join(os.Getenv("HOME"), "Library", "Application Support", appname)
	default:
		configDir = "./"
	}

	if err := os.MkdirAll(configDir, os.ModePerm); err != nil {
		log.Fatalln("Failed to create dir:", err)
	}

	return filepath.Join(configDir, configFile)
}

type Config struct {
	StudentId string
	Login     login
	Task      task
}

type task struct {
	Name string
	Id   string
}

type login struct {
	Password      string
	Authorization string
}
