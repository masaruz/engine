package common

import (
	"log"
	"os"
)

// Print only when stage is devevelopment
func Print(a ...interface{}) {
	if os.Getenv("stage") == "dev" {
		log.Println(a)
	}
}

// Printf only when stage is devevelopment
func Printf(format string, v ...interface{}) {
	if os.Getenv("stage") == "dev" {
		log.Printf(format, v)
	}
}
