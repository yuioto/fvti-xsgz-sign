package config

import (
	"log"
	"os"

	"fvti-xsgz-sign/pkg/set"

	"github.com/pelletier/go-toml/v2"
)

func CreateDefaultConfig(filename string) {
	defaultConfig := Config{
		StudentId: "{fvti_student_id}",
		Login: login{
			Password:      "{fvti_xsgz_password}",
			Authorization: "",
		},
		Task: task{
			Name:   "24级新生晚点名",
			Id:     "",
			SignId: "No need to fill in anything, auto-populates on request",
		},
		Nofy: set.NotyId,
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
