package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Dotfiles struct {
	Filepath        []string `json:"dotfiles"`
	DestinationPath string   `json:"destination-path"`
}

func (d *Dotfiles) ReadConfig() error {
	// Get home folder
	homeFolder, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("error getting user home folder %v", err)
	}

	// Read json config file
	jsonConfigPath := filepath.Join(homeFolder, ".config/dotty/config.json")
	byteValue, err := ReadFile(&jsonConfigPath)
	if err != nil {
		return err
	}

	// Parse the json file
	err = json.Unmarshal(byteValue, d)
	if err != nil {
		return fmt.Errorf("error unmarshalling json: %v", err)
	}

	return nil
}
