package cli

import "fmt"

// make error centralize, for convinient in error uniformity.
var (
	ErrNoCommandProvided = fmt.Errorf("no command provided")
)

func ErrCommandNotRegistered(command string) error {
	return fmt.Errorf("command %s not registered", command)
}
