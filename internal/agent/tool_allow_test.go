package agent

import (
	"os"
	"path/filepath"
	"testing"
)

func TestToolAllow_IdempotentAndRefusesMalformed(t *testing.T) {
	path := filepath.Join(t.TempDir(), "settings.json")
	rules := ToolAllowRules("claude")
	if len(rules) != len(mcpReadOnlyTools) {
		t.Fatalf("rules = %v", rules)
	}
	if changed, err := ConfigureToolAllow(path, rules); err != nil || !changed {
		t.Fatalf("configure: changed=%v err=%v", changed, err)
	}
	if changed, err := ConfigureToolAllow(path, rules); err != nil || changed {
		t.Errorf("re-configure: changed=%v err=%v; want no change", changed, err)
	}
	if changed, err := RemoveToolAllow(path, rules); err != nil || !changed {
		t.Fatalf("remove: changed=%v err=%v", changed, err)
	}

	for _, bad := range []string{`{"permissions":[]}`, `{"permissions":{"allow":"x"}}`} {
		if err := os.WriteFile(path, []byte(bad), 0o644); err != nil {
			t.Fatal(err)
		}
		if _, err := ConfigureToolAllow(path, rules); err == nil {
			t.Errorf("%s: want error", bad)
		}
		if b, _ := os.ReadFile(path); string(b) != bad {
			t.Errorf("%s: rewritten to %s", bad, b)
		}
	}
	if ToolAllowFile("copilot", "user", t.TempDir()) != "" || ToolAllowFile("antigravity", "project", t.TempDir()) != "" {
		t.Error("unexpected allow file for unsupported agent/scope")
	}
}
