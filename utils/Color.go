package utils

import "github.com/fatih/color"

func ResetColor() {
	color.Unset()
	color.Set(color.FgWhite)
}
