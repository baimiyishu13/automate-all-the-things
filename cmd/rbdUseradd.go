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
	Short: "rbdUseradd <user> <password> <email> <cluster> ",
	Long: `
	"store-test": "http://172.21.14.149:7070",
	"ent4-store": "https://rbdwg.hwwt2.com",
	"ent3-store": "https://rbdent3.hwwt2.com",
	"ent3-enterprise": "https://etrbd-prd-tc.hwwt2.com",
	"ent4-enterprise": "https://etrbd-prd.hwwt2.com",`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 3 {
			fmt.Println("rbdUseradd <user> <password> <email> <cluster> ")
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
