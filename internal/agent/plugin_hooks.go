package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
)

// Prompt hooks for agents whose hook format differs from the
// Claude/Codex/Gemini settings shape (see prompt_hook.go).

// ── Antigravity ──────────────────────────────────────────────────────────────
//
// hooks.json maps a hook name to its events:
//
//	{"grimoire": {"SessionStart": [{"type": "command", ...}], "PreInvocation": [...]}}
//
// grimoire owns the "grimoire" entry and leaves every other hook alone.

// agyHookName is the hook entry grimoire owns in Antigravity's hooks.json.
const agyHookName = "grimoire"

// agyHookTimeout bounds each hook (seconds).
const agyHookTimeout = 10

// AntigravityHooksFile returns Antigravity's hooks.json for scope.
func AntigravityHooksFile(scope, projectDir string) string {
	if scope == "project" {
		return filepath.Join(projectDir, ".agents", "hooks.json")
	}
	return filepath.Join(mustHome(), ".gemini", "config", "hooks.json")
}

func agyHookEntry(bin string) map[string]any {
	handler := func(sub string) []any {
		return []any{map[string]any{
			"type":    "command",
			"command": shellQuote(bin) + " hook " + sub,
			"timeout": agyHookTimeout,
		}}
	}
	return map[string]any{
		"SessionStart":  handler("antigravity-session"),
		"PreInvocation": handler("antigravity-prompt"),
	}
}

// ConfigureAntigravityHooks writes grimoire's hook entry into the hooks.json
// at path. Reports whether the file changed.
func ConfigureAntigravityHooks(path, bin string) (bool, error) {
	cfg, err := readSettings(path)
	if err != nil {
		return false, err
	}
	want := agyHookEntry(bin)
	// Compare through JSON so numbers read back as float64 still match.
	if cur, ok := cfg[agyHookName]; ok && jsonEqual(cur, want) {
		return false, nil
	}
	cfg[agyHookName] = want
	return true, writeSettings(path, cfg)
}

// RemoveAntigravityHooks drops grimoire's hook entry from the hooks.json at
// path. Reports whether the file changed.
func RemoveAntigravityHooks(path string) (bool, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false, nil
	}
	cfg, err := readSettings(path)
	if err != nil {
		return false, err
	}
	if _, ok := cfg[agyHookName]; !ok {
		return false, nil
	}
	delete(cfg, agyHookName)
	return true, writeSettings(path, cfg)
}

func jsonEqual(a, b any) bool {
	ja, errA := json.Marshal(a)
	jb, errB := json.Marshal(b)
	if errA != nil || errB != nil {
		return false
	}
	var va, vb any
	if json.Unmarshal(ja, &va) != nil || json.Unmarshal(jb, &vb) != nil {
		return false
	}
	return reflect.DeepEqual(va, vb)
}

// ── OpenCode ─────────────────────────────────────────────────────────────────
//
// OpenCode loads every JS/TS file in its plugins directory at startup. The
// grimoire plugin runs `grimoire hook prompt` on each new user message
// (chat.message) and adds the result to the system prompt for that turn
// (experimental.chat.system.transform), so nothing is written into the
// session history.

// openCodePluginMarker identifies a plugin file grimoire wrote; uninstall
// only deletes a file that carries it.
const openCodePluginMarker = "// grimoire prompt hook — written by `grimoire mcp install`"

// OpenCodePluginFile returns the grimoire plugin path for scope.
func OpenCodePluginFile(scope, projectDir string) string {
	if scope == "project" {
		return filepath.Join(projectDir, ".opencode", "plugins", "grimoire.js")
	}
	return filepath.Join(OpenCodeConfigDir(), "plugins", "grimoire.js")
}

const openCodePluginTemplate = `%s; ` + "`grimoire mcp uninstall`" + ` removes it.
// OpenCode has no session-start hook, so this plugin plays both of grimoire's
// hook roles: it keeps grimoire's workflow router (start-best-practice) in
// the system prompt of every turn, adds the per-prompt routing reminder and
// skill candidates, and shows grimoire's status lines ("grimoire ▸ …") as
// toasts. Nothing is written into the session history.
const BIN = %s
const TIMEOUT_MS = 10000
const contextBySession = new Map()
const routedSessions = new Set()
let router = null // Promise of the session hook output, shared by all sessions

async function hook(sub, input) {
  try {
    const proc = Bun.spawn([BIN, "hook", sub], {
      stdin: new Blob([JSON.stringify(input)]),
      stdout: "pipe",
      stderr: "ignore",
    })
    const timer = setTimeout(() => proc.kill(), TIMEOUT_MS)
    const out = await new Response(proc.stdout).text()
    clearTimeout(timer)
    if (!out.trim()) return {}
    const o = JSON.parse(out)
    return { context: o.hookSpecificOutput?.additionalContext ?? "", status: o.systemMessage ?? "" }
  } catch {
    return {}
  }
}

async function toast(client, message) {
  if (!message) return
  try {
    await client?.tui?.showToast?.({ body: { message, variant: "info" } })
  } catch {}
}

export const GrimoirePromptHook = async ({ client } = {}) => ({
  "chat.message": async (input, output) => {
    router ??= hook("session", { hook_event_name: "SessionStart" })
    if (!routedSessions.has(input.sessionID)) {
      routedSessions.add(input.sessionID)
      await toast(client, (await router).status)
    }
    const prompt = (output.parts ?? [])
      .filter((p) => p.type === "text" && !p.synthetic)
      .map((p) => p.text)
      .join("\n")
    const { context = "", status = "" } = prompt
      ? await hook("prompt", { prompt, hook_event_name: "UserPromptSubmit" })
      : {}
    contextBySession.set(input.sessionID, context)
    await toast(client, status)
  },
  "experimental.chat.system.transform": async (input, output) => {
    const { context: routerCtx = "" } = router ? await router : {}
    if (routerCtx) output.system.push(routerCtx)
    const ctx = input.sessionID ? contextBySession.get(input.sessionID) : ""
    if (ctx) output.system.push(ctx)
  },
})
`

func openCodePlugin(bin string) []byte {
	quoted, _ := json.Marshal(bin) // a JSON string is a valid JS string literal
	return []byte(fmt.Sprintf(openCodePluginTemplate, openCodePluginMarker, quoted))
}

// ConfigureOpenCodePlugin writes the grimoire plugin to path. A file at path
// that grimoire did not write is left alone (error). Reports whether it changed.
func ConfigureOpenCodePlugin(path, bin string) (bool, error) {
	want := openCodePlugin(bin)
	cur, err := os.ReadFile(path)
	switch {
	case err == nil && bytes.Equal(cur, want):
		return false, nil
	case err == nil && !bytes.HasPrefix(cur, []byte(openCodePluginMarker)):
		return false, fmt.Errorf("%s exists and was not written by grimoire — remove or rename it, then rerun", path)
	case err != nil && !os.IsNotExist(err):
		return false, err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return false, err
	}
	return true, os.WriteFile(path, want, 0o644)
}

// RemoveOpenCodePlugin deletes the grimoire plugin at path if grimoire wrote
// it. Reports whether it was removed.
func RemoveOpenCodePlugin(path string) (bool, error) {
	cur, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	if !bytes.HasPrefix(cur, []byte(openCodePluginMarker)) {
		return false, nil
	}
	return true, os.Remove(path)
}
