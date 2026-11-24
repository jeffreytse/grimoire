package cmd

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/mark3labs/mcp-go/mcp"

	"github.com/jeffreytse/grimoire/internal/config"
	"github.com/jeffreytse/grimoire/internal/manifest"
	"github.com/jeffreytse/grimoire/internal/profiles"
	"github.com/jeffreytse/grimoire/internal/search"
	"github.com/jeffreytse/grimoire/internal/skills"
	"github.com/jeffreytse/grimoire/internal/tui"
)

// skillScope is what a project's grimoire.toml declares: its [dependencies]
// skills, its profiles' skills and its pinned practices. By default MCP mode
// prefers them — they rank first and are labeled — while the whole library
// stays available; with [standards] scope = "strict" it serves only them.
// grimoire's meta workflows always load by name, since routing needs them.
type skillScope struct {
	declared map[string]bool // skill names; nil = the project declares nothing
	strict   bool            // [standards] scope = "strict": serve only declared skills
	missing  []string        // dependency packages the project needs but are not on disk
}

// hasDeclaration reports whether the project declares any practices.
func (s *skillScope) hasDeclaration() bool { return s != nil && s.declared != nil }

// restricted reports whether the scope hides undeclared skills (strict).
func (s *skillScope) restricted() bool { return s.hasDeclaration() && s.strict }

// isDeclared reports whether a skill is one of the project's declared practices.
func (s *skillScope) isDeclared(name string) bool { return s.hasDeclaration() && s.declared[name] }

// filter returns the skills to search, list and offer: all of them, or in
// strict mode only the declared ones (meta workflows load by name instead;
// in a small strict scope they would only pad the results).
func (s *skillScope) filter(all []skills.Skill) []skills.Skill {
	if !s.restricted() {
		return all
	}
	out := make([]skills.Skill, 0, len(s.declared))
	for i := range all {
		if s.declared[all[i].Name] {
			out = append(out, all[i])
		}
	}
	return out
}

// declaredBoost multiplies a declared skill's search score in prefer mode:
// enough to lead among comparably relevant hits, not to outrank a much
// better match.
const declaredBoost = 1.5

// prefer re-ranks hits so declared skills lead among comparable matches.
func (s *skillScope) prefer(hits []search.Result) {
	if !s.hasDeclaration() || s.strict {
		return
	}
	for i := range hits {
		if s.declared[hits[i].Skill.Name] {
			hits[i].Score *= declaredBoost
		}
	}
	sort.SliceStable(hits, func(a, b int) bool { return hits[a].Score > hits[b].Score })
}

// size is the number of declared practices.
func (s *skillScope) size() int {
	if !s.hasDeclaration() {
		return 0
	}
	return len(s.declared)
}

// strictCopy returns the scope as strict, for a call that asks for the
// declared practices only (scope "project").
func (s *skillScope) strictCopy() *skillScope {
	if !s.hasDeclaration() {
		return s
	}
	c := *s
	c.strict = true
	return &c
}

var (
	scopeMu    sync.Mutex
	scopeCache = map[string]cachedScope{}
)

type cachedScope struct {
	stamp string
	scope *skillScope
}

// projectScope returns the skill scope of the project at dir, cached until
// its grimoire.toml changes. A project without [dependencies] or profiles
// has no restriction — MCP mode then serves the whole library.
func projectScope(dir string) *skillScope {
	path := manifest.ProjectPath(dir)
	stamp := ""
	if fi, err := os.Stat(path); err == nil {
		stamp = fi.ModTime().UTC().Format(time.RFC3339Nano)
	}
	scopeMu.Lock()
	defer scopeMu.Unlock()
	if c, ok := scopeCache[dir]; ok && c.stamp == stamp {
		return c.scope
	}
	s := resolveProjectScope(dir)
	scopeCache[dir] = cachedScope{stamp, s}
	return s
}

// resolveProjectScope computes what a project declares. Dependencies and
// the strict setting come from the project file only, as in a project
// install (runInstallFromManifest); profiles and pinned practices from the
// merged config, as in resolveInstallPool and preferencesContext.
func resolveProjectScope(dir string) *skillScope {
	s := &skillScope{}
	add := func(name string) {
		if s.declared == nil {
			s.declared = map[string]bool{}
		}
		s.declared[name] = true
	}
	if mf, err := manifest.ParseFile(manifest.ProjectPath(dir)); err == nil {
		s.strict = strings.EqualFold(strings.TrimSpace(mf.Standards.Scope), "strict")
	}

	if keys := projectDepKeys(dir); len(keys) > 0 {
		regs := skills.AllPackages()
		refs := make([]config.PackageRef, 0, len(keys))
		for _, k := range keys {
			ref := config.ParsePackageRef(k)
			refs = append(refs, ref)
			if ref.IsLocal() || ref.IsOfficialRepoPath() || ref.PackageName == "" || hasPackage(regs, ref.PackageName) {
				continue
			}
			// A project-only package: serve it from its clone, or report it missing.
			home := skills.PackageHome(ref.PackageName)
			if dirExists(home) {
				regs = append(regs, skills.PackageEntry{Name: ref.PackageName, Home: home})
			} else {
				s.missing = append(s.missing, ref.PackageName)
			}
		}
		pool, _, _ := skills.SkillsMatchingAnyMeta(refs, regs)
		for i := range pool {
			add(pool[i].Name)
		}
		if s.declared == nil {
			s.declared = map[string]bool{} // declared, though nothing resolved yet
		}
	}

	if r, err := config.Load(dir); err == nil {
		if len(r.Core.Profiles) > 0 {
			if eff, err := profiles.ResolveEffectiveSkills(r.Core.Profiles, dir, nil, nil); err == nil {
				for _, sk := range eff {
					add(sk.Name)
				}
			}
		}
		for _, key := range r.SectionKeys() {
			for _, p := range r.ResolveSection(key).Practices {
				add(p)
			}
		}
	}
	return s
}

func hasPackage(regs []skills.PackageEntry, name string) bool {
	for i := range regs {
		if regs[i].Name == name {
			return true
		}
	}
	return false
}

// requestScope is the scope a tool call asks for: the project's (default),
// the declared practices only ("project"), or the whole library with no
// preference ("all").
func requestScope(request mcp.CallToolRequest) *skillScope { //nolint:gocritic
	switch strings.ToLower(strings.TrimSpace(request.GetString("scope", ""))) {
	case "all":
		return nil
	case "project":
		return projectScope(getProjectDir()).strictCopy()
	}
	return projectScope(getProjectDir())
}

// scopeNote tells the model how the project's declaration shaped a result,
// and how to change that; "" when the project declares nothing.
func scopeNote(s *skillScope) string {
	if !s.hasDeclaration() {
		return ""
	}
	var n string
	if s.strict {
		n = fmt.Sprintf("Limited to this project's %d declared practices (grimoire.toml scope \"strict\"; grimoire's workflows still load by name). Pass scope \"all\" to search the whole library.", s.size())
	} else {
		n = fmt.Sprintf("Ranked by this project's grimoire.toml: its %d declared practices come first (project_standard). The rest of the library is available; when you recommend a skill that is not a project standard, say so. Pass scope \"project\" for declared practices only.", s.size())
	}
	if len(s.missing) > 0 {
		n += " Missing packages (run grimoire update): " + strings.Join(s.missing, ", ") + "."
	}
	return n
}

// withScopeNote appends the scope note to a tool result as a second text
// item, leaving the JSON payload unchanged.
func withScopeNote(res *mcp.CallToolResult, err error) func(*skillScope) (*mcp.CallToolResult, error) {
	return func(s *skillScope) (*mcp.CallToolResult, error) {
		if err != nil || res == nil {
			return res, err
		}
		if n := scopeNote(s); n != "" {
			res.Content = append(res.Content, mcp.NewTextContent(n))
		}
		return res, nil
	}
}

// outOfScopeNote marks a skill loaded by name that the project does not
// declare. It still loads — it was asked for by name.
func outOfScopeNote(name, path string) string {
	s := projectScope(getProjectDir())
	if !s.hasDeclaration() || s.declared[name] || strings.Contains(filepath.ToSlash(path), "/meta/") {
		return ""
	}
	if s.strict {
		return name + " is not in this project's grimoire.toml; it was loaded because it was asked for by name."
	}
	return name + " is not one of this project's declared practices."
}

// scopeStatus is the session status line about the project's declaration:
// missing dependency packages are a problem, shown at every visible level;
// the declaration itself is detail, shown at full.
func scopeStatus(l displayLevel) string {
	s := projectScope(getProjectDir())
	switch {
	case len(s.missing) > 0:
		return statusLine(l, "this project needs "+strings.Join(s.missing, ", ")+" — run `grimoire update`")
	case s.hasDeclaration() && l == displayFull:
		mode := "preferred"
		if s.strict {
			mode = "strict"
		}
		return statusLine(l, fmt.Sprintf("project standards: %d declared practices (%s)", s.size(), mode))
	}
	return ""
}

// joinLines joins the non-empty lines with newlines.
func joinLines(lines ...string) string {
	var out []string
	for _, l := range lines {
		if l != "" {
			out = append(out, l)
		}
	}
	return strings.Join(out, "\n")
}

// cloneMissingPackages clones each dependency's package that is not on disk
// yet. Local and official refs need no clone. A failed clone is reported and
// skipped.
func cloneMissingPackages(keys []string, out, errOut io.Writer) {
	for _, key := range keys {
		ref := config.ParsePackageRef(key)
		if ref.IsLocal() || ref.IsOfficialRepoPath() || ref.PackageURL == "" {
			continue
		}
		if dirExists(skills.PackageHome(ref.PackageName)) {
			continue
		}
		fmt.Fprintf(out, "  Cloning %s…\n", ref.PackageName)
		buf := &bytes.Buffer{}
		if err := updateNamedPackage(ref.PackageName, ref.PackageURL, ref.Tag, buf); err != nil {
			fmt.Fprintf(errOut, "%s  dependency %q: %v\n", tui.IconWarn, key, err)
			continue
		}
		_, _ = out.Write(buf.Bytes())
	}
}

// projectDepKeys returns the dependency keys of the project's grimoire.toml.
func projectDepKeys(dir string) []string {
	mf, err := manifest.ParseFile(manifest.ProjectPath(dir))
	if err != nil {
		return nil
	}
	keys := make([]string, 0, len(mf.Deps)+len(mf.DevDeps))
	for k := range mf.Deps {
		keys = append(keys, k)
	}
	for k := range mf.DevDeps {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// fetchProjectPackages clones the project's missing dependency packages, so
// MCP mode can serve the skills its grimoire.toml declares. grimoire update
// and grimoire mcp install call it; the MCP server itself never fetches.
func fetchProjectPackages() {
	cloneMissingPackages(projectDepKeys(getProjectDir()), os.Stdout, os.Stderr)
}
