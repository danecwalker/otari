package utils

import (
	"fmt"

	"github.com/fatih/color"
)

func Success(msg string) string {
	return fmt.Sprintf("%s %s", color.New(color.FgGreen, color.Bold).Sprint("[✔]"), msg)
}

func Error(msg string) string {
	return fmt.Sprintf("%s %s", color.New(color.FgRed, color.Bold).Sprint("[✘]"), msg)
}

func Info(msg string) string {
	return fmt.Sprintf("%s %s", color.New(color.FgCyan, color.Bold).Sprint("[i]"), msg)
}
