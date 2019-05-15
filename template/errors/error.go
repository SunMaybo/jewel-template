package errors

type HttpError struct {
	Code int
	Message string
}

func (e HttpError)Error() string {
	return e.Message
}
func New(code int,text string) error {
	return &HttpError{Code:code,Message:text}
}






