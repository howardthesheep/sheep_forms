package parser

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

func Parse(bytes []byte) error {
	var (
		style  string
		output string
		form   Form
	)

	// Define our regex that will be used in parsing
	styleRx := regexp.MustCompile(`(?i)(style)(?-i):([a-zA-Z]+)`)
	outputRx := regexp.MustCompile(`(?i)(output)(?-i):([a-zA-Z]+)`)

	// Evaluate on line-by-line basis
	lineTokens := strings.Split(string(bytes), "\n")
	maxIdx := len(lineTokens) - 1
	for idx, line := range lineTokens {
		// Skip blank lines or lines only containing return characters
		if line == "" || line == "\r" {
			continue
		}

		// Look for `style:` assignment if provided & not already found
		if style == "" {
			styleMatches := styleRx.FindStringSubmatch(line)
			if len(styleMatches) > 0 {
				style = styleMatches[len(styleMatches)-1]
				fmt.Println("Using provided style: " + style)
				continue
			}
		}

		// Look for `output:` assignment if provided & not already found
		if output == "" {
			outputMatches := outputRx.FindStringSubmatch(line)
			if len(outputMatches) > 0 {
				output = outputMatches[len(outputMatches)-1]
				fmt.Println("Using provided output format: " + output)
				continue
			}
		}

		// Only look for header if we haven't set it yet, and we've not identified any form sections or form inputs
		// since the form header will always precede form sections & inputs
		if form.Header == "" && len(form.Sections) == 0 && len(form.Inputs) == 0 {
			// Handle if idx+1 would cause arr out of bounds, then we know we have our header,
			// and we don't need to see if this is a section header
			if idx+1 > maxIdx {
				form.Header = line
				continue
			}

			// If proceeded by `\t<`, this is a section header
			if strings.HasPrefix(lineTokens[idx+1], "\t<") {
				form.Sections = append(form.Sections, FormSection{
					Header: line,
				})
			} else {
				form.Header = line
			}
			continue
		}

		// Now look for root level form inputs if they exist. Root level meaning they are not nested within a section
		if strings.HasPrefix(line, "<") {
			fInput, err := ParseInput(line) // TODO: Support multiple returned Inputs
			if err != nil {
				return err
			}

			form.Inputs = append(form.Inputs, *fInput)
			continue
		}

		// Look for form input options attached to root level form inputs. If there is no root level form inputs,
		// then we know that this is a section input
		if strings.HasPrefix(line, "\t<") {
			if strings.HasPrefix(lineTokens[idx-1], "<") {
				recentInputOptions := &form.Inputs[len(form.Inputs)-1].Options
				*recentInputOptions = append(*recentInputOptions, strings.Trim(line, "\t<>"))
			} else {
				recentSection := len(form.Sections) - 1
				rows := &form.Sections[recentSection].InputRows
				fInput, err := ParseInput(strings.Trim(line, "\t")) // TODO: Support multiple returned Inputs
				if err != nil {
					return err
				}

				*rows = append(*rows, FormInputRow{
					Inputs: []FormInput{fInput},
				})
			}
			continue
		}

		// The "\t\t<" prefix denotes that  this is a section input option
		if strings.HasPrefix(line, "\t\t<") {
			// Get reference to the latest Section > Input Row > Input
			newestSectIdx := len(form.Sections) - 1
			newestInputRowIdx := len(form.Sections[newestSectIdx].InputRows) - 1
			newestInputIdx := len(form.Sections[newestSectIdx].InputRows[newestInputRowIdx].Inputs) - 1
			options := &form.Sections[newestSectIdx].InputRows[newestInputRowIdx].Inputs[newestInputIdx].Options

			// Append our option to form data structure
			*options = append(*options, strings.Trim(line, "\r\n\t<>"))
			continue
		}

		form.Sections = append(form.Sections, FormSection{
			Header:    line,
			InputRows: nil,
		})
	}

	// If file did not define a style, use the default (Material)
	if style == "" {
		style = string(StyleMaterial)
		fmt.Printf("No style provided, using default (%s)\n", style)
	}

	// If file did not define an output, use the default (Flutter)
	if output == "" {
		output = string(OutputFlutter)
		fmt.Printf("No output format provided, using default (%s)\n", output)
	}

	form.Style = FormStyle(style)
	form.Output = FormOutput(output)

	DebugForm(form)

	return nil
}

// ParseInput returns the FormInput that was parsed from a line of text
// TODO: Support if multiple Inputs defined per line
func ParseInput(line string) (*FormInput, error) {
	var fi FormInput
	noCarats := strings.Trim(line, "<>")
	inputTokens := strings.Split(noCarats, ",")
	switch len(inputTokens) {
	case 1:
		// No type give, default to text
		fi = FormInput{
			Title: inputTokens[0],
			Type:  text,
		}
		break
	case 2:
		inputType, err := ResolveType(inputTokens[1])
		if err != nil {
			return nil, err
		}
		fi = FormInput{
			Title: inputTokens[0],
			Type:  *inputType,
		}
		break
	case 3:
		inputType, err := ResolveType(inputTokens[1])
		if err != nil {
			return nil, err
		}
		fi = FormInput{
			Title:      inputTokens[0],
			Type:       *inputType,
			Attributes: ParseAttributes(inputTokens[2:]),
		}
	default:
	}

	return &fi, nil
}

func ParseAttributes(input []string) map[string]string {
	return map[string]string{}
}

// ResolveType takes a string representation of a FormInputType and resolve that to a valid FormInputType
func ResolveType(input string) (*FormInputType, error) {
	var iType FormInputType
	trimmed := strings.TrimSpace(input)
	lowered := strings.ToLower(trimmed)
	switch lowered {
	case "text":
		iType = text
		break
	case "int":
		iType = integer
		break
	// TODO: support all cases of FormInputTypes
	default:
		return nil, errors.New("Unrecognized FormInputType provided: " + lowered)
	}

	return &iType, nil
}

// DebugForm prints the contents of the form to stdout
func DebugForm(form Form) {
	fmt.Println("header: " + form.Header)
	fmt.Println("style: " + form.Style)
	fmt.Println("output: " + form.Output)

	for _, row := range form.InputRows {
		fmt.Printf("row: %s\n", row.Inputs)
	}

	for _, section := range form.Sections {
		fmt.Printf("section: %s\n", section.Header)
		for _, row := range section.InputRows {
			fmt.Printf("\t row: ")
			for _, input := range row.Inputs {
				fmt.Print("row: " + input.Title + "\n")
				for _, option := range input.Options {
					fmt.Printf("\t option: %s\n", option)
				}
			}
		}
	}
}
