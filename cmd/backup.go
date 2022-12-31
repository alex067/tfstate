/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"alex067/tfstate/utils"
	"log"

	"github.com/spf13/cobra"
)

// backupCmd represents the backup command
var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Backup the current state file",
	Long:  `Backup the current state file into the tfstate folder.`,
	Run: func(cmd *cobra.Command, args []string) {
		if VersionFlag {
			log.Printf(outputVersion())
			return
		}

		if IsRemote {
			utils.StateDownload(CurrentWorkingDirectory)
		}
	},
}

func init() {
	rootCmd.AddCommand(backupCmd)
}
