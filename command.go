package cli

import (
	"flag"
	"strings"
)

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
			err := printHelp(app, cmd)
			return err
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
