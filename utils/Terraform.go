package utils

import (
	"fmt"
	"os"
)

func IsTerraformInit(currentWorkingDir string) bool {
	fullPath := fmt.Sprintf("%s/%s", currentWorkingDir, ".terraform")
	_, err := os.Stat(fullPath)
	if err != nil || os.IsNotExist(err) {
		return false
	}
	return true
}

func IsTfstateInit(currentWorkingDir string) bool {
	fullPath := fmt.Sprintf("%s/%s/tfstate", currentWorkingDir, ".terraform")
	_, err := os.Stat(fullPath)
	if err != nil || os.IsNotExist(err) {
		return false
	}
	return true
}
