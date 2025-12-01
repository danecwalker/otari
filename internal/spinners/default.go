package spinners

import (
	"time"

	"github.com/danecwalker/otari/pkg/spinner"
	"github.com/fatih/color"
)

func DefaultSpinner() *spinner.Spinner {
	sp := spinner.NewCustom(
		[]string{
			"[|]", "[/]", "[-]", "[\\]",
		},
		spinner.WithSuccessSymbol(color.New(color.FgGreen, color.Bold).Sprint("[✔]")),
		spinner.WithErrorSymbol(color.New(color.FgRed, color.Bold).Sprint("[✘]")),
		spinner.WithInfoSymbol(color.New(color.FgCyan, color.Bold).Sprint("[i]")),
	)
	sp.Enable(100 * time.Millisecond)
	return sp
}
