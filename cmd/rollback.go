/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"alex067/tfstate/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"runtime"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var rollbackLatest bool

// rollbackCmd represents the rollback command
var rollbackCmd = &cobra.Command{
	Use:   "rollback",
	Short: "Rollback to an earlier State file version",
	Long:  `Rollback your State file to an earlier version by providing the backup file name located inside the tfstate folder.`,
	Run: func(cmd *cobra.Command, args []string) {
		if VersionFlag {
			log.Printf(outputVersion())
			return
		}

		if !IsRemote {
			log.Println("Rollback is disabled when using local state")
			return
		}

		rollbackLatest, _ := cmd.Flags().GetBool("latest")

		if !rollbackLatest && len(args) != 1 {
			log.Fatal("Must provide valid state file to restore")
		}

		currentStateOutput := utils.StateOutput(IsRemote, CurrentWorkingDirectory)
		var currentStateFile map[string]interface{}
		json.Unmarshal([]byte(currentStateOutput), &currentStateFile)

		currentStateSerial := int(currentStateFile["serial"].(float64))
		previousStateSerial := currentStateSerial - 1

		var restoreStateFile string

		if rollbackLatest {
			stateFullPath := fmt.Sprintf("%s/.terraform/tfstate", CurrentWorkingDirectory)
			statePrevious := fmt.Sprintf("state-%d", previousStateSerial)

			files, err := ioutil.ReadDir(stateFullPath)
			if err != nil {
				log.Fatal(err)
			}

			for _, f := range files {
				if len(f.Name()) >= len("state-00") && f.Name()[0:len(statePrevious)] == statePrevious {
					restoreStateFile = f.Name()
					break
				}
			}
		} else {
			restoreStateFile = args[0]
		}
		// TODO: add support to just provide serial and parse the latest state file

		// Read backup version of state file
		backupStateFile := utils.ReadStateFile(CurrentWorkingDirectory, restoreStateFile, IsRemote)
		var backupStatePayload State
		json.Unmarshal(backupStateFile, &backupStatePayload)

		backupStateSerial := backupStatePayload.Serial

		color.Set(color.Bold)
		log.Printf("Current state version: %d, Restore state version: %d", currentStateSerial, backupStateSerial)
		utils.ResetColor()

		if currentStateSerial == backupStateSerial {
			log.Println("Rollback command canceled.")
		}

		if currentStateSerial-backupStateSerial != 1 {
			log.Printf("Restore state version is older than the previous version of the current state, which is %d.", previousStateSerial)
		}

		var confirm string

		fmt.Println("Are you sure you wish to rollback the current state file? Only 'yes' will be accepted to confirm.")
		color.Set(color.Bold)
		fmt.Print("Enter a value: ")
		fmt.Scan(&confirm)

		utils.ResetColor()

		if confirm != "yes" {
			log.Println("Rollback command canceled.")
			return
		}

		fullRestoreStatePath := fmt.Sprintf("%s/%s/%s", CurrentWorkingDirectory, TfstateFullPath, restoreStateFile)
		if runtime.GOOS == "windows" {
			restoreCmd := exec.Command("cmd", "/c", "terraform", "state", "push", "--force", fullRestoreStatePath)

			var outb, errb bytes.Buffer
			restoreCmd.Stdout = &outb
			restoreCmd.Stderr = &errb

			if err := restoreCmd.Run(); err != nil {
				log.Fatal("Error running rollback command: ", errb.String(), err)
			}

			log.Printf("State rolled back to serial: %d", backupStateSerial)
			log.Printf("Restored file: %s", fullRestoreStatePath)
		}
	},
}

func init() {
	rootCmd.AddCommand(rollbackCmd)
	rollbackCmd.Flags().Bool("latest", false, "Rollback to the previous state version")
}
