package main

import (
	"fmt"
	"os"
	"strings"
	"text/template"
)

// CLIFlag represents a CLI flag
type Flags interface {
	TypeFlag() string // Type returns the type of the flag
	NameFlag() string // Name returns the name of the flag
	ValueFlag() interface{}
}

// StringFlag represents a string flag
type String struct {
	Name  string
	Value string
}

// Type returns the type of the flag
func (f String) TypeFlag() string {
	return "string"
}

// Name returns the name of the flag
func (f String) NameFlag() string {
	return f.Name
}

func (f String) ValueFlag() interface{} {
	return f.Value
}

// BoolFlag represents a boolean flag
type Bool struct {
	Name  string
	Value bool
}

// Type returns the type of the flag
func (f Bool) TypeFlag() string {
	return "bool"
}

// Name returns the name of the flag
func (f Bool) NameFlag() string {
	return f.Name
}

func (f Bool) ValueFlag() interface{} {
	return f.Value
}

// Context represents the context passed to the Action function of a command
type Context struct {
	Flags map[string]interface{} // Flags map to store flag values
}

// Command represents a CLI command
type Command struct {
	Name        string
	Usage       string
	Help        string
	Description string
	Aliases     []string
	Flags       []Flags // Slice of CLIFlag for storing flags
	Action      func(Context)
}

// New creates a new CLI instance
type New struct {
	Name     string
	Usage    string
	Version  string
	flags    map[any]Flags
	commands map[string]Command
}

// AddCommand registers a new command
func (c *New) AddCommand(cmd Command) {
	c.commands[cmd.Name] = cmd
	for _, alias := range cmd.Aliases {
		c.commands[alias] = cmd
	}
}

// Exec executes a command by name
func (c *New) Exec(name string, args []string) {
	cmd, ok := c.commands[name]
	if !ok {
		fmt.Println("Unknown command:", name)
		c.PrintUsage()
		return
	}

	// Check if help flag is provided
	if len(args) > 0 && (args[0] == "-h" || args[0] == "--help") {
		c.PrintCommandHelp(name)
		return
	}

	// Parse command flags
	flagValues := make(map[string]interface{})
	i := 0
	for i < len(args) {
		arg := args[i]

		if strings.HasPrefix(arg, "-") && i+1 < len(args) {
			flagName := arg[1:]
			flagValue := args[i+1]
			i += 2

			// Check if the flag is valid for the command
			validFlag := false
			for _, flag := range cmd.Flags {
				if flagName == flag.NameFlag() {
					validFlag = true
					// If the flag requires a value, ensure it's provided
					_, ok := flag.(String)
					if ok && flagValue == "" {
						fmt.Printf("Flag -%s requires a value\n", flagName)
						return
					}
					break
				}
			}

			if validFlag {
				flagValues[flagName] = flagValue
			} else {
				// Flag is not valid, suggest correct flags
				fmt.Printf("Unknown flag: -%s\n", flagName)
				fmt.Println("Valid flags:")
				for _, flag := range cmd.Flags {
					fmt.Printf("-%s\n", flag.NameFlag())
				}
				return
			}
		} else {
			i++ // Skip the next argument (flag value)
		}
	}

	// Set default values for flags
	for _, flag := range cmd.Flags {
		defaultValue := flag.ValueFlag()
		if defaultValue != "" {
			if _, exists := flagValues[flag.NameFlag()]; !exists {
				flagValues[flag.NameFlag()] = defaultValue
			}
		}
	}

	// Set default values for flags
	for _, flag := range cmd.Flags {
		defaultValue := flag.ValueFlag()
		if defaultValue != "" {
			_, exists := flagValues[flag.NameFlag()]
			if !exists {
				flagValues[flag.NameFlag()] = defaultValue
			}
		}
	}

	context := Context{Flags: flagValues}

	// Call the Action function with the context
	cmd.Action(context)
}

// PrintUsage prints usage information for the CLI program
func (c *New) PrintUsage() {
	fmt.Printf("Usage: %s %s \n\n", c.Name, c.Usage)
	fmt.Println("Commands:")
	printedCmds := make(map[string]bool)
	for _, cmd := range c.commands {
		if !printedCmds[cmd.Name] {
			fmt.Printf("  %s\t%s\n", cmd.Name, cmd.Usage)
			printedCmds[cmd.Name] = true
		}

		// I don't know but this is to circumvent commands that are printed twice
		for _, alias := range cmd.Aliases {
			if !printedCmds[alias] {
				fmt.Print()
				printedCmds[alias] = true
			}
		}
	}
}

// PrintCommandHelp prints help information for a specific command
func (c *New) PrintCommandHelp(name string) {
	cmd, ok := c.commands[name]
	if !ok {
		fmt.Println("Unknown command:", name)
		c.PrintUsage()
		return
	}

	// Define the help template
	tmpl := `Usage: {{.CLIName}} {{.CmdName}} {{.Usage}}
{{if .Aliases}}Aliases: {{join .Aliases ", "}}{{end}}
{{if .Description}}{{.Description}}{{end}}
{{.Help}}`

	// Define the template functions
	funcMap := template.FuncMap{
		"join": strings.Join,
	}

	// Parse the template
	helpTemplate, err := template.New("help").Funcs(funcMap).Parse(tmpl)
	if err != nil {
		fmt.Println("Error parsing help template:", err)
		return
	}

	// Execute the template with command data
	err = helpTemplate.Execute(os.Stdout, map[string]interface{}{
		"CLIName":     c.Name,
		"CmdName":     cmd.Name,
		"Usage":       cmd.Usage,
		"Aliases":     cmd.Aliases,
		"Description": cmd.Description,
		"Help":        cmd.Help,
	})
	if err != nil {
		fmt.Println("Error executing help template:", err)
		return
	}
}

func main() {
	cli := New{
		Name:     "test",
		Usage:    "[commands] [flag]",
		Version:  "experimental",
		commands: make(map[string]Command),
	}

	echoCmd := Command{
		Name:    "echo",
		Usage:   "[-m <message>]",
		Aliases: []string{"say"},
		Help:    "Echoes the provided message",
		Flags: []Flags{
			String{
				Name:  "m",
				Value: "cute",
			},
			String{
				Name:  "cute",
				Value: "mumumu",
			},
		},
		Action: func(c Context) {
			message, _ := c.Flags["m"].(string)
			msg, _ := c.Flags["cute"].(string)
			fmt.Println("Echo:", message)
			fmt.Println("Echo:", msg)
		},
	}

	cli.AddCommand(echoCmd)

	// Parse command line arguments
	if len(os.Args) < 2 {
		cli.PrintUsage()
		return
	}

	// Execute the specified command
	cmdName := os.Args[1]
	cli.Exec(cmdName, os.Args[2:])
}
