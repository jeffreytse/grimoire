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
