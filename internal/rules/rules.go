package rules

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jeffreytse/grimoire/internal/agent"
	"github.com/jeffreytse/grimoire/internal/compliance"
	"github.com/jeffreytse/grimoire/internal/settings"
	"github.com/jeffreytse/grimoire/internal/skills"
)

const ruleSource = "grimoire-rules"

func diag(uri, code, message string, severity int) compliance.Diagnostic {
	return compliance.Diagnostic{
		URI:      uri,
		Severity: severity,
		Code:     code,
		Source:   ruleSource,
		Message:  message,
		Status:   "fail",
	}
}

// checkSkillHasSkillMd reports Error for any skill directory missing SKILL.md.
func checkSkillHasSkillMd(sources []skills.SkillsSource) []compliance.Diagnostic {
	var out []compliance.Diagnostic
	for _, src := range sources {
		domains, err := skills.ListDomains(src.Root)
		if err != nil {
			continue
		}
		for _, d := range domains {
			dirs := skillDirs(d)
			for _, dir := range dirs {
				if _, err := os.Stat(filepath.Join(dir, "SKILL.md")); err != nil {
					out = append(out, diag(
						"file://"+dir,
						"skill-has-skill-md",
						fmt.Sprintf("skill directory has no SKILL.md: %s", filepath.Base(dir)),
						1,
					))
				}
			}
		}
	}
	return out
}

// checkSkillMdFrontmatter reports Warning when a SKILL.md lacks name or tags in frontmatter.
func checkSkillMdFrontmatter(sources []skills.SkillsSource) []compliance.Diagnostic {
	var out []compliance.Diagnostic
	for _, src := range sources {
		allSkills, err := skills.ListAllSkills(src.Root)
		if err != nil {
			continue
		}
		for _, sk := range allSkills {
			skillMd := filepath.Join(sk.Path, "SKILL.md")
			data, err := os.ReadFile(skillMd)
			if err != nil {
				continue
			}
			content := string(data)
			uri := "file://" + skillMd

			hasFrontmatter, fmName, fmTags := parseFrontmatterFields(content)
			if !hasFrontmatter {
				out = append(out, diag(uri, "skill-md-has-frontmatter",
					fmt.Sprintf("%s: missing YAML frontmatter block", filepath.Base(sk.Path)), 2))
				continue
			}
			if fmName == "" {
				out = append(out, diag(uri, "skill-md-has-name",
					fmt.Sprintf("%s: frontmatter missing name: field", filepath.Base(sk.Path)), 2))
			}
			if len(fmTags) == 0 {
				out = append(out, diag(uri, "skill-md-has-tags",
					fmt.Sprintf("%s: frontmatter missing tags: field", filepath.Base(sk.Path)), 2))
			}
		}
	}
	return out
}

// parseFrontmatterFields extracts name and tags directly from SKILL.md content.
// Returns (hasFrontmatter, name, tags).
func parseFrontmatterFields(content string) (bool, string, []string) {
	if !strings.HasPrefix(content, "---") {
		return false, "", nil
	}
	rest := content[3:]
	end := strings.Index(rest, "\n---")
	if end == -1 {
		return false, "", nil
	}
	fm := rest[:end]

	var name string
	var tags []string
	inTagsBlock := false

	for _, line := range strings.Split(fm, "\n") {
		trimmed := strings.TrimRight(line, " \t\r")
		if inTagsBlock {
			if strings.HasPrefix(trimmed, "  - ") || strings.HasPrefix(trimmed, "- ") {
				tag := strings.Trim(strings.TrimLeft(trimmed, " -"), `"'`)
				if tag != "" {
					tags = append(tags, tag)
				}
				continue
			}
			if !strings.HasPrefix(trimmed, " ") && !strings.HasPrefix(trimmed, "\t") {
				inTagsBlock = false
			}
		}
		if after, ok := strings.CutPrefix(trimmed, "name:"); ok {
			name = strings.Trim(after, ` "'`)
			continue
		}
		if after, ok := strings.CutPrefix(trimmed, "tags:"); ok {
			after = strings.TrimSpace(after)
			if strings.HasPrefix(after, "[") {
				inner := strings.Trim(after, "[]")
				for _, t := range strings.Split(inner, ",") {
					if t = strings.Trim(t, ` "'`); t != "" {
						tags = append(tags, t)
					}
				}
			} else if after == "" {
				inTagsBlock = true
			}
		}
	}
	return true, name, tags
}

// checkAgentBrokenSymlinks reports Error for broken symlinks in agent skill dirs.
func checkAgentBrokenSymlinks() []compliance.Diagnostic {
	var out []compliance.Diagnostic
	for _, ag := range agent.All {
		dir := agent.SkillsDir(ag)
		if dir == "" {
			continue
		}
		entries, err := os.ReadDir(dir)
		if err != nil {
			continue
		}
		for _, e := range entries {
			if e.Name() == "" || e.Name()[0] == '.' {
				continue
			}
			full := filepath.Join(dir, e.Name())
			if isBrokenSymlink(full) {
				out = append(out, diag(
					"file://"+full,
					"agent-no-broken-symlinks",
					fmt.Sprintf("%s: broken skill symlink in %s skills dir", e.Name(), ag),
					1,
				))
			}
		}
	}
	return out
}

func isBrokenSymlink(path string) bool {
	if _, err := os.Lstat(path); err != nil {
		return false
	}
	_, err := os.Stat(path)
	return err != nil
}

// checkSettingsParseable reports Error if .grimoire/settings.toml exists but fails to parse.
func checkSettingsParseable(projectDir string) []compliance.Diagnostic {
	path := filepath.Join(projectDir, ".grimoire", "settings.toml")
	if _, err := os.Stat(path); err != nil {
		return nil // absent is fine
	}
	if _, err := settings.ParseFile(path); err != nil {
		return []compliance.Diagnostic{diag(
			"file://"+path,
			"settings-toml-parseable",
			fmt.Sprintf("settings.toml failed to parse: %v", err),
			1,
		)}
	}
	return nil
}

// skillDirs returns all skill entry directories (not the skills/ parent) for a domain.
func skillDirs(d skills.Domain) []string {
	var dirs []string
	if d.Nested {
		subs, err := skills.ListSubdomains(d.Path)
		if err != nil {
			return nil
		}
		for _, s := range subs {
			dirs = append(dirs, leafSkillDirs(s.Path)...)
		}
	} else {
		dirs = leafSkillDirs(d.Path)
	}
	return dirs
}

func leafSkillDirs(parentDir string) []string {
	skillsDir := filepath.Join(parentDir, "skills")
	entries, err := os.ReadDir(skillsDir)
	if err != nil {
		return nil
	}
	var dirs []string
	for _, e := range entries {
		if !e.IsDir() || strings.HasPrefix(e.Name(), ".") {
			continue
		}
		dirs = append(dirs, filepath.Join(skillsDir, e.Name()))
	}
	return dirs
}
