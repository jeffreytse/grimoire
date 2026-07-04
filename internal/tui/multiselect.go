package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type multiselectModel struct {
	title    string
	items    []string
	selected []bool
	cursor   int
	offset   int
	maxVis   int
	quit     bool
}

func newMultiselectModel(title string, items []string, preSelected []bool) multiselectModel {
	sel := make([]bool, len(items))
	if preSelected != nil {
		for i := range sel {
			if i < len(preSelected) {
				sel[i] = preSelected[i]
			}
		}
	}
	return multiselectModel{
		title:    title,
		items:    items,
		selected: sel,
		maxVis:   12,
	}
}

func (m *multiselectModel) Init() tea.Cmd { return nil }

func (m *multiselectModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.maxVis = msg.Height - 8
		if m.maxVis < 3 {
			m.maxVis = 3
		}
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
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
		case " ":
			m.selected[m.cursor] = !m.selected[m.cursor]
		case "a", "A":
			anyOff := false
			for _, s := range m.selected {
				if !s {
					anyOff = true
					break
				}
			}
			for i := range m.selected {
				m.selected[i] = anyOff
			}
		case "enter":
			return m, tea.Quit
		}
		// scroll
		if m.cursor < m.offset {
			m.offset = m.cursor
		}
		if m.cursor >= m.offset+m.maxVis {
			m.offset = m.cursor - m.maxVis + 1
		}
	}
	return m, nil
}

func (m *multiselectModel) View() string {
	var sb strings.Builder
	sb.WriteString("\n")
	sb.WriteString(StyleTitle.Render(m.title))
	sb.WriteString("\n")
	sb.WriteString(StyleHint.Render("  ↑↓ navigate   SPACE toggle   A select all   ENTER confirm   ESC cancel"))
	sb.WriteString("\n\n")

	end := m.offset + m.maxVis
	if end > len(m.items) {
		end = len(m.items)
	}

	selCount := 0
	for _, s := range m.selected {
		if s {
			selCount++
		}
	}

	for i := m.offset; i < end; i++ {
		var mark, name string
		if m.selected[i] {
			mark = StyleCheck.Render("✓")
			name = StyleSelected.Render(m.items[i])
		} else {
			mark = StyleCircle.Render("○")
			name = m.items[i]
		}
		if i == m.cursor {
			fmt.Fprintf(&sb, "  %s %s %s\n",
				StyleCursor.Render("▶"), mark, name)
		} else {
			fmt.Fprintf(&sb, "    %s %s\n", mark, name)
		}
	}

	if len(m.items) > m.maxVis {
		fmt.Fprintf(&sb, "\n  %s\n",
			StyleDim.Render(fmt.Sprintf("(%d/%d shown)  %d selected",
				end-m.offset, len(m.items), selCount)))
	}

	return sb.String()
}

func (m *multiselectModel) chosen() []string {
	var out []string
	for i, item := range m.items {
		if m.selected[i] {
			out = append(out, item)
		}
	}
	return out
}

// RunMultiselect shows a multiselect menu. Returns selected items and ok=false if cancelled.
func RunMultiselect(title string, items []string, preSelected []bool) ([]string, bool) {
	m := newMultiselectModel(title, items, preSelected)
	p := tea.NewProgram(&m, ttyOpts()...)
	final, err := runProgram(p)
	if err != nil {
		return nil, false
	}
	res, ok := final.(*multiselectModel)
	if !ok {
		return nil, false
	}
	if res.quit {
		return nil, false
	}
	return res.chosen(), true
}
