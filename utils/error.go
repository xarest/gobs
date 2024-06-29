package utils

import "errors"

var (
	ErrorEndOfProcessing = errors.New("end of processing")
	ErrorServiceNotFound = errors.New("service not found")
	ErrorServiceRan      = errors.New("service has already run")
	ErrorServiceNotReady = errors.New("service is not ready")
)

func WrapCommonError(err error) error {
	switch err {
	case ErrorServiceRan:
		return nil
	default:
		return err
	}
}
