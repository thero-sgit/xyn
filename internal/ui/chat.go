package ui

import (
	"fmt"
	"os/user"

	"github.com/charmbracelet/lipgloss"
)

var currentUser, currentUserRrr = user.Current()
var username = func() string {
		if currentUserRrr != nil {
			return "user"
		}
		return currentUser.Username
}()

func newUserPrompt(prompt string, width int) string {
	return lipgloss.NewStyle().
		MarginBottom(1).
		Width(width).
		PaddingLeft(1).
		PaddingRight(1).
		Border(lipgloss.NormalBorder(), false, false, false, true).
		Render(
			lipgloss.JoinHorizontal(
				lipgloss.Left,
				subtleStyle.Italic(true).Render(username) + lipgloss.NewStyle().Foreground(lipgloss.Color("#E47753")).Render(" $ "),
				lipgloss.NewStyle().Render(prompt),
			),
		)
}

type agentBackgroundActivityLabel struct {
	rawString      string
	loaderIndex    int
	loaderFrames   []string
	loader         string
	prettyString   string
	rawStringStyle lipgloss.Style
	loaderStyle    lipgloss.Style
}

func newAgentBackgroundActivity(activity string) agentBackgroundActivityLabel {
	frames := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	loader := frames[0]

	rawStringStyle := lipgloss.NewStyle() 
	loaderStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#E47753")).Bold(true)


	prettyString := fmt.Sprintf(
		"%s %s...",
		loaderStyle.Render(loader),
		rawStringStyle.Render(activity),
	)

	return agentBackgroundActivityLabel {
		rawString: activity,
		loaderIndex: 0,
		loaderFrames: frames,
		loader: loader,
		prettyString: prettyString,
		rawStringStyle: rawStringStyle,
		loaderStyle: loaderStyle,
	}
}

func (abal *agentBackgroundActivityLabel) animate() {
	abal.loaderIndex = (abal.loaderIndex + 1) % len(abal.loaderFrames)

	abal.prettyString = fmt.Sprintf(
		"%s %s...",
		abal.loaderStyle.Render(abal.loaderFrames[abal.loaderIndex]),
		abal.rawStringStyle.Render(abal.rawString),
	)
}