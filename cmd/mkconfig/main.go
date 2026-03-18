// Package main provides a tool to create an example configuration file.
package main

import (
	"log"

	"github.com/yuioto/fvti-xsgz-sign/internal/config"
)

// CreateExampleConfig creates an example configuration file.
func CreateExampleConfig(filename string) {
	if err := config.CreateDefaultConfig(filename); err != nil {
		log.Fatalf("Failed to create example config: %v", err)
	}
}

func main() {
	filename := "config.example.kdl"
	CreateExampleConfig(filename)
}
