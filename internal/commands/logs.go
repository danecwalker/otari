package commands

import (
	"context"
	"fmt"
	"os"

	"github.com/danecwalker/otari/internal/definition"
	"github.com/danecwalker/otari/internal/systemd"
	"github.com/danecwalker/otari/internal/utils"
	"github.com/fatih/color"
)

func Logs(ctx context.Context, stackPath string, containerName string) {
	if stackPath == "" {
		stackPath = utils.DefaultStackPath()
	}
	c, err := os.ReadFile(stackPath)
	if err != nil {
		fmt.Println(utils.Error("Failed to read " + stackPath))
		color.New(color.FgWhite).Println("    " + err.Error())
		os.Exit(1)
	}

	stack, err := definition.Parse(c)
	if err != nil {
		fmt.Println(utils.Error("Failed to parse stack definition"))
		color.New(color.FgWhite).Println("    " + err.Error())
		os.Exit(1)
	}

	// check if container name is valid
	if containerName != "" {
		if _, exists := stack.Containers[containerName]; !exists {
			fmt.Println(utils.Error("Container '" + containerName + "' not found in stack definition"))
			os.Exit(1)
		}
	}

	logs, err := systemd.GetLogs(containerName)
	if err != nil {
		fmt.Println(utils.Error("Failed to get logs"))
		color.New(color.FgWhite).Println("    " + err.Error())
		os.Exit(1)
	}

	fmt.Println(string(logs))
}
