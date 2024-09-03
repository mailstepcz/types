package iface

import (
	"reflect"
	"unsafe"
)

// Copiable is a type supporting custom copying.
type Copiable interface {
	CanCopyTo(reflect.Type) bool
	Copy(reflect.Type) unsafe.Pointer
}
