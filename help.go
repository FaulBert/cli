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
const defaultCmdHelpTemplate = `Usage: {{.Name}} {{.CmdName}} {{ if .Usage}}{{.Usage}}{{end}}{{if .Short}}

   {{.Short}}{{end}}{{if .Alias}}

Aliases: {{join .Alias ", "}}{{end}}{{if .Description}}

Description:
   {{.Description}}{{end}}`

func printHelp(app *App, data interface{}) {
	switch d := data.(type) {
	case *Command:
		if d.HelpTemplate == "" {
			cmdHelpParser(app, d)
		} else {
			tmpl, err := template.New("help").Parse(d.HelpTemplate)
			if err != nil {
				fmt.Println("Error parsing help template:", err)
				return
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
				fmt.Println("Error parsing help template:", err)
				return
			}
			err = tmpl.Execute(os.Stdout, d)
			if err != nil {
				fmt.Println("Error executing help template:", err)
			}
		}
	default:
		fmt.Println("Unknown type:", reflect.TypeOf(data))
	}
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

func cmdHelpParser(app *App, cmd *Command) {
	funcMap := template.FuncMap{
		"join": strings.Join,
	}

	tmpl, err := template.New("help").Funcs(funcMap).Parse(defaultCmdHelpTemplate)
	if err != nil {
		fmt.Println("Error parsing help template:", err)
		return
	}
	err = tmpl.Execute(os.Stdout, map[string]interface{}{
		"Name":        app.Name,
		"Alias":       cmd.Alias,
		"CmdName":     cmd.Name,
		"Usage":       cmd.Usage,
		"Description": cmd.Description,
	})
	if err != nil {
		fmt.Println("Error executing help template:", err)
		return
	}
}
