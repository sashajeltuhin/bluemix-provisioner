package main

import (
	"os"

	"github.com/sashajeltuhin/bluemix-provisioner/provision/softlayer"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "provision",
	Short: "Provision infrastructure for ACP",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
	SilenceUsage: true,
}

func init() {
	rootCmd.AddCommand(softlayer.Cmd())
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}
