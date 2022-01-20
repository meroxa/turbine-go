package platform

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/meroxa/meroxa-go/pkg/meroxa"
	"net/http"
	"strings"
)

const FunctionsBasePath = "/v1/functions"

type FunctionState string

const (
	FunctionStatePending  FunctionState = "pending"
	FunctionStateStarting FunctionState = "starting"
	FunctionStateError    FunctionState = "error"
	FunctionStateReady    FunctionState = "ready"
)

type CreateFunctionInput struct {
	InputStream string             `json:"input_stream"`
	Image       string             `json:"image"`
	EnvVars     map[string]string  `json:"env_vars"`
	Args        []string           `json:"args"`
	Pipeline    PipelineIdentifier `json:"pipeline"`
}

type PipelineIdentifier struct {
	Name string `json:"name"`
}

type FunctionStatus struct {
	State   FunctionState `json:"state"`
	Details string        `json:"details"`
}

type Function struct {
	UUID         string            `json:"uuid"`
	Name         string            `json:"name"`
	InputStream  string            `json:"input_stream"`
	OutputStream string            `json:"output_stream"`
	Image        string            `json:"image"`
	Command      []string          `json:"command"`
	Args         []string          `json:"args"`
	EnvVars      map[string]string `json:"env_vars"`
	Status       FunctionStatus    `json:"status"`
	Pipeline     meroxa.Pipeline   `json:"pipeline"`
}

func (c *Client) CreateFunction(ctx context.Context, input *CreateFunctionInput) (*Function, error) {
	resp, err := c.MakeRequest(ctx, http.MethodPost, FunctionsBasePath, input, nil)
	if err != nil {
		return nil, err
	}

	err = handleAPIErrors(resp)
	if err != nil {
		return nil, err
	}

	var f Function
	err = json.NewDecoder(resp.Body).Decode(&f)
	if err != nil {
		return nil, err
	}

	return &f, nil
}

func handleAPIErrors(resp *http.Response) error {
	if resp.StatusCode > 204 {
		apiError, err := parseErrorFromBody(resp)

		// err if there was a problem decoding the resp.Body as the `errResponse` struct
		if err != nil {
			return err
		}

		// API error returned by Meroxa Platform API
		return apiError
	}
	return nil
}

type errResponse struct {
	Code    string              `json:"code,omitempty"`
	Message string              `json:"message,omitempty"`
	Details map[string][]string `json:"details,omitempty"`
}

func (err *errResponse) Error() string {
	msg := err.Message

	if errCount := len(err.Details); errCount > 0 {
		msg = fmt.Sprintf("%s. %d %s occurred:%s",
			msg,
			errCount,
			func() string {
				if errCount > 1 {
					return "problems"
				}
				return "problem"
			}(),
			mapToString(err.Details),
		)
	}
	return msg
}

func mapToString(m map[string][]string) string {
	s := ""
	count := 1
	for k, v := range m {
		s = fmt.Sprintf("%s\n%d. %s: \"%s\"", s, count, k, strings.Join(v, `", "`))
		count++
	}
	return s
}

func parseErrorFromBody(resp *http.Response) (error, error) {
	var er errResponse
	var body = resp.Body
	err := json.NewDecoder(body).Decode(&er)
	if err != nil {
		// In cases we didn't receive a proper JSON response
		if _, ok := err.(*json.SyntaxError); ok {
			return nil, errors.New(fmt.Sprintf("%s %s", resp.Proto, resp.Status))
		}

		return nil, err
	}

	return &er, nil
}
