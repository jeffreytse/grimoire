package cmd

import (
	"os"
	"path/filepath"
	"sort"
	"testing"
)

func TestIsAnalyzableFile(t *testing.T) {
	dir := t.TempDir()

	writeTmp := func(name, content string) string {
		p := filepath.Join(dir, name)
		if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
		return p
	}
	writeBinary := func(name string) string {
		p := filepath.Join(dir, name)
		if err := os.WriteFile(p, []byte("data\x00more"), 0o644); err != nil {
			t.Fatal(err)
		}
		return p
	}

	tests := []struct {
		name string
		path string
		want bool
	}{
		{"png ext", writeTmp("img.png", "data"), false},
		{"zip ext", writeTmp("arch.zip", "data"), false},
		{"exe ext", writeTmp("bin.exe", "data"), false},
		{"woff ext", writeTmp("font.woff", "data"), false},
		{"mp3 ext", writeTmp("sound.mp3", "data"), false},
		{"go ext text", writeTmp("main.go", "package main"), true},
		{"md ext text", writeTmp("README.md", "# hello"), true},
		{"yaml ext text", writeTmp("cfg.yaml", "key: val"), true},
		{"null byte in content", writeBinary("binary.go"), false},
		{"text file no known ext", writeTmp("Makefile", "all:\n\tgo build"), true},
		{"nonexistent path", filepath.Join(dir, "missing.go"), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isAnalyzableFile(tt.path); got != tt.want {
				t.Errorf("isAnalyzableFile(%q) = %v, want %v", tt.path, got, tt.want)
			}
		})
	}
}

func TestExtractExtHints(t *testing.T) {
	tests := []struct {
		name    string
		rubrics []skillRubric
		want    []string
	}{
		{
			name:    "empty rubrics",
			rubrics: nil,
			want:    nil,
		},
		{
			name:    "no ext patterns",
			rubrics: []skillRubric{{Name: "a", Body: "apply solid principles to all code"}},
			want:    nil,
		},
		{
			name:    "single pattern",
			rubrics: []skillRubric{{Name: "a", Body: "applies to *.go files"}},
			want:    []string{".go"},
		},
		{
			name: "multiple patterns",
			rubrics: []skillRubric{
				{Name: "a", Body: "check *.go and *.ts files"},
			},
			want: []string{".go", ".ts"},
		},
		{
			name: "patterns across rubrics",
			rubrics: []skillRubric{
				{Name: "a", Body: "applies to *.go"},
				{Name: "b", Body: "applies to *.py"},
			},
			want: []string{".go", ".py"},
		},
		{
			name: "duplicates deduplicated",
			rubrics: []skillRubric{
				{Name: "a", Body: "*.go files"},
				{Name: "b", Body: "also *.go"},
			},
			want: []string{".go"},
		},
		{
			name:    "uppercase pattern lowercased",
			rubrics: []skillRubric{{Name: "a", Body: "applies to *.GO files"}},
			want:    []string{".go"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractExtHints(tt.rubrics)
			if tt.want == nil {
				if got != nil {
					t.Errorf("expected nil, got %v", got)
				}
				return
			}
			sort.Strings(got)
			sort.Strings(tt.want)
			if len(got) != len(tt.want) {
				t.Errorf("got %v, want %v", got, tt.want)
				return
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("got %v, want %v", got, tt.want)
					return
				}
			}
		})
	}
}

func TestMatchesExcludePatterns(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		patterns []string
		want     bool
	}{
		// doublestar glob
		{"vendor dir glob", "vendor/foo/bar.go", []string{"vendor/**"}, true},
		{"vendor not matched outside", "cmd/vendor.go", []string{"vendor/**"}, false},
		{"double-star pb.go", "internal/foo.pb.go", []string{"**/*.pb.go"}, true},
		{"double-star pb.go root", "foo.pb.go", []string{"**/*.pb.go"}, true},
		// basename match (glob without slash)
		{"bare ext pattern", "internal/foo.pb.go", []string{"*.pb.go"}, true},
		{"bare ext no match", "internal/foo.go", []string{"*.pb.go"}, false},
		// plain name — directory prefix
		{"plain dir no slash", "testdata/input.go", []string{"testdata"}, true},
		{"plain dir with slash", "testdata/sub/input.go", []string{"testdata/"}, true},
		{"plain dir no match", "cmd/testdata.go", []string{"testdata"}, false},
		// exact file match via plain name
		{"exact file match", "go.sum", []string{"go.sum"}, true},
		// multiple patterns — any match wins
		{"multi patterns first", "vendor/foo.go", []string{"vendor/**", "**/*.pb.go"}, true},
		{"multi patterns second", "internal/x.pb.go", []string{"vendor/**", "**/*.pb.go"}, true},
		{"multi patterns none", "cmd/check.go", []string{"vendor/**", "**/*.pb.go"}, false},
		// empty patterns
		{"empty patterns", "anything.go", nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := matchesExcludePatterns(tt.path, tt.patterns); got != tt.want {
				t.Errorf("matchesExcludePatterns(%q, %v) = %v, want %v", tt.path, tt.patterns, got, tt.want)
			}
		})
	}
}
