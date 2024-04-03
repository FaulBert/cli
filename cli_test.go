package cli

import (
	"errors"
	"fmt"
	"testing"
)

func TestSubCommandApp(t *testing.T) {
	// initiate dummy app
	dummyApp := App{}
	dummyApp.Name = "uwe"
	cmdAlias := []string{"r", "rnu", "nur"}
	cmdFlags := []Flag{
		&StringFlag{
			Name:  "m",
			Value: "moe",
		},
	}
	cmd := &Command{
		Name:        "run",
		Usage:       "",
		Alias:       cmdAlias,
		Flags:       cmdFlags,
		Description: "simply run",
		Action: func(ctx Context) {
			flagValue := ctx.String().Get("m")

			fmt.Printf("run command invoked with m flag value %s \n", flagValue)
		},
	}
	dummyApp.AddCommand(cmd)

	// testing scenarios
	okArgsTest := [][]string{
		{""},
		// Flag test
		{"", "-h"},

		// Subcommand Alias test
		{"", "run"},
		{"", "r"},
		{"", "rnu"},
		{"", "nur"},

		// Subcommand With Flags

		{"", "run", "-m", "32"},
	}

	for _, args := range okArgsTest {
		if err := dummyApp.Run(args); err != nil {
			t.Error(err)
			t.Fail()
		}
	}

	notOkArgsTest := []struct {
		Args      []string
		ShouldErr error
	}{
		{
			[]string{"", "nazan-cute-uwu"},
			ErrCommandNotRegistered("nazan-cute-uwu"),
		},
	}

	for _, test := range notOkArgsTest {

		if appErr := dummyApp.Run(test.Args); errors.Is(appErr, test.ShouldErr) {
			fmt.Println("error get : ", appErr)
			fmt.Println("error should: ", test.ShouldErr)
			t.Fail()
		}
	}

}
