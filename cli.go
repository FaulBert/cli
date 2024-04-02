package cli

import (
	"flag"
	"fmt"
	"os"
	"reflect"
	"strings"
	"text/template"
)

// Context represents the context of a command execution.
type Context struct {
	Args  Args
	Flags map[string]interface{}
}

type ActionFunc func(Context)

type Command struct {
	Name         string
	Usage        string
	Description  string
	Action       ActionFunc
	Flags        []Flag
	Subcommands  []*Command
	Help         string
	HelpTemplate string
}

type App struct {
	Name         string
	Version      string
	Commands     []*Command
	Help         string
	HelpTemplate string
}

func (app *App) AddCommand(cmd *Command) {
	app.Commands = append(app.Commands, cmd)
}

func findCommand(cmd *Command, args []string) *Command {
	if len(args) == 0 || len(cmd.Subcommands) == 0 {
		return cmd
	}

	nextCmdName := args[0]
	for _, subCmd := range cmd.Subcommands {
		if subCmd.Name == nextCmdName {
			return findCommand(subCmd, args[1:])
		}
	}

	return cmd
}

func (app *App) Run(args []string) {
	help := len(args) < 2 || args[1] == "help" || args[1] == "-h" || args[1] == "--help"
	if help {
		if app.Help != "" {
			fmt.Print(app.Help)
			return
		}
		renderText(app)
		return
	}

	// Find the deepest subcommand
	var cmd *Command
	if len(args) >= 2 {
		cmd = findCommand(&Command{Name: app.Name, Subcommands: app.Commands}, args[1:])
	}

	// Execute the action of the deepest subcommand
	flagSet := flag.NewFlagSet(cmd.Name, flag.ExitOnError)
	for _, flag := range cmd.Flags {
		flag.Parse(flagSet)
	}

	if len(args) >= 3 && !strings.HasPrefix(args[2], "-") {
		flagSet.Parse(args[len(args)-len(flagSet.Args()):])
		cmd.Action(Context{
			Args:  flagSet.Args(),
			Flags: parseFlags(flagSet, cmd.Flags),
		})
	} else {
		flagSet.Parse(args[2:])
		cmd.Action(Context{
			Args:  flagSet.Args(),
			Flags: parseFlags(flagSet, cmd.Flags),
		})
	}
}

func renderText(data interface{}) {
	switch d := data.(type) {
	case *Command:
		tmpl, err := template.New("help").Parse(d.HelpTemplate)
		if err != nil {
			fmt.Println("Error parsing help template:", err)
			return
		}
		err = tmpl.Execute(os.Stdout, d)
		if err != nil {
			fmt.Println("Error executing help template:", err)
		}
	case *App:
		tmpl, err := template.New("help").Parse(d.HelpTemplate)
		if err != nil {
			fmt.Println("Error parsing help template:", err)
			return
		}
		err = tmpl.Execute(os.Stdout, d)
		if err != nil {
			fmt.Println("Error executing help template:", err)
		}
	default:
		fmt.Println("Unknown type:", reflect.TypeOf(data))
	}
}
