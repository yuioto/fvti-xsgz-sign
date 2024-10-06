package config

import (
	"log"
	"os"

	"github.com/pelletier/go-toml/v2"
)

func CreateDefaultConfig(filename string) {
	defaultConfig := Config{
		StudentId: "{fvti_sso_id}",
		Login: login{
			Password:      "{fvti_sso_password}",
			Authorization: "Bearer {token}",
		},
		Task: task{
			Name: "24级新生晚点名",
			Id:   "3ebb9f9d-a8fe-49b4-9e61-9079e09868be",
		},
	}

	file, err := os.Create(filename)
	if err != nil {
		log.Fatalln("Failed to create file:", err)
	}
	defer file.Close()

	encoder := toml.NewEncoder(file)
	if err := encoder.Encode(defaultConfig); err != nil {
		log.Fatalln("Failed to write default config:", err)
	}
	log.Println("Create default config file successfully")
}
