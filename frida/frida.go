package frida

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"unsafe"
)

/*
#include "frida-core.h"
void cgo_on_message(FridaScript *script, const gchar *message, GBytes *data, gpointer user_data);
*/
import "C"

var pointers = &pointerMap{}

func Start(ctx context.Context, udid, bundleID, source string, callback func(string, []byte)) error {
	if len(strings.TrimSpace(bundleID)) == 0 {
		return errors.New("BundleID 不能为空")
	}

	if _, err := C.frida_init(); err != nil {
		return fmt.Errorf("frida_init error: %v", err)
	}
	defer C.frida_deinit()

	loop := C.g_main_loop_new(nil, C.int(1))
	defer C.g_main_loop_unref(loop)

	devManager := C.frida_device_manager_new()
	defer func() {
		C.frida_device_manager_close_sync(devManager, nil, nil)
		C.frida_unref(C.gpointer(devManager))
	}()

	var gerr *C.GError
	devList := C.frida_device_manager_enumerate_devices_sync(devManager, nil, &gerr)
	if gerr != nil {
		return fmt.Errorf("enumerate devices error: %v", C.GoString(gerr.message))
	}

	var device *C.FridaDevice
	count := int(C.frida_device_list_size(devList))
	for i := 0; i < count; i++ {
		device = C.frida_device_list_get(devList, C.int(i))
		id := C.GoString(C.frida_device_get_id(device))
		if id != udid {
			C.g_object_unref(C.gpointer(device))
			device = nil
			continue
		}
	}
	C.frida_unref(C.gpointer(devList))

	if device == nil {
		return errors.New("frida device not found")
	}
	defer C.frida_unref(C.gpointer(device))

	bid := C.CString(bundleID)
	defer C.free(unsafe.Pointer(bid))
	pid := C.frida_device_spawn_sync(device, bid, nil, nil, &gerr)
	if gerr != nil {
		return fmt.Errorf("spawn target error: %v", C.GoString(gerr.message))
	}

	session := C.frida_device_attach_sync(device, pid, C.FRIDA_REALM_NATIVE, nil, &gerr)
	if gerr != nil {
		return fmt.Errorf("attach target error: %v", C.GoString(gerr.message))
	}
	defer func() {
		C.frida_session_detach_sync(session, nil, nil)
		C.frida_unref(C.gpointer(session))
	}()

	C.frida_device_resume_sync(device, pid, nil, nil)

	jssrc := C.CString(source)
	defer C.free(unsafe.Pointer(jssrc))
	script := C.frida_session_create_script_sync(session, jssrc, nil, nil, &gerr)
	if gerr != nil {
		return fmt.Errorf("create script error: %v", C.GoString(gerr.message))
	}
	defer func() {
		C.frida_script_unload_sync(script, nil, nil)
		C.frida_unref(C.gpointer(script))
	}()

	msg := C.CString("message")
	defer C.free(unsafe.Pointer(msg))
	C.g_signal_connect_data(C.gpointer(script), msg, C.GCallback(C.cgo_on_message),
		C.gpointer(pointers.Store(callback)), nil, 0)

	C.frida_script_load_sync(script, nil, &gerr)
	if gerr != nil {
		return fmt.Errorf("load script error: %v", C.GoString(gerr.message))
	}

	go func() {
		<-ctx.Done()
		C.g_main_loop_quit(loop)
	}()

	ok := C.g_main_loop_is_running(loop)
	if ok == 1 {
		C.g_main_loop_run(loop)
	}

	return nil
}

//export onMessage
func onMessage(script *C.FridaScript, message *C.gchar, data *C.GBytes, userData C.gpointer) {
	callback, ok := pointers.Load(unsafe.Pointer(userData)).(func(string, []byte))
	if !ok {
		return
	}

	bs := make([]byte, 0)
	if !IsNullCPointer(unsafe.Pointer(data)) {
		var size C.ulong
		buf := C.g_bytes_get_data(data, &size)
		bs = C.GoBytes(unsafe.Pointer(buf), C.int(size))
	}

	if callback != nil {
		callback(C.GoString(message), bs)
	}
}
