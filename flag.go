package cli

import (
	"flag"
	"fmt"
	"reflect"
)

type Flag interface {
	Parse(*flag.FlagSet)
	GetName() string
}

// String flag
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

// Bool flag
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
			switch fl := flag.(type) {
			case *StringFlag:
				if f.Name == fl.Name {
					flagValues[f.Name] = f.Value.String()
					break
				}
			case *BoolFlag:
				if f.Name == fl.Name {
					flagValues[f.Name] = f.Value.String() == "true"
					break
				}
			default:
				fmt.Println("Unknown flag type:", reflect.TypeOf(flag))
			}
		}
	})
	return flagValues
}
