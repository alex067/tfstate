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
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func handleActionRm(args []string) {
	stateFileName := "terraform.tfstate"

	if IsRemote {
		stateFileName = utils.StateDownload(CurrentWorkingDirectory)
	}

	stateFile := utils.ReadStateFile(CurrentWorkingDirectory, stateFileName, IsRemote)

	var payload State
	affectedDependencies := []string{}

	json.Unmarshal(stateFile, &payload)

	for _, arg := range args {
		isModule := arg[0:len("module")] == "module"

		for resource, _ := range payload.Resources {
			resourceAddress := fmt.Sprintf("%s.%s", payload.Resources[resource].Type, payload.Resources[resource].Name)

			if payload.Resources[resource].Module != "" {
				resourceAddress = fmt.Sprintf("%s.%s", payload.Resources[resource].Module, resourceAddress)
			}

			for instance, _ := range payload.Resources[resource].Instances {
				dependencies := payload.Resources[resource].Instances[instance].Dependencies

				if isModule && utils.ModuleContains(dependencies, arg) {
					affectedDependencies = append(affectedDependencies, resourceAddress)
				} else if !isModule && utils.Contains(dependencies, arg) {
					affectedDependencies = append(affectedDependencies, resourceAddress)
				}
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

	rmCommand := []string{}
	rmCommand = append(rmCommand, args...)
	rmCommand = append([]string{"rm"}, rmCommand...)
	rmCommand = append([]string{"state"}, rmCommand...)
	rmCommand = append([]string{"terraform"}, rmCommand...)

	if runtime.GOOS == "windows" {
		rmCommand = append([]string{"/c"}, rmCommand...)
		stateCmd := exec.Command("cmd", rmCommand...)

		stateCmd.Stdout = &outb
		stateCmd.Stderr = &errb

		if err := stateCmd.Run(); err != nil {
			log.Fatal("Error running state command: ", outb.String())
		}
	} else if runtime.GOOS == "darwin" {
		// pop terraform
		rmCommand = rmCommand[1:]
		stateCmd := exec.Command("terraform", rmCommand...)
		stateCmd.Stdout = &outb
		stateCmd.Stderr = &errb

		if err := stateCmd.Run(); err != nil {
			log.Fatal("Error running state command: ", outb.String())
		}
	}

	log.Println(outb.String())
}

var rmCmd = &cobra.Command{
	Use:   "rm",
	Short: "A wrapper around terraform state rm",
	Long: `A wrapper around terraform state rm, adding automatic state backup to easily rollback changes,
	and a confirmation step to first review affected resources.`,
	Run: func(cmd *cobra.Command, args []string) {
		if VersionFlag {
			log.Printf(outputVersion())
			return
		}

		if !utils.IsTerraformInit(CurrentWorkingDirectory) {
			log.Fatal("Run terraform init to initialize .terraform used by tfstate")
		}

		if len(args) == 0 {
			log.Fatal("Must provide valid terraform state rm command")
		}

		if IsRemote && !utils.IsTfstateInit(CurrentWorkingDirectory) {
			dirFullPath := fmt.Sprintf("%s/.terraform/tfstate", CurrentWorkingDirectory)
			err := os.Mkdir(dirFullPath, 0777)
			if err != nil {
				log.Fatal("Error initializing tfstate directory: ", err)
			}
			log.Println("Initializing tfstate directory")
		}

		handleActionRm(args)
	},
}

func init() {
	rootCmd.AddCommand(rmCmd)
}
