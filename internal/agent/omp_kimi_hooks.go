package agent

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pelletier/go-toml/v2"
)

// ── OMP (oh-my-pi) ───────────────────────────────────────────────────────────
//
// OMP loads TS/JS extensions from <agent dir>/extensions (user) and
// <cwd>/.omp/extensions (project; OMP has no trust step). Its
// before_agent_start event carries the current system prompt as a string
// array and takes a replacement array back, so the extension appends
// grimoire's router and per-prompt context to it, keeping them out of the
// transcript. A "revive" (a parked subagent waking) is skipped: OMP discards
// its result. API checked against oh-my-pi's
// packages/coding-agent/src/extensibility/extensions/types.ts.

const ompExtensionMarker = "// grimoire prompt hook — written by `grimoire mcp install`"

// OMPExtensionFile returns the grimoire extension path for scope.
func OMPExtensionFile(scope, projectDir string) string {
	if scope == "project" {
		return filepath.Join(projectDir, ".omp", "extensions", "grimoire.ts")
	}
	return filepath.Join(ompAgentDir(mustHome()), "extensions", "grimoire.ts")
}

const ompExtensionTemplate = `%s; ` + "`grimoire mcp uninstall`" + ` removes it.
// OMP has no settings-file hooks, so this extension plays grimoire's two hook
// roles: before each user turn it adds grimoire's workflow router
// (start-best-practice) and the per-prompt routing reminder and skill
// candidates to the system prompt, and shows grimoire's status lines
// ("grimoire ▸ …") as notifications.
import { spawn } from "node:child_process";
import type { ExtensionAPI } from "@oh-my-pi/pi-coding-agent";

const BIN = %s;
const TIMEOUT_MS = 10000;

type HookOut = { context?: string; status?: string };

function hook(sub: string, input: unknown, cwd: string): Promise<HookOut> {
	return new Promise((resolve) => {
		let out = "";
		let done = false;
		const finish = (v: HookOut) => {
			if (!done) {
				done = true;
				resolve(v);
			}
		};
		try {
			const proc = spawn(BIN, ["hook", sub], { cwd, stdio: ["pipe", "pipe", "ignore"] });
			const timer = setTimeout(() => {
				proc.kill();
				finish({});
			}, TIMEOUT_MS);
			proc.stdout.on("data", (d) => (out += d));
			proc.on("error", () => finish({}));
			proc.on("close", () => {
				clearTimeout(timer);
				if (!out.trim()) return finish({});
				try {
					const o = JSON.parse(out);
					finish({
						context: o.hookSpecificOutput?.additionalContext ?? "",
						status: o.systemMessage ?? "",
					});
				} catch {
					finish({});
				}
			});
			proc.stdin.end(JSON.stringify(input));
		} catch {
			finish({});
		}
	});
}

export default function (pi: ExtensionAPI) {
	let router: Promise<HookOut> | null = null;
	let announced = false;

	pi.on("before_agent_start", async (event, ctx) => {
		if (event.reason !== "prompt" || !event.prompt) return;
		const cwd = ctx.cwd ?? process.cwd();

		router ??= hook("session", { hook_event_name: "SessionStart", cwd }, cwd);
		const r = await router;
		if (!announced && r.status) {
			announced = true;
			ctx.ui?.notify?.(r.status, "info");
		}

		let session_id: string | undefined;
		try {
			session_id = ctx.sessionManager?.getSessionId?.();
		} catch {}
		const p = await hook("prompt", { prompt: event.prompt, hook_event_name: "UserPromptSubmit", cwd, session_id }, cwd);
		if (p.status) ctx.ui?.notify?.(p.status, "info");

		const extra = [r.context, p.context].filter(Boolean);
		if (extra.length) return { systemPrompt: [...event.systemPrompt, ...extra] };
	});
}
`

func ompExtension(bin string) []byte {
	quoted, _ := json.Marshal(bin) // a JSON string is a valid TS string literal
	return []byte(fmt.Sprintf(ompExtensionTemplate, ompExtensionMarker, quoted))
}

// ConfigureOMPExtension writes the grimoire extension to path. A file at path
// that grimoire did not write is left alone (error). Reports whether it changed.
func ConfigureOMPExtension(path, bin string) (bool, error) {
	return writeOwnedFile(path, ompExtension(bin), ompExtensionMarker)
}

// RemoveOMPExtension deletes the grimoire extension at path if grimoire wrote
// it. Reports whether it was removed.
func RemoveOMPExtension(path string) (bool, error) {
	return removeOwnedFile(path, ompExtensionMarker)
}

// ── Kimi Code ────────────────────────────────────────────────────────────────
//
// Kimi Code reads hooks only from [[hooks]] in its user config.toml (there is
// no project hook file). A UserPromptSubmit hook's stdout is appended to the
// context as text; SessionStart is observation-only, so its output never
// reaches the model. grimoire therefore registers only the prompt hook, in
// text form, and its reminder tells the model to load the router with
// grimoire_get_skill. Each entry may hold only event, matcher, command and
// timeout (seconds); any other field stops the config from loading.
// (https://moonshotai.github.io/kimi-code/en/customization/hooks)

// kimiHookMarker precedes grimoire's [[hooks]] entry in config.toml; uninstall
// removes exactly the marked entry.
const kimiHookMarker = "# grimoire prompt hook — written by `grimoire mcp install`; `grimoire mcp uninstall` removes it"

const kimiHookTimeout = 10 // seconds

// KimiConfigFile returns Kimi Code's user config.toml.
func KimiConfigFile() string {
	return filepath.Join(kimiHome(mustHome()), "config.toml")
}

func kimiHookBlock(bin string) string {
	cmd, _ := json.Marshal(shellQuote(bin) + " " + promptHookArgs + " --format text") // a JSON string is a valid TOML basic string
	return fmt.Sprintf("%s\n[[hooks]]\nevent = \"UserPromptSubmit\"\ncommand = %s\ntimeout = %d\n", kimiHookMarker, cmd, kimiHookTimeout)
}

// ConfigureKimiHook adds grimoire's prompt hook to Kimi's config.toml at path,
// replacing an older grimoire entry and keeping the rest of the file as it
// was. Reports whether the file changed.
func ConfigureKimiHook(path, bin string) (bool, error) {
	cur, err := os.ReadFile(path)
	if err != nil && !os.IsNotExist(err) {
		return false, err
	}
	block := kimiHookBlock(bin)
	if strings.Count(string(cur), kimiHookMarker) == 1 && strings.Contains(string(cur), block) {
		return false, nil
	}
	rest := withoutKimiHook(string(cur))
	if rest != "" && !strings.HasSuffix(rest, "\n") {
		rest += "\n"
	}
	if rest != "" {
		rest += "\n"
	}
	next := rest + block
	// Never leave Kimi with a config it can't load.
	var probe map[string]any
	if err := toml.Unmarshal([]byte(next), &probe); err != nil {
		return false, fmt.Errorf("%s: adding the grimoire hook would not parse (%v) — fix the file by hand", path, err)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return false, err
	}
	return true, os.WriteFile(path, []byte(next), 0o644)
}

// RemoveKimiHook drops grimoire's marked [[hooks]] entry from path. Reports
// whether the file changed.
func RemoveKimiHook(path string) (bool, error) {
	cur, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	rest := withoutKimiHook(string(cur))
	if rest == string(cur) {
		return false, nil
	}
	return true, os.WriteFile(path, []byte(strings.TrimRight(rest, "\n")+"\n"), 0o644)
}

// HasKimiHook reports whether path holds grimoire's prompt hook.
func HasKimiHook(path string) bool {
	b, err := os.ReadFile(path)
	return err == nil && strings.Contains(string(b), kimiHookMarker)
}

// withoutKimiHook removes each marked grimoire entry: the marker line and the
// [[hooks]] table after it, up to the next blank line, comment or table.
func withoutKimiHook(s string) string {
	if !strings.Contains(s, kimiHookMarker) {
		return s
	}
	lines := strings.Split(s, "\n")
	out := make([]string, 0, len(lines))
	for i := 0; i < len(lines); i++ {
		if strings.TrimSpace(lines[i]) != kimiHookMarker {
			out = append(out, lines[i])
			continue
		}
		j := i + 1
		if j < len(lines) && strings.TrimSpace(lines[j]) == "[[hooks]]" {
			j++
			for j < len(lines) {
				t := strings.TrimSpace(lines[j])
				if t == "" || strings.HasPrefix(t, "[") || strings.HasPrefix(t, "#") {
					break
				}
				j++
			}
		}
		i = j - 1
		// Drop the blank line grimoire added before its entry.
		if n := len(out); n > 0 && strings.TrimSpace(out[n-1]) == "" {
			out = out[:n-1]
		}
	}
	return strings.Join(out, "\n")
}
