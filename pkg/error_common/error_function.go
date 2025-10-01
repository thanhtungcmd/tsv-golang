package error_common

import "errors"

type ErrorResult struct {
	ErrorCode int    `json:"errorCode"`
	Message   string `json:"message"`
	Language  string `json:"language"`
}

func getErrors(code int) (map[string]string, error) {
	err, ok := messages[code]
	if !ok {
		return err, errors.New("error code not found")
	}
	return err, nil
}

func Error(code int) *ErrorResult {
	var e ErrorResult
	result, _ := getErrors(code)
	e.Language = "vi"
	e.ErrorCode = code
	e.Message = result[e.Language]
	return &e
}

func (e *ErrorResult) GetText() string {
	return e.Message
}
