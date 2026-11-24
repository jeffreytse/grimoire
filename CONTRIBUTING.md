# Contributing to grimoire

Thanks for helping. This repository is the **grimoire CLI and MCP server**. To add or improve a **skill** (a best practice), contribute to [grimoire-core](https://github.com/jeffreytse/grimoire-core/blob/main/CONTRIBUTING.md) instead.

## Before you start

- **Bugs and small fixes:** open a pull request directly.
- **New features or behavior changes:** open an issue first, so we can agree on the approach before you write code.
- **Where new features go:** MCP mode is grimoire's primary mode. New capabilities land in MCP mode first. Skills mode is fully supported but gets fixes and compatibility work, not features that only work there.

## Set up

You need Go 1.26.3 (see `go.mod`), git, and [golangci-lint](https://golangci-lint.run) for linting.

```bash
git clone https://github.com/jeffreytse/grimoire.git
cd grimoire
make build   # build the binary into bin/
make test    # run all tests
make lint    # golangci-lint, the same check CI runs
```

`make help` lists the other targets (`vet`, `fmt`, `tidy`, `dist`, `install`).

Testing an install? Point grimoire at a scratch home, so your real setup isn't touched. Set both, because grimoire reads `XDG_CONFIG_HOME` too:

```bash
HOME=$(mktemp -d) XDG_CONFIG_HOME=$(mktemp -d) ./bin/grimoire mcp install --target claude
```

## Making a change

1. **Write a test.** Bug fixes come with a test that fails before the fix. Features come with tests for their behavior.
2. **Keep it passing.** `make test` and `make lint` must pass. CI runs golangci-lint on every push and pull request.
3. **Update the docs** when behavior changes: `README.md` for anything a new user meets, and `docs/` for reference material.
4. **Commit messages** follow [Conventional Commits](https://www.conventionalcommits.org/), e.g. `fix(mcp): …`, `feat(install): …`, `docs(readme): …`, `test(eval): …`. Explain *why* in the body when it isn't obvious.

## Evals (optional, costs API credits)

Some behavior is measured against a real agent, not unit tests. These evals call Claude Code and are skipped unless you opt in:

| Eval | Run |
| --- | --- |
| Display | `GRIMOIRE_DISPLAY_EVAL=1 go test ./cmd -run TestDisplayEval -v -timeout 60m` |
| Keyword candidates (no API cost) | `GRIMOIRE_HOOK_EVAL=1 go test ./cmd -run TestHookEval -v` |
| Skill selection | `GRIMOIRE_SELECTION_EVAL=1 go test ./cmd -run TestSelectionEval -v -timeout 60m` |
| Outcome (with vs. without grimoire) | `GRIMOIRE_OUTCOME_EVAL=1 go test ./cmd -run TestOutcomeEval -v -timeout 120m` |

Each test's doc comment lists its options. Run them when you change routing, the display, or skill search, and include the before-and-after numbers in your pull request.

## Code of Conduct

By taking part you agree to follow our [Code of Conduct](./CODE_OF_CONDUCT.md).
