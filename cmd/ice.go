/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/felixqin/p2p-tunnel/session"
	"github.com/spf13/cobra"
)

// iceCmd represents the ice command
var iceCmd = &cobra.Command{
	Use:   "ice",
	Short: "ice servers manager command",
	Long:  "",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// iceAddCmd represents the ice add command
var iceAddCmd = &cobra.Command{
	Use:   "add",
	Short: "add ice server to list",
	Long:  "",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// fmt.Println("ice add called!", args)
		server := args[0]
		session.AddIceServer(server)
	},
}

// iceListCmd represents the ice list command
var iceListCmd = &cobra.Command{
	Use:   "list",
	Short: "show all ice servers",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		// fmt.Println("ice list called")
		servers := session.IceServers()
		for _, server := range *servers {
			fmt.Printf("%v\n", server)
		}
	},
}

func init() {
	iceCmd.AddCommand(iceAddCmd)
	iceCmd.AddCommand(iceListCmd)
	registerCommand(iceCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// iceCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// iceCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
