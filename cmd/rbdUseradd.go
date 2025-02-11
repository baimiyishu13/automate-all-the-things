/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"paddle/internal/rbd"
)

// rbdUseraddCmd represents the rbdUseradd command
var rbdUseraddCmd = &cobra.Command{
	Use:   "rbdUseradd",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 3 {
			fmt.Println("Usage: rbdUseradd <user> <password> <cluster> ")
			return
		}

		user := args[0]
		password := args[1]
		email := args[2]
		cluster := args[3]

		rbd.CreateUser(user, password, email, cluster)
	},
}

func init() {
	rootCmd.AddCommand(rbdUseraddCmd)
}
