package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/gofmt/itool/idevice"
	"github.com/gofmt/itool/idevice/installation"
	"github.com/gookit/gcli/v3"
)

var AppCmd = &gcli.Command{
	Name: "apps",
	Desc: "显示设备应用列表",
	Config: func(c *gcli.Command) {
		c.AddArg("name", "应用名称，用于过滤")
	},
	Func: func(c *gcli.Command, args []string) error {
		device, err := idevice.GetDefaultDevice()
		if err != nil {
			return err
		}

		cli, err := installation.NewClient(device.UDID)
		if err != nil {
			return err
		}
		defer func(cli *installation.Client) {
			_ = cli.Close()
		}(cli)

		apps, err := cli.InstalledApps()
		if err != nil {
			return err
		}

		w := new(tabwriter.Writer)
		w.Init(os.Stdout, 0, 0, 1, ' ', 0)
		_, _ = fmt.Fprintln(w, "--------------------------------------------------------------")
		for i, app := range apps {
			if c.Arg("name").HasValue() && c.Arg("name").String() != app.CFBundleDisplayName {
				continue
			}

			_, _ = fmt.Fprintln(w, fmt.Sprintf("Number\t: %d", i))
			_, _ = fmt.Fprintln(w, fmt.Sprintf("Name\t: %s", app.CFBundleDisplayName))
			_, _ = fmt.Fprintln(w, fmt.Sprintf("BundleId\t: %s", app.CFBundleIdentifier))
			_, _ = fmt.Fprintln(w, fmt.Sprintf("Version\t: %s", app.CFBundleShortVersionString))
			_, _ = fmt.Fprintln(w, fmt.Sprintf("Executable\t: %s", app.CFBundleExecutable))
			_, _ = fmt.Fprintln(w, fmt.Sprintf("Container\t: %s", app.Container))
			_, _ = fmt.Fprintln(w, fmt.Sprintf("Path\t: %s", app.Path))
			_, _ = fmt.Fprintln(w, "--------------------------------------------------------------")
		}

		return w.Flush()
	},
}

var InstallCmd = &gcli.Command{
	Name:     "install",
	Desc:     "安装应用到设备",
	Examples: "{$binName} {$cmd} /path/example.ipa",
	Config: func(c *gcli.Command) {
		c.AddArg("ipa", "ipa 文件路径", true)
	},
	Func: func(c *gcli.Command, args []string) error {
		device, err := idevice.GetDefaultDevice()
		if err != nil {
			return err
		}

		cli, err := installation.NewClient(device.UDID)
		if err != nil {
			return err
		}
		defer func(cli *installation.Client) {
			_ = cli.Close()
		}(cli)

		ipa := c.Arg("ipa").String()
		f, err := os.Open(ipa)
		if err != nil {
			return err
		}
		defer func(f *os.File) {
			_ = f.Close()
		}(f)

		return cli.CopyAndInstall(ipa, nil)
	},
}

var UninstallCmd = &gcli.Command{
	Name: "uninstall",
	Desc: "卸载设备应用",
	Config: func(c *gcli.Command) {
		c.AddArg("bundleID", "应用 BundleID", true)
	},
	Func: func(c *gcli.Command, args []string) error {
		device, err := idevice.GetDefaultDevice()
		if err != nil {
			return err
		}

		cli, err := installation.NewClient(device.UDID)
		if err != nil {
			return err
		}
		defer func(cli *installation.Client) {
			_ = cli.Close()
		}(cli)

		return cli.Uninstall(c.Arg("bundleID").String(), nil)
	},
}
