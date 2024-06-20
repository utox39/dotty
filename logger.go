package main

import (
	"log"
	"os"
)

// LogErr is a custom logger
func LogErr(err error) {
	logger := log.New(os.Stderr, "dotty: ", 0)
	logger.Fatalln(err)
}
