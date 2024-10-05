package main

import (
	"fmt"
	"os"

	"github.com/pelletier/go-toml/v2"
)

func main() {
	const configFileName = "Config"
	const configFile = configFileName + ".toml"
	var config Config

	if _, err := os.Stat(configFile); err != nil {
		if os.IsNotExist(err) {
			fmt.Println("Create default config")
			createDefaultConfig(configFile)
			panic("Please write your config file")
		}
	} else {
		cfg, err := os.ReadFile(configFile)
		if err != nil {
			panic("Read file faild")
		}
		if err := toml.Unmarshal([]byte(cfg), &config); err != nil {
			panic("Read toml file faild")
		}
		fmt.Println("Read file Sueecuful")
	}

	fmt.Println("Config: \n", config) // debug
}

func createDefaultConfig(filename string) {
	defaultConfig := Config{
		StudentId: "{fvti sso id}",
		Login: login{
			Password:      "{fvti sso password}",
			Authorization: "Bearer {token}",
		},
		Task: task{
			Name: "24级新生晚点名",
			Id:   "3ebb9f9d-a8fe-49b4-9e61-9079e09868be",
		},
	}

	file, err := os.Create(filename)
	if err != nil {
		panic("Create file faild")
	}
	defer file.Close()

	encoder := toml.NewEncoder(file)
	if err := encoder.Encode(defaultConfig); err != nil {
		panic("Write default config faild")
	}
	fmt.Println("Create default config file succerful")
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
