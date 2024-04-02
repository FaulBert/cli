package cli

import "flag"

type UintFlag struct {
	Name    string
	Value   uint
	Default uint
	Usage   string
}

type Uint struct {
	Context
}

func (c Context) Uint() Uint {
	return Uint{c}
}

func (f *UintFlag) Parse(flagSet *flag.FlagSet) {
	flagSet.UintVar(&f.Value, f.Name, f.Default, f.Usage)
}

func (f *UintFlag) GetName() string {
	return f.Name
}

func (f *UintFlag) GetValue() interface{} {
	return f.Value
}

func (u Uint) Get(name string) uint {
	if val, ok := u.Flags[name]; ok {
		if v, ok := val.(uint); ok {
			return v
		}
	}
	return 0
}

func (u Uint) Set(name string, value uint) uint {
	u.Flags[name] = value
	return value
}
