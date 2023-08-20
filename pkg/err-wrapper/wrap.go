package ewrap

import "github.com/ztrue/tracerr"

func Errorf(message string, args ...interface{}) error {
	return tracerr.Errorf(message, args...)
}

func Print(err error) {
	tracerr.Print(err)
}

func New(msg string) error {
	return tracerr.New(msg)
}

func Sprint(err error) string {
	return tracerr.Sprint(err)
}
