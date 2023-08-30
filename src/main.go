package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"sheep_forms/parser"
)

func main() {
	var (
		err      error
		response string
		files    []string
	)

	if len(os.Args) == 1 {
		// Input Mode -- User did not supply any .sheepform in initial args
		fmt.Printf("Where are your input .sheepforms: ")
		_, err = fmt.Scanln(&response)
		if err != nil {
			fmt.Printf("error reading in response: %v", err)
			return
		}

		files = strings.Split(response, " ")
	} else {
		// Straight to it Mode -- User supplied .sheepform in initial args
		files = os.Args[1:]
	}

	err = ParseFiles(files)
	if err != nil {
		log.Fatalf("error while parsing files: %v", err)
	}
}

// ParseFiles takes a list of files, reads the content, then attempts to parse to create desired output
func ParseFiles(files []string) error {
	var (
		fileBytes []byte
		err       error
	)

	for _, file := range files {
		fileBytes, err = os.ReadFile(file)
		if err != nil {
			return err
		}

		err = parser.Parse(fileBytes)
		if err != nil {
			return err
		}
	}
	return nil
}
