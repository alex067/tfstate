/*
Copyright © 2022 NAME HERE panayi067@gmail.com

*/
package cmd

import (
	"alex067/tfstate/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"reflect"
	"runtime"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// mvCmd represents the mv command
var mvCmd = &cobra.Command{
	Use:   "mv",
	Short: "A wrapper around terraform state mv",
	Long: `A wrapper around terraform state mv, adding automatic state backup to easily rollback changes,
and a confirmation step to first review affected resources.`,
	Run: func(cmd *cobra.Command, args []string) {
		if VersionFlag {
			log.Printf(outputVersion())
			return
		}

		if len(args) != 2 {
			log.Fatal("Must provide valid terraform state mv command")
		}

		stateFileName := "terraform.tfstate"

		if IsRemote {
			stateFileName = utils.StateDownload(CurrentWorkingDirectory)
		}

		stateFile := utils.ReadStateFile(CurrentWorkingDirectory, stateFileName, IsRemote)

		var payload State
		affectedDependencies := []string{}

		json.Unmarshal(stateFile, &payload)

		isModule := args[0][0:len("module")] == "module"

		for resource, _ := range payload.Resources {
			resourceAddress := fmt.Sprintf("%s.%s", payload.Resources[resource].Type, payload.Resources[resource].Name)
			if reflect.TypeOf(payload.Resources[resource].Module).Kind() == reflect.String {
				resourceAddress = fmt.Sprintf("%s.%s", payload.Resources[resource].Module, resourceAddress)
			}

			for instance, _ := range payload.Resources[resource].Instances {
				dependencies := payload.Resources[resource].Instances[instance].Dependencies

				if isModule && utils.ModuleContains(dependencies, args[0]) {
					affectedDependencies = append(affectedDependencies, resourceAddress)
				} else if !isModule && utils.Contains(dependencies, args[0]) {
					affectedDependencies = append(affectedDependencies, resourceAddress)
				}
			}
		}

		affectedDependenciesString := ""

		if len(affectedDependencies) > 0 {
			affectedDependenciesString = fmt.Sprintf("%s", strings.Join(affectedDependencies, "\n"))
		}

		log.Println("Calculating affected resources...")
		if affectedDependenciesString != "" {
			fmt.Println(affectedDependenciesString)
		}

		color.Set(color.Bold)
		log.Printf("Affected resources: %d", len(affectedDependencies))
		fmt.Println("Do you really want to alter your state file?")

		utils.ResetColor()

		fmt.Println("Changes to your state file will affect the resources listed above.")
		fmt.Println("Undo changes by running the rollout command, and supplying the back up state file generated earlier. Only 'yes' will be accepted to confirm.")

		var confirmation string
		color.Set(color.Bold)
		fmt.Print("Enter a value: ")
		fmt.Scanln(&confirmation)

		utils.ResetColor()

		if confirmation != "yes" {
			log.Println("State command canceled.")
			return
		}

		var outb, errb bytes.Buffer

		if runtime.GOOS == "windows" {
			stateCmd := exec.Command("cmd", "/c", "terraform", "state", "mv", args[0], args[1])
			stateCmd.Stdout = &outb
			stateCmd.Stderr = &errb

			if err := stateCmd.Run(); err != nil {
				log.Fatal("Error running state command: ", outb.String())
			}
		} else if runtime.GOOS == "darwin" {
			stateCmd := exec.Command("terraform", "state", "mv", args[0], args[1])
			stateCmd.Stdout = &outb
			stateCmd.Stderr = &errb

			if err := stateCmd.Run(); err != nil {
				log.Fatal("Error running state command: ", outb.String())
			}
		}
		log.Println(outb.String())
	},
}

func init() {
	rootCmd.AddCommand(mvCmd)
}
