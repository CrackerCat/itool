package cmd

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gofmt/itool/idevice"
	"github.com/gookit/gcli/v3"
)

var Scp2Cmd = &gcli.Command{
	Name: "scp2",
	Desc: "通过SSH传递文件, 兼容 dropbear",
	Config: func(c *gcli.Command) {
		c.IntOpt(&rport, "rport", "r", 22, "设备SSH端口")
		c.AddArg("src", "源路径", true)
		c.AddArg("dst", "目标路径", true)
	},
	Examples: "local file to device: {$binName} {$cmd} /path/example.js :/tmp/example." +
		"js\ndevice file to local: {$binName} {$cmd} :/tmp/example.js /path/example.js",
	Func: func(c *gcli.Command, args []string) error {
		var (
			remotePath   = ""
			localPath    = ""
			copyToRemote = false
		)

		if strings.HasPrefix(args[0], ":") {
			remotePath = strings.TrimLeft(args[0], ":")
			localPath = args[1]
			copyToRemote = false
		} else if strings.HasPrefix(args[1], ":") {
			remotePath = strings.TrimLeft(args[1], ":")
			localPath = args[0]
			copyToRemote = true
		}

		if len(remotePath) == 0 || len(localPath) == 0 {
			return errors.New("SCP参数错误")
		}

		device, err := idevice.GetDefaultDevice()
		if err != nil {
			return err
		}

		if copyToRemote {
			idx := strings.LastIndex(localPath, "/")
			if err := push(device, localPath, localPath[idx:]); err != nil {
				return err
			}
			rpath := fmt.Sprintf("/private/var/mobile/Media/%s", localPath[idx:])
			_, err := shellRun(rport, fmt.Sprintf("cp -r %s %s", rpath, remotePath))
			return err
		}

		idx := strings.LastIndex(remotePath, "/")
		rpath := fmt.Sprintf("/private/var/mobile/Media/%s", remotePath[idx:])
		result, err := shellRun(rport, fmt.Sprintf("cp -r %s %s", remotePath, rpath))
		if err != nil {
			return err
		}
		if strings.Contains(string(result), "No such file or directory") {
			return errors.New(string(result))
		}

		return pull(device, rpath, localPath)
	},
}
