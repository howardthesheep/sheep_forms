package parser

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

func Parse(bytes []byte) (*Form, error) {
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
		if form.Header == "" && len(form.Sections) == 0 && len(form.InputRows) == 0 {
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
			fInputs, err := ParseInputs(line, idx)
			if err != nil {
				return nil, err
			}

			form.InputRows = append(form.InputRows, FormInputRow{Inputs: fInputs})
			continue
		}

		// Look for form input options attached to root level form inputs. If there is no root level form inputs,
		// then we know that this is a section input
		if strings.HasPrefix(line, "\t<") {
			if strings.HasPrefix(lineTokens[idx-1], "<") {
				// Get the most recently created input row
				recentInputRow := &form.InputRows[len(form.InputRows)-1]
				recentInputs := recentInputRow.Inputs
				recentInputOptions := &recentInputs[len(recentInputs)-1].Options
				*recentInputOptions = append(*recentInputOptions, strings.Trim(line, "\t<>"))
			} else {
				recentSection := len(form.Sections) - 1
				rows := &form.Sections[recentSection].InputRows
				fInputs, err := ParseInputs(strings.Trim(line, "\t"), idx)
				if err != nil {
					return nil, err
				}

				*rows = append(*rows, FormInputRow{
					Inputs: fInputs,
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
			if newestInputIdx > 0 {
				// Throw error if user specified two form inputs, with one being dropdown. This causes ambiguity
				// when figuring out where the proceeding FormInputOptions belong...
				// Maybe this could be done better, but I'm deciding this for now.
				errorStr := fmt.Sprintf("line #%d: cannot have two from inputs on the same line if one is dropdown", idx)
				return nil, errors.New(errorStr)
			}

			// Append our option to form data structure, cleaning up unnecessary syntax characters
			*options = append(*options, strings.Trim(line, "\r\n\t<>"))
			continue
		}

		// If none of the above conditionals triggered, we know were left with a Section header
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
	form.Style = FormStyle(style)

	// If file did not define an output, use the default (Flutter)
	if output == "" {
		output = string(OutputFlutter)
		fmt.Printf("No output format provided, using default (%s)\n", output)
	}
	form.Output = FormOutput(output)

	// DebugForm(form)

	return &form, nil
}

// ParseInputs returns an []FormInput that was parsed from a line of text
func ParseInputs(line string, idx int) ([]FormInput, error) {
	var fis []FormInput
	dirtyInputs := strings.Split(line, "<")
	for _, input := range dirtyInputs {
		if input == "\t" || input == "" {
			continue
		}

		var fi FormInput
		leftTrimmed := strings.TrimLeft(input, " <")
		trimmed := strings.TrimRight(leftTrimmed, "> \r\n")
		inputTokens := strings.Split(trimmed, ",")
		switch len(inputTokens) {
		case 1:
			// No type give, default to text
			fi = FormInput{
				Title: inputTokens[0],
				Type:  text,
			}
			break
		case 2:
			inputType, err := ResolveType(inputTokens[1], idx)
			if err != nil {
				return nil, err
			}
			fi = FormInput{
				Title: inputTokens[0],
				Type:  *inputType,
			}
			break
		case 3:
			inputType, err := ResolveType(inputTokens[1], idx)
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

		fis = append(fis, fi)
	}

	return fis, nil
}

// ParseAttributes TODO: parses all of the potential user created key=value pairs present in a FormInput
func ParseAttributes(input []string) map[string]string {
	return map[string]string{}
}

// ResolveType takes a string representation of a FormInputType and resolve that to a valid FormInputType
func ResolveType(input string, idx int) (*FormInputType, error) {
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
	case "double":
		iType = double
		break
	case "phone":
		iType = phone
		break
	case "email":
		iType = email
		break
	case "dropdown":
		iType = dropdown
		break
	case "date":
		iType = date
		break
	case "date range":
		iType = dateRange
		break
	case "checkbox":
		iType = checkbox
		break
	case "tri-state box":
		iType = tristateBox
		break
	case "files":
		iType = files
		break
	case "images":
		iType = images
		break
	case "time":
		iType = time
		break
	case "rich text":
		iType = richText
		break
	case "slider":
		iType = slider
		break
	case "captcha":
		iType = captcha
		break
	case "color":
		iType = color
		break
	case "credit card":
		iType = creditCard
		break
	case "address":
		iType = address
		break
	case "search and select":
		iType = searchAndSelect
		break
	case "progress bar":
		iType = progressBar
		break
	default:
		errorStr := fmt.Sprintf("line #%d: unrecognized FormInputType provided: %s", idx, lowered)
		return nil, errors.New(errorStr)
	}

	return &iType, nil
}

// DebugForm prints the contents of the form to stdout
func DebugForm(form Form) {
	fmt.Println("header: " + form.Header)
	fmt.Println("style: " + form.Style)
	fmt.Println("output: " + form.Output)

	// Print root level form inputs
	for _, row := range form.InputRows {
		for _, input := range row.Inputs {
			fmt.Print("input: " + input.Title + " ")
			for _, option := range input.Options {
				fmt.Printf("\n\t\t option: %s", option)
			}
		}
		fmt.Print("\n")
	}

	// Print form sections
	for _, section := range form.Sections {
		fmt.Printf("section: %s\n", section.Header)
		for _, row := range section.InputRows {
			fmt.Printf("\t row: ")
			for _, input := range row.Inputs {
				fmt.Print("input: " + input.Title + " ")
				for _, option := range input.Options {
					fmt.Printf("\n\t\t option: %s", option)
				}
			}
			fmt.Print("\n")
		}
	}
}
