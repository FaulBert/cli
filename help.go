package cli

import (
	"fmt"
	"os"
	"reflect"
	"strings"
	"text/template"
)

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

func flagsFunc(cmd *Command) map[string]interface{} {
	return getFlags(cmd)
}

func printHelp(app *App, data interface{}) (err error) {
	switch d := data.(type) {
	case *Command:
		if d.HelpTemplate == "" {
			err := cmdHelpParser(app, d)
			if err != nil {
				return err
			}
		} else {
			tmpl, err := template.New("help").Parse(d.HelpTemplate)
			if err != nil {
				return ErrParsingHelpTemplate
			}
			err = tmpl.Execute(os.Stdout, d)
			if err != nil {
				fmt.Println("Error executing help template:", err)
			}
		}
	case *App:
		if d.HelpTemplate == "" {
			appHelpParser(d)
		} else {
			tmpl, err := template.New("help").Parse(d.HelpTemplate)
			if err != nil {
				return ErrParsingHelpTemplate
			}
			err = tmpl.Execute(os.Stdout, d)
			if err != nil {
				fmt.Println("Error executing help template:", err)
			}
		}
	default:
		fmt.Println("Unknown type:", reflect.TypeOf(data))
	}

	return
}

func appHelpParser(app *App) {
	tmpl, err := template.New("help").Parse(defaultAppHelpTemplate)
	if err != nil {
		fmt.Println("Error parsing app's help template:", err)
		return
	}
	err = tmpl.Execute(os.Stdout, map[string]interface{}{
		"Name":    app.Name,
		"Version": app.Version,
	})
	if err != nil {
		fmt.Println("Error executing app's help template:", err)
		return
	}
}

func cmdHelpParser(app *App, cmd *Command) error {
	funcMap := template.FuncMap{
		"join":  strings.Join,
		"flags": flagsFunc,
	}

	data := struct {
		App *App
		Cmd *Command
	}{
		App: app,
		Cmd: cmd,
	}

	tmpl, err := template.New("help").Funcs(funcMap).Parse(defaultCmdHelpTemplate)
	if err != nil {
		return ErrParsingHelpTemplate
	}
	err = tmpl.Execute(os.Stdout, data)
	if err != nil {
		return ErrParsingHelpTemplate
	}

	return nil
}
