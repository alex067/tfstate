package utils

import (
	"fmt"
	"strings"
)

func Contains(arr []string, lookup string) bool {
	for _, val := range arr {
		if strings.Contains(val, lookup) {
			return true
		}
	}
	return false
}

func ModuleContains(arr []string, lookup string) bool {
	isRootModule := strings.Count(lookup, ".") == 1 && len(strings.Split(lookup, ".")) == 2

	for _, val := range arr {
		if isRootModule && val[0:len(lookup)+1] == fmt.Sprintf("%s.", lookup) {
			return true
		} else if strings.Contains(val, lookup) {
			return true
		}
	}
	return false
}
