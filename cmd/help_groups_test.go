package cmd

import (
	"bytes"
	"strings"
	"testing"
)

// Every visible top-level command must sit in a help group, or it falls
// into cobra's "Additional Commands" and drifts out of the curated help.
func TestHelpGroupsCoverEveryCommand(t *testing.T) {
	groupCommands(rootCmd)
	listed := map[string]bool{}
	for _, g := range helpGroups {
		for _, name := range g.commands {
			if listed[name] {
				t.Errorf("%q is in more than one help group", name)
			}
			listed[name] = true
		}
	}
	for _, c := range rootCmd.Commands() {
		if c.Hidden || c.Name() == "help" || c.Name() == "completion" {
			continue
		}
		if !listed[c.Name()] {
			t.Errorf("command %q has no help group; add it to helpGroups", c.Name())
		}
	}
}

func TestRootHelpLeadsWithMCP(t *testing.T) {
	groupCommands(rootCmd)
	var out bytes.Buffer
	rootCmd.SetOut(&out)
	defer rootCmd.SetOut(nil)
	if err := rootCmd.Help(); err != nil {
		t.Fatal(err)
	}
	help := out.String()
	start := strings.Index(help, "Get started:")
	if start < 0 || !strings.Contains(help[start:], "mcp") {
		t.Fatalf("help should open with a Get started group containing mcp:\n%s", help)
	}
	if strings.Contains(help, "Additional Commands") {
		t.Errorf("an ungrouped command leaked into Additional Commands:\n%s", help)
	}
	if strings.Count(help, "  install ") != 1 {
		t.Errorf("install should be listed once:\n%s", help)
	}
}
