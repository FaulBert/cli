package cli

import (
	"flag"
)

type Flag interface {
	Parse(*flag.FlagSet)
	GetName() string
	GetValue() interface{}
}

func parseFlags(flagSet *flag.FlagSet, flags []Flag) map[string]interface{} {
	flagValues := make(map[string]interface{})
	flagSet.VisitAll(func(f *flag.Flag) {
		for _, fl := range flags {
			if f.Name == fl.GetName() {
				flagValues[fl.GetName()] = fl.GetValue()
				break
			}
		}
	})
	return flagValues
}

// Assume cmd is an instance of Command
func getFlags(cmd *Command) map[string]interface{} {
	flagsMap := make(map[string]interface{})
	for _, flag := range cmd.Flags {
		switch f := flag.(type) {
		case *StringFlag:
			flagsMap[f.Name] = f
		case *IntFlag:
			flagsMap[f.Name] = f
		}
	}
	return flagsMap
}

// Assume cmd is an instance of Command
func getFlagValue(cmd *Command, flagName string) (interface{}, error) {
	for _, flag := range cmd.Flags {
		switch f := flag.(type) {
		case *StringFlag:
			if f.Name == flagName {
				return f.Value, nil
			}
		case *IntFlag:
			if f.Name == flagName {
				return f.Value, nil
			}
		}
	}
	return nil, ErrFlagNotFound(flagName)
}
