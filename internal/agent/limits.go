package agent

// SkillLoadCap returns the known default skill-loading cap for an agent, or 0 if unknown.
// These caps are agent-internal limits on how many skills are read from disk.
func SkillLoadCap(ag string) int {
	switch ag {
	case "openclaw":
		return 200
	default:
		return 0
	}
}

// SkillListingBudget returns how many characters of skill metadata (names and
// descriptions) an agent loads into its context to discover skills, or 0 if
// unknown. Skills past the budget are truncated or dropped from the listing,
// so the agent's own discovery cannot find them. The values are the agents'
// documented fallbacks — Claude Code scales its budget with 1–2% of the
// context window (8,000–16,000 characters as fallback), Codex uses 2% or
// 8,000 characters — so estimates against them are conservative.
func SkillListingBudget(ag string) int {
	switch ag {
	case "claude", "codex":
		return 8000
	default:
		return 0
	}
}

// SkillListingCost is the characters one skill takes in an agent's skill
// listing: its name, its description and the entry's formatting.
func SkillListingCost(name, description string) int {
	return len(name) + len(description) + 8
}
