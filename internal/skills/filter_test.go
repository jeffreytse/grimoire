package skills

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGlobMatch(t *testing.T) {
	cases := []struct {
		pattern string
		path    string
		want    bool
	}{
		// empty pattern matches all
		{"", "anything", true},
		// exact match
		{"engineering", "engineering", true},
		{"engineering", "security", false},
		// .gitignore style: bare name (no '/') matches at any depth
		{"engineering", "abc/engineering", true},
		// ** prefix = anywhere (explicit)
		{"**/engineering", "abc/engineering", true},
		{"**/engineering", "engineering", true},
		// nested path
		{"engineering/tdd", "engineering/tdd", true},
		{"engineering/tdd", "engineering/bdd", false},
		// ** in middle
		{"engineering/**/apply-*", "engineering/tdd/apply-tdd", true},
		{"engineering/**/apply-*", "engineering/apply-solid", true},
		{"engineering/**/apply-*", "engineering/tdd/bdd/apply-solid", true},
		// alternation
		{"{engineering,security}", "engineering", true},
		{"{engineering,security}", "security", true},
		{"{engineering,security}", "devops", false},
		// character class
		{"[es]ngineering", "engineering", true},
		{"[es]ngineering", "sngineering", true},
		{"[es]ngineering", "bngineering", false},
		// negated class
		{"[!es]ngineering", "bngineering", true},
		{"[!es]ngineering", "engineering", false},
		// POSIX class
		{"[[:alpha:]]ngineering", "engineering", true},
		{"[[:digit:]]bc", "1bc", true},
		{"[[:digit:]]bc", "abc", false},
		{"[[:alnum:]]ngineering", "engineering", true},
		{"[[:alnum:]]ngineering", "1ngineering", true},
		// wildcard *
		{"engineering/*", "engineering/tdd", true},
		{"engineering/*", "engineering/tdd/sub", false},
	}

	for _, c := range cases {
		got := GlobMatch(c.pattern, c.path)
		if got != c.want {
			t.Errorf("GlobMatch(%q, %q) = %v, want %v", c.pattern, c.path, got, c.want)
		}
	}
}

func TestWalkSkills(t *testing.T) {
	// Build a temp dir tree:
	// root/
	//   engineering/tdd/apply-tdd/SKILL.md
	//   engineering/bdd/apply-bdd/SKILL.md
	//   security/owasp/SKILL.md
	root := t.TempDir()
	must := func(err error) {
		t.Helper()
		if err != nil {
			t.Fatal(err)
		}
	}
	mkSkill := func(relDir string) {
		dir := filepath.Join(root, filepath.FromSlash(relDir))
		must(os.MkdirAll(dir, 0o755))
		must(os.WriteFile(filepath.Join(dir, "SKILL.md"), []byte("---\nname: "+filepath.Base(dir)+"\n---\n"), 0o644))
	}
	mkSkill("engineering/tdd/apply-tdd")
	mkSkill("engineering/bdd/apply-bdd")
	mkSkill("security/owasp")

	skills, err := WalkSkills(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(skills) != 3 {
		t.Errorf("expected 3 skills, got %d", len(skills))
	}
}

func TestSkillsMatchingGlob(t *testing.T) {
	root := t.TempDir()
	must := func(err error) {
		t.Helper()
		if err != nil {
			t.Fatal(err)
		}
	}
	mkSkill := func(relDir string) {
		dir := filepath.Join(root, filepath.FromSlash(relDir))
		must(os.MkdirAll(dir, 0o755))
		must(os.WriteFile(filepath.Join(dir, "SKILL.md"), []byte("---\nname: "+filepath.Base(dir)+"\n---\n"), 0o644))
	}
	mkSkill("engineering/tdd/apply-tdd")
	mkSkill("engineering/bdd/apply-bdd")
	mkSkill("security/owasp")

	t.Run("empty glob returns all", func(t *testing.T) {
		got, err := SkillsMatchingGlob(root, "")
		if err != nil {
			t.Fatal(err)
		}
		if len(got) != 3 {
			t.Errorf("expected 3, got %d", len(got))
		}
	})

	t.Run("top-level dir", func(t *testing.T) {
		got, err := SkillsMatchingGlob(root, "engineering/**")
		if err != nil {
			t.Fatal(err)
		}
		if len(got) != 2 {
			t.Errorf("expected 2, got %d", len(got))
		}
	})

	t.Run("exact path", func(t *testing.T) {
		got, err := SkillsMatchingGlob(root, "engineering/tdd/apply-tdd")
		if err != nil {
			t.Fatal(err)
		}
		if len(got) != 1 {
			t.Errorf("expected 1, got %d", len(got))
		}
	})

	t.Run("no match", func(t *testing.T) {
		got, err := SkillsMatchingGlob(root, "devops/**")
		if err != nil {
			t.Fatal(err)
		}
		if len(got) != 0 {
			t.Errorf("expected 0, got %d", len(got))
		}
	})
}
