package octokit

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

type gitHubError struct {
	Resource string      `json:"resource"`
	Field    string      `json:"field"`
	Value    interface{} `json:"value"`
	Code     string      `json:"code"`
	Message  string      `json:"message"`
}

type gitHubErrors struct {
	Message string        `json:"message"`
	Errors  []gitHubError `json:"errors"`
}

func handleErrors(body []byte) error {
	var githubErrors gitHubErrors
	err := json.Unmarshal(body, &githubErrors)
	if err != nil {
		return err
	}

	msg := buildErrorMessage(githubErrors)

	return errors.New(msg)
}

func buildErrorMessage(githubErrors gitHubErrors) string {
	errorMessages := make([]string, len(githubErrors.Errors))
	for _, e := range githubErrors.Errors {
		switch e.Code {
		case "custom":
			errorMessages = append(errorMessages, e.Message)
		case "missing_field":
			msg := fmt.Sprintf("Missing field: %s", e.Field)
			errorMessages = append(errorMessages, msg)
		case "invalid":
			msg := fmt.Sprintf("Invalid value for %s: %v", e.Field, e.Value)
			errorMessages = append(errorMessages, msg)
		case "unauthorized":
			msg := fmt.Sprintf("Not allow to change field %s", e.Field)
			errorMessages = append(errorMessages, msg)
		}
	}

	msg := strings.Join(errorMessages, "\n")
	if msg == "" {
		msg = githubErrors.Message
	}

	return msg
}
