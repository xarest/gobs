package gobs

import (
	"fmt"
	"reflect"

	"github.com/xarest/gobs/common"
)

type Dependencies []IService

func (d Dependencies) Assign(pTargets ...IService) error {
	for id := 0; id < len(pTargets) && id < len(d); id++ {
		dep := d[id]
		dst := pTargets[id]
		if dst == nil {
			continue
		}

		dstType := reflect.TypeOf(dst).Elem()
		// If the destination is a pointer to a struct, we can directly assign pointers address to the pointer's value
		if dstType.Kind() == reflect.Ptr {
			depType := reflect.TypeOf(dep)
			if !depType.AssignableTo(dstType) {
				return fmt.Errorf("require %s but got %s %w",
					depType.String(), dstType.String(),
					common.ErrorInvalidType)
			}
			reflect.ValueOf(dst).Elem().Set(reflect.ValueOf(dep))
			continue
		}

		dstValue := reflect.ValueOf(dst)
		// If the destination is a pointer to an interface, we can assign the interface to the pointer's value
		if dstValue.Kind() != reflect.Ptr || dstValue.IsNil() {
			return fmt.Errorf("destination must be a non-nil pointer: %w", common.ErrorInvalidType)
		}

		dstElemType := dstValue.Elem().Type()
		if dstElemType.Kind() != reflect.Interface {
			return fmt.Errorf("destination must be a pointer to an interface: %w", common.ErrorInvalidType)
		}

		depValue := reflect.ValueOf(dep)
		if !depValue.Type().Implements(dstElemType) {
			return fmt.Errorf("dependency does not implement the interface: %w", common.ErrorInvalidType)
		}

		dstValue.Elem().Set(depValue)
	}
	return nil
}
