package errors

import "net/http"

var (
	// Discovery
	ErrorHeaders    = &Error{http.StatusBadRequest, "ErrorHeaders", "An internal error has occurred. Retry your request, but if the problem persists, contact us with details by posting a message on the Qchat forums"}
	ErrorParameters = &Error{http.StatusBadRequest, "ErrorParameters", "please check your parameters"}
	InternalError   = &Error{http.StatusInternalServerError, "ServerInternalError", "An internal error has occurred. Retry your request, but if the problem persists, contact us with details by posting a message on the Qchat forums"}
	InvalidArgument = &Error{http.StatusBadRequest, "InvalidArgument", "Invalid Argument"}
	MalformedJSON   = &Error{http.StatusBadRequest, "MalformedJSON", "The JSON you provided was not well-formed or did not validate against our published schema."}
	NotFoundError   = &Error{http.StatusNotFound, "NotFoundError", "not found"}
	AlreadyExist    = &Error{http.StatusBadRequest, "AlreadyExist", "Already exist!"}

	// UserId
	ErrorUserIdInvalid = &Error{http.StatusBadRequest, "ErrorUserIdInvalid", "userId invalid"}
	ErrorUserNotFound  = &Error{http.StatusBadRequest, "ErrorUserNotFound", "user not found"}

	// Group
	ErrorDeleteAdmin = &Error{http.StatusBadRequest, "ErrorDeleteAdmin", "you should not delete admin when someone exists"}
)

func BadRequestError(msg string) *Error {
	return &Error{
		Code:    http.StatusBadRequest,
		Name:    "BadRequestParams",
		Message: msg,
	}
}

func UnHandleError(msg string) *Error {
	return &Error{
		Code:    http.StatusInternalServerError,
		Name:    "UnHandleError",
		Message: msg,
	}
}
