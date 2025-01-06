package main

import (
	cfgset "fvti-xsgz-sign/internal/config"
	"fvti-xsgz-sign/pkg/sign"
	"log"
)

type Config cfgset.Config

func main() {
	config, err := sign.LoadConfig()
	if err != nil {
		log.Fatalln(err)
	}

	if err := sign.SignAndNotify(config); err != nil {
		log.Fatalln(err)
	}
}
