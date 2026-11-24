package agent

import (
	"os"
	"path/filepath"
	"testing"
)

func TestStatusLine_SetChainRestore(t *testing.T) {
	path := filepath.Join(t.TempDir(), "settings.json")

	// Empty: badge is set; second run is a no-op.
	if changed, other, err := ConfigureStatusLine(path, "/usr/local/bin/grimoire", false); err != nil || !changed || other != "" {
		t.Fatalf("set: changed=%v other=%q err=%v", changed, other, err)
	}
	if changed, _, _ := ConfigureStatusLine(path, "/usr/local/bin/grimoire", false); changed {
		t.Error("second set changed the file")
	}
	if st, _ := StatusLineState(path); st != StatusLineOurs {
		t.Errorf("state = %q", st)
	}
	if changed, err := RemoveStatusLine(path); err != nil || !changed {
		t.Fatalf("remove: %v %v", changed, err)
	}
	if st, _ := StatusLineState(path); st != StatusLineNone {
		t.Errorf("after remove state = %q", st)
	}

	// Someone else's: left alone without chain, chained with it, restored on remove.
	orig := `bash "/home/u/.claude/hooks/it's-mine.sh"`
	if err := os.WriteFile(path, []byte(`{"statusLine":{"type":"command","command":`+jsonString(orig)+`,"padding":1}}`), 0o644); err != nil {
		t.Fatal(err)
	}
	if changed, other, _ := ConfigureStatusLine(path, "grimoire", false); changed || other != orig {
		t.Errorf("no-chain: changed=%v other=%q", changed, other)
	}
	if changed, _, err := ConfigureStatusLine(path, "grimoire", true); err != nil || !changed {
		t.Fatalf("chain: %v %v", changed, err)
	}
	if st, then := StatusLineState(path); st != StatusLineChained || then != orig {
		t.Errorf("chained state = %q then = %q", st, then)
	}
	if changed, _, _ := ConfigureStatusLine(path, "grimoire", true); changed {
		t.Error("re-chaining changed the file (would nest)")
	}
	if changed, err := RemoveStatusLine(path); err != nil || !changed {
		t.Fatalf("remove chained: %v %v", changed, err)
	}
	cfg, _ := readSettings(path)
	sl, ok := cfg["statusLine"].(map[string]any)
	if !ok {
		t.Fatalf("statusLine = %v, want object", cfg["statusLine"])
	}
	if sl["command"] != orig || sl["padding"] != float64(1) {
		t.Errorf("restored statusLine = %v", sl)
	}
}

func jsonString(s string) string {
	b := []byte{'"'}
	for _, r := range s {
		if r == '"' || r == '\\' {
			b = append(b, '\\')
		}
		b = append(b, string(r)...)
	}
	return string(append(b, '"'))
}
