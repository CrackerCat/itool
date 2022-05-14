package cmd

import (
	"context"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofmt/itool/idevice"
	"github.com/gofmt/itool/idevice/pcap"
	"github.com/gookit/gcli/v3"
)

var PcapCmd = &gcli.Command{
	Name: "pcap",
	Desc: "对设备应用进行网络抓包, 抓包结束后执行 wireshark.sh",
	Func: func(c *gcli.Command, args []string) error {
		device, err := idevice.GetDefaultDevice()
		if err != nil {
			return err
		}

		app, err := GetSelectApp(device)
		if err != nil {
			return err
		}

		name := app.CFBundleDisplayName
		fmt.Println("["+name+"]", "正在抓包,[CTRL+C]停止抓包...")

		execName := app.CFBundleExecutable
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
