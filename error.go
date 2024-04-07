package cli

import "fmt"

// make error centralize, for convinient in error uniformity.
var (
	ErrNoCommandProvided   = fmt.Errorf("no command provided")
	ErrParsingHelpTemplate = fmt.Errorf("Error parsing help template")
)

func ErrFlagNotFound(flag string) error {
	return fmt.Errorf("flag '%s' not found", flag)
}

func ErrCommandNotRegistered(command string) error {
	return fmt.Errorf("Command %s not registered", command)
}

func ErrCommandNotFound(command string) error {
	return fmt.Errorf("Command %s not found", command)
}
