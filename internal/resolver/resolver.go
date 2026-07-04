// Package resolver resolves grimoire dependency graphs and produces lockfile entries.
// It performs flat (Cargo-style) resolution: one version per skill name across all packages.
// Skill identity is the final path component of the dep key (the skill name).
// Version conflicts are hard errors, not silent duplicates.
package resolver

import (
	"fmt"
	"path"
	"sort"
	"strings"

	"github.com/Masterminds/semver/v3"

	"github.com/jeffreytse/grimoire/internal/lock"
	"github.com/jeffreytse/grimoire/internal/manifest"
)

// SkillMeta holds metadata for one skill fetched from its package.
// The caller (install command) populates this by cloning packages and reading SKILL.toml/SKILL.md.
type SkillMeta struct {
	Name         string            // skill name, e.g. "apply-solid-principles"
	Version      string            // semver from SKILL.toml/SKILL.md; "" = no version info
	Source       string            // dep key prefix, e.g. "acmecorp/practices"
	Resolved     string            // full git URL, e.g. "https://github.com/acmecorp/practices.git"
	Commit       string            // git commit SHA
	Checksum     string            // sha256:<hex>
	Dependencies map[string]string // transitive deps: skill name → semver constraint
}

// ConstraintSource links a version constraint to the dep that required it.
type ConstraintSource struct {
	Constraint string // e.g. "^1.0.0"
	RequiredBy string // dep key that imposed this constraint
}

// Conflict describes irreconcilable version constraints for one skill.
type Conflict struct {
	Skill       string
	Version     string             // actual resolved version (from SkillMeta)
	Constraints []ConstraintSource // all constraints that couldn't be satisfied
}

// ErrConflict is returned when one or more skills have unsatisfiable constraints.
type ErrConflict struct {
	Conflicts []Conflict
}

func (e *ErrConflict) Error() string {
	var sb strings.Builder
	sb.WriteString("dependency conflicts:\n")
	for _, c := range e.Conflicts {
		fmt.Fprintf(&sb, "  %s (version %q):\n", c.Skill, c.Version)
		for _, cs := range c.Constraints {
			fmt.Fprintf(&sb, "    %s (required by %s)\n", cs.Constraint, cs.RequiredBy)
		}
	}
	return sb.String()
}

// Resolver resolves a dependency graph against pre-fetched skill metadata.
type Resolver struct {
	// Meta maps skill name to its fetched metadata.
	// Skills absent from Meta are treated as versionless (version = "").
	Meta map[string]SkillMeta
}

// New creates a Resolver with the given skill metadata.
func New(meta map[string]SkillMeta) *Resolver {
	return &Resolver{Meta: meta}
}

type constraint struct {
	raw        string
	requiredBy string
}

// Resolve performs flat dependency resolution starting from root deps.
// Returns lock entries sorted by skill name, or ErrConflict on version mismatch.
// Skills with version "" (no metadata) satisfy any constraint; non-wildcard constraints
// against a versionless skill emit a warning in the entry but do not fail resolution.
func (r *Resolver) Resolve(deps map[string]manifest.DepSpec) ([]lock.Entry, error) {
	// collected[skillName] = list of constraints imposed on that skill
	collected := make(map[string][]constraint)

	// BFS over dep graph
	type workItem struct {
		depKey     string
		spec       manifest.DepSpec
		requiredBy string
	}

	queue := make([]workItem, 0, len(deps))
	for key, spec := range deps {
		queue = append(queue, workItem{key, spec, "grimoire.toml"})
	}

	visited := make(map[string]bool) // keyed by skill name

	for len(queue) > 0 {
		w := queue[0]
		queue = queue[1:]

		skillName := skillNameFromKey(w.depKey)

		if w.spec.Version != "" && w.spec.Version != "*" {
			collected[skillName] = append(collected[skillName], constraint{w.spec.Version, w.requiredBy})
		}

		if visited[skillName] {
			continue
		}
		visited[skillName] = true

		// Expand transitive deps if metadata available
		if meta, ok := r.Meta[skillName]; ok {
			for transDep, transConstraint := range meta.Dependencies {
				queue = append(queue, workItem{transDep, manifest.DepSpec{Version: transConstraint}, skillName})
			}
		}
	}

	// Resolve constraints for each skill
	var conflicts []Conflict
	entries := make([]lock.Entry, 0, len(visited))

	for skillName := range visited {
		meta, hasMeta := r.Meta[skillName]

		skillConstraints := collected[skillName]

		if len(skillConstraints) > 0 && hasMeta && meta.Version != "" {
			v, err := semver.NewVersion(meta.Version)
			if err == nil {
				var violated []ConstraintSource
				for _, sc := range skillConstraints {
					c, err := semver.NewConstraint(sc.raw)
					if err != nil {
						continue // malformed constraint — skip
					}
					if !c.Check(v) {
						violated = append(violated, ConstraintSource{sc.raw, sc.requiredBy})
					}
				}
				if len(violated) > 0 {
					conflicts = append(conflicts, Conflict{skillName, meta.Version, violated})
					continue
				}
			}
		}

		entry := lock.Entry{Name: skillName}
		if hasMeta {
			entry.Version = meta.Version
			entry.Source = meta.Source
			entry.Resolved = meta.Resolved
			entry.Commit = meta.Commit
			entry.Checksum = meta.Checksum
		}
		entries = append(entries, entry)
	}

	if len(conflicts) > 0 {
		sort.Slice(conflicts, func(i, j int) bool { return conflicts[i].Skill < conflicts[j].Skill })
		return nil, &ErrConflict{Conflicts: conflicts}
	}

	sort.Slice(entries, func(i, j int) bool { return entries[i].Name < entries[j].Name })
	return entries, nil
}

// skillNameFromKey extracts the skill name from a dep key.
//
//	"apply-solid"                          → "apply-solid"
//	"acmecorp/practices:apply-tdd"         → "apply-tdd"
//	"github.com/myteam/skills:domain/skill"→ "skill"
func skillNameFromKey(key string) string {
	if idx := strings.Index(key, ":"); idx >= 0 {
		return path.Base(key[idx+1:])
	}
	return path.Base(key)
}
