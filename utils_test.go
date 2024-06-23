package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestValidatePathValidPath(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "example-*.txt")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	// Remove the temp file
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			t.Fatalf("Failed to remove temporary file: %v", err)
		}
	}(tmpFile.Name())

	// Get the temp file path
	path := tmpFile.Name()

	// Add some extra slashes to test cleaning
	path = filepath.Join(path, "//")

	// Validate the path
	err = ValidatePath(&path)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Check if the path is cleaned
	expectedPath := filepath.Clean(tmpFile.Name())
	if path != expectedPath {
		t.Errorf("Expected path %v, got %v", expectedPath, path)
	}
}

func TestValidatePathInvalidPath(t *testing.T) {
	// Define a non-existent path
	path := "non/existent/path"

	// Validate the path
	err := ValidatePath(&path)
	if err == nil {
		t.Errorf("Expected an error, got nil")
	}

	// Check if the error is of type os.IsNotExist
	if !os.IsNotExist(err) {
		t.Errorf("Expected os.IsNotExist error, got %v", err)
	}
}

func TestReplaceTildeStringWithTilde(t *testing.T) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("error getting user home directory: %v", err)
	}

	path := "~/test/path"
	expectedPath := strings.Replace(path, "~", homeDir, 1)

	err = ReplaceTilde(&path)
	if err != nil {
		t.Fatalf("ReplaceTilde returned an error: %v", err)
	}

	if path != expectedPath {
		t.Errorf("expected %s, got %s", expectedPath, path)
	}
}

func TestReplaceTildeStringWithoutTilde(t *testing.T) {
	path := "test/path"
	expectedPath := "test/path"

	err := ReplaceTilde(&path)
	if err != nil {
		t.Fatalf("ReplaceTilde returned an error: %v", err)
	}

	if path != expectedPath {
		t.Errorf("expected %s, got %s", expectedPath, path)
	}
}

func TestReplaceTildeWithEmptyString(t *testing.T) {
	path := ""
	expectedPath := ""

	err := ReplaceTilde(&path)
	if err != nil {
		t.Fatalf("ReplaceTilde returned an error: %v", err)
	}

	if path != expectedPath {
		t.Errorf("expected %s, got %s", expectedPath, path)
	}
}

func TestReadFileExistentFile(t *testing.T) {
	content := []byte("This is a test file.\n")

	tmpFile, err := os.CreateTemp("", "example-*.txt")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	// Remove the temp file
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			t.Fatalf("Failed to remove temporary file: %v", err)
		}
	}(tmpFile.Name())

	// Writes the test content to the temp file
	if _, err := tmpFile.Write(content); err != nil {
		t.Fatalf("Failed to write to temporary file: %v", err)
	}

	err = tmpFile.Close()
	if err != nil {
		t.Fatalf("Failed to close temporary file: %v", err)
	}

	path := tmpFile.Name()

	fileContent, err := ReadFile(&path)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if !bytes.Equal(fileContent, content) {
		t.Errorf("Expected %s, got %s", content, fileContent)
	}
}

func TestReadFileNonExistentPath(t *testing.T) {
	path := "testdata/non-existent-file.txt"

	_, err := ReadFile(&path)
	if err == nil {
		t.Fatalf("Expected an error, got nil")
	}

	if !os.IsNotExist(err) {
		t.Errorf("Expected a non-existent file error, got %v", err)
	}
}
