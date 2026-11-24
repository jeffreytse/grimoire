package cmd

import (
	"sort"

	"github.com/spf13/cobra"
)

// helpGroups orders `grimoire --help` by what a user is trying to do,
// starting with MCP setup, grimoire's primary mode. Each top-level command
// is listed in exactly one group; TestHelpGroupsCoverEveryCommand keeps it
// that way as commands are added.
var helpGroups = []struct {
	id, title string
	commands  []string
}{
	{"start", "Get started:", []string{"mcp", "wizard", "init", "enable", "disable"}},
	{"skills", "Find and install skills:", []string{"search", "list", "info", "install", "uninstall", "update", "profile"}},
	{"standards", "Project standards:", []string{"check", "status", "watch", "lock"}},
	{"maintain", "Configure and maintain:", []string{"config", "package", "doctor", "clean", "self-update", "version"}},
	{"build", "Integrations and authoring:", []string{"context", "lsp", "statusline", "generate", "validate", "completion", "help"}},
}

// groupCommands assigns every visible top-level command to its help group.
// Called at Execute time, after every file's init has registered its
// commands. Safe to call more than once.
func groupCommands(root *cobra.Command) {
	if len(root.Groups()) == 0 {
		for _, g := range helpGroups {
			root.AddGroup(&cobra.Group{ID: g.id, Title: g.title})
		}
	}
	root.SetHelpCommandGroupID("build")
	root.SetCompletionCommandGroupID("build")
	byName := map[string]string{}
	for _, g := range helpGroups {
		for _, name := range g.commands {
			byName[name] = g.id
		}
	}
	for _, c := range root.Commands() {
		if id, ok := byName[c.Name()]; ok {
			c.GroupID = id
		}
	}
	// List each group in the order above, not alphabetically, so the
	// first command under "Get started" is the one to run first.
	cobra.EnableCommandSorting = false
	cmds := root.Commands()
	rank := map[string]int{}
	i := 0
	for _, g := range helpGroups {
		for _, name := range g.commands {
			rank[name] = i
			i++
		}
	}
	sort.SliceStable(cmds, func(a, b int) bool {
		ra, okA := rank[cmds[a].Name()]
		rb, okB := rank[cmds[b].Name()]
		if !okA || !okB {
			return okA && !okB
		}
		return ra < rb
	})
	root.RemoveCommand(cmds...)
	root.AddCommand(cmds...)
}
