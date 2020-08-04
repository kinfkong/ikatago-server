package errors

type HttpError struct {
	StatusCode int    `json:"statusCode"`
	Msg        string `json:"msg"`
	Key        string `json:"key"`
}

func (e *HttpError) Error() string {
	return e.Msg
}

func CreateErrorWithMsg(status int, key string, msg string) error {
	return &HttpError{StatusCode: status, Msg: msg, Key: key}
}
func CreateError(status int, key string) error {
	return &HttpError{StatusCode: status, Msg: key, Key: key}
}
