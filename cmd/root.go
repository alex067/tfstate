/*
Copyright © 2022 NAME HERE panayi067@gmail.com

*/
package cmd

import (
	"alex067/tfstate/utils"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
)

var Version string
var VersionFlag bool
var CurrentWorkingDirectory string
var TfstateFullPath = ".terraform/tfstate"
var IsRemote bool

// Define remot state json lookup
type Dependencies struct {
	Dependencies []string
}

type Instances struct {
	Attributes     map[string]struct{}
	Dependencies   []string
	Schema_Version int
}

type Resources struct {
	Module    string
	Mode      string
	Type      string
	Name      string
	Provider  string
	Instances []Instances
}

type State struct {
	Serial    int
	Resources []Resources
}

type AffectedDependencies struct {
	Address string
}

func outputVersion() string {
	var terraformVersion string
	var stdout []byte
	var err error

	if runtime.GOOS == "windows" {
		versionCmd := exec.Command("cmd", "/c", "terraform", "--version")
		stdout, err = versionCmd.Output()

	} else if runtime.GOOS == "darwin" {
		versionCmd := exec.Command("terraform", "--version")
		stdout, err = versionCmd.Output()
	}

	if err != nil || len(stdout) == 0 {
		log.Fatal("Unable to execute terraform --version")
		os.Exit(1)
	}

	versionCmdOutput := string([]byte(stdout))
	terraformVersion = strings.Split(versionCmdOutput, "\n")[0]

	if terraformVersion[0:len("Terraform")] != "Terraform" {
		log.Fatal("Unable to read the terraform version")
		os.Exit(1)
	}

	return fmt.Sprintf("%s, Tfstate v%s", terraformVersion, Version)
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "tfstate",
	Short: "A wrapper around terraform state",
	Long:  `tfstate provides simple guard rails and automatic backup recovery when running state commands.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(Version)
		if len(args) == 0 && !VersionFlag {
			cmd.Help()
			return
		}

		if VersionFlag {
			log.Printf(outputVersion())
			return
		}
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	workingDirectory, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	CurrentWorkingDirectory = workingDirectory
	IsRemote = utils.IsRemoteState(CurrentWorkingDirectory)

	rootCmd.PersistentFlags().BoolVarP(&VersionFlag, "version", "v", false, "get the current version of tfstate and terraform")
}
