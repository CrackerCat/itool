package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/gofmt/itool/idevice"
	"github.com/gookit/gcli/v3"
)

var DebugCmd = &gcli.Command{
	Name: "debug",
	Desc: "启动 debugserver 调试目标应用",
	Func: func(c *gcli.Command, args []string) error {
		device, err := idevice.GetDefaultDevice()
		if err != nil {
			return err
		}

		app, err := GetSelectApp(device)
		if err != nil {
			return err
		}

		execPath := filepath.Join(app.Path, app.CFBundleExecutable)
		fmt.Println(execPath)

		return nil
	},
}
