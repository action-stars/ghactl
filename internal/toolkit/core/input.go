package core

import (
	"fmt"
	"os"
	"strings"
)

// InputOptions are options for the GetInput, GetMultilineInput and GetBooleanInput functions.
type InputOptions struct {
	// Required indicates if the input is required.
	// If required and not present, GetInput will return an error.
	Required bool
	// TrimWhitespace indicates if the input should be trimmed.
	// If nil, the input is trimmed by default.
	TrimWhitespace *bool
}

// GetInput gets the value of an input.
// Unless TrimWhitespace is set to false in InputOptions, the value is also trimmed.
// Returns an empty string if the value is not defined.
func GetInput(name string, opts InputOptions) (string, error) {
	envKey := "INPUT_" + strings.ToUpper(strings.ReplaceAll(name, " ", "_"))
	val := os.Getenv(envKey)

	if opts.Required && val == "" {
		return "", fmt.Errorf("input required and not supplied: %s", name)
	}

	if opts.TrimWhitespace != nil && !*opts.TrimWhitespace {
		return val, nil
	}

	return strings.TrimSpace(val), nil
}

// GetMultilineInput gets the values of a multiline input.
// Each value is also trimmed unless TrimWhitespace is set to false.
func GetMultilineInput(name string, opts InputOptions) ([]string, error) {
	val, err := GetInput(name, opts)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(val, "\n")
	result := make([]string, 0, len(lines))
	for _, line := range lines {
		if opts.TrimWhitespace != nil && !*opts.TrimWhitespace {
			if line != "" {
				result = append(result, line)
			}
		} else {
			trimmed := strings.TrimSpace(line)
			if trimmed != "" {
				result = append(result, trimmed)
			}
		}
	}

	return result, nil
}

// GetBooleanInput gets the input value of the boolean type in the YAML 1.2 "core schema" specification.
// Supported boolean input list: true | True | TRUE | false | False | FALSE.
func GetBooleanInput(name string, opts InputOptions) (bool, error) {
	val, err := GetInput(name, opts)
	if err != nil {
		return false, err
	}

	switch val {
	case "true", "True", "TRUE":
		return true, nil
	case "false", "False", "FALSE":
		return false, nil
	default:
		return false, fmt.Errorf(
			"input does not meet YAML 1.2 \"Core Schema\" specification: %s\n"+
				"Support boolean input list: `true | True | TRUE | false | False | FALSE`",
			name,
		)
	}
}
