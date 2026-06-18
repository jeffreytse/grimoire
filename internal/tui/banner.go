package tui

import (
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
)

func PrintBanner(version string) {
	gold := lipgloss.NewStyle().Foreground(lipgloss.Color("178"))
	star := lipgloss.NewStyle().Foreground(lipgloss.Color("227"))
	line := lipgloss.NewStyle().Foreground(lipgloss.Color("100"))
	dim := lipgloss.NewStyle().Faint(true)
	bold := lipgloss.NewStyle().Bold(true)
	ver := lipgloss.NewStyle().Faint(true)
	cyan := lipgloss.NewStyle().Foreground(lipgloss.Color("36"))

	fmt.Fprintln(os.Stderr)
	fmt.Fprintf(os.Stderr, " %s%s%s  %s %s\n",
		star.Render("✦"),
		gold.Render("▗▄▄▄▄▄▄▄▄▖"),
		star.Render("✦"),
		bold.Render("grimoire"),
		ver.Render("v"+version),
	)
	fmt.Fprintf(os.Stderr, "  %s%s%s   %s\n",
		gold.Render("▐"),
		line.Render("▬▬▬│▬▬▬▬"),
		gold.Render("▌"),
		dim.Render("The world's best practices for AI assistants"),
	)
	fmt.Fprintf(os.Stderr, "  %s%s%s   %s\n",
		gold.Render("▐"),
		line.Render("▬▬ ✦ ▬▬▬"),
		gold.Render("▌"),
		cyan.Render("https://github.com/jeffreytse/grimoire"),
	)
	fmt.Fprintf(os.Stderr, "  %s%s%s\n",
		gold.Render("▐"),
		line.Render("▬▬▬│▬▬▬▬"),
		gold.Render("▌"),
	)
	fmt.Fprintf(os.Stderr, " %s%s%s  ⭐ %s  💖 %s  🐛 %s\n",
		star.Render("✦"),
		gold.Render("▝▀▀▀▀▀▀▀▀▘"),
		star.Render("✦"),
		osc8("https://github.com/jeffreytse/grimoire", "Star"),
		osc8("https://github.com/sponsors/jeffreytse", "Sponsor"),
		osc8("https://github.com/jeffreytse/grimoire/issues", "Issues"),
	)
	fmt.Fprintln(os.Stderr)
}

func osc8(url, text string) string {
	return "\033]8;;" + url + "\007" + text + "\033]8;;\007"
}
