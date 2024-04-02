package cli

import "flag"

type BoolFlag struct {
	Name  string
	Value bool
	Usage string
}

func (f *BoolFlag) GetName() string {
	return f.Name
}

func (f *BoolFlag) GetValue() interface{} {
	return f.Value
}

func (f *BoolFlag) Parse(flagSet *flag.FlagSet) {
	flagSet.BoolVar(&f.Value, f.Name, f.Value, f.Usage)
}

type Bool struct {
	Context
}

func (c Context) Bool() Bool {
	return Bool{c}
}

func (b Bool) Get(name string) bool {
	if val, ok := b.Flags[name]; ok {
		if v, ok := val.(bool); ok {
			return v
		}
	}
	return false
}

func (b Bool) Set(name string, value bool) bool {
	b.Flags[name] = value
	return value
}
