package errors

import "fmt"

type HttpError struct {
	Status int
	Message string
}

func (e HttpError)Error() string {
	 return fmt.Sprintf(`{"status":%d,"message":"%s"}`, e.Status, e.Message)
}
func New(code int,text string) error {
	return &HttpError{Status:code,Message:text}
}




