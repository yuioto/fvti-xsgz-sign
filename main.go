package main

import (
	"fmt"
	"log"

	cfgset "fvti-xsgz-sign/utils/config"
	"fvti-xsgz-sign/utils/savestusignin"
)

type Config cfgset.Config

func main() {
	// Read config file(Don't need Close())
	configFile := cfgset.GetConfigFilePath("fvti-xsgz-sign")
	config, err := cfgset.LoadConfig(configFile)
	if err != nil {
		log.Fatalln(err)
	}

	// while config.Task.Id == nil, jump GetTaskId
	if config.Task.Id == "" {
		config.Task.Id, _ = savestusignin.GetTaskId(config.Task.Name, config.Login.Authorization)
		if err != nil {
			log.Fatalln("Failed to requeset Task.Id:", err)
		}
	}

	//fmt.Println("Config:\n", config) // debug

	fmt.Println("StudentId:", config.StudentId)
	fmt.Println("Password:", config.Login.Password)
	fmt.Println("Authorization:", config.Login.Authorization)
	fmt.Println("Task.Name:", config.Task.Name)
	fmt.Println("Task.Id:", config.Task.Id)
}
