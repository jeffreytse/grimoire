package cmd

import "testing"

func TestNormalizeScope(t *testing.T) {
	cases := []struct{ in, want string }{
		{"user", "user"}, {"global", "user"}, {"GLOBAL", "user"},
		{"project", "project"}, {"  project", "project"},
	}
	for _, c := range cases {
		got, err := normalizeScope(c.in)
		if err != nil || got != c.want {
			t.Errorf("normalizeScope(%q) = %q, %v; want %q", c.in, got, err, c.want)
		}
	}
	if _, err := normalizeScope("everywhere"); err == nil {
		t.Error("unknown scope should error")
	}
}

func TestInstallUserScope(t *testing.T) {
	orig, origG := flagInstallScope, flagInstallGlobal
	t.Cleanup(func() { flagInstallScope, flagInstallGlobal = orig, origG })
	for scope, want := range map[string]bool{"project": false, "user": true, "global": true} {
		flagInstallScope, flagInstallGlobal = scope, false
		if got := installUserScope(); got != want {
			t.Errorf("--scope %s: installUserScope() = %v, want %v", scope, got, want)
		}
	}
	flagInstallScope, flagInstallGlobal = "project", true
	if !installUserScope() {
		t.Error("--user should select the user scope")
	}
}

func TestValidateTarget(t *testing.T) {
	for _, ok := range []string{"", "auto", "all", "claude", "codex"} {
		if err := validateTarget(ok); err != nil {
			t.Errorf("validateTarget(%q): %v", ok, err)
		}
	}
	if err := validateTarget("bogus"); err == nil {
		t.Error("unknown target should error")
	}
}

// A bad --target must fail before install touches packages or disk.
func TestInstall_RejectsUnknownTargetFirst(t *testing.T) {
	orig := flagInstallTarget
	t.Cleanup(func() { flagInstallTarget = orig })
	flagInstallTarget = "bogus"
	err := runInstall(nil, nil)
	if err == nil || err == errNoSkillLibrary {
		t.Fatalf("want unknown-target error, got %v", err)
	}
}
