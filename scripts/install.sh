#!/usr/bin/env bash
set -euo pipefail

# ── Colors ────────────────────────────────────────────────────────────────────
# Disabled when NO_COLOR is set or stdout is not a TTY.
if [ -t 1 ] && [ -z "${NO_COLOR:-}" ]; then
  BOLD=$'\033[1m'; DIM=$'\033[2m'; RESET=$'\033[0m'
  GREEN=$'\033[32m'; CYAN=$'\033[36m'; YELLOW=$'\033[33m'; RED=$'\033[31m'
  GOLD=$'\033[38;5;178m'; STAR=$'\033[38;5;227m'; LINE=$'\033[38;5;100m'
else
  BOLD=""; DIM=""; RESET=""; GREEN=""; CYAN=""; YELLOW=""; RED=""
  GOLD=""; STAR=""; LINE=""
fi

info()    { printf "  ${CYAN}→${RESET}  %s\n"        "$*"; }
success() { printf "  ${GREEN}✓${RESET}  %s\n"       "$*"; }
warn()    { printf "  ${YELLOW}!${RESET}  %s\n" "$*" >&2; }
error()   { printf "  ${RED}✗${RESET}  %s\n"   "$*" >&2; exit 1; }

osc8() { printf "\033]8;;%s\007%s\033]8;;\007" "$1" "$2"; }  # clickable hyperlink

# ── Banner ────────────────────────────────────────────────────────────────────
printf "\n"
printf " ${STAR}✦${GOLD}▗▄▄▄▄▄▄▄▄▖${STAR}✦${RESET}  ${BOLD}grimoire${RESET}  ${DIM}installer${RESET}\n"
printf "  ${GOLD}▐${LINE}▬▬▬│▬▬▬▬${GOLD}▌${RESET}   ${DIM}The world's best practices for AI assistants${RESET}\n"
printf "  ${GOLD}▐${LINE}▬▬ ✦ ▬▬▬${GOLD}▌${RESET}   ${CYAN}https://github.com/jeffreytse/grimoire${RESET}\n"
printf "  ${GOLD}▐${LINE}▬▬▬│▬▬▬▬${GOLD}▌${RESET}\n"
printf " ${STAR}✦${GOLD}▝▀▀▀▀▀▀▀▀▘${STAR}✦${RESET}  ⭐ $(osc8 "https://github.com/jeffreytse/grimoire" "Star")  💖 $(osc8 "https://github.com/sponsors/jeffreytse" "Sponsor")  🐛 $(osc8 "https://github.com/jeffreytse/grimoire/issues" "Issues")\n"
printf "\n"

# ── Platform detection ────────────────────────────────────────────────────────
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)
case "$ARCH" in
  x86_64)        ARCH=amd64 ;;
  aarch64|arm64) ARCH=arm64 ;;
  *) error "unsupported architecture: $ARCH" ;;
esac

info "Platform: ${OS}/${ARCH}"

# ── Install destination ───────────────────────────────────────────────────────
INSTALL_DIR="${GRIMOIRE_INSTALL_DIR:-/usr/local/bin}"
DEST="${INSTALL_DIR}/grimoire"

# Existing install — show current version
if command -v grimoire &>/dev/null; then
  CURRENT=$(grimoire --version 2>/dev/null | sed 's/grimoire version //' || echo "unknown")
  info "Replacing existing install ${DIM}(${CURRENT})${RESET}"
fi

# Write permission / sudo
SUDO=""
if [ ! -w "$INSTALL_DIR" ]; then
  if command -v sudo &>/dev/null; then
    warn "No write access to ${INSTALL_DIR} — using sudo"
    SUDO="sudo"
  else
    error "Cannot write to ${INSTALL_DIR}. Set GRIMOIRE_INSTALL_DIR to a writable path."
  fi
fi

# ── Download ──────────────────────────────────────────────────────────────────
BIN="grimoire-${OS}-${ARCH}"
URL="https://github.com/jeffreytse/grimoire/releases/latest/download/${BIN}"
TMP=$(mktemp)
trap 'rm -f "$TMP"' EXIT

info "Downloading ${DIM}${BIN}${RESET}"
if ! curl -fsSL "$URL" -o "$TMP"; then
  error "Download failed. Check https://github.com/jeffreytse/grimoire/releases"
fi

# ── Install ───────────────────────────────────────────────────────────────────
$SUDO install -m 755 "$TMP" "$DEST"

VERSION=$("$DEST" --version 2>/dev/null | sed 's/grimoire version //' || echo "")
VERSION_LABEL="${VERSION:+ ${DIM}(${VERSION})${RESET}}"

success "Installed grimoire${VERSION_LABEL}  →  ${DEST}"

# PATH warning
case ":${PATH}:" in
  *":${INSTALL_DIR}:"*) ;;
  *) warn "${INSTALL_DIR} is not in your PATH — add it to your shell profile" ;;
esac

# ── Next steps ────────────────────────────────────────────────────────────────
printf "\n"
printf "  ${BOLD}Get started:${RESET}\n"
printf "\n"
printf "    ${CYAN}grimoire wizard${RESET}     ${DIM}# interactive setup${RESET}\n"
printf "\n"
printf "  ${DIM}Or manually:${RESET}\n"
printf "\n"
printf "    ${CYAN}grimoire update${RESET}     ${DIM}# fetch the official skill library${RESET}\n"
printf "    ${CYAN}grimoire install${RESET}    ${DIM}# install skills to all detected AI agents${RESET}\n"
printf "    ${CYAN}grimoire check${RESET}      ${DIM}# run a compliance check on your project${RESET}\n"
printf "\n"
