/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/felixqin/p2p-tunnel/contacts"
	"github.com/felixqin/p2p-tunnel/session"
	"github.com/spf13/cobra"
)

// tunnelCmd represents the tunnel command
var tunnelCmd = &cobra.Command{
	Use:   "tunnel",
	Short: "tunnel manager command",
	Long:  "",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// tunnelConnectCmd represents the tunnel connect command
var tunnelConnectCmd = &cobra.Command{
	Use:   "connect",
	Short: "create p2p tunnel to connect server node",
	Long:  "",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		nodeName := args[0]
		fmt.Println("create tunnel to connect", nodeName)
		contact, err := contacts.FindContact(nodeName)
		if err != nil {
			fmt.Printf("node %v not found!\n", nodeName)
			return
		}

		err = session.Connect(contact.ClientId)
		if err != nil {
			fmt.Printf("connect %v failed!\n", nodeName)
			return
		}
	},
}

// tunnelCloseCmd represents the tunnel close command
var tunnelCloseCmd = &cobra.Command{
	Use:   "close",
	Short: "close p2p tunnel",
	Long:  "",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("tunnel close called")
	},
}

// tunnelListCmd represents the tunnel list command
var tunnelListCmd = &cobra.Command{
	Use:   "list",
	Short: "show all tunnels",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("tunnel list called")
	},
}

func init() {
	tunnelCmd.AddCommand(tunnelConnectCmd)
	tunnelCmd.AddCommand(tunnelCloseCmd)
	tunnelCmd.AddCommand(tunnelListCmd)
	registerCommand(tunnelCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// tunnelCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// tunnelCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
