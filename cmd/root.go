/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var commands = []*cobra.Command{}

func registerCommand(cmd *cobra.Command) {
	commands = append(commands, cmd)
}

func Execute(cmdline string) {
	args := strings.Fields(cmdline)
	if len(args) == 0 {
		return
	}

	for _, cmd := range commands {
		if cmd.Use == args[0] {
			cmd.SetArgs(args[1:])
			err := cmd.Execute()
			if err != nil {
				fmt.Printf("execute cmdline(%v) failed! err(%v)\n", cmdline, err)
			}
			break
		}
	}
}
