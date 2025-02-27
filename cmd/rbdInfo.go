/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/baimiyishu13/automate-all-the-things/internal/rbd/rbdmysql"
	"github.com/spf13/cobra"
)

// rbdInfoCmd represents the rbdInfo command
var rbdInfoCmd = &cobra.Command{
	Use:   "rbdInfo",
	Short: "A brief description of your command",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 5 {
			fmt.Println("rbdInfo sip port user 'password' env")
			return
		}
		sip := args[0]
		port := args[1]
		user := args[2]
		password := args[3]
		env := args[4]

		rbdmysql.RunSQL(sip, port, user, password, env)

	},
}

func init() {
	rootCmd.AddCommand(rbdInfoCmd)
}
