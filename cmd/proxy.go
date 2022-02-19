/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/felixqin/p2p-tunnel/session"
	"github.com/spf13/cobra"
)

var (
	flagListenPort string
)

// proxyCmd represents the proxy command
var proxyCmd = &cobra.Command{
	Use:   "proxy",
	Short: "proxy manager command",
	Long:  "",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// proxyCreateCmd represents the proxy create command
var proxyCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "create proxy to tunnel stub",
	Long:  "",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		tunnelName := args[0]
		stub := args[1]
		fmt.Printf("create proxy over tunnel(%v)\n", tunnelName)
		client := session.FindClient(tunnelName)
		if client == nil {
			fmt.Printf("tunnel %v not found!\n", tunnelName)
			return
		}

		proxy := session.NewProxy(client, stub)
		go proxy.ListenAndServe(flagListenPort)
	},
}

// proxyCloseCmd represents the proxy close command
var proxyCloseCmd = &cobra.Command{
	Use:   "close",
	Short: "close proxy",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("proxy called")
	},
}

// proxyListCmd represents the proxy list command
var proxyListCmd = &cobra.Command{
	Use:   "list",
	Short: "show all proxies",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("proxy called")
	},
}

func init() {
	proxyCmd.AddCommand(proxyCreateCmd)
	proxyCmd.AddCommand(proxyCloseCmd)
	proxyCmd.AddCommand(proxyListCmd)
	registerCommand(proxyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// proxyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// proxyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	proxyCreateCmd.Flags().StringVarP(&flagListenPort, "listen", "l", "", "listen port")
}
