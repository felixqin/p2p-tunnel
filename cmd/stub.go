/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/felixqin/p2p-tunnel/session"
	"github.com/spf13/cobra"
)

// stubCmd represents the proxy command
var stubCmd = &cobra.Command{
	Use:   "stub",
	Short: "stub manager command",
	Long:  "",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// stubEnableCmd represents the stub enable command
var stubEnableCmd = &cobra.Command{
	Use:   "enable",
	Short: "enable stub upstream",
	Long:  "",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("stub enable called")
		name := args[0]
		upstream := args[1]
		session.EnableStub(name, upstream)
	},
}

// stubDisableCmd represents the stub disable command
var stubDisableCmd = &cobra.Command{
	Use:   "disable",
	Short: "disable stub upstream",
	Long:  "",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("stub disable called")
		name := args[0]
		session.DisableStub(name)
	},
}

// stubListCmd represents the stub list command
var stubListCmd = &cobra.Command{
	Use:   "list",
	Short: "list stub",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("stub list called")
		session.DumpStubs()
	},
}

func init() {
	stubCmd.AddCommand(stubEnableCmd)
	stubCmd.AddCommand(stubDisableCmd)
	stubCmd.AddCommand(stubListCmd)
	registerCommand(stubCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// stubCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// stubCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
