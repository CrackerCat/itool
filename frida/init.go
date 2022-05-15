package frida

/*
 #cgo darwin CFLAGS: -g -O2 -w -I./macos
 #cgo linux CFLAGS: -g -O2 -w -I./linux
 #cgo darwin LDFLAGS: -framework Foundation -framework AppKit -lbsm -lresolv -L./macos -lfrida-core
 #cgo linux LDFLAGS: -static-libgcc -L./linux -lfrida-core -ldl -lm -lrt -lresolv -lpthread -Wl,--export-dynamic
 #include "frida-core.h"

 void cgo_on_message(FridaScript *script, const gchar *message, GBytes *data, gpointer user_data) {
	onMessage(script, message, data, user_data);
 }
*/
import "C"
import (
	"unsafe"
)

func IsNullCPointer(ptr unsafe.Pointer) bool {
	return uintptr(ptr) == uintptr(0)
}
