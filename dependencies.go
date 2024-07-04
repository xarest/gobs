package gobs

import (
	"fmt"
	"reflect"

	"github.com/traphamxuan/gobs/common"
)

type Dependencies []IService

func (d Dependencies) Assign(pTargets ...any) error {
	for id := 0; id < len(pTargets) && id < len(d); id++ {
		dep := d[id]
		dst := pTargets[id]
		if dst == nil {
			continue
		}
		dstType := reflect.TypeOf(dst).Elem()
		if dstType.Kind() != reflect.Ptr {
			return fmt.Errorf("require pointer type. But got %T %w", dst, common.ErrorInvalidType)
		}

		// Check if the dependency is assignable to the target
		depType := reflect.TypeOf(dep)
		if !depType.AssignableTo(dstType) {
			return fmt.Errorf("require %s but got %s %w",
				depType.String(), dstType.String(),
				common.ErrorInvalidType)
		}
		reflect.ValueOf(dst).Elem().Set(reflect.ValueOf(dep))
	}
	return nil
}
