package main

import (
	"flag"
	"fmt"
	"os"
)

type Command struct {
	Name        string
	Usage       string
	Description string
	Action      func(Context)
	Flags       []Flag
}

type Flag interface {
	Parse(*flag.FlagSet)
	GetName() string
}

type Context struct {
	Args  []string
	Flags map[string]interface{}
}

func main() {
	app := App{
		Name:    "cli",
		Version: "1.0.0",
	}

	echoCmd := Command{
		Name:        "echo",
		Usage:       "[option] [arg]",
		Description: "Print text",
		Action: func(c Context) {
			text := c.Flags["m"].(string)
			fmt.Println(text)
		},
		Flags: []Flag{
			&StringFlag{
				Name:  "m",
				Usage: "print special message",
			},
			&BoolFlag{
				Name:  "dl",
				Usage: "download",
			},
		},
	}

	app.AddCommand(echoCmd)
	app.Run(os.Args)
}

type App struct {
	Name     string
	Version  string
	Commands []Command
}

func (app *App) AddCommand(cmd Command) {
	app.Commands = append(app.Commands, cmd)
}

func (app *App) Run(args []string) {
	if len(args) < 2 {
		fmt.Println("Usage: " + app.Name + " [command]")
		return
	}

	commandName := args[1]
	for _, cmd := range app.Commands {
		if cmd.Name == commandName {
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

type StringFlag struct {
	Name  string
	Usage string
}

func (f *StringFlag) GetName() string {
	return f.Name
}

func (f *StringFlag) Parse(flagSet *flag.FlagSet) {
	flagSet.String(f.Name, "", f.Usage)
}

type BoolFlag struct {
	Name  string
	Usage string
}

func (f *BoolFlag) GetName() string {
	return f.Name
}

func (f *BoolFlag) Parse(flagSet *flag.FlagSet) {
	flagSet.Bool(f.Name, false, f.Usage)
}

func parseFlags(flagSet *flag.FlagSet, flags []Flag) map[string]interface{} {
	flagValues := make(map[string]interface{})
	flagSet.VisitAll(func(f *flag.Flag) {
		for _, flag := range flags {
			if f.Name == flag.GetName() {
				flagValues[f.Name] = f.Value.String()
				break
			}
		}
	})
	return flagValues
}
