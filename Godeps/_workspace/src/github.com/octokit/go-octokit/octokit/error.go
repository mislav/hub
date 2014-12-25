package octokit

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/jingweno/go-sawyer"
)

type ResponseErrorType int

const (
	ErrorClientError             ResponseErrorType = iota // 400-499
	ErrorBadRequest              ResponseErrorType = iota // 400
	ErrorUnauthorized            ResponseErrorType = iota // 401
	ErrorOneTimePasswordRequired ResponseErrorType = iota // 401
	ErrorForbidden               ResponseErrorType = iota // 403
	ErrorTooManyRequests         ResponseErrorType = iota // 403
	ErrorTooManyLoginAttempts    ResponseErrorType = iota // 403
	ErrorNotFound                ResponseErrorType = iota // 404
	ErrorNotAcceptable           ResponseErrorType = iota // 406
	ErrorUnsupportedMediaType    ResponseErrorType = iota // 414
	ErrorUnprocessableEntity     ResponseErrorType = iota // 422
	ErrorServerError             ResponseErrorType = iota // 500-599
	ErrorInternalServerError     ResponseErrorType = iota // 500
	ErrorNotImplemented          ResponseErrorType = iota // 501
	ErrorBadGateway              ResponseErrorType = iota // 502
	ErrorServiceUnavailable      ResponseErrorType = iota // 503
	ErrorMissingContentType      ResponseErrorType = iota
	ErrorUnknownError            ResponseErrorType = iota
)

type ErrorObject struct {
	Resource string `json:"resource,omitempty"`
	Code     string `json:"code,omitempty"`
	Field    string `json:"field,omitempty"`
	Message  string `json:"message,omitempty"`
}

func (e *ErrorObject) Error() string {
	err := fmt.Sprintf("%v error", e.Code)
	if e.Field != "" {
		err = fmt.Sprintf("%v caused by %v field", err, e.Field)
	}
	err = fmt.Sprintf("%v on %v resource", err, e.Resource)
	if e.Message != "" {
		err = fmt.Sprintf("%v: %v", err, e.Message)
	}

	return err
}

type ResponseError struct {
	Response         *http.Response    `json:"-"`
	Type             ResponseErrorType `json:"-"`
	Message          string            `json:"message,omitempty"`
	Err              string            `json:"error,omitempty"`
	Errors           []ErrorObject     `json:"errors,omitempty"`
	DocumentationURL string            `json:"documentation_url,omitempty"`
}

func (e *ResponseError) Error() string {
	return fmt.Sprintf("%v %v: %d - %s",
		e.Response.Request.Method, e.Response.Request.URL,
		e.Response.StatusCode, e.errorMessage())
}

func (e *ResponseError) errorMessage() string {
	messages := []string{}

	if e.Message != "" {
		messages = append(messages, e.Message)
	}

	if e.Err != "" {
		m := fmt.Sprintf("Error: %s", e.Err)
		messages = append(messages, m)
	}

	if len(e.Errors) > 0 {
		m := []string{}
		m = append(m, "\nError summary:")
		for _, e := range e.Errors {
			m = append(m, fmt.Sprintf("\t%s", e.Error()))
		}
		messages = append(messages, strings.Join(m, "\n"))
	}

	if e.DocumentationURL != "" {
		messages = append(messages, fmt.Sprintf("// See: %s", e.DocumentationURL))
	}

	return strings.Join(messages, "\n")
}

func NewResponseError(resp *sawyer.Response) (err *ResponseError) {
	err = &ResponseError{}

	e := resp.Decode(&err)
	if e != nil {
		err.Message = fmt.Sprintf("Problems parsing error message: %s", e)
	}

	err.Response = resp.Response
	err.Type = getResponseErrorType(err)
	return
}

func getResponseErrorType(err *ResponseError) ResponseErrorType {
	code := err.Response.StatusCode
	header := err.Response.Header

	switch {
	case code == http.StatusBadRequest:
		return ErrorBadRequest

	case code == http.StatusUnauthorized:
		otp := header.Get("X-GitHub-OTP")
		r := regexp.MustCompile(`(?i)required; (\w+)`)
		if r.MatchString(otp) {
			return ErrorOneTimePasswordRequired
		}

		return ErrorUnauthorized

	case code == http.StatusForbidden:
		msg := err.Message
		rr := regexp.MustCompile("(?i)rate limit exceeded")
		if rr.MatchString(msg) {
			return ErrorTooManyRequests
		}
		lr := regexp.MustCompile("(?i)login attempts exceeded")
		if lr.MatchString(msg) {
			return ErrorTooManyLoginAttempts
		}

		return ErrorForbidden

	case code == http.StatusNotFound:
		return ErrorNotFound

	case code == http.StatusNotAcceptable:
		return ErrorNotAcceptable

	case code == http.StatusUnsupportedMediaType:
		return ErrorUnsupportedMediaType

	case code == 422:
		return ErrorUnprocessableEntity

	case code >= 400 && code <= 499:
		return ErrorClientError

	case code == http.StatusInternalServerError:
		return ErrorInternalServerError

	case code == http.StatusNotImplemented:
		return ErrorNotImplemented

	case code == http.StatusBadGateway:
		return ErrorBadGateway

	case code == http.StatusServiceUnavailable:
		return ErrorServiceUnavailable

	case code >= 500 && code <= 599:
		return ErrorServerError
	}

	return ErrorUnknownError
}
