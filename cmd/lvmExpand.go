/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/baimiyishu13/automate-all-the-things/internal/lvm"
	"github.com/spf13/cobra"
	"strconv"
)

// lvmExpandCmd represents the lvmExpand command
var lvmExpandCmd = &cobra.Command{
	Use:   "lvmExpand",
	Short: "lvmExpand <isLargeDisk> <expansionTarget> <selectedDisk>",
	Long:  `eg: lvmExpand false /test sde`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 3 {
			fmt.Println("Usage: lvmExpand <isLargeDisk> <expansionTarget> <selectedDisk>")
			return
		}

		isLargeDisk, err := strconv.ParseBool(args[0])
		if err != nil {
			fmt.Printf("Invalid value for isLargeDisk: %v\n", err)
			return
		}

		expansionTarget := args[1]
		selectedDisk := args[2]

		lvm.LvmExpand(isLargeDisk, expansionTarget, selectedDisk)
	},
}

func init() {
	rootCmd.AddCommand(lvmExpandCmd)
}
