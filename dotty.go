package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

type Dotfiles struct {
	Filename []string `json:"dotfiles"`
}

func CopyFile(fileName string) {
	homeFolder, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	destinationPath := filepath.Join(homeFolder, "dotfiles/")

	fileAbsPath := filepath.Join(homeFolder, fileName)

	// Open the file and save the content
	file, err := os.Open(fileAbsPath)
	if err != nil {
		log.Fatal(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatalf("Error closing the file %v: %v", fileAbsPath, err)
		}
	}(file)

	// Create a copy of the file to copy
	newFile, err := os.Create(filepath.Join(destinationPath, fileName))
	if err != nil {
		log.Fatal(err)
	}
	defer func(newFile *os.File) {
		err := newFile.Close()
		if err != nil {
			log.Fatalf("Error closing the file %v: %v", fileAbsPath, err)
		}
	}(newFile)

	// Write the content of the file to copy in the copy file
	bytesWritten, err := io.Copy(newFile, file)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Copied %d bytes.\n", bytesWritten)

	err = newFile.Sync()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%v copied successfully\n\n", fileName)
}

func main() {
	fmt.Println("Hello world from dotty :)")

	// Get home folder
	homeFolder, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	// Read json config file
	jsonConfigPath := filepath.Join(homeFolder, ".config/dotty/config.json")
	jsonFile, err := os.Open(jsonConfigPath)
	if err != nil {
		log.Fatal(err)
	}
	defer func(jsonFile *os.File) {
		err := jsonFile.Close()
		if err != nil {
			log.Fatalf("Error closing the file %v: %v", jsonFile, err)
		}
	}(jsonFile)

	byteValue, _ := io.ReadAll(jsonFile)

	// Parse the json file
	var dotfile Dotfiles
	err = json.Unmarshal(byteValue, &dotfile)
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < len(dotfile.Filename); i++ {
		fmt.Printf("Saving %v...\n", dotfile.Filename[i])

		if _, err := os.Stat(filepath.Join(homeFolder, dotfile.Filename[i])); err == nil {
			fmt.Printf("copying %v...\n", dotfile.Filename[i])
			CopyFile(dotfile.Filename[i])
		} else if errors.Is(err, os.ErrNotExist) {
			fmt.Printf("dotty: %v does not exists", filepath.Join(homeFolder, dotfile.Filename[i]))
			continue
		} else {
			log.Fatalf("dotty: error: %v", err)
		}
	}
}
