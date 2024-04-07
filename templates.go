package cli

const defaultAppHelpTemplate = `{{.Name}}{{if .Version}}

Version:
   {{.Version}}{{end}}`
const defaultCmdHelpTemplate = `Usage: {{ with .App}}{{.Name}}{{ end }} {{.Cmd.Name}}{{ if .Cmd.Usage}} {{.Cmd.Usage}}{{end}}{{if .Cmd.Short}}

   {{.Cmd.Short}}{{end}}{{if .Cmd.Alias}}

Aliases: {{join .Cmd.Alias ", "}}{{end}}{{if .Cmd.Description}}

Description:
   {{.Cmd.Description}}{{end}}{{if .Cmd.Flags}}
{{ $flags := flags .Cmd }}
Flags:{{ range $name, $flag := $flags }}
   {{ printf "%-15s" $name }}    {{ $flag.Usage }}{{ end }}{{ end }}`
