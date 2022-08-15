package jorge

import (
	"fmt"

	log "github.com/sirupsen/logrus"
)

type EncapsulatedError struct {
	OriginalErr error
	Message     string
	Solution    string
	Code        int
}

func NewEncapsulatedError(err error, message string, solution string) *EncapsulatedError {
	return &EncapsulatedError{
		OriginalErr: err,
		Message:     message,
		Solution:    solution,
	}
}

func (receiver EncapsulatedError) print(debug bool) {
	if debug || len(receiver.Message) == 0 {
		log.Debug(receiver.OriginalErr.Error())
	}

	if len(receiver.Message) > 0 {
		fmt.Println(receiver.Message)
	}

	if len(receiver.Solution) > 0 {
		fmt.Println(receiver.Solution)
	}
}
