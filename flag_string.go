package cli

import "flag"

type StringFlag struct {
	Name  string
	Value string
	Usage string
}

func (f *StringFlag) GetName() string {
	return f.Name
}

func (f *StringFlag) GetValue() interface{} {
	return f.Value
}

func (f *StringFlag) Parse(flagSet *flag.FlagSet) {
	flagSet.StringVar(&f.Value, f.Name, f.Value, f.Usage)
}

type String struct {
	Context
}

func (c Context) String() String {
	return String{c}
}

func (s String) Get(name string) string {
	if val, ok := s.Flags[name]; ok {
		if v, ok := val.(string); ok {
			return v
		}
	}
	return ""
}

func (s String) Set(name, value string) string {
	s.Flags[name] = value
	return value
}
