package converter

import (
	"errors"
	"fmt"

	"sheep_forms/parser"
)

// Convert converts a parser.Form to the actual code implementation of the form
func Convert(form parser.Form, outputDir string) error {
	var (
		fileBytes []byte
		err       error
	)

	fmt.Printf("converting %s to %s", form.Header, string(form.Output))

	switch form.Output {
	case parser.OutputFlutter:
		fileBytes, err = ConvertFlutter(form)
		break
	case parser.OutputHtml:
		fileBytes, err = ConvertHtml(form)
		break
	default:
		return errors.New("unrecognized output type " + string(form.Output))
	}

	if err != nil {
		return err
	}

	fmt.Printf("created file with contents: %s", string(fileBytes))

	return nil
}

func ConvertFlutter(form parser.Form) ([]byte, error) {
	return nil, nil
}

func ConvertHtml(form parser.Form) ([]byte, error) {
	return nil, nil
}
