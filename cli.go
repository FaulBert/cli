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
	Short        string
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
	Description  string
	Action       ActionFunc
	Flags        []Flag
	Commands     []*Command
	HelpTemplate string
}

func (app *App) AddCommand(cmd *Command) {
	app.Commands = append(app.Commands, cmd)
}

func (app *App) Run(args []string) (err error) {
	flagSet := flag.NewFlagSet(app.Name, flag.ExitOnError)
	flag.Usage = func() {
		printHelp(app, app)
	}

	if len(args) <= 1 || strings.HasPrefix(args[1], "-"){
		// args[1] = strings.ReplaceAll(args[1], "-", "")
		err = flagSet.Parse(args[0:])
		if err != nil {
			return err
		}
		runApp(args, app, flagSet)
		return
	}
	if len(args) >= 2 {
		err = flagSet.Parse(args[1:])
		if err != nil {
			return err
		}
		if strings.HasPrefix(args[1], "-") || app.Commands == nil {
			err := runApp(args, app, flagSet)
			if err != nil {
				return err
			}
		}
		if !strings.HasPrefix(args[1], "-") && app.Commands != nil {
			err := runCmd(args, app, flagSet)
			if err != nil {
				return err
			}
		}

		if args[1] == "-h" || args[1] == "--help" {
			printHelp(app, app)
			return
		}
	}

	return
}

func runApp(args []string, app *App, flagSet *flag.FlagSet) (err error) {
	if len(flagSet.Args()) == 0 {
		return ErrNoCommandProvided
	}

	// MAGIC!!!, must odd ["","-m","nazan"], ["-m","nazan"] wouldn't work
	// this pattern happen because, os.Args always return [<binary name>, "flag","flag value"],
	// the flag package follow those doctrine.
	args = append([]string{""}, args...)

	for _, flag := range app.Flags {
		flag.Parse(flagSet)
	}

	for i, arg := range flagSet.Args() {
		if strings.HasPrefix(arg, "-") {
			flagSet.Parse(args[1+i:])
		}
	}

	if app.Action == nil {
		return ErrAppActionNotProvided
	}

	app.Action(Context{
		Args:  flag.Args(),
		Flags: parseFlags(flagSet, app.Flags),
	})
	return
}

func runCmd(args []string, app *App, flagSet *flag.FlagSet) error {
	if len(flagSet.Args()) == 0 {
		return ErrNoCommandProvided
	}

	// Find the deepest subcommand
	cmd, cmdName := findCommand(&Command{Name: app.Name, Subcommands: app.Commands}, flagSet.Args())

	if cmdName != "" {
		return ErrCommandNotFound(cmdName)
	}

	// Check if --help flag is provided
	for _, arg := range flagSet.Args() {
		if arg == "--help" || arg == "-h" {
			printHelp(app, cmd)
			return nil
		}
	}

	// if subcommand not found, action will be empty.
	if cmd.Action == nil {
		return ErrCommandNotRegistered(cmdName)
	}

	for _, flag := range cmd.Flags {
		flag.Parse(flagSet)
	}

	for i, arg := range flagSet.Args() {
		if strings.HasPrefix(arg, "-") {
			flagSet.Parse(args[1+i:])
		}
	}

	cmd.Action(Context{
		Args:  flagSet.Args(),
		Flags: parseFlags(flagSet, cmd.Flags),
	})
	return nil
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
