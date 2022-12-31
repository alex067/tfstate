/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"alex067/tfstate/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"runtime"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// rmCmd represents the rm command
var rmCmd = &cobra.Command{
	Use:   "rm",
	Short: "A wrapper around terraform state rm",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if VersionFlag {
			log.Printf(outputVersion())
			return
		}

		if len(args) != 1 {
			log.Fatal("Must provide valid terraform state rm command")
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

			if payload.Resources[resource].Module != "" {
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

		rmCommand := []string{}
		rmCommand = append([]string{args[0]}, rmCommand...)
		rmCommand = append([]string{"rm"}, rmCommand...)
		rmCommand = append([]string{"state"}, rmCommand...)
		rmCommand = append([]string{"terraform"}, rmCommand...)

		initCommand := ""

		if runtime.GOOS == "windows" {
			rmCommand = append([]string{"/c"}, rmCommand...)
			initCommand = "cmd"
		}

		stateCmd := exec.Command(initCommand, rmCommand...)

		var outb, errb bytes.Buffer
		stateCmd.Stdout = &outb
		stateCmd.Stderr = &errb

		if err := stateCmd.Run(); err != nil {
			log.Fatal("Error running state command: ", outb.String())
		}

		log.Println(outb.String())
	},
}

func init() {
	rootCmd.AddCommand(rmCmd)
}
