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

// CLI represents the command line interface
type CLI struct {
	Name     string
	Usage    string
	Version  string
	flags    map[any]Flags
	commands map[string]Command
}

// NewCLI creates a new CLI instance
func NewCLI(name, usage, version string) *CLI {
	return &CLI{
		Name:     name,
		Usage:    usage,
		Version:  version,
		commands: make(map[string]Command),
	}
}

// AddCommand registers a new command
func (c *CLI) AddCommand(cmd Command) {
	c.commands[cmd.Name] = cmd
	for _, alias := range cmd.Aliases {
		c.commands[alias] = cmd
	}
}

// Exec executes a command by name
func (c *CLI) Exec(name string, args []string) {
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
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if strings.HasPrefix(arg, "-") && i+1 < len(args) {
			flagName := arg[1:]
			flagValue := args[i+1]
			// Check if the flag is valid for the command
			validFlag := false
			for _, flag := range cmd.Flags {
				if flag.NameFlag() == flagName {
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
			i++ // Skip the next argument (flag value)
		}
	}

	// Set default values for flags
	// for _, flag := range cmd.Flags {
	// 	X == ?
	// 	if defaultValue := X; defaultValue != "" {
	// 		if _, exists := flagValues[flag.NameFlag()]; !exists {
	// 			flagValues[flag.NameFlag()] = defaultValue
	// 		}
	// 	}
	// }

	context := Context{Flags: flagValues}

	// Call the Action function with the context
	cmd.Action(context)
}

// PrintUsage prints usage information for the CLI program
func (c *CLI) PrintUsage() {
	fmt.Printf("Usage: %s %s \n\n", c.Name, c.Usage)
	fmt.Println("Commands:")
	for _, cmd := range c.commands {
		fmt.Printf("  %s\t%s\n", cmd.Name, cmd.Usage)
	}
}

// PrintCommandHelp prints help information for a specific command
func (c *CLI) PrintCommandHelp(name string) {
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
	cli := NewCLI("mycli", "<command> [options]", "1.0.0")

	echoCmd := Command{
		Name:    "echo",
		Usage:   "[-m <message>]",
		Aliases: []string{"say"},
		Help:    "Echoes the provided message",
		Flags: []Flags{
			String{
				Name:  "message",
				Value: "default",
			},
		},
		Action: func(c Context) {
			// Access flag values from the Context
			message, _ := c.Flags["message"].(string)
			fmt.Println("Echo:", message)
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
