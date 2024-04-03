package cli

import (
	"flag"
	"strings"
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
	Alias        []string
	Action       ActionFunc
	Flags        []Flag
	Subcommands  []*Command
	HelpTemplate string
}

type App struct {
	Name         string
	Version      string
	Commands     []*Command
	HelpTemplate string
}

func (app *App) AddCommand(cmd *Command) {
	app.Commands = append(app.Commands, cmd)
}

func findCommand(cmd *Command, args []string) (*Command, string) {
	if len(args) == 0 || len(cmd.Subcommands) == 0 {
		return cmd, ""
	}

	nextCmdName := args[0]
	for _, subCmd := range cmd.Subcommands {
		if subCmd.Name == nextCmdName || contains(subCmd.Alias, nextCmdName) {
			return findCommand(subCmd, args[1:])
		}
	}

	return cmd, nextCmdName
}

func contains(slice []string, name string) bool {
	for _, s := range slice {
		if s == name {
			return true
		}
	}
	return false
}

func (app *App) Run(args []string) (err error) {
	if len(args) <= 1 || args[1] == "-h" || args[1] == "--help" {
		printHelp(app, app)
		return
	}

	flagSet := flag.NewFlagSet(app.Name, flag.ExitOnError)
	flag.Usage = func() {
		printHelp(app, app)
	}

	err = flagSet.Parse(args[1:])
	if err != nil {
		return err
	}

	if len(flagSet.Args()) == 0 {
		return ErrNoCommandProvided
	}

	// Find the deepest subcommand
	cmd, cmdName := findCommand(&Command{Name: app.Name, Subcommands: app.Commands}, flagSet.Args())

	// Check if --help flag is provided
	for _, arg := range flagSet.Args() {
		if arg == "--help" || arg == "-h" {
			printHelp(app, cmd)
			return
		}
	}

	for _, flag := range cmd.Flags {
		flag.Parse(flagSet)
	}

	for i, arg := range flagSet.Args() {
		if strings.HasPrefix(arg, "-") {
			flagSet.Parse(args[1+i:])
		}
	}

	// if subcommand not found, action will be empty.
	if cmd.Action == nil {
		return ErrCommandNotRegistered(cmdName)
	}

	cmd.Action(Context{
		Args:  flagSet.Args(),
		Flags: parseFlags(flagSet, cmd.Flags),
	})
	return
}
