/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/felixqin/p2p-tunnel/contacts"
	"github.com/spf13/cobra"
)

var (
	flagNodeName string
	flagUser     string
	flagPassword string
)

// connectCmd represents the connect command
var connectCmd = &cobra.Command{
	Use:   "connect",
	Short: "connect message server",
	Long:  "connect mqtt message broker.",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("connect called!", args)
		server := args[0]
		contacts.Open(&contacts.Option{
			Name:     flagNodeName,
			Server:   server,
			Username: flagUser,
			Password: flagPassword,
		})
	},
}

func init() {
	registerCommand(connectCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// connectCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// connectCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	connectCmd.Flags().StringVarP(&flagNodeName, "node", "n", "", "node name")
	connectCmd.Flags().StringVarP(&flagUser, "user", "u", "", "user name")
	connectCmd.Flags().StringVarP(&flagPassword, "password", "s", "", "password of account")
}
