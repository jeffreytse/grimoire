package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type selectModel struct {
	title  string
	items  []string
	cursor int
	chosen string
	quit   bool
}

func (m *selectModel) Init() tea.Cmd { return nil }

func (m *selectModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			m.quit = true
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.items)-1 {
				m.cursor++
			}
		case "enter", " ":
			m.chosen = m.items[m.cursor]
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m *selectModel) View() string {
	var sb strings.Builder
	sb.WriteString("\n")
	sb.WriteString(StyleTitle.Render(m.title))
	sb.WriteString("\n")
	sb.WriteString(StyleHint.Render("  ↑↓ navigate   ENTER confirm   q quit"))
	sb.WriteString("\n\n")
	for i, item := range m.items {
		if i == m.cursor {
			fmt.Fprintf(&sb, "  %s %s\n",
				StyleCursor.Render("▶"),
				StyleSelected.Render(item))
		} else {
			fmt.Fprintf(&sb, "    %s\n", StyleDim.Render(item))
		}
	}
	return sb.String()
}

// RunSelect shows a single-select menu. Returns chosen item and ok=false if cancelled.
func RunSelect(title string, items []string) (string, bool) {
	m := selectModel{title: title, items: items}
	p := tea.NewProgram(&m, tea.WithOutput(selectOutput()))
	final, err := p.Run()
	if err != nil {
		return "", false
	}
	res, ok := final.(*selectModel)
	if !ok {
		return "", false
	}
	if res.quit || res.chosen == "" {
		return "", false
	}
	return res.chosen, true
}
