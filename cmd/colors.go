package cmd

import (
	"fmt"
	"os"
	"strings"
)

// grimoire's color roles. Every color grimoire draws itself comes from here;
// the hues are the brand colors of assets/icon-book.svg.
//
//	brand.gold  #c9933a   the book's gilt
//	brand.ink   #1a0d06   the book's leather, used as text on gold
//
// Roles, not raw values, are used in code:
//
//	badge.bg = brand.gold   badge.fg = brand.ink   (ink on gold: 7.0:1)
//
// The badge is a filled chip so its contrast does not depend on the
// terminal's theme: gold text alone is 2.5–2.8:1 on light backgrounds,
// below WCAG's 4.5:1 for text.

type rgb struct{ r, g, b uint8 }

var (
	brandGold = rgb{0xc9, 0x93, 0x3a}
	brandInk  = rgb{0x1a, 0x0d, 0x06}

	// xterm-256 stand-ins, for terminals without 24-bit color: 173 is the
	// nearest to brand.gold (#d7875f), 16 is black (black on 173: 7.5:1).
	brandGold256 = 173
	brandInk256  = 16
)

// colorMode is how grimoire may color its own output.
type colorMode int

const (
	colorNone colorMode = iota // NO_COLOR, TERM=dumb: plain text
	color256
	colorTrue
)

// terminalColorMode follows the CLI conventions: NO_COLOR (no-color.org)
// and TERM=dumb turn color off; COLORTERM=truecolor|24bit allows 24-bit.
func terminalColorMode() colorMode {
	if os.Getenv("NO_COLOR") != "" || os.Getenv("TERM") == "dumb" {
		return colorNone
	}
	switch strings.ToLower(os.Getenv("COLORTERM")) {
	case "truecolor", "24bit":
		return colorTrue
	}
	return color256
}

// renderBadge draws label as grimoire's badge: an ink-on-gold chip, or
// [label] in plain text when color is off.
func renderBadge(label string) string {
	switch terminalColorMode() {
	case colorTrue:
		return fmt.Sprintf("\033[1;38;2;%d;%d;%d;48;2;%d;%d;%dm %s \033[0m",
			brandInk.r, brandInk.g, brandInk.b, brandGold.r, brandGold.g, brandGold.b, label)
	case color256:
		return fmt.Sprintf("\033[1;38;5;%d;48;5;%dm %s \033[0m", brandInk256, brandGold256, label)
	}
	return "[" + label + "]"
}
