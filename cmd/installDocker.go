package cmd

import (
	"fmt"
	"github.com/baimiyishu13/automate-all-the-things/internal/docker"
	"github.com/spf13/cobra"
	"strings"
)

// installDockerCmd represents the installDocker command
var installDockerCmd = &cobra.Command{
	Use:   "installDocker <insecureRegistries> <base>",
	Short: "Install Docker with specific configurations",
	Long:  `Install Docker with specific configurations including sysctl parameters, rebuilding RPM database, and setting up daemon.json.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			fmt.Println("installDocker <insecureRegistries> <base>")
			return
		}
		insecureRegistries := strings.Split(strings.TrimSpace(args[0]), ",")
		base := strings.TrimSpace(args[1])

		// Call InstallDocker with user inputs
		if err := docker.InstallDocker(insecureRegistries, base); err != nil {
			fmt.Printf("Error installing Docker: %v\n", err)
		} else {
			fmt.Println("Docker installed successfully")
		}
	},
}

func init() {
	rootCmd.AddCommand(installDockerCmd)
}
