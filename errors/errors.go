package errors

import "fmt"

type Error struct {
	Code    int    `json:"code"`
	Name    string `json:"name"`
	Message string `json:"message"`
}

func NewError(code int, name string, message string) *Error {
	return &Error{
		Code:    code,
		Name:    name,
		Message: message,
	}
}

func (e *Error) Error() string {
	return fmt.Sprintf("code:%d name:%s message:%s", e.Code, e.Name, e.Message)
}

// IsNotFoundError notfound是一个比较重要的状态码，这里独立出来判断
func IsNotFoundError(err error) bool {
	if e, ok := err.(*Error); ok {
		return e.Code == 404
	}
	return false
}
