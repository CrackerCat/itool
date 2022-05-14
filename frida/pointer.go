package frida

// #include <stdlib.h>
import "C"
import (
	"sync"
	"unsafe"
)

type pointerMap struct {
	sync.Map
}

func (p *pointerMap) Store(v interface{}) unsafe.Pointer {
	if v == nil {
		return nil
	}

	var ptr unsafe.Pointer = C.malloc(C.size_t(1))
	if ptr == nil {
		panic("can't alloc 'cgo-pointer hack index pointer': ptr == nil")
	}

	p.Map.Store(ptr, v)

	return ptr
}

func (p *pointerMap) Load(ptr unsafe.Pointer) (v interface{}) {
	if ptr == nil {
		return nil
	}

	v, _ = p.Map.Load(ptr)

	return
}

func (p *pointerMap) Delete(ptr unsafe.Pointer) {
	if ptr == nil {
		return
	}

	p.Map.Delete(ptr)

	C.free(ptr)
}
