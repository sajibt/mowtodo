package pprint

import (
	"fmt"
	"github.com/fatih/color"
)

func Print(s string, colors ...color.Attribute) {
	color.Set(colors...)
	fmt.Print(s)
	color.Unset()
}

func Error(msg string) {
	Print("Error: ", color.FgRed, color.Bold)
	Print(msg, color.FgWhite)
	fmt.Println()
}

func Success(msg string) {
	Print("Success: ", color.FgGreen, color.Bold)
	Print(msg, color.FgWhite)
	fmt.Println()
}

