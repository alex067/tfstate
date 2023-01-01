package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"time"
)

type Backend struct {
	Type string
}

type RemoteState struct {
	Backend Backend
}

func IsRemoteState(currentWorkingDirectory string) bool {
	stateFile, err := ioutil.ReadFile(fmt.Sprintf("%s/.terraform/terraform.tfstate", currentWorkingDirectory))
	if err != nil {
		return false
	}

	var payload RemoteState
	json.Unmarshal(stateFile, &payload)

	return true
}

func StateDownload(currentWorkingDir string) string {
	stateFileOutput := StateOutput(true, currentWorkingDir)
	var currentStateFile map[string]interface{}
	json.Unmarshal([]byte(stateFileOutput), &currentStateFile)

	stateSerial := int(currentStateFile["serial"].(float64))
	currentTime := strconv.FormatInt(time.Now().Unix(), 10)

	// Check if state serial backup already exists
	stateFullPath := fmt.Sprintf("%s/.terraform/tfstate", currentWorkingDir)
	files, err := ioutil.ReadDir(stateFullPath)
	if err != nil {
		log.Fatal(err)
	}

	stateFileName := fmt.Sprintf("state-%d-%s.json", stateSerial, currentTime)
	fullPath := fmt.Sprintf("%s/.terraform/tfstate/%s", currentWorkingDir, stateFileName)

	if runtime.GOOS == "windows" {
		cmd := exec.Command("cmd", "/c", "terraform", "state", "pull", ">", fullPath)

		_, err := cmd.Output()
		if err != nil {
			log.Fatal("Error performing state file backup: ", err)
		}
	} else if runtime.GOOS == "darwin" {
		cmd := exec.Command("terraform", "state", "pull", ">", fullPath)

		_, err := cmd.Output()
		if err != nil {
			log.Fatal("Error performing state file backup: ", err)
		}
	}

	log.Printf("State file backup: %s", stateFileName)

	// Remove existing state backups
	stateBackupName := fmt.Sprintf("state-%d", stateSerial)
	stateBackupsToDelete := []string{}

	for _, f := range files {
		if len(f.Name()) >= len(stateBackupName) && f.Name()[0:len(stateBackupName)] == stateBackupName {
			if f.Name() == stateFileName {
				continue
			}
			stateBackupsToDelete = append(stateBackupsToDelete, f.Name())
		}
	}

	if len(stateBackupsToDelete) > 0 {
		for _, val := range stateBackupsToDelete {
			os.Remove(fmt.Sprintf("%s/.terraform/tfstate/%s", currentWorkingDir, val))
		}
	}

	return stateFileName
}

func StateOutput(isRemote bool, currentWorkingDir string) string {
	var stdout []byte
	var err error

	if runtime.GOOS == "windows" {
		if isRemote {
			cmd := exec.Command("cmd", "/c", "terraform", "state", "pull")

			stdout, err = cmd.Output()
			if err != nil {
				log.Fatal("Error reading state file: ", err)
			}
		}
	} else if runtime.GOOS == "darwin" {
		if isRemote {
			cmd := exec.Command("terraform", "state", "pull")

			stdout, err = cmd.Output()
			if err != nil {
				log.Fatal("Error reading state file: ", err)
			}
		}
	}
	statefile := string([]byte(stdout))
	return statefile
}

func ReadStateFile(currentWorkingDir string, stateFileName string, isRemote bool) []byte {
	// Read backup version of state file
	remoteDir := ""
	if isRemote {
		remoteDir = ".terraform/tfstate/"
	}

	stateFile, err := ioutil.ReadFile(fmt.Sprintf("%s/%s%s", currentWorkingDir, remoteDir, stateFileName))
	if err != nil {
		log.Fatalf("Error reading %s: %s", stateFileName, err)
	}
	return stateFile
}
