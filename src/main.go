package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"sheep_forms/converter"
	"sheep_forms/parser"
)

func main() {
	var (
		err         error
		response    string
		outResponse string
		files       []string
		forms       []parser.Form
	)

	if len(os.Args) == 1 {
		// Input Mode -- User did not supply any .sheepform in initial args
		fmt.Printf("Where are your input .sheepforms: ")
		_, err = fmt.Scanln(&response)
		if err != nil {
			log.Fatalf("error reading in response: %v", err)
			return
		}
		files = strings.Split(response, " ")

		fmt.Printf("Where is your output directory: ")
		_, err = fmt.Scanln(&outResponse)
		if err != nil {
			log.Fatalf("error while parsing output dir: %v", err)
			return
		}
	} else {
		// Straight to it Mode -- User supplied .sheepform in initial args
		files = os.Args[1:]
	}

	forms, err = ParseFiles(files)
	if err != nil {
		log.Fatalf("error while parsing %v", err)
	}

	print(len(forms))

	err = ConvertForms(forms, outResponse)
	if err != nil {
		log.Fatalf("error while coding form %v", err)
	}
}

// ParseFiles takes a list of files, reads the content, then attempts to parse to create desired output
func ParseFiles(files []string) ([]parser.Form, error) {
	var (
		forms     []parser.Form
		form      *parser.Form
		fileBytes []byte
		err       error
	)

	for _, file := range files {
		fileBytes, err = os.ReadFile(file)
		if err != nil {
			return nil, fmt.Errorf("file %s %v", file, err)
		}

		form, err = parser.Parse(fileBytes)
		if err != nil {
			return nil, fmt.Errorf("file %s %v", file, err)
		}

		forms = append(forms, *form)
	}
	return forms, nil
}

func ConvertForms(forms []parser.Form, outputDir string) error {
	for formIdx, form := range forms {
		err := converter.Convert(form, outputDir)
		if err != nil {
			return fmt.Errorf("form #%d (%s) %v", formIdx, form.Header, err)
		}
	}

	return nil
}
