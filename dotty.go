package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/urfave/cli/v2"
)

func Backup() error {
	fmt.Println("Hello world from dotty :)")

	dotfiles := &Dotfiles{}

	// Get home folder
	homeFolder, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("error getting user home folder %v", err)
	}

	// Read json config file
	jsonConfigPath := filepath.Join(homeFolder, ".config/dotty/config.json")

	err = dotfiles.ReadConfig(jsonConfigPath)
	if err != nil {
		return err
	}

	// The destination path specified in the config file
	destinationPath := dotfiles.DestinationPath

	// Replace the ~ with the home folder path
	err = ReplaceTilde(&destinationPath)
	if err != nil {
		return err
	}

	// Check if the destination path exists
	err = ValidatePath(&destinationPath)
	if err != nil {
		return err
	}

	for i := 0; i < len(dotfiles.Dotfiles); i++ {
		filePath := dotfiles.Dotfiles[i]

		// Replace the ~ with the home folder path
		err = ReplaceTilde(&filePath)
		if err != nil {
			return err
		}

		err = ValidatePath(&filePath)
		if errors.Is(err, os.ErrNotExist) {
			fmt.Printf("dotty: %v does not exist", filePath)
			continue
		} else if err != nil {
			return err
		}

		if !strings.HasPrefix(filepath.Base(filePath), ".") {
			fmt.Printf("dotty: %v is not a dotfile\n", filePath)
			continue
		}

		fmt.Printf("- Copying %v...\n", filePath)
		if err := CopyFile(filePath, destinationPath); err != nil {
			return err
		}
	}

	return nil
}

func CopyFile(filePath string, destinationPath string) error {
	// Open the file and save the content
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("could not open file %v: %v", filePath, err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			LogErr(fmt.Errorf("error closing the file %v: %v", filePath, err))
		}
	}(file)

	// Create a copy of the file to copy
	newFile, err := os.Create(filepath.Join(destinationPath, filepath.Base(filePath)))
	if err != nil {
		return fmt.Errorf("could not create file %v: %v", newFile, err)
	}
	defer func(newFile *os.File) {
		err := newFile.Close()
		if err != nil {
			LogErr(fmt.Errorf("error closing the file %v: %v", filepath.Base(filePath), err))
		}
	}(newFile)

	// Write the content of the file to copy in the copy file
	bytesWritten, err := io.Copy(newFile, file)
	if err != nil {
		return fmt.Errorf("could not copy file %v to %v: %v", file, newFile, err)
	}
	fmt.Printf("Copied %d bytes.\n", bytesWritten)

	err = newFile.Sync()
	if err != nil {
		return fmt.Errorf("error syncing file: %v", err)
	}

	fmt.Printf("%v copied successfully\n\n", filePath)

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
					// Get home folder
					homeFolder, err := os.UserHomeDir()
					if err != nil {
						return fmt.Errorf("error getting user home folder %v", err)
					}

					// Read json config file
					jsonConfigPath := filepath.Join(homeFolder, ".config/dotty/config.json")

					err = AddFile(cCtx.Args().Get(0), jsonConfigPath)
					if err != nil {
						LogErr(err)
					}
					return nil
				},
			},
			{
				Name:  "backup",
				Usage: "backup your dotfiles",
				Action: func(cCtx *cli.Context) error {
					err := Backup()
					if err != nil {
						LogErr(err)
					}
					return nil
				},
			},
		},
		Action:  nil,
		Version: "1.2.0",
	}

	if err := app.Run(os.Args); err != nil {
		LogErr(err)
	}
}
