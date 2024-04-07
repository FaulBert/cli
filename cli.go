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

	if app.Action == nil && len(args) <= 1 {
		printHelp(app, app)
		return
	}

	if len(args) <= 1 && strings.HasPrefix(args[0], "-") {
		err = flagSet.Parse(args[0:])
		if err != nil {
			return err
		}
		err := runApp(args, app, flagSet)
		if err != nil {
			return err
		}
	}

	if len(args) >= 2 {
		err = flagSet.Parse(args[1:])
		if err != nil {
			return err
		}

		if args[1] == "-h" || args[1] == "--help" {
			printHelp(app, app)
			return
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
	}

	return
}

func runApp(args []string, app *App, flagSet *flag.FlagSet) (err error) {
	if len(flagSet.Args()) == 0 {
		return ErrNoCommandProvided
	}

	for _, flag := range app.Flags {
		flag.Parse(flagSet)
	}

	for i, arg := range flagSet.Args() {
		if strings.HasPrefix(arg, "-") {
			flagSet.Parse(args[1+i:])
		}
	}

	app.Action(Context{
		Args:  flag.Args(),
		Flags: parseFlags(flagSet, app.Flags),
	})
	return
}
