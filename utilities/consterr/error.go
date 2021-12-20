package consterr

import "fmt"

type Error string

func (err Error) Error() string {
	return string(err)
}

func Newf(format string, args ...interface{}) Error {
	return Error(fmt.Sprintf(format, args...))
}
