package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gofmt/itool/idevice"
	"github.com/gofmt/itool/idevice/installation"
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

		appCli, err := installation.NewClient(device.UDID)
		if err != nil {
			return err
		}
		defer func(appCli *installation.Client) {
			_ = appCli.Close()
		}(appCli)

		apps, err := appCli.InstalledApps()
		if err != nil {
			return err
		}

		fmt.Println("应用列表：")
		fmt.Println("--------------------------------------------------------------")
		for i, app := range apps {
			if app.CFBundleDisplayName != "" {
				fmt.Println(i, "\t|", app.CFBundleDisplayName, "["+app.CFBundleIdentifier+"]["+app.CFBundleExecutable+"]")
			}
		}

		fmt.Println("--------------------------------------------------------------")
		fmt.Println("输入应用编号开始调试：")
		var input string
		_, err = fmt.Scan(&input)
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}

		idx, err := strconv.Atoi(input)
		if err != nil {
			fmt.Printf("'%s' 不是正确的应用ID\n", input)
			os.Exit(-1)
		}

		if idx > len(apps)-1 {
			fmt.Printf("'%d' 应用ID不存在\n", idx)
			os.Exit(-1)
		}

		app := apps[idx]
		execPath := filepath.Join(app.Path, app.CFBundleExecutable)
		fmt.Println(execPath)

		return nil
	},
}
