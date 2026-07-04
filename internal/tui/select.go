package tui

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// ── Basic select ──────────────────────────────────────────────────────────────

type selectModel struct {
	title  string
	items  []string
	cursor int
	chosen string
	quit   bool
	ctrlC  bool
}

func (m *selectModel) Init() tea.Cmd { return nil }

func (m *selectModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		switch msg.String() {
		case "ctrl+c":
			m.quit = true
			m.ctrlC = true
			return m, tea.Quit
		case "q", "esc":
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
	p := tea.NewProgram(&m, ttyOpts()...)
	final, err := runProgram(p)
	if err != nil {
		return "", false
	}
	res, ok := final.(*selectModel)
	if !ok {
		return "", false
	}
	if res.ctrlC {
		os.Exit(0) //nolint:revive // intentional: TUI exit on Ctrl+C
	}
	if res.quit || res.chosen == "" {
		return "", false
	}
	return res.chosen, true
}

// ── Profile select (viewport + sections + annotations + sentinels) ────────────

const profileViewport = 10

// Sentinels — non-profile values returned by RunProfileSelect.
const (
	ProfileSelectOther   = "\x00other"    // user wants to type a custom name
	ProfileSelectNone    = "\x00none"     // user wants no profile
	ProfileSectionPrefix = "\x00section:" // non-selectable section header
)

// Sentinels — returned by RunProfileSelect when used as a preset picker.
const (
	PresetSelectCustom = "\x00preset-custom" // user wants the custom wizard
	PresetSelectSkip   = "\x00preset-skip"   // user wants no setup
)

func isSectionItem(s string) bool {
	return strings.HasPrefix(s, ProfileSectionPrefix)
}

func isSelectableItem(s string) bool {
	return !isSectionItem(s)
}

type profileSelectModel struct {
	items       []string
	annotations map[string]string
	cursor      int
	offset      int
	chosen      string
	quit        bool
	ctrlC       bool
}

func (m *profileSelectModel) Init() tea.Cmd { return nil }

func (m *profileSelectModel) moveCursor(delta int) {
	next := m.cursor + delta
	for next >= 0 && next < len(m.items) && isSectionItem(m.items[next]) {
		next += delta
	}
	if next < 0 || next >= len(m.items) {
		return
	}
	m.cursor = next
	if m.cursor < m.offset {
		m.offset = m.cursor
	} else if m.cursor >= m.offset+profileViewport {
		m.offset = m.cursor - profileViewport + 1
	}
}

func (m *profileSelectModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		switch msg.String() {
		case "ctrl+c":
			m.quit = true
			m.ctrlC = true
			return m, tea.Quit
		case "esc":
			m.quit = true
			return m, tea.Quit
		case "up", "k":
			m.moveCursor(-1)
		case "down", "j":
			m.moveCursor(1)
		case "enter", " ":
			if isSelectableItem(m.items[m.cursor]) {
				m.chosen = m.items[m.cursor]
				return m, tea.Quit
			}
		}
	}
	return m, nil
}

func (m *profileSelectModel) View() string {
	var sb strings.Builder
	sb.WriteString("\n")
	sb.WriteString(StyleHint.Render("  ↑↓ navigate   ENTER confirm   ESC cancel"))
	sb.WriteString("\n\n")

	if m.offset > 0 {
		fmt.Fprintf(&sb, "    %s\n", StyleHint.Render(fmt.Sprintf("↑ %d more", m.offset)))
	}

	end := m.offset + profileViewport
	if end > len(m.items) {
		end = len(m.items)
	}

	for i := m.offset; i < end; i++ {
		item := m.items[i]

		// Section header — non-selectable separator
		if isSectionItem(item) {
			title := strings.TrimPrefix(item, ProfileSectionPrefix)
			fmt.Fprintf(&sb, "\n  %s\n", StyleHint.Render("── "+title+" ──"))
			continue
		}

		ann := m.annotations[item]
		var label string
		switch item {
		case ProfileSelectOther:
			label = "other (type name)…"
		case ProfileSelectNone:
			label = "(no profile)"
		case PresetSelectCustom:
			label = "custom setup…"
		case PresetSelectSkip:
			label = "skip"
		default:
			label = item
		}

		if i == m.cursor {
			line := StyleSelected.Render(label)
			if ann != "" {
				line += " " + StyleHint.Render(ann)
			}
			fmt.Fprintf(&sb, "  %s %s\n", StyleCursor.Render("▶"), line)
		} else {
			var rendered string
			switch item {
			case ProfileSelectOther, ProfileSelectNone, PresetSelectCustom, PresetSelectSkip:
				rendered = StyleHint.Render(label)
			default:
				rendered = StyleDim.Render(label)
			}
			if ann != "" {
				rendered += " " + StyleHint.Render(ann)
			}
			fmt.Fprintf(&sb, "    %s\n", rendered)
		}
	}

	remaining := len(m.items) - end
	if remaining > 0 {
		fmt.Fprintf(&sb, "    %s\n", StyleHint.Render(fmt.Sprintf("↓ %d more", remaining)))
	}

	return sb.String()
}

// ── Scrollable select (configurable viewport + section headers) ───────────────

type scrollSelectModel struct {
	title      string
	items      []string
	maxVisible int
	cursor     int
	offset     int
	chosen     string
	quit       bool
	ctrlC      bool
}

func (m *scrollSelectModel) Init() tea.Cmd { return nil }

func (m *scrollSelectModel) moveCursor(delta int) {
	next := m.cursor + delta
	for next >= 0 && next < len(m.items) && isSectionItem(m.items[next]) {
		next += delta
	}
	if next < 0 || next >= len(m.items) {
		return
	}
	m.cursor = next
	if m.cursor < m.offset {
		m.offset = m.cursor
	} else if m.cursor >= m.offset+m.maxVisible {
		m.offset = m.cursor - m.maxVisible + 1
	}
}

func (m *scrollSelectModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		switch msg.String() {
		case "ctrl+c":
			m.quit = true
			m.ctrlC = true
			return m, tea.Quit
		case "q", "esc":
			m.quit = true
			return m, tea.Quit
		case "up", "k":
			m.moveCursor(-1)
		case "down", "j":
			m.moveCursor(1)
		case "enter", " ":
			if isSelectableItem(m.items[m.cursor]) {
				m.chosen = m.items[m.cursor]
				return m, tea.Quit
			}
		}
	}
	return m, nil
}

func (m *scrollSelectModel) View() string {
	var sb strings.Builder
	sb.WriteString("\n")
	sb.WriteString(StyleTitle.Render(m.title))
	sb.WriteString("\n")
	sb.WriteString(StyleHint.Render("  ↑↓ navigate   ENTER confirm   q quit"))
	sb.WriteString("\n\n")

	if m.offset > 0 {
		fmt.Fprintf(&sb, "    %s\n", StyleHint.Render(fmt.Sprintf("↑ %d more", m.offset)))
	}

	end := m.offset + m.maxVisible
	if end > len(m.items) {
		end = len(m.items)
	}

	for i := m.offset; i < end; i++ {
		item := m.items[i]
		if isSectionItem(item) {
			title := strings.TrimPrefix(item, ProfileSectionPrefix)
			fmt.Fprintf(&sb, "\n  %s\n", StyleHint.Render("── "+title+" ──"))
			continue
		}
		if i == m.cursor {
			fmt.Fprintf(&sb, "  %s %s\n", StyleCursor.Render("▶"), StyleSelected.Render(item))
		} else {
			fmt.Fprintf(&sb, "    %s\n", StyleDim.Render(item))
		}
	}

	remaining := len(m.items) - end
	if remaining > 0 {
		fmt.Fprintf(&sb, "    %s\n", StyleHint.Render(fmt.Sprintf("↓ %d more", remaining)))
	}

	return sb.String()
}

// RunSelectScrollable shows a single-select menu with a configurable viewport.
// items may include ProfileSectionPrefix ("\x00section:") headers — non-selectable,
// rendered as "── title ──". maxVisible controls visible rows before scrolling.
func RunSelectScrollable(title string, items []string, maxVisible int) (string, bool) {
	cursor := 0
	for cursor < len(items) && isSectionItem(items[cursor]) {
		cursor++
	}
	m := &scrollSelectModel{
		title:      title,
		items:      items,
		maxVisible: maxVisible,
		cursor:     cursor,
	}
	p := tea.NewProgram(m, ttyOpts()...)
	final, err := runProgram(p)
	if err != nil {
		return "", false
	}
	res, ok := final.(*scrollSelectModel)
	if !ok {
		return "", false
	}
	if res.ctrlC {
		os.Exit(0) //nolint:revive // intentional: TUI exit on Ctrl+C
	}
	if res.quit || res.chosen == "" {
		return "", false
	}
	return res.chosen, true
}

// RunProfileSelect shows a scrollable (10-item viewport) profile picker.
// items may include ProfileSectionPrefix headers (non-selectable), ProfileSelectNone,
// and ProfileSelectOther sentinels — caller controls the full list.
// defaultItem sets the initial cursor position (section headers are skipped).
func RunProfileSelect(items []string, annotations map[string]string, defaultItem string) (string, bool) {
	// Find starting cursor — skip section headers
	cursor := 0
	for i, it := range items {
		if it == defaultItem {
			cursor = i
			break
		}
	}
	// If cursor landed on a section header, advance to next selectable
	for cursor < len(items) && isSectionItem(items[cursor]) {
		cursor++
	}

	offset := 0
	if cursor >= profileViewport {
		offset = cursor - profileViewport + 1
	}

	m := &profileSelectModel{
		items:       items,
		annotations: annotations,
		cursor:      cursor,
		offset:      offset,
	}
	p := tea.NewProgram(m, ttyOpts()...)
	final, err := runProgram(p)
	if err != nil {
		return "", false
	}
	res, ok := final.(*profileSelectModel)
	if !ok {
		return "", false
	}
	if res.ctrlC {
		os.Exit(0) //nolint:revive // intentional: TUI exit on Ctrl+C
	}
	if res.quit || res.chosen == "" {
		return "", false
	}
	return res.chosen, true
}
