package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Dotfiles struct {
	Filepath        []string `json:"dotfiles"`
	DestinationPath string   `json:"destination-path"`
}

func CopyFile(filePath string, destinationPath string) {
	// Open the file and save the content
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatalf("Error closing the file %v: %v", filePath, err)
		}
	}(file)

	// Create a copy of the file to copy
	newFile, err := os.Create(filepath.Join(destinationPath, filepath.Base(filePath)))
	if err != nil {
		log.Fatal(err)
	}
	defer func(newFile *os.File) {
		err := newFile.Close()
		if err != nil {
			log.Fatalf("Error closing the file %v: %v", filepath.Base(filePath), err)
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

	fmt.Printf("%v copied successfully\n\n", filePath)
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

	destinationPath := dotfile.DestinationPath
	// Replace the ~ with the home folder path
	destinationPath = strings.Replace(destinationPath, "~", homeFolder, 1)

	// Check if the destination path exists
	if _, err := os.Stat(destinationPath); errors.Is(err, os.ErrNotExist) {
		fmt.Printf("dotty: %v: the directory does not exists.\n", destinationPath)
		os.Exit(3)
	}

	for i := 0; i < len(dotfile.Filepath); i++ {
		filePath := dotfile.Filepath[i]
		// Replace the ~ with the home folder path
		filePath = strings.Replace(filePath, "~", homeFolder, 1)

		if !strings.HasPrefix(filepath.Base(filePath), ".") {
			fmt.Printf("dotty: %v is not a dotfile\n", filePath)
			continue
		}

		if _, err := os.Stat(filePath); err == nil {
			fmt.Printf("- Copying %v...\n", filePath)
			CopyFile(filePath, destinationPath)
		} else if errors.Is(err, os.ErrNotExist) {
			fmt.Printf("dotty: %v does not exists", filePath)
			continue
		} else {
			log.Fatalf("dotty: error: %v", err)
		}
	}
}
