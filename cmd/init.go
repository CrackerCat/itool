package cmd

import (
	"github.com/gookit/gcli/v3"
)

func init() {
	gcli.AppHelpTemplate = `{{.Desc}} (版本: <info>{{.Version}}</>)
-----------------------------------------------------{{range $cmdName, $c := .Cs}}
  <info>{{$c.Name | paddingName }}</> {{$c.HelpDesc}}{{if $c.Aliases}} (别名: <green>{{ join $c.Aliases ","}}</>){{end}}{{end}}
  <info>{{ paddingName "help" }}</> 显示帮助信息

使用 "<cyan>{$binName} COMMAND -h</>" 查看命令的其他帮助信息
`

	gcli.CmdHelpTemplate = `{{.Desc}}
	
<comment>用法:</>
  {$binName} {{if .Cmd.NotStandalone}}<cyan>{{.Cmd.Path}}</> {{end}}[--选项 ...] [参数 ...]{{ if .Subs }}
  {$binName} {{if .Cmd.NotStandalone}}<cyan>{{.Cmd.Path}}</> {{end}}<cyan>SUBCOMMAND</> [--选项 ...] [参数 ...]{{end}}
{{if .Options}}
<comment>选项:</>
{{.Options}}{{end}}{{if .Cmd.Args}}
<comment>参数:</>{{range $a := .Cmd.Args}}
  <info>{{$a.HelpName | printf "%-12s"}}</>{{$a.Desc | ucFirst}}{{if $a.Required}}<red>*</>{{end}}{{end}}
{{end}}{{ if .Subs }}
<comment>子命令列表:</>{{range $n,$c := .Subs}}
  <info>{{$c.Name | paddingName }}</> {{$c.HelpDesc}}{{if $c.Aliases}} (alias: <green>{{ join $c.Aliases ","}}</>){{end}}{{end}}
{{end}}{{if .Cmd.Examples}}
<comment>用例:</>
{{.Cmd.Examples}}{{end}}{{if .Cmd.Help}}
<comment>帮助:</>
{{.Cmd.Help}}{{end}}`
}
