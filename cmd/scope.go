package cmd

import (
	"fmt"
	"slices"
	"strings"

	"github.com/jeffreytse/grimoire/internal/agent"
)

// Every command names its scopes the same way: "user" (all your projects)
// and "project" (this project), plus "system" for config. "global" is still
// accepted as an older name for "user", so existing scripts keep working.
const scopeHelp = `"project" (this project) or "user" (all your projects)`

// normalizeScope maps a --scope value to "user" or "project".
func normalizeScope(s string) (string, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "user", "global":
		return "user", nil
	case "project":
		return "project", nil
	}
	return "", fmt.Errorf(`--scope must be "user" or "project", got %q`, s)
}

// installUserScope reports whether install acts on the user scope.
func installUserScope() bool {
	s, _ := normalizeScope(flagInstallScope)
	return flagInstallGlobal || s == "user"
}

// uninstallUserScope reports whether uninstall acts on the user scope.
func uninstallUserScope() bool {
	s, _ := normalizeScope(flagUninstallScope)
	return flagUninstallGlobal || s == "user"
}

// validateTarget rejects an unknown --target before any work starts.
func validateTarget(target string) error {
	switch target {
	case "", "auto", "all":
		return nil
	}
	if slices.Contains(agent.All, target) {
		return nil
	}
	return fmt.Errorf("unknown --target %q — choose one of: %s, all, auto", target, strings.Join(agent.All, ", "))
}
