#!/usr/bin/env bash
set -eo pipefail

SCRIPT_PATH="${BASH_SOURCE[0]:-}"
if [[ -n "${SCRIPT_PATH}" && -f "${SCRIPT_PATH}" ]]; then
  REPO_ROOT="$(cd "$(dirname "${SCRIPT_PATH}")/.." && pwd)"
  _TMPDIR=""
else
  _TMPDIR="$(mktemp -d)"
  trap 'rm -rf "${_TMPDIR}"' EXIT
  echo "Downloading grimoire..."
  git clone --depth 1 --quiet https://github.com/jeffreytse/grimoire.git "${_TMPDIR}"
  REPO_ROOT="${_TMPDIR}"
fi
SKILLS_ROOT="${REPO_ROOT}/skills"
CLAUDE_SKILLS_DIR="${HOME}/.claude/skills"
AGENTS_SKILLS_DIR="${HOME}/.agents/skills"
GEMINI_SKILLS_DIR="${HOME}/.gemini/skills"

usage() {
  cat <<EOF
Usage: install.sh [OPTIONS]

Options:
  --domain <name>       Install/uninstall all skills for a domain
  --subdomain <name>    Restrict to one sub-domain within a domain
  --skill <path>        Install/uninstall one skill (e.g. engineering/development/propose-conventional-commit)
  --target <agent>      Target: claude, codex, gemini, all
  --uninstall           Remove skills instead of installing
  --list                List available domains, sub-domains, and skills
  --yes                 Non-interactive: install all skills to all detected agents
  --help                Show this help

Examples:
  install.sh                                                # Interactive TUI
  install.sh --yes                                          # Install everything, no prompts
  install.sh --domain engineering --target claude
  install.sh --skill engineering/development/propose-conventional-commit
  install.sh --uninstall --domain engineering --target claude
  install.sh --uninstall --skill engineering/development/propose-conventional-commit
EOF
}

# ── Banner ────────────────────────────────────────────────────────────────────
print_banner() {
  local version
  version="$(grep -oE '"version": *"[^"]+"' "${REPO_ROOT}/.claude-plugin/plugin.json" 2>/dev/null | grep -oE '[0-9]+\.[0-9]+\.[0-9]+' | head -1)"
  version="${version:-1.0.0}"

  # All colors use $'...' so they contain actual ESC bytes (works in printf %s and variables)
  local e=$'\033'
  local D="${e}[2m"          local W="${e}[1;37m"       local V="${e}[0;36m"
  local R="${e}[0m"          local GD="${e}[38;5;178m"  local ST="${e}[38;5;214m"
  local SK="${e}[38;5;227m"  local LN="${e}[38;5;100m"

  local _bel=$'\007'
  local _hs="${e}]8;;"
  local _URL="https://github.com/jeffreytse/grimoire"
  local _URL_SPONSOR="https://github.com/sponsors/jeffreytse"
  local _URL_ISSUES="https://github.com/jeffreytse/grimoire/issues"

  # Build OSC8 links using actual ESC+BEL bytes — one write per line
  local _star="${_hs}${_URL}${_bel}Star${_hs}${_bel}"
  local _spon="${_hs}${_URL_SPONSOR}${_bel}Sponsor${_hs}${_bel}"
  local _iss="${_hs}${_URL_ISSUES}${_bel}Issues${_hs}${_bel}"

  # Open book: two pages with spine │ in center, magical star ✦ glowing from middle
  printf '\n' >/dev/tty
  printf ' %s✦%s▗▄▄▄▄▄▄▄▄▖%s✦%s  %sgrimoire%s %sv%s%s\n' \
    "${SK}" "${GD}" "${SK}" "${R}" "${W}" "${R}" "${D}" "${version}" "${R}" >/dev/tty
  printf '  %s▐%s▬▬▬%s│%s▬▬▬▬%s▌%s   %sThe world'\''s best practices for AI assistants%s\n' \
    "${GD}" "${LN}" "${GD}" "${LN}" "${GD}" "${R}" "${D}" "${R}" >/dev/tty
  printf '  %s▐%s▬▬ %s✦%s ▬▬▬%s▌%s   %s\n' \
    "${GD}" "${LN}" "${ST}" "${LN}" "${GD}" "${R}" "${_URL}" >/dev/tty
  printf '  %s▐%s▬▬▬%s│%s▬▬▬▬%s▌%s\n' "${GD}" "${LN}" "${GD}" "${LN}" "${GD}" "${R}" >/dev/tty
  printf ' %s✦%s▝▀▀▀▀▀▀▀▀▘%s✦%s  ⭐ %s  💖 %s  🐛 %s\n' \
    "${SK}" "${GD}" "${SK}" "${R}" "${_star}" "${_spon}" "${_iss}" >/dev/tty
  printf '\n' >/dev/tty
}
# ── TUI multiselect ────────────────────────────────────────────────────────────
# Usage: multiselect RESULT_VAR "Prompt" option1 option2 ...
# Sets RESULT_VAR to a space-separated list of selected options.
multiselect() {
  local _ms_var="$1"; shift
  local prompt="$1";  shift
  local -a opts=("$@")
  local -a sel=()
  local cur=0 offset=0
  local n=${#opts[@]}

  # Cap visible rows so draw_height never causes a scroll
  local term_rows
  term_rows="$(tput lines 2>/dev/null)" || term_rows=24
  local MAX_VIS=$(( term_rows - 8 ))
  (( MAX_VIS < 3 )) && MAX_VIS=3
  local vis=$(( n < MAX_VIS ? n : MAX_VIS ))
  # Total lines drawn per render (constant → cursor-up is reliable)
  local draw_height=$vis
  (( n > vis )) && draw_height=$(( vis + 1 ))  # +1 for scroll indicator

  for (( i=0; i<n; i++ )); do
    if [[ "${opts[$i]}" == +* ]]; then
      sel+=("true")
      opts[$i]="${opts[$i]#+}"
    else
      sel+=("false")
    fi
  done

  local _stty_saved
  _stty_saved="$(stty -g </dev/tty 2>/dev/null)" || _stty_saved=""
  [[ -n "$_stty_saved" ]] && stty -echo </dev/tty 2>/dev/null || true
  tput civis >/dev/tty 2>/dev/null || true
  trap '
    tput cnorm >/dev/tty 2>/dev/null || true
    [[ -n "${_stty_saved}" ]] && stty "${_stty_saved}" </dev/tty 2>/dev/null || true
  ' RETURN

  printf "\n\033[1m%s\033[0m\n" "${prompt}" >/dev/tty
  printf "  \033[2m↑↓ navigate   SPACE toggle   A select all   ENTER confirm\033[0m\n\n" >/dev/tty

  local first=1
  while true; do
    (( cur < offset )) && offset=$cur
    (( cur >= offset + vis )) && offset=$(( cur - vis + 1 ))

    # Cursor-up by exact draw_height to overwrite previous render in place
    if [[ $first -eq 0 ]]; then
      printf "\033[%dA" "$draw_height" >/dev/tty
    fi
    first=0

    for (( i=offset; i<offset+vis; i++ )); do
      local mark
      [[ "${sel[$i]}" == "true" ]] && mark="✅" || mark="⬜"
      if [[ $i -eq $cur ]]; then
        printf "\r\033[K  👉 %s \033[1m%s\033[0m\n" "${mark}" "${opts[$i]}" >/dev/tty
      else
        printf "\r\033[K     %s %s\n"               "${mark}" "${opts[$i]}" >/dev/tty
      fi
    done
    if (( n > vis )); then
      local sel_count
      sel_count="$(printf '%s\n' "${sel[@]}" | grep -c true)" || sel_count=0
      printf "\r\033[K  \033[2m(%d/%d)  %d selected\033[0m\n" \
        "$(( cur + 1 ))" "$n" "$sel_count" >/dev/tty
    fi

    local key seq
    IFS= read -rsn1 key </dev/tty
    if [[ "${key}" == $'\x1b' ]]; then
      IFS= read -rsn2 -t 1 seq </dev/tty || true
      key="${key}${seq}"
    fi

    case "${key}" in
      $'\x1b[A') (( cur = (cur - 1 + n) % n )) || true ;;
      $'\x1b[B') (( cur = (cur + 1)     % n )) || true ;;
      ' ')
        if [[ "${sel[$cur]}" == "true" ]]; then sel[$cur]="false"
        else sel[$cur]="true"; fi
        ;;
      a|A)
        local any_off=0
        for s in "${sel[@]}"; do
          if [[ "$s" == "false" ]]; then any_off=1; break; fi
        done
        for (( i=0; i<n; i++ )); do
          if [[ $any_off -eq 1 ]]; then sel[$i]="true"; else sel[$i]="false"; fi
        done
        ;;
      ''|$'\n'|$'\r') break ;;
    esac
  done

  while IFS= read -rn1 -t 0 _x </dev/tty 2>/dev/null; do :; done || true

  printf "\n" >/dev/tty
  tput cnorm >/dev/tty 2>/dev/null || true
  [[ -n "$_stty_saved" ]] && stty "${_stty_saved}" </dev/tty 2>/dev/null || true

  local -a chosen=()
  for (( i=0; i<n; i++ )); do
    [[ "${sel[$i]}" == "true" ]] && chosen+=("${opts[$i]}")
  done
  eval "${_ms_var}=(\"\${chosen[@]}\")"
}

# Usage: select_one RESULT_VAR "Prompt" option1 option2 ...
# Sets RESULT_VAR to the chosen option. cur=0 is default.
select_one() {
  local _so_var="$1"; shift
  local prompt="$1";  shift
  local -a opts=("$@")
  local cur=0 n=${#opts[@]}

  local _stty_saved
  _stty_saved="$(stty -g </dev/tty 2>/dev/null)" || _stty_saved=""
  [[ -n "$_stty_saved" ]] && stty -echo </dev/tty 2>/dev/null || true
  tput civis >/dev/tty 2>/dev/null || true
  trap '
    tput cnorm >/dev/tty 2>/dev/null || true
    [[ -n "${_stty_saved}" ]] && stty "${_stty_saved}" </dev/tty 2>/dev/null || true
  ' RETURN

  printf "\n\033[1m%s\033[0m\n" "${prompt}" >/dev/tty
  printf "  \033[2m↑↓ navigate   ENTER confirm\033[0m\n\n" >/dev/tty

  local first=1
  while true; do
    if [[ $first -eq 0 ]]; then
      printf "\033[%dA" "$n" >/dev/tty
    fi
    first=0

    for (( i=0; i<n; i++ )); do
      if [[ $i -eq $cur ]]; then
        printf "\r\033[K  👉 \033[1m%s\033[0m\n" "${opts[$i]}" >/dev/tty
      else
        printf "\r\033[K     %s\n"               "${opts[$i]}" >/dev/tty
      fi
    done

    local key seq
    IFS= read -rsn1 key </dev/tty
    if [[ "${key}" == $'\x1b' ]]; then
      IFS= read -rsn2 -t 1 seq </dev/tty || true
      key="${key}${seq}"
    fi

    case "${key}" in
      $'\x1b[A') (( cur = (cur - 1 + n) % n )) || true ;;
      $'\x1b[B') (( cur = (cur + 1)     % n )) || true ;;
      ''|$'\n'|$'\r') break ;;
    esac
  done

  while IFS= read -rn1 -t 0 _x </dev/tty 2>/dev/null; do :; done || true

  printf "\n" >/dev/tty
  tput cnorm >/dev/tty 2>/dev/null || true
  [[ -n "$_stty_saved" ]] && stty "${_stty_saved}" </dev/tty 2>/dev/null || true

  eval "${_so_var}=\"\${opts[$cur]}\""
}
# ──────────────────────────────────────────────────────────────────────────────

is_nested() {
  local domain_dir="$1"
  [[ ! -d "${domain_dir}/skills" ]] || [[ -z "$(ls -A "${domain_dir}/skills/" 2>/dev/null)" ]]
}

list_skills() {
  for domain_dir in "${SKILLS_ROOT}"/*/; do
    domain=$(basename "${domain_dir}")
    [[ "${domain}" == .* ]] && continue
    [[ ! -d "${domain_dir}" ]] && continue
    if is_nested "${domain_dir}"; then
      echo "Domain: ${domain} (sub-domains)"
      for sub_dir in "${domain_dir}"/*/; do
        sub=$(basename "${sub_dir}")
        [[ "${sub}" == .* ]] && continue
        [[ ! -d "${sub_dir}/skills" ]] && continue
        echo "  Sub-domain: ${domain}/${sub}"
        for skill_dir in "${sub_dir}/skills"/*/; do
          [[ -f "${skill_dir}/SKILL.md" ]] && echo "    ${domain}/${sub}/$(basename "${skill_dir}")"
        done
      done
    else
      echo "Domain: ${domain} (flat)"
      for skill_dir in "${domain_dir}/skills"/*/; do
        [[ -f "${skill_dir}/SKILL.md" ]] && echo "  ${domain}/$(basename "${skill_dir}")"
      done
    fi
  done
}

install_skill_dir() {
  local src="$1" dest_dir="$2" skill_name
  skill_name=$(basename "${src}")
  mkdir -p "${dest_dir}/${skill_name}"
  cp -r "${src}/." "${dest_dir}/${skill_name}/"
  echo "  installed: ${skill_name} -> ${dest_dir}/${skill_name}" >/dev/tty
  (( _installed_count++ )) || true
}

detect_agents() {
  local detected=()
  [[ -d "${HOME}/.claude" ]] && detected+=("claude")
  [[ -d "${HOME}/.agents" ]] && detected+=("codex")
  [[ -d "${HOME}/.gemini" ]] && detected+=("gemini")
  echo "${detected[@]:-}"
}

do_install() {
  local src="$1" target="$2"
  case "${target}" in
    claude) install_skill_dir "${src}" "${CLAUDE_SKILLS_DIR}" ;;
    codex)  install_skill_dir "${src}" "${AGENTS_SKILLS_DIR}" ;;
    gemini) install_skill_dir "${src}" "${GEMINI_SKILLS_DIR}" ;;
    all)
      install_skill_dir "${src}" "${CLAUDE_SKILLS_DIR}"
      install_skill_dir "${src}" "${AGENTS_SKILLS_DIR}"
      install_skill_dir "${src}" "${GEMINI_SKILLS_DIR}"
      ;;
  esac
}

install_subdomain() {
  local sub_dir="$1" target="$2"
  [[ ! -d "${sub_dir}/skills" ]] && return
  local found=0
  for skill_dir in "${sub_dir}/skills"/*/; do
    [[ -f "${skill_dir}/SKILL.md" ]] || continue
    [[ ${found} -eq 0 ]] && echo "  Installing sub-domain: $(basename "${sub_dir}")" >/dev/tty
    found=1
    do_install "${skill_dir}" "${target}"
  done
}

install_domain() {
  local domain="$1" subdomain="$2" target="$3"
  local domain_dir="${SKILLS_ROOT}/${domain}"
  [[ ! -d "${domain_dir}" ]] && echo "Domain not found: ${domain}" >/dev/tty && exit 1
  echo "Installing domain: ${domain}" >/dev/tty
  if is_nested "${domain_dir}"; then
    if [[ -n "${subdomain}" ]]; then
      install_subdomain "${domain_dir}/${subdomain}" "${target}"
    else
      for sub_dir in "${domain_dir}"/*/; do
        [[ "${sub_dir}" == */.claude-plugin* ]] && continue
        [[ "$(basename "${sub_dir}")" == .* ]] && continue
        install_subdomain "${sub_dir}" "${target}"
      done
    fi
  else
    for skill_dir in "${domain_dir}/skills"/*/; do
      [[ -f "${skill_dir}/SKILL.md" ]] && do_install "${skill_dir}" "${target}"
    done
  fi
}

install_to_agents() {
  local domain="$1" subdomain="$2"
  shift 2
  local agents=("$@")
  for agent in "${agents[@]}"; do
    install_domain "${domain}" "${subdomain}" "${agent}"
  done
}

_in_array() {
  local needle="$1"; shift
  local item
  for item in "$@"; do [[ "${item}" == "${needle}" ]] && return 0; done
  return 1
}

is_skill_installed() {
  local skill_name="$1"
  [[ -d "${CLAUDE_SKILLS_DIR}/${skill_name}" ]] && return 0
  [[ -d "${AGENTS_SKILLS_DIR}/${skill_name}" ]] && return 0
  [[ -d "${GEMINI_SKILLS_DIR}/${skill_name}" ]] && return 0
  return 1
}

is_subdomain_installed() {
  local sub_dir="$1"
  [[ -d "${sub_dir}/skills" ]] || return 1
  for skill_dir in "${sub_dir}/skills"/*/; do
    [[ -f "${skill_dir}/SKILL.md" ]] || continue
    is_skill_installed "$(basename "${skill_dir}")" && return 0
  done
  return 1
}

is_domain_installed() {
  local domain_dir="$1"
  if is_nested "${domain_dir}"; then
    for sub_dir in "${domain_dir}"/*/; do
      [[ "${sub_dir}" == */.claude-plugin* ]] && continue
      [[ "$(basename "${sub_dir}")" == .* ]] && continue
      is_subdomain_installed "${sub_dir}" && return 0
    done
  else
    for skill_dir in "${domain_dir}/skills"/*/; do
      [[ -f "${skill_dir}/SKILL.md" ]] || continue
      is_skill_installed "$(basename "${skill_dir}")" && return 0
    done
  fi
  return 1
}

uninstall_skill_dir() {
  local skill_name="$1" target="$2"
  local dest_dir
  case "${target}" in
    claude) dest_dir="${CLAUDE_SKILLS_DIR}" ;;
    codex)  dest_dir="${AGENTS_SKILLS_DIR}" ;;
    gemini) dest_dir="${GEMINI_SKILLS_DIR}" ;;
    *) return ;;
  esac
  if [[ -d "${dest_dir}/${skill_name}" ]]; then
    rm -rf "${dest_dir}/${skill_name}"
    echo "  uninstalled: ${skill_name} from ${dest_dir}" >/dev/tty
    (( _uninstalled_count++ )) || true
  fi
}

do_uninstall() {
  local skill_name="$1" target="$2"
  case "${target}" in
    claude) uninstall_skill_dir "${skill_name}" "claude" ;;
    codex)  uninstall_skill_dir "${skill_name}" "codex"  ;;
    gemini) uninstall_skill_dir "${skill_name}" "gemini" ;;
    all)
      uninstall_skill_dir "${skill_name}" "claude"
      uninstall_skill_dir "${skill_name}" "codex"
      uninstall_skill_dir "${skill_name}" "gemini"
      ;;
  esac
}

uninstall_subdomain() {
  local sub_dir="$1" target="$2"
  [[ -d "${sub_dir}/skills" ]] || return
  local found=0
  for skill_dir in "${sub_dir}/skills"/*/; do
    [[ -f "${skill_dir}/SKILL.md" ]] || continue
    [[ ${found} -eq 0 ]] && echo "  Uninstalling sub-domain: $(basename "${sub_dir}")" >/dev/tty
    found=1
    do_uninstall "$(basename "${skill_dir}")" "${target}"
  done
}

uninstall_domain() {
  local domain="$1" subdomain="$2" target="$3"
  local domain_dir="${SKILLS_ROOT}/${domain}"
  [[ -d "${domain_dir}" ]] || return
  echo "Uninstalling domain: ${domain}" >/dev/tty
  if is_nested "${domain_dir}"; then
    if [[ -n "${subdomain}" ]]; then
      uninstall_subdomain "${domain_dir}/${subdomain}" "${target}"
    else
      for sub_dir in "${domain_dir}"/*/; do
        [[ "${sub_dir}" == */.claude-plugin* ]] && continue
        [[ "$(basename "${sub_dir}")" == .* ]] && continue
        uninstall_subdomain "${sub_dir}" "${target}"
      done
    fi
  else
    for skill_dir in "${domain_dir}/skills"/*/; do
      [[ -f "${skill_dir}/SKILL.md" ]] && do_uninstall "$(basename "${skill_dir}")" "${target}"
    done
  fi
}

# ── Argument parsing ───────────────────────────────────────────────────────────
DOMAIN="" SUBDOMAIN="" SKILL="" TARGET="" YES=0 UNINSTALL=0

while [[ $# -gt 0 ]]; do
  case "$1" in
    --domain)     DOMAIN="$2";    shift 2 ;;
    --subdomain)  SUBDOMAIN="$2"; shift 2 ;;
    --skill)      SKILL="$2";     shift 2 ;;
    --target)     TARGET="$2";    shift 2 ;;
    --yes|-y)     YES=1;          shift   ;;
    --uninstall)  UNINSTALL=1;    shift   ;;
    --list)       list_skills; exit 0 ;;
    --help)       usage; exit 0 ;;
    *) echo "Unknown option: $1"; usage; exit 1 ;;
  esac
done

# ── Interactive TUI (when no flags given) ──────────────────────────────────────
if [[ -z "${DOMAIN}" && -z "${SKILL}" && -z "${TARGET}" && ${YES} -eq 0 && ${UNINSTALL} -eq 0 ]]; then
  print_banner

  # Mode selection (default: Install)
  select_one _tui_mode "⚙️  What would you like to do?" "📥 Install" "🗑️  Uninstall"

  read -ra _detected <<< "$(detect_agents)"

  # Agent selection
  if [[ ${#_detected[@]} -eq 0 ]]; then
    echo "No agents detected. Defaulting to Claude Code."
    _agent_list=("claude")
  else
    if [[ "${_tui_mode}" == *Uninstall* ]]; then
      multiselect _agent_list "🤖 Which agents to uninstall from?" "${_detected[@]}"
    else
      multiselect _agent_list "🤖 Which agents to install to?" "${_detected[@]}"
    fi
    if [[ ${#_agent_list[@]} -eq 0 ]]; then
      echo "No agents selected. Exiting."
      exit 0
    fi
  fi

  if [[ "${_tui_mode}" == *Uninstall* ]]; then
    # ── Uninstall mode ──────────────────────────────────────────────────────────
    _uninst_domain_opts=()
    _all_uninst_domains=()
    for d in "${SKILLS_ROOT}"/*/; do
      name=$(basename "${d}")
      [[ "${name}" == .* ]] && continue
      [[ -d "${d}" ]] || continue
      _all_uninst_domains+=("${name}")
      _uninst_domain_opts+=("${name}")
    done

    multiselect _domain_to_rm "🗑️  Which domains to uninstall?" "${_uninst_domain_opts[@]}"
    if [[ ${#_domain_to_rm[@]} -eq 0 ]]; then
      echo "Nothing selected. Exiting."
      exit 0
    fi

    # Sub-domain selection for nested domains to uninstall
    for domain in "${_domain_to_rm[@]}"; do
      domain_dir="${SKILLS_ROOT}/${domain}"
      if is_nested "${domain_dir}"; then
        _uninst_sub_opts=()
        for sub_dir in "${domain_dir}"/*/; do
          [[ "${sub_dir}" == */.claude-plugin* ]] && continue
          sub=$(basename "${sub_dir}")
          [[ "${sub}" == .* ]] && continue
          [[ -d "${sub_dir}/skills" ]] || continue
          _uninst_sub_opts+=("+${sub}")
        done
        if [[ ${#_uninst_sub_opts[@]} -gt 0 ]]; then
          multiselect _sub_to_rm "🗑️  ${domain}: which sub-domains to uninstall?" "${_uninst_sub_opts[@]}"
          _varname="_urmsubs_${domain//-/_}"
          eval "${_varname}=(\"\${_sub_to_rm[@]}\")"
        fi
      fi
    done

    _B=$'\033[1m' _D=$'\033[2m' _G=$'\033[38;5;178m' _C=$'\033[0;36m' _R=$'\033[0m'
    printf '%s🤖 Agents%s  %s%s%s\n' "${_B}" "${_R}" "${_C}" "${_agent_list[*]}" "${_R}" >/dev/tty
    printf '%s🗑️  Mode%s    %sUninstall%s\n' "${_B}" "${_R}" "${_G}" "${_R}" >/dev/tty
    for domain in "${_domain_to_rm[@]}"; do
      _varname="_urmsubs_${domain//-/_}"
      eval "_preview_subs=(\"\${${_varname}[@]}\")"
      if [[ ${#_preview_subs[@]} -gt 0 ]]; then
        printf '%s📂 Domain%s  %s%s%s %s[%s]%s\n' "${_B}" "${_R}" "${_G}" "${domain}" "${_R}" "${_D}" "${_preview_subs[*]}" "${_R}" >/dev/tty
      else
        printf '%s📂 Domain%s  %s%s%s\n' "${_B}" "${_R}" "${_G}" "${domain}" "${_R}" >/dev/tty
      fi
    done
    printf '\n' >/dev/tty
    echo "Uninstalling..." >/dev/tty
    _uninstalled_count=0
    for domain in "${_domain_to_rm[@]}"; do
      domain_dir="${SKILLS_ROOT}/${domain}"
      if is_nested "${domain_dir}"; then
        _varname="_urmsubs_${domain//-/_}"
        eval "_usubs=(\"\${${_varname}[@]}\")"
        if [[ ${#_usubs[@]} -gt 0 ]]; then
          for sub in "${_usubs[@]}"; do
            for agent in "${_agent_list[@]}"; do
              uninstall_subdomain "${domain_dir}/${sub}" "${agent}"
            done
          done
        else
          for agent in "${_agent_list[@]}"; do
            uninstall_domain "${domain}" "" "${agent}"
          done
        fi
      else
        for agent in "${_agent_list[@]}"; do
          uninstall_domain "${domain}" "" "${agent}"
        done
      fi
    done
    _skill_count_unique=$(( _uninstalled_count / ${#_agent_list[@]} ))
    printf '\n' >/dev/tty
    printf '🗑️  %d skills uninstalled' "${_skill_count_unique}" >/dev/tty
    [[ ${#_agent_list[@]} -gt 1 ]] && printf ' × %d agents (%d total)' "${#_agent_list[@]}" "${_uninstalled_count}" >/dev/tty
    printf ' → %s\n' "${_agent_list[*]}" >/dev/tty

  else
    # ── Install mode ─────────────────────────────────────────────────────────────
    _all_domains=()
    for d in "${SKILLS_ROOT}"/*/; do
      name=$(basename "${d}")
      [[ "${name}" == .* ]] && continue
      [[ -d "${d}" ]] || continue
      _all_domains+=("${name}")
    done

    multiselect _domain_list "📚 Which domains to install?" "${_all_domains[@]}"
    if [[ ${#_domain_list[@]} -eq 0 ]]; then
      echo "No domains selected. Exiting."
      exit 0
    fi

    # Sub-domain selection for nested domains
    for domain in "${_domain_list[@]}"; do
      domain_dir="${SKILLS_ROOT}/${domain}"
      if is_nested "${domain_dir}"; then
        _subs=()
        for sub_dir in "${domain_dir}"/*/; do
          [[ "${sub_dir}" == */.claude-plugin* ]] && continue
          sub=$(basename "${sub_dir}")
          [[ "${sub}" == .* ]] && continue
          [[ -d "${sub_dir}/skills" ]] || continue
          _subs+=("+${sub}")
        done
        if [[ ${#_subs[@]} -gt 0 ]]; then
          multiselect _sub_list "📂 ${domain}: which sub-domains?" "${_subs[@]}"
          _varname="_dsubs_${domain//-/_}"
          eval "${_varname}=(\"\${_sub_list[@]}\")"
        fi
      fi
    done

    _B=$'\033[1m' _D=$'\033[2m' _G=$'\033[38;5;178m' _C=$'\033[0;36m' _R=$'\033[0m'
    printf '%s🤖 Agents%s  %s%s%s\n' "${_B}" "${_R}" "${_C}" "${_agent_list[*]}" "${_R}" >/dev/tty
    printf '%s📥 Mode%s    %sInstall%s\n' "${_B}" "${_R}" "${_G}" "${_R}" >/dev/tty
    for domain in "${_domain_list[@]}"; do
      _varname="_dsubs_${domain//-/_}"
      eval "_preview_subs=(\"\${${_varname}[@]}\")"
      if [[ ${#_preview_subs[@]} -gt 0 ]]; then
        printf '%s📂 Domain%s  %s%s%s %s[%s]%s\n' "${_B}" "${_R}" "${_G}" "${domain}" "${_R}" "${_D}" "${_preview_subs[*]}" "${_R}" >/dev/tty
      else
        printf '%s📂 Domain%s  %s%s%s\n' "${_B}" "${_R}" "${_G}" "${domain}" "${_R}" >/dev/tty
      fi
    done
    printf '\n' >/dev/tty
    echo "Installing..." >/dev/tty
    _installed_count=0
    for domain in "${_domain_list[@]}"; do
      domain_dir="${SKILLS_ROOT}/${domain}"
      if is_nested "${domain_dir}"; then
        _varname="_dsubs_${domain//-/_}"
        eval "_stored=(\"\${${_varname}[@]}\")"
        if [[ ${#_stored[@]} -gt 0 ]]; then
          for sub in "${_stored[@]}"; do
            for agent in "${_agent_list[@]}"; do
              install_domain "${domain}" "${sub}" "${agent}"
            done
          done
        else
          echo "  No sub-domains selected for ${domain}, skipping." >/dev/tty
        fi
      else
        for agent in "${_agent_list[@]}"; do
          install_domain "${domain}" "" "${agent}"
        done
      fi
    done
    _skill_count_unique=$(( _installed_count / ${#_agent_list[@]} ))
    printf '\n' >/dev/tty
    printf '✅ %d skills installed' "${_skill_count_unique}" >/dev/tty
    [[ ${#_agent_list[@]} -gt 1 ]] && printf ' × %d agents (%d total)' "${#_agent_list[@]}" "${_installed_count}" >/dev/tty
    printf ' → %s\n' "${_agent_list[*]}" >/dev/tty
  fi

# ── Non-interactive / flag-driven ─────────────────────────────────────────────
else
  if [[ -z "${TARGET}" || "${TARGET}" == "auto" ]]; then
    read -ra _detected <<< "$(detect_agents)"
    if [[ ${#_detected[@]} -eq 0 ]]; then
      echo "No agents detected. Defaulting to Claude Code."
      TARGET="claude"
    else
      if [[ ${UNINSTALL} -eq 1 ]]; then
        echo "Uninstalling from: ${_detected[*]}"
      else
        echo "Installing to: ${_detected[*]}"
      fi
      TARGET="auto"
    fi
  fi

  if [[ ${UNINSTALL} -eq 1 ]]; then
    if [[ -n "${SKILL}" ]]; then
      IFS='/' read -ra parts <<< "${SKILL}"
      if [[ ${#parts[@]} -eq 3 ]]; then
        skill_path="${SKILLS_ROOT}/${parts[0]}/${parts[1]}/skills/${parts[2]}"
      elif [[ ${#parts[@]} -eq 2 ]]; then
        skill_path="${SKILLS_ROOT}/${parts[0]}/skills/${parts[1]}"
      else
        echo "Invalid skill path: ${SKILL}"; exit 1
      fi
      skill_name="$(basename "${skill_path}")"
      echo "Uninstalling skill: ${SKILL}"
      if [[ "${TARGET}" == "auto" ]]; then
        for agent in "${_detected[@]}"; do do_uninstall "${skill_name}" "${agent}"; done
      else
        do_uninstall "${skill_name}" "${TARGET}"
      fi
    elif [[ -n "${DOMAIN}" ]]; then
      if [[ "${TARGET}" == "auto" ]]; then
        for agent in "${_detected[@]}"; do uninstall_domain "${DOMAIN}" "${SUBDOMAIN}" "${agent}"; done
      else
        uninstall_domain "${DOMAIN}" "${SUBDOMAIN}" "${TARGET}"
      fi
    else
      for domain_dir in "${SKILLS_ROOT}"/*/; do
        domain=$(basename "${domain_dir}")
        [[ "${domain}" == .* ]] && continue
        [[ ! -d "${domain_dir}" ]] && continue
        if [[ "${TARGET}" == "auto" ]]; then
          for agent in "${_detected[@]}"; do uninstall_domain "${domain}" "" "${agent}"; done
        else
          uninstall_domain "${domain}" "" "${TARGET}"
        fi
      done
    fi
  else
    if [[ -n "${SKILL}" ]]; then
      IFS='/' read -ra parts <<< "${SKILL}"
      if [[ ${#parts[@]} -eq 3 ]]; then
        skill_path="${SKILLS_ROOT}/${parts[0]}/${parts[1]}/skills/${parts[2]}"
      elif [[ ${#parts[@]} -eq 2 ]]; then
        skill_path="${SKILLS_ROOT}/${parts[0]}/skills/${parts[1]}"
      else
        echo "Invalid skill path: ${SKILL}"; exit 1
      fi
      [[ ! -f "${skill_path}/SKILL.md" ]] && echo "Skill not found: ${SKILL}" && exit 1
      echo "Installing skill: ${SKILL}"
      if [[ "${TARGET}" == "auto" ]]; then
        for agent in "${_detected[@]}"; do do_install "${skill_path}" "${agent}"; done
      else
        do_install "${skill_path}" "${TARGET}"
      fi
    elif [[ -n "${DOMAIN}" ]]; then
      if [[ "${TARGET}" == "auto" ]]; then
        for agent in "${_detected[@]}"; do install_domain "${DOMAIN}" "${SUBDOMAIN}" "${agent}"; done
      else
        install_domain "${DOMAIN}" "${SUBDOMAIN}" "${TARGET}"
      fi
    else
      for domain_dir in "${SKILLS_ROOT}"/*/; do
        domain=$(basename "${domain_dir}")
        [[ "${domain}" == .* ]] && continue
        [[ ! -d "${domain_dir}" ]] && continue
        if [[ "${TARGET}" == "auto" ]]; then
          for agent in "${_detected[@]}"; do install_domain "${domain}" "" "${agent}"; done
        else
          install_domain "${domain}" "" "${TARGET}"
        fi
      done
    fi
  fi
fi

_c_bold=$'\033[1m'
_c_dim=$'\033[2m'
_c_gold=$'\033[38;5;178m'
_c_cyan=$'\033[0;36m'
_c_reset=$'\033[0m'
printf '\n'
printf '%s💡 Also available via marketplace:%s\n' "${_c_bold}" "${_c_reset}"
printf '   %s🖱️  Cursor%s   %s/add-plugin grimoire%s  %s(in Agent chat)%s\n' \
  "${_c_gold}" "${_c_reset}" "${_c_cyan}" "${_c_reset}" "${_c_dim}" "${_c_reset}"
printf '   %s🐙 Copilot%s  %scopilot plugin install grimoire%s\n' \
  "${_c_gold}" "${_c_reset}" "${_c_cyan}" "${_c_reset}"
printf '   %s📖 Docs%s     %shttps://github.com/jeffreytse/grimoire#-install%s\n' \
  "${_c_gold}" "${_c_reset}" "${_c_dim}" "${_c_reset}"
printf '\n'
