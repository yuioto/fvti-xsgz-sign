package main

import (
	"log"
	"os"

	cfgset "fvti-xsgz-sign/utils/config"
	"fvti-xsgz-sign/utils/notify"
	"fvti-xsgz-sign/utils/savestusignin"
)

type Config cfgset.Config

func main() {
	// Read config file(Don't need Close())
	configFile := cfgset.GetConfigFilePath("fvti-xsgz-sign")
	if os.Getenv("FvtiSign") != "" {
		configFile = os.Getenv("FvtiSign")
	}
	config, err := cfgset.LoadConfig(configFile)
	if err != nil {
		log.Fatalln(err)
	}

	//fmt.Println(config)

	if config.Login.Authorization == "" {
		if config.Login.Authorization, err = savestusignin.GetAuthorization(config.StudentId, config.Login.Password); err != nil {
			log.Fatalln("Failed to requeset Login.Authorization:", err)
		}
	}

	//fmt.Println(config.Login.Authorization)

	// while config.Task.Id == nil, jump GetTaskId
	if config.Task.Id == "" {
		if config.Task.Id, err = savestusignin.GetTaskId(config.Task.Name, config.Login.Authorization); err != nil {
			log.Fatalln("Failed to requeset Task.Id:", err)
		}
	}

	//fmt.Println(config.Login.Authorization, config.Task.Id) // debug: check id update

	if err := savestusignin.PostStuSignIn(config.StudentId, config.Task.Id, config.Login.Authorization); err != nil {
		notify.SendNtfyMessage(config.Nofy, "max", "Sign Failed", "Sign")
		log.Fatalln("Failed sign:", err)
	}
	//log.Println("Sign successfully.")

	msg := "StudentId: " + config.StudentId + " sign " + config.Task.Id + " SignId is " + savestusignin.GetSignId(config.Task.Id, config.Login.Authorization)
	notify.SendNtfyMessage(config.Nofy, "high", "Sign Done", msg)
}
