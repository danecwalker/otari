package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/danecwalker/otari/internal/commands"
	"github.com/danecwalker/otari/internal/podman"
	"github.com/danecwalker/otari/internal/spinners"
	"github.com/danecwalker/otari/internal/systemd"
	"github.com/fatih/color"
	"github.com/urfave/cli/v3"
)

var version = "dev"
var commit = "none"
var date = "unknown"

func main() {
	printLogo()

	cmd := &cli.Command{
		Name: "otari",
		Commands: []*cli.Command{
			{
				Name:  "version",
				Usage: "Print the version information",
				Action: func(ctx context.Context, c *cli.Command) error {
					fmt.Printf("Version: %s\nCommit: %s\nDate: %s\n", version, commit, date)
					return nil
				},
			},
			{
				Name:  "start",
				Usage: "Start the stack",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "file",
						Value:   "",
						Usage:   "Path to the stack definition file",
						Aliases: []string{"f"},
					},
				},
				Action: func(ctx context.Context, c *cli.Command) error {
					stackPath := c.String("file")
					systemCheck()
					commands.Start(ctx, stackPath)
					return nil
				},
			},
			{
				Name:  "stop",
				Usage: "Stop the stack",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "file",
						Value:   "",
						Usage:   "Path to the stack definition file",
						Aliases: []string{"f"},
					},
				},
				Action: func(ctx context.Context, c *cli.Command) error {
					stackPath := c.String("file")
					systemCheck()
					commands.Stop(ctx, stackPath)
					return nil
				},
			},
			{
				Name:  "remove",
				Usage: "Remove the stack",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "file",
						Value:   "",
						Usage:   "Path to the stack definition file",
						Aliases: []string{"f"},
					},
				},
				Action: func(ctx context.Context, c *cli.Command) error {
					stackPath := c.String("file")
					systemCheck()
					commands.Remove(ctx, stackPath)
					return nil
				},
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

func printLogo() {
	logo := `
 ██████╗ ████████╗ █████╗ ██████╗ ██╗
██╔═══██╗╚══██╔══╝██╔══██╗██╔══██╗██║
██║   ██║   ██║   ███████║██████╔╝██║
██║   ██║   ██║   ██╔══██║██╔══██╗██║
╚██████╔╝   ██║   ██║  ██║██║  ██║██║
 ╚═════╝    ╚═╝   ╚═╝  ╚═╝╚═╝  ╚═╝╚═╝
`

	c := color.New(color.FgMagenta, color.Bold)
	c.Println(logo)

	fmt.Println(color.WhiteString("A modern container orchestration tool"))
	fmt.Println()
}

func systemCheck() {
	sp := spinners.DefaultSpinner()
	sp.SetMessage("Checking for Podman...")

	// Check if Podman is installed and get version
	installed, podmanVersion := podman.PodmanVersion()
	if !installed {
		sp.FinishWithError("Podman is not installed. Please install Podman to use Otari.")
		os.Exit(1)
	}

	c := color.New(color.FgWhite)
	sp.Println(c.Sprint(" >  Podman intalled OK"))

	sp.SetMessage("Checking Podman version...")

	podmanMajorVersion, podmanMinorVersion, _ := podman.ParsePodmanVersion(podmanVersion)
	// must be greater than or equal to 4.4.0
	if podmanMajorVersion < 4 || (podmanMajorVersion == 4 && podmanMinorVersion < 4) {
		sp.FinishWithError("Podman version 4.4.0 or higher is required. Please upgrade Podman to use Otari.")
		os.Exit(1)
	}

	sp.Println(c.Sprintf(" >  Podman version %s OK", podmanVersion))

	sp.SetMessage("Checking if systemd is running...")
	// check if systemd is running
	if !systemd.IsSystemdRunning() {
		sp.FinishWithError("systemd is not running. Otari requires systemd to manage containers.")
		os.Exit(1)
	}

	sp.Println(c.Sprint(" >  systemd is running OK"))

	sp.SetMessage("Checking if user lingering is enabled...")

	// check for user lingering
	lingeringEnabled, err := systemd.IsUserLingeringEnabled()
	if err != nil {
		sp.FinishWithError(fmt.Sprintf("Failed to check user lingering: %v", err))
		os.Exit(1)
	}

	if !lingeringEnabled {
		sp.FinishWithError("User lingering is not enabled. Please enable user lingering to use Otari.")

		fmt.Println()

		fmt.Println("To enable user lingering, run the following command:")
		user := os.Getenv("USER")
		fmt.Println()
		fmt.Println(color.YellowString("    sudo loginctl enable-linger %s", user))
		fmt.Println()

		os.Exit(1)
	}

	sp.Println(c.Sprint(" >  User lingering is enabled OK"))

	// All checks passed
	sp.FinishWithSuccess("System checks passed!")
}
