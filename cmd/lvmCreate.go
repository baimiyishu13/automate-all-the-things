/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/baimiyishu13/automate-all-the-things/internal/lvm"
	"github.com/spf13/cobra"
)

// lvmCreateCmd represents the lvmCreate command
var lvmCreateCmd = &cobra.Command{
	Use:   "lvmCreate",
	Short: "lvmCreate <disk> <mountPoint> <fsType>",
	Long:  `eg：lvmCreate sdb /mnt ext4`,
	Run: func(cmd *cobra.Command, args []string) {
		disk := args[0]
		mountPoint := args[1]
		fsType := args[2]

		lvm.LvmCreate(disk, mountPoint, fsType)
	},
}

func init() {
	rootCmd.AddCommand(lvmCreateCmd)
}
