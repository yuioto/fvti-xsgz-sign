package main

import (
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

	//if savestusignin.GetTaskList(config.Login.Authorization) != nil {}

	// while config.Task.Id == nil, jump GetTaskId
	if config.Task.Id == "" {
		if config.Task.Id, err = savestusignin.GetTaskId(config.Task.Name, config.Login.Authorization); err != nil {
			log.Fatalln("Failed to requeset Task.Id:", err)
		}
	}
	//fmt.Println(config.Login.Authorization, config.Task.Id) // debug: check id update

	if err := savestusignin.PostStuSignIn(config.StudentId, config.Task.Id, config.Login.Authorization); err != nil {
		log.Fatalln("Failed sign:", err)
	}
	log.Println("Sign successfully.")

	//if a, err := checklogin.GetSignHistory(config.Task.Id, "2", config.Login.Authorization); err != nil {
	//	log.Fatalln("Failed to GetSignHistory")
	//	fmt.Println(a)
	//}
}
