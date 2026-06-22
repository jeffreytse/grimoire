package tui

import "github.com/charmbracelet/lipgloss"

var (
	StyleTitle    = lipgloss.NewStyle().Bold(true)
	StyleHint     = lipgloss.NewStyle().Faint(true)
	StyleCursor   = lipgloss.NewStyle().Foreground(lipgloss.Color("214"))
	StyleSelected = lipgloss.NewStyle().Bold(true)
	StyleCheck    = lipgloss.NewStyle().Foreground(lipgloss.Color("82"))
	StyleCircle   = lipgloss.NewStyle().Faint(true)
	StyleDim      = lipgloss.NewStyle().Faint(true)
	StyleGold     = lipgloss.NewStyle().Foreground(lipgloss.Color("178"))
	StyleCyan     = lipgloss.NewStyle().Foreground(lipgloss.Color("36"))
	StyleGreen    = lipgloss.NewStyle().Foreground(lipgloss.Color("82"))
	StyleRed      = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
	StyleYellow   = lipgloss.NewStyle().Foreground(lipgloss.Color("214"))
	StyleBold     = lipgloss.NewStyle().Bold(true)

	IconOK    = StyleGreen.Render("✅")
	IconWarn  = StyleYellow.Render("⚠️ ")
	IconFail  = StyleRed.Render("❌")
	IconSkip  = "⬜"
	IconDone  = StyleGreen.Render("✓")
	IconError = StyleRed.Render("✗")
)
