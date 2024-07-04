package common

import "errors"

var (
	ErrorEndOfProcessing = errors.New("end of processing")
	ErrorServiceNotFound = errors.New("service not found")
	ErrorServiceRan      = errors.New("service has already run")
	ErrorServiceNotReady = errors.New("service is not ready")
	ErrorInvalidLength   = errors.New("invalid length")
	ErrorInvalidType     = errors.New("invalid type")
)
