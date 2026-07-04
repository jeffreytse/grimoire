// Package context provides project-level context detection for grimoire.
// Context drives profile selection; profiles drive skill membership.
// Skills are pure knowledge — they have no activation conditions of their own.
package context

import (
	"github.com/jeffreytse/grimoire/internal/config"
	"github.com/jeffreytse/grimoire/internal/detect"
)

// ProjectContext holds the detected context for a project directory.
// Profile is the only field that affects which skills are active:
// the detected profile's skill list determines the active set.
type ProjectContext struct {
	Dir     string // absolute path to the project root
	Profile string // detected or configured profile name; "" = no profile
}

// Detect returns the ProjectContext for dir.
// Profile comes from detect.Profile() file-signal heuristics;
// if no signal is found it falls back to the first entry in [standards] profiles
// from the effective settings.
func Detect(dir string) ProjectContext {
	profile := detect.Profile(dir)
	if profile == "" {
		// fall back to user-configured profile
		if r, err := config.Load(dir); err == nil && len(r.Core.Profiles) > 0 {
			profile = r.Core.Profiles[0]
		}
	}
	return ProjectContext{Dir: dir, Profile: profile}
}
