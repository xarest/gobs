package utils

import (
	"context"
	"reflect"
	"strings"
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

func EmptyFunc(ctx context.Context) error {
	return nil
}
