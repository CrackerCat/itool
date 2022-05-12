package cmd

import (
	"context"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/gofmt/itool/idevice"
	"github.com/gofmt/itool/idevice/installation"
	"github.com/gofmt/itool/idevice/pcap"
	"github.com/gookit/gcli/v3"
)

var PcapCmd = &gcli.Command{
	Name: "pcap",
	Desc: "设备网络抓包，可过滤进程",
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
		fmt.Println("输入应用编号开始抓包：")
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

		name := apps[idx].CFBundleDisplayName
		fmt.Println("["+name+"]", "正在抓包,[CTRL+C]停止抓包...")

		execName := apps[idx].CFBundleExecutable
		ctx, cancel := signal.NotifyContext(context.Background(), os.Kill, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGABRT)

		pcapClient, err := pcap.NewClient(device.UDID)
		if err != nil {
			fmt.Println("创建PCAP客户端错误:", err)
			os.Exit(-1)
		}
		defer func(pcapClient *pcap.Client) {
			_ = pcapClient.Close()
		}(pcapClient)

		pcapFile, err := os.Create(name + ".pcap")
		if err != nil {
			fmt.Println("创建PCAP文件错误:", err)
			os.Exit(-1)
		}
		defer func(pcapFile *os.File) {
			_ = pcapFile.Close()
		}(pcapFile)

		go func() {
			<-ctx.Done()
			fmt.Println("正在停止抓包，封包数据回写有点慢，请等待几秒出现抓包结束提示...")
		}()

		err = pcapClient.ReadPacket(ctx, execName, pcapFile, func(data []byte) {
			fmt.Println(hex.Dump(data))
		})
		if err != nil {
			fmt.Println("读取网络封包错误:", err)
			os.Exit(-1)
		}

		cancel()

		// wireshark -r xxx.pcap -o "tls.keylog_file:./xxx.keylog"
		wiresharkParam := fmt.Sprintf(`wireshark -r %s.pcap`, name)
		_ = ioutil.WriteFile("wireshark.sh", []byte(wiresharkParam), os.ModePerm)

		fmt.Println("["+name+"]", "抓包结束")

		return nil
	},
}
