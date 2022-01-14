package job

import (
	"github.com/tyagip966/utility/fraazoError"
)

const JobInvalidInputError = "JOB_ERR_INVALID_INPUT"
const JobFailedError = "JOB_ERR_FAILED"

func BadRequestErrorFunc(errorMsg string) *fraazoError.Error {
	return &fraazoError.Error{
		ErrorCode:    JobInvalidInputError,
		ErrorMessage: errorMsg,
	}
}

func JobFailedErrorFunc(errorMsg string) *fraazoError.Error {
	return &fraazoError.Error{
		ErrorCode:    JobFailedError,
		ErrorMessage: errorMsg,
	}
}
