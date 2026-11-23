package cmd

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"

	"github.com/spf13/cobra"

	"github.com/jeffreytse/grimoire/internal/agent"
)

// Antigravity hooks (~/.gemini/config/hooks.json) use their own protocol:
// camelCase input with lastUserInput, and {"injectSteps": [...]} output. It
// has no per-prompt event — PreInvocation fires before every model call — and
// no user-level instruction file, so grimoire registers two hooks:
//
//	SessionStart   grimoire hook antigravity-session  injects the instruction line
//	PreInvocation  grimoire hook antigravity-prompt   suggests skills, once per prompt
//
// Like `grimoire hook prompt`, both exit 0 with no output on any failure.

var hookAgySessionCmd = &cobra.Command{
	Use:           "antigravity-session",
	Short:         "Antigravity SessionStart hook: inject the grimoire instruction",
	Args:          cobra.NoArgs,
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(_ *cobra.Command, _ []string) error {
		runAgySession(os.Stdin, os.Stdout)
		return nil
	},
}

var hookAgyPromptCmd = &cobra.Command{
	Use:           "antigravity-prompt",
	Short:         "Antigravity PreInvocation hook: suggest skills matching the prompt",
	Args:          cobra.NoArgs,
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(_ *cobra.Command, _ []string) error {
		runAgyPrompt(os.Stdin, os.Stdout, agyStateDir())
		return nil
	},
}

func init() {
	hookCmd.AddCommand(hookAgySessionCmd, hookAgyPromptCmd)
}

type agyInput struct {
	ConversationID string `json:"conversationId"`
	LastUserInput  string `json:"lastUserInput"`
}

type agyStep struct {
	SystemMessage    *agySystemMessage `json:"systemMessage,omitempty"`
	EphemeralMessage string            `json:"ephemeralMessage,omitempty"`
}

type agySystemMessage struct {
	SystemMessage string `json:"systemMessage"`
}

type agyOutput struct {
	InjectSteps []agyStep `json:"injectSteps"`
}

func writeAgy(out io.Writer, step agyStep) {
	b, err := json.Marshal(agyOutput{InjectSteps: []agyStep{step}})
	if err == nil {
		fmt.Fprintln(out, string(b))
	}
}

// runAgySession injects the instruction line that other agents read from
// their instruction file (Antigravity has no user-level one), followed by
// grimoire's router.
func runAgySession(in io.Reader, out io.Writer) {
	_, _ = io.Copy(io.Discard, io.LimitReader(in, 1<<20))
	ctx, _ := sessionTurn(traceEnabled())
	msg := agent.MCPTriggerLine + "\n\n" + ctx
	writeAgy(out, agyStep{SystemMessage: &agySystemMessage{SystemMessage: msg}})
}

// runAgyPrompt suggests skills for the latest user prompt. PreInvocation runs
// before every model call in a turn, so the context is injected only when
// lastUserInput differs from the one last handled for this conversation
// (remembered under stateDir).
func runAgyPrompt(in io.Reader, out io.Writer, stateDir string) {
	raw, err := io.ReadAll(io.LimitReader(in, 1<<20))
	if err != nil {
		return
	}
	var input agyInput
	if json.Unmarshal(raw, &input) != nil || input.LastUserInput == "" {
		return
	}
	if stateDir != "" && input.ConversationID != "" {
		sum := sha256.Sum256([]byte(input.LastUserInput))
		hash := hex.EncodeToString(sum[:])
		state := filepath.Join(stateDir, agyStateName(input.ConversationID))
		if prev, err := os.ReadFile(state); err == nil && string(prev) == hash {
			return
		}
		if os.MkdirAll(stateDir, 0o700) == nil {
			_ = os.WriteFile(state, []byte(hash), 0o600)
		}
	}
	// Antigravity has no user-visible hook field, so the status line is
	// dropped; the model's trace lines still show grimoire's decisions.
	if ctx, _ := hookTurn(input.LastUserInput, traceEnabled()); ctx != "" {
		writeAgy(out, agyStep{EphemeralMessage: ctx})
	}
}

var agyUnsafeName = regexp.MustCompile(`[^A-Za-z0-9_-]`)

// agyStateName keeps a conversation id usable as a file name.
func agyStateName(id string) string {
	return agyUnsafeName.ReplaceAllString(id, "_")
}

// agyStateDir is where runAgyPrompt remembers the last prompt per
// conversation, or "" when no cache dir is available.
func agyStateDir() string {
	dir, err := os.UserCacheDir()
	if err != nil {
		return ""
	}
	return filepath.Join(dir, "grimoire", "antigravity-hook")
}
