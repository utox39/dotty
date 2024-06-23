package main

import (
	"encoding/json"
	"fmt"
	"github.com/tidwall/sjson"
	"os"
	"path/filepath"
)

type Dotfiles struct {
	Filepath        []string `json:"dotfiles"`
	DestinationPath string   `json:"destination-path"`
}

func (d *Dotfiles) ReadConfig(jsonConfigPath string) error {
	// Read json config file
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

func AddFile(newFilePath, jsonConfigPath string) error {
	absFilePath, err := filepath.Abs(newFilePath)
	if err != nil {
		return fmt.Errorf("could not determine absolute path of file %v: %v", newFilePath, err)
	}

	// Replace the ~ with the home folder path
	err = ReplaceTilde(&absFilePath)
	if err != nil {
		return err
	}

	err = ValidatePath(&absFilePath)
	if err != nil {
		return err
	}

	byteValue, err := ReadFile(&jsonConfigPath)

	// Add the new dotfile to the end of the dotfiles array
	value, err := sjson.Set(string(byteValue), "dotfiles.-1", absFilePath)
	if err != nil {
		return fmt.Errorf("error adding the new dotfile to the dotfiles array %v: %v", jsonConfigPath, err)
	}
	err = os.WriteFile(jsonConfigPath, []byte(value), 0644)
	if err != nil {
		return fmt.Errorf("error writing the file %v: %v", jsonConfigPath, err)
	}

	return nil
}
