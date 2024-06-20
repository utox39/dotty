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

	"github.com/tidwall/sjson"
	"github.com/urfave/cli/v2"
)

type Dotfiles struct {
	Filepath        []string `json:"dotfiles"`
	DestinationPath string   `json:"destination-path"`
}

func LogError(err error) {
	log.Fatalf("dotty: error: %v\n", err)
}

func Backup() error {
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
			log.Fatalf("dotty: error closing the file %v: %v", jsonFile, err)
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
		fmt.Printf("dotty: %v: the directory does not exist.\n", destinationPath)
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
			fmt.Printf("dotty: %v does not exist", filePath)
			continue
		} else {
			LogError(err)
		}
	}

	return nil
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
			log.Fatalf("dotty: error closing the file %v: %v", filePath, err)
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
			log.Fatalf("dotty: error closing the file %v: %v", filepath.Base(filePath), err)
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

func AddFile(newFilePath string) error {
	if _, err := os.Stat(newFilePath); errors.Is(err, os.ErrNotExist) {
		fmt.Printf("dotty: %v does not exist\n", newFilePath)
	} else if err != nil {
		LogError(err)
	}

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
			log.Fatalf("dotty: error closing the file %v: %v", jsonFile, err)
		}
	}(jsonFile)

	byteValue, _ := io.ReadAll(jsonFile)

	// Add the new dotfile to the end of the dotfiles array
	value, _ := sjson.Set(string(byteValue), "dotfiles.-1", newFilePath)

	err = os.WriteFile(jsonConfigPath, []byte(value), 0644)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func main() {
	app := &cli.App{
		Name:  "dotty",
		Usage: "backup your dotfiles of choice in a folder",
		Commands: []*cli.Command{
			{
				Name:  "add",
				Usage: "add a dotfile to the list",
				Action: func(cCtx *cli.Context) error {
					err := AddFile(cCtx.Args().Get(0))
					if err != nil {
						LogError(err)
					}
					return nil
				},
			},
		},
		Action: func(*cli.Context) error {
			err := Backup()
			if err != nil {
				LogError(err)
			}
			return nil
		},
		Version: "0.1",
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
