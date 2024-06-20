package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// ValidatePath clean and check the path
func ValidatePath(path *string) error {
	// Replace the ~ with the home folder path
	if strings.HasPrefix(*path, "~") {
		homeFolder, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("error getting user home folder %v", err)
		}

		*path = strings.Replace(*path, "~", homeFolder, 1)
	}

	// Clean the path
	*path = filepath.Clean(*path)

	// Check that the path exists
	if _, err := os.Stat(*path); os.IsNotExist(err) {
		return err
	}

	return nil
}

// ReadFile returns the contents of a file
func ReadFile(path *string) ([]byte, error) {
	// Validate the path
	if err := ValidatePath(path); err != nil {
		return nil, err
	}

	// Open the file
	f, err := os.Open(*path)
	if err != nil {
		return nil, err
	}

	// Close the file
	defer func(f *os.File) {
		if err := f.Close(); err != nil {
			LogErr(err)
		}
	}(f)

	// Read the file content
	byteValue, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	return byteValue, nil
}
