package rest

// ErrorResponse is a json struct for error response
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Error returns and error in a format, that most handler need
func Error(code int, msg string) (int, interface{}) {
	return code, &ErrorResponse{Code: code, Message: msg}
}
