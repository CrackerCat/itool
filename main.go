package main

import (
	"fmt"
	"os"

	"github.com/gofmt/itool/cmd"
	"github.com/gookit/gcli/v3"
)

func main() {
	gcli.DefaultVerb = gcli.VerbQuiet

	app := gcli.NewApp(func(app *gcli.App) {
		app.Version = "v1.0.1-beta"
		app.Desc = "itool"
		app.ExitOnEnd = false
		app.On(gcli.EvtAppInit, func(data ...interface{}) (stop bool) {
			return false
		})
	})

	app.On(gcli.EvtAppRunError, func(data ...interface{}) (stop bool) {
		fmt.Println(data[1])
		return false
	})

	app.On(gcli.EvtCmdNotFound, func(data ...interface{}) (stop bool) {
		return false
	})

	app.On(gcli.EvtAppCmdNotFound, func(data ...interface{}) (stop bool) {
		return false
	})

	app.On(gcli.EvtCmdRunError, func(data ...interface{}) (stop bool) {
		return false
	})

	app.Add(
		cmd.InfoCmd,
		cmd.RestartCmd,
		cmd.ShutdownCmd,
		cmd.SyslogCmd,
		cmd.AppCmd,
		cmd.InstallCmd,
		cmd.UninstallCmd,
		cmd.ShellCmd,
		cmd.ForwardCmd,
		cmd.PcapCmd,
		cmd.DebugCmd,
		cmd.FileCmd,
		cmd.ScreenShotCmd,
	)

	os.Exit(app.Run(nil))
}
