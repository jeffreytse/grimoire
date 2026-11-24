# Editor integration (LSP)

Live compliance diagnostics in your editor.

← [Back to the README](../README.md)

`grimoire lsp` is a [language server](https://microsoft.github.io/language-server-protocol/): any editor that supports language servers (Neovim, Helix, VS Code with a generic client, …) shows compliance findings inline as you work — no plugin needed beyond a one-time setup.

**How it works:** on every file save the server runs `grimoire check` in the background and shows the findings in your editor. Pass/hint items are suppressed; only errors, warnings, and info appear.

| Editor         | Setup                                                          |
| -------------- | -------------------------------------------------------------- |
| Neovim         | `lspconfig` custom server (see below)                          |
| VSCode         | grimoire extension, or `tasks.json` with a custom server entry |
| Helix          | `languages.toml` custom language server entry                  |
| Any LSP client | Point `cmd` at `grimoire lsp`                                  |

**Neovim (nvim-lspconfig):**

```lua
local lspconfig = require('lspconfig')
local configs = require('lspconfig.configs')

if not configs.grimoire then
  configs.grimoire = {
    default_config = {
      cmd = { 'grimoire', 'lsp' },
      filetypes = { 'go', 'python', 'javascript', 'typescript', 'rust', 'ruby' },
      root_dir = lspconfig.util.root_pattern('grimoire.toml', '.git'),
      single_file_support = true,
    },
  }
end

lspconfig.grimoire.setup {}
```

**Helix (`languages.toml`):**

```toml
[[language-server]]
name = "grimoire"
command = "grimoire"
args = ["lsp"]

[[language]]
name = "go"
language-servers = ["gopls", "grimoire"]
```

**VSCode (`settings.json` — requires grimoire extension or a custom extension):**

```json
{
  "grimoire.lsp.enable": true,
  "grimoire.lsp.command": "grimoire",
  "grimoire.lsp.args": ["lsp"]
}
```
