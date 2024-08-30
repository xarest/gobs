package utils

import (
	"reflect"
	"strings"

	"github.com/xarest/gobs/common"
)

func CompactName(name string) string {
	names := strings.Split(name, "/")
	return names[len(names)-1]
}

func DefaultServiceName(s any) string {
	t := reflect.TypeOf(s)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t.PkgPath() + "." + t.Name()
}

func WrapCommonError(err error) error {
	switch err {
	case common.ErrorServiceRan:
		return nil
	default:
		return err
	}
}
