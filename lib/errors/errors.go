package errors

import (
	"errors"
	"log"
	"os"

	"github.com/go-playground/validator/v10"
)

func ParseError(err error) (errMsgs []string) {
	if !errors.As(err, &validator.ValidationErrors{}) {
		errMsgs = append(errMsgs, err.Error())
		return errMsgs
	}

	for _, err := range err.(validator.ValidationErrors) {
		errMsgs = append(errMsgs, err.Error())
	}

	return errMsgs
}

var (
	Error = log.New(os.Stdout, "ERROR: ", log.Llongfile)
	Warn  = log.New(os.Stdout, "WARN: ", 0)
	Info  = log.New(os.Stdout, "INFO: ", 0)
)

type LogWriter struct{}

func (f LogWriter) Write(p []byte) (n int, err error) {
	log.Printf("%s", p)
	return len(p), nil
}
