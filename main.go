package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/thero-sgit/xyn/internal/ai"
	"github.com/thero-sgit/xyn/internal/config"
	"github.com/thero-sgit/xyn/internal/ui"
)

func main() {
	ai.InitSession()
	config.Init()
	

	program := tea.NewProgram(
		ui.InitialModel(),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	if _, err := program.Run(); err != nil {
		fmt.Println("error: ", err)
		os.Exit(1)
	}
}
