package cli

import "fmt"

// make error centralize, for convinient in error uniformity.
var (
	ErrNoCommandProvided = fmt.Errorf("no command provided")
)

func ErrCommandNotRegistered(command string) error {
	return fmt.Errorf("Command %s not registered", command)
}

func ErrCommandNotFound(command string) error {
	return fmt.Errorf("Command %s not found", command)
}
