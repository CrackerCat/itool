package debugserver

import (
	"io"
	"os"
	"testing"

	"github.com/gofmt/itool/idevice"
	"github.com/gofmt/itool/idevice/installation"
)

func TestProcess_Start(t *testing.T) {
	conn, err := idevice.NewConn()
	if err != nil {
		t.Fatal(err)
	}
	defer func(conn *idevice.Conn) {
		_ = conn.Close()
	}(conn)

	devices, err := conn.ListDevices()
	if err != nil {
		t.Fatal(err)
	}

	if len(devices) < 1 {
		t.Fatal("not device")
	}

	device := devices[0]
	cli, err := installation.NewClient(device.UDID)
	if err != nil {
		t.Fatal(err)
	}
	path, err := cli.LookupExePath("com.meituan.itakeaway")
	if err != nil {
		t.Fatal(err)
	}
	var appEnv []string
	if os.Getenv("IDE_DISABLED_OS_ACTIVITY_DT_MODE") == "" {
		appEnv = append(appEnv, "OS_ACTIVITY_DT_MODE=enable")
	}
	proc, err := NewProcess(device.UDID, []string{path}, appEnv)
	if err != nil {
		t.Fatal(err)
	}
	defer func(proc *Process) {
		_ = proc.Kill()
	}(proc)

	go func() {
		_, _ = io.Copy(os.Stdout, proc.Stdout())
	}()
	if err := proc.Start(); err != nil {
		t.Fatal(err)
	}

	select {}
}
