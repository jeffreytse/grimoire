package rules

import (
	"github.com/jeffreytse/grimoire/internal/compliance"
	"github.com/jeffreytse/grimoire/internal/skills"
)

// Engine runs deterministic structural rules against the local grimoire setup.
// Findings are returned as compliance.Diagnostic with Source="grimoire-rules".
type Engine struct {
	SkillsRegistries []skills.SkillsRegistry
	ProjectDir       string
}

// Run executes all rules and returns all findings.
func (e *Engine) Run() []compliance.Diagnostic {
	var all []compliance.Diagnostic
	all = append(all, checkSkillHasSkillMd(e.SkillsRegistries)...)
	all = append(all, checkSkillMdFrontmatter(e.SkillsRegistries)...)
	all = append(all, checkAgentBrokenSymlinks()...)
	all = append(all, checkSettingsParseable(e.ProjectDir)...)
	return all
}
