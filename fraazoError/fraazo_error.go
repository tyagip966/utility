package fraazoError

type additionalData interface{}

type Error struct {
	ErrorCode      string         `json:"errorCode"`
	ErrorMessage   string         `json:"errorMessage"`
	AdditionalData additionalData `json:"additionalData,omitempty"`
}

type MFAError struct {
	HttpStatusCode int
	OptimusError   Error
}

func (r Error) Error() string {
	return "ErrorCode: " + r.ErrorCode + " ErrorMessage: " + r.ErrorMessage
}

func New(errorCode, errorMessage string, additionalData interface{}) Error {
	return Error{
		ErrorCode:      errorCode,
		ErrorMessage:   errorMessage,
		AdditionalData: additionalData,
	}
}
