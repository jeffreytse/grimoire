package cmd

import (
	"errors"
	"os"

	"github.com/spf13/cobra"
)

var lspCmd = &cobra.Command{
	Use:   "lsp",
	Short: "Start the grimoire Language Server (LSP) for editor integration",
	Long: `Starts a Language Server Protocol (LSP) server over stdio.

Connect any LSP-capable editor (VSCode, Neovim, Helix, ...) by pointing its
language client at:

  grimoire lsp

The server publishes compliance diagnostics in real time as you edit. On each
file save, grimoire check runs in the background and pushes findings to the
editor as standard LSP diagnostics — visible in the gutter, Problems panel, or
inline alongside your code.`,
	RunE: runLSP,
	Args: cobra.NoArgs,
}

func runLSP(_ *cobra.Command, _ []string) error {
	err := serveLSP(os.Stdin, os.Stdout)
	var exitErr *lspExitError
	if errors.As(err, &exitErr) {
		os.Exit(exitErr.code) //nolint:revive // LSP spec requires immediate exit here; cobra RunE cannot encode exit codes
	}
	return err
}
