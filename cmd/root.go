/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"os"
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

	switch args[0] {
	case "exit":
		os.Exit(0)
	case "help":
		for _, cmd := range commands {
			fmt.Printf("%-12v%v\n", cmd.Use, cmd.Short)
		}
	default:
		for _, cmd := range commands {
			if args[0] == cmd.Use {
				cmd.SetArgs(args[1:])
				err := cmd.Execute()
				if err != nil {
					fmt.Printf("execute cmdline(%v) failed! err(%v)\n", cmdline, err)
				}
				break
			}
		}
	}
}
