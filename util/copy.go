package util

import (
	"reflect"
	"unsafe"
)

// FastCopyString copies strings with pointers which are much faster.
func FastCopyString(s string) string {
	var b []byte
	h := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	h.Data = (*reflect.StringHeader)(unsafe.Pointer(&s)).Data
	h.Len = len(s)
	h.Cap = len(s)
	return string(b)
}
