package cmd

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/jeffreytse/grimoire/internal/agent"
)

// An agent in MCP mode is reported as such — never told to run
// `grimoire install`.
func TestDoctorMCPChecks(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("XDG_CONFIG_HOME", filepath.Join(home, ".config"))
	t.Setenv("COPILOT_HOME", "")
	pkg := t.TempDir()

	if _, found := doctorMCPChecks("copilot", "Copilot", "1.0", pkg); found {
		t.Fatal("copilot reported in MCP mode before install")
	}
	if _, err := agent.ConfigureMCPInstruction(agent.MCPInstructionFile("copilot", "user", getProjectDir())); err != nil {
		t.Fatal(err)
	}
	if _, err := agent.ConfigureCopilotHooks(agent.CopilotHooksFile("user", getProjectDir()), "grimoire"); err != nil {
		t.Fatal(err)
	}
	checks, found := doctorMCPChecks("copilot", "Copilot", "1.0", pkg)
	if !found {
		t.Fatal("copilot not detected in MCP mode")
	}
	var all strings.Builder
	for _, c := range checks {
		all.WriteString(c.Status + " " + c.Detail + "\n")
		if c.Status == "warn" || c.Status == "error" {
			t.Errorf("unexpected %s: %s", c.Status, c.Detail)
		}
	}
	if strings.Contains(all.String(), "grimoire install") || !strings.Contains(all.String(), "MCP mode") {
		t.Errorf("checks = %s", all.String())
	}
}
