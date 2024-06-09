package gobs

import "errors"

var (
	ErrorServiceNotFound = errors.New("service not found")
	ErrorServiceRan      = errors.New("service has already run")
	ErrorServiceNotReady = errors.New("service is not ready")
)

func wrapCommonError(err error) error {
	switch err {
	case ErrorServiceRan, ErrorServiceNotReady:
		return nil
	default:
		return err
	}
}
