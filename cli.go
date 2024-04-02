package cli

import (
	"flag"
	"fmt"
	"os"
	"reflect"
	"text/template"
)

// Context represents the context of a command execution.
type Context struct {
	Args  Args
	Flags map[string]interface{}
}

type ActionFunc func(Context)

type Command struct {
	Name        string
	Usage       string
	Description string
	Action      ActionFunc
	Flags       []Flag
	Subcommands []*Command
	Help        string
}

type App struct {
	Name     string
	Version  string
	Commands []*Command
	Help     string
}

func (app *App) AddCommand(cmd *Command) {
	app.Commands = append(app.Commands, cmd)
}

func (app *App) Run(args []string) {
	if len(args) < 2 || args[1] == "help" {
		printHelp(app)
		return
	}

	commandName := args[1]
	for _, cmd := range app.Commands {
		if cmd.Name == commandName {
			if len(args) == 3 && args[2] == "help" {
				printHelp(cmd)
				return
			}

			flagSet := flag.NewFlagSet(commandName, flag.ExitOnError)
			for _, flag := range cmd.Flags {
				flag.Parse(flagSet)
			}
			flagSet.Parse(args[2:])
			cmd.Action(Context{
				Args:  flagSet.Args(),
				Flags: parseFlags(flagSet, cmd.Flags),
			})
			return
		}
	}

	fmt.Println("Command not found:", commandName)
}

func printHelp(data interface{}) {
	switch d := data.(type) {
	case *Command:
		tmpl, err := template.New("help").Parse(d.Help)
		if err != nil {
			fmt.Println("Error parsing help template:", err)
			return
		}
		err = tmpl.Execute(os.Stdout, d)
		if err != nil {
			fmt.Println("Error executing help template:", err)
		}
	case *App:
		tmpl, err := template.New("help").Parse(d.Help)
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
