package idevice

import (
	"errors"
)

func GetDefaultDevice() (*DeviceAttachment, error) {
	conn, err := NewConn()
	if err != nil {
		return nil, err
	}
	defer func(conn *Conn) {
		_ = conn.Close()
	}(conn)

	devices, err := conn.ListDevices()
	if err != nil {
		return nil, err
	}

	if len(devices) < 1 {
		return nil, errors.New("not device")
	}

	return devices[0], nil
}
