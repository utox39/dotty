package main

import (
	"os"
	"testing"
)

func TestReadConfigValidJSONFile(t *testing.T) {
	// Create a temporary JSON file
	tmpFile, err := os.CreateTemp("", "config*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	// Remove the temp file
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			t.Fatalf("Failed to remove temporary file: %v", err)
		}
	}(tmpFile.Name())

	// Write valid JSON content to the temp file
	jsonContent := `
	{
  	"dotfiles" : [
    	".test"
  	],
  	"destination-path" : "destination/path/test"
	}`

	if _, err := tmpFile.Write([]byte(jsonContent)); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	if err := tmpFile.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	// Initialize Dotfiles struct
	var dotfiles Dotfiles

	// Invoke the code under test
	err = dotfiles.ReadConfig(tmpFile.Name())
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Validate the result
	if dotfiles.Dotfiles[0] != ".test" {
		t.Errorf("Expected dotfiles[0] to be '.test', got %v", dotfiles.Dotfiles[0])
	}

	if dotfiles.DestinationPath != "destination/path/test" {
		t.Errorf("Expected destination-path to be 'destination/path/test', got %v", dotfiles.Dotfiles[0])
	}
}

func TestReadConfigInvalidJSONFile(t *testing.T) {
	// Create a temporary JSON file
	tmpFile, err := os.CreateTemp("", "config*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	// Remove the temp file
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			t.Fatalf("Failed to remove temporary file: %v", err)
		}
	}(tmpFile.Name())

	// Write valid JSON content to the temp file
	jsonContent := `
	{
  	"invalid" : [
    	".test"
  	],
  	"invalid" : "destination/path/test"
	}`

	if _, err := tmpFile.Write([]byte(jsonContent)); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	if err := tmpFile.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	// Initialize Dotfiles struct
	var dotfiles Dotfiles

	// Invoke the code under test
	err = dotfiles.ReadConfig(tmpFile.Name())

	// Check dotfiles.Dotfiles[]
	if dotfiles.Dotfiles != nil {
		t.Fatalf("Expected to be nil, got %v", dotfiles.Dotfiles)
	}

	// Check dotfiles.DestinationPath
	if dotfiles.DestinationPath != "" {
		t.Fatalf("Expected to be empty, got %v", dotfiles.Dotfiles)
	}
}

func TestReadConfigNonExistentFilePath(t *testing.T) {
	// Initialize Dotfiles struct
	var dotfiles Dotfiles

	// Use a non-existent file path
	nonExistentPath := "/path/to/nonexistent/file.json"

	// Invoke the code under test
	err := dotfiles.ReadConfig(nonExistentPath)

	// Validate the result
	if err == nil {
		t.Errorf("Expected an error, got nil")
	}
}
