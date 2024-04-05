package cli

import (
	"flag"
)

type Flag interface {
	Parse(*flag.FlagSet)
	GetName() string
	GetValue() interface{}
	GetAction() *ActionFunc
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
