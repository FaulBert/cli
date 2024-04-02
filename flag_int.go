package cli

import "flag"

type IntFlag struct {
	Name  string
	Value int
	Usage string
}

type Int struct {
	Context
}

func (c Context) Int() Int {
	return Int{c}
}

func (f *IntFlag) GetName() string {
	return f.Name
}

func (f *IntFlag) GetValue() interface{} {
	return f.Value
}

func (f *IntFlag) Parse(flagSet *flag.FlagSet) {
	flagSet.IntVar(&f.Value, f.Name, f.Value, f.Usage)
}

func (ic Int) Get(name string) int {
	if val, ok := ic.Flags[name]; ok {
		if v, ok := val.(int); ok {
			return v
		}
	}
	return 0
}

func (ic Int) Set(name string, value int) int {
	ic.Flags[name] = value
	return value
}
