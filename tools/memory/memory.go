package memory

/*
#include <stdlib.h>
*/
import "C"
import "unsafe"

func Alloc(size uintptr) *byte {
	return (*byte)(C.malloc(C.size_t(size)))
}

func Free(ptr *byte, size uintptr) {
	C.free(unsafe.Pointer(ptr))
}
