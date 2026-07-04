package config

import (
	"testing"
)

func TestParsePackageRef(t *testing.T) {
	cases := []struct {
		input       string
		local       bool
		localPath   string
		host        string
		owner       string
		repo        string
		tag         string
		path        string
		packageName string
		official    bool
	}{
		// official repo path shorthand (single segment, no '/')
		{
			input:    "engineering",
			path:     "engineering",
			official: true,
		},
		// local Unix absolute path — no glob
		{
			input:     "/local/path/to/skills",
			local:     true,
			localPath: "/local/path/to/skills",
		},
		// local Unix path + glob
		{
			input:     "/local/path/to/skills:engineering/bdd",
			local:     true,
			localPath: "/local/path/to/skills",
			path:      "engineering/bdd",
		},
		// relative path
		{
			input:     "./myskills:engineering",
			local:     true,
			localPath: "./myskills",
			path:      "engineering",
		},
		// owner/repo only (all skills)
		{
			input:       "jeffreytse/grimoire-core",
			host:        "github.com",
			owner:       "jeffreytse",
			repo:        "grimoire-core",
			packageName: "github.com/jeffreytse/grimoire-core@latest",
		},
		// owner/repo + glob path
		{
			input:       "jeffreytse/grimoire-core:engineering/bdd",
			host:        "github.com",
			owner:       "jeffreytse",
			repo:        "grimoire-core",
			path:        "engineering/bdd",
			packageName: "github.com/jeffreytse/grimoire-core@latest",
		},
		// owner/repo + version + glob
		{
			input:       "acmecorp/standards@v0.1.0:engineering/tdd",
			host:        "github.com",
			owner:       "acmecorp",
			repo:        "standards",
			tag:         "v0.1.0",
			path:        "engineering/tdd",
			packageName: "github.com/acmecorp/standards@v0.1.0",
		},
		// explicit non-github host
		{
			input:       "gitlab.com/google/standards",
			host:        "gitlab.com",
			owner:       "google",
			repo:        "standards",
			packageName: "gitlab.com/google/standards@latest",
		},
		// explicit host + version + glob
		{
			input:       "github.com/acmecorp/new-standards@v0.2.1:engineering/**/apply-*",
			host:        "github.com",
			owner:       "acmecorp",
			repo:        "new-standards",
			tag:         "v0.2.1",
			path:        "engineering/**/apply-*",
			packageName: "github.com/acmecorp/new-standards@v0.2.1",
		},
	}

	for _, c := range cases {
		t.Run(c.input, func(t *testing.T) {
			got := ParsePackageRef(c.input)
			if got.Local != c.local {
				t.Errorf("Local: got %v, want %v", got.Local, c.local)
			}
			if got.LocalPath != c.localPath {
				t.Errorf("LocalPath: got %q, want %q", got.LocalPath, c.localPath)
			}
			if got.Host != c.host {
				t.Errorf("Host: got %q, want %q", got.Host, c.host)
			}
			if got.Owner != c.owner {
				t.Errorf("Owner: got %q, want %q", got.Owner, c.owner)
			}
			if got.Repo != c.repo {
				t.Errorf("Repo: got %q, want %q", got.Repo, c.repo)
			}
			if got.Tag != c.tag {
				t.Errorf("Tag: got %q, want %q", got.Tag, c.tag)
			}
			if got.Path != c.path {
				t.Errorf("Path: got %q, want %q", got.Path, c.path)
			}
			if got.PackageName != c.packageName {
				t.Errorf("PackageName: got %q, want %q", got.PackageName, c.packageName)
			}
			if c.official && !got.IsOfficialRepoPath() {
				t.Errorf("expected IsOfficialRepoPath()=true")
			}
		})
	}
}
