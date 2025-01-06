package sign

import (
	"fmt"
	"os"

	cfgset "fvti-xsgz-sign/internal/config"
	"fvti-xsgz-sign/pkg/notify"
	"fvti-xsgz-sign/pkg/savestusignin"
)

type Config cfgset.Config

func LoadConfig() (cfgset.Config, error) {
	configFile := cfgset.GetConfigFilePath("fvti-xsgz-sign")
	if envConfigFile := os.Getenv("FvtiSign"); envConfigFile != "" {
		configFile = envConfigFile
	}
	return cfgset.LoadConfig(configFile)
}

func SignAndNotify(config cfgset.Config) error {
	config, err := Sign(config)
	if err != nil {
		notify.SendNtfyMessage(config.Nofy, "max", "Sign Failed", err.Error())
		return err
	}

	msg := fmt.Sprintf("StudentId: %s Task.Name: %s Task.Id: %s Task.SignId: %s", config.StudentId, config.Task.Name, config.Task.Id, config.Task.SignId)
	notify.SendNtfyMessage(config.Nofy, "high", "Sign Done", msg)
	return nil
}

func Sign(config cfgset.Config) (cfgset.Config, error) {
	var err error

	if config.Login.Authorization == "" {
		config.Login.Authorization, err = savestusignin.GetAuthorization(config.StudentId, config.Login.Password)
		if err != nil {
			return config, fmt.Errorf("failed to request Login.Authorization: %w", err)
		}
	}

	if config.Task.Id == "" {
		config.Task.Id, err = savestusignin.GetTaskId(config.Task.Name, config.Login.Authorization)
		if err != nil {
			return config, fmt.Errorf("failed to request Task.Id: %w", err)
		}
	}

	err = savestusignin.PostStuSignIn(config.StudentId, config.Task.Id, config.Login.Authorization)
	if err != nil {
		return config, fmt.Errorf("failed to sign: %w", err)
	}

	config.Task.SignId, err = savestusignin.GetSignId(config.Task.Id, config.Login.Authorization)
	if err != nil {
		return config, fmt.Errorf("get SignId failed: %w", err)
	}

	return config, nil
}
