package main

import (
	"fmt"
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

	if err := SignNtfy(config); err != nil {
		log.Fatalln(err)
	}

}

func SignNtfy(config cfgset.Config) error {
	if config, err := Sign(config); err != nil {
		notify.SendNtfyMessage(config.Nofy, "max", "Sign Failed", err.Error())
		return err
	} else {
		msg := "StudentId: " + config.StudentId + " Task.Name: " + config.Task.Name + " Task.Id: " + config.Task.Id + " Task.SignId: " + config.Task.SignId
		notify.SendNtfyMessage(config.Nofy, "high", "Sign Done", msg)
		return nil
	}
}

func Sign(config cfgset.Config) (cfgset.Config, error) {
	var err error

	if config.Login.Authorization == "" {
		if config.Login.Authorization, err = savestusignin.GetAuthorization(config.StudentId, config.Login.Password); err != nil {
			return config, fmt.Errorf("failed to requeset Login.Authorization: %w", err)
		}
	}

	//fmt.Println(config.Login.Authorization)

	// while config.Task.Id == nil, jump GetTaskId
	if config.Task.Id == "" {
		if config.Task.Id, err = savestusignin.GetTaskId(config.Task.Name, config.Login.Authorization); err != nil {
			return config, fmt.Errorf("failed to requeset Task.Id: %w", err)
		}
	}

	//fmt.Println(config.Login.Authorization, config.Task.Id) // debug: check id update

	if err := savestusignin.PostStuSignIn(config.StudentId, config.Task.Id, config.Login.Authorization); err != nil {
		return config, fmt.Errorf("failed sign: %w", err)
	}
	//log.Println("Sign successfully.")

	config.Task.SignId, err = savestusignin.GetSignId(config.Task.Id, config.Login.Authorization)
	if err != nil {
		return config, fmt.Errorf("get SignId failed: %w", err)
	}

	return config, nil
}
