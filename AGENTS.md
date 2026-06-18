# grimoire CLI

Go CLI for installing and managing grimoire skill libraries across AI agents.

- Skills library: https://github.com/jeffreytse/grimoire-skills
- Website: https://grimoire.jeffreytse.net

## Structure

```
cmd/            # Cobra commands: install, update, check, clean, doctor,
                #   init, list, uninstall, version, self-update
internal/
  agent/        # Agent detection and linking (Claude, Codex, Gemini, etc.)
  compliance/   # Skill compliance checking
  detect/       # OS and agent auto-detection
  git/          # Git operations (clone, pull, commit detection)
  skills/       # Skill path resolution, root.go (GrimoireRepo constant)
  tui/          # Terminal UI helpers (icons, color, prompts)
main.go         # Entrypoint
go.mod          # module github.com/jeffreytse/grimoire
scripts/
  install.sh    # Curl-pipe installer (macOS/Linux)
  install.ps1   # PowerShell installer (Windows)
```

## Key constants

- `internal/skills/root.go` — `GrimoireRepo` (points to grimoire-skills repo)

## Build

```bash
go build ./...
make build      # builds local binary
make test       # runs all tests
make dist       # cross-compile all platforms to dist/
```
