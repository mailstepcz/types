package iface

import (
	"database/sql/driver"
	"reflect"
	"unsafe"
)

// Copiable is a type supporting custom copying.
type Copiable interface {
	CanCopyTo(reflect.Type) bool
	Copy(reflect.Type) unsafe.Pointer
}

// ClosedEnum is a closed enum.
type ClosedEnum interface {
	driver.Valuer
	EnumValueIsValid() bool
	DefaultValue() string
}
