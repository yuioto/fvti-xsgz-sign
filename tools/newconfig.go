package main

import cfg "fvti-xsgz-sign/internal/config"

func CreateExampleConfig(filename string) {
	cfg.CreateDefaultConfig(filename)
}

func main() {
	filename := "Config.example.toml"
	CreateExampleConfig(filename)
}
