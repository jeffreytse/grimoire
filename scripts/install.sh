#!/usr/bin/env bash
set -uo pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
SKILLS_ROOT="${REPO_ROOT}/skills"
CLAUDE_SKILLS_DIR="${HOME}/.claude/skills"
AGENTS_SKILLS_DIR="${HOME}/.agents/skills"
GEMINI_SKILLS_DIR="${HOME}/.gemini/skills"

usage() {
  cat <<EOF
Usage: install.sh [OPTIONS]

Options:
  --domain <name>       Install all skills for a domain (e.g. engineering, photography)
  --subdomain <name>    Restrict to one sub-domain within a domain (e.g. development)
  --skill <path>        Install one skill (e.g. engineering/development/propose-conventional-commit)
  --target <agent>      Target: claude (default), codex, gemini, all
  --list                List available domains, sub-domains, and skills
  --help                Show this help

Examples:
  install.sh                                                # All skills, Claude Code
  install.sh --domain engineering                           # All engineering sub-domains
  install.sh --domain engineering --subdomain development   # One sub-domain
  install.sh --skill engineering/development/propose-conventional-commit
  install.sh --domain engineering --target all              # All agents
EOF
}

# Detect if a domain dir is flat (has skills/) or nested (has sub-domain dirs)
is_nested() {
  local domain_dir="$1"
  [[ ! -d "${domain_dir}/skills" ]]
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
  local src="$1"
  local dest_dir="$2"
  local skill_name
  skill_name=$(basename "${src}")
  mkdir -p "${dest_dir}/${skill_name}"
  cp -r "${src}/." "${dest_dir}/${skill_name}/"
  echo "  installed: ${skill_name} -> ${dest_dir}/${skill_name}"
}

do_install() {
  local src="$1"
  local target="$2"
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
  local sub_dir="$1"
  local target="$2"
  [[ ! -d "${sub_dir}/skills" ]] && return
  local found=0
  for skill_dir in "${sub_dir}/skills"/*/; do
    [[ -f "${skill_dir}/SKILL.md" ]] || continue
    [[ ${found} -eq 0 ]] && echo "  Installing sub-domain: $(basename "${sub_dir}")"
    found=1
    do_install "${skill_dir}" "${target}"
  done
}

install_domain() {
  local domain="$1"
  local subdomain="$2"
  local target="$3"
  local domain_dir="${SKILLS_ROOT}/${domain}"

  [[ ! -d "${domain_dir}" ]] && echo "Domain not found: ${domain}" && exit 1

  echo "Installing domain: ${domain}"
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

# Defaults
DOMAIN=""
SUBDOMAIN=""
SKILL=""
TARGET="claude"

while [[ $# -gt 0 ]]; do
  case "$1" in
    --domain)    DOMAIN="$2";    shift 2 ;;
    --subdomain) SUBDOMAIN="$2"; shift 2 ;;
    --skill)     SKILL="$2";     shift 2 ;;
    --target)    TARGET="$2";    shift 2 ;;
    --list)      list_skills; exit 0 ;;
    --help)      usage; exit 0 ;;
    *) echo "Unknown option: $1"; usage; exit 1 ;;
  esac
done

if [[ -n "${SKILL}" ]]; then
  # Format: domain/subdomain/skill-name OR domain/skill-name
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
  do_install "${skill_path}" "${TARGET}"
elif [[ -n "${DOMAIN}" ]]; then
  install_domain "${DOMAIN}" "${SUBDOMAIN}" "${TARGET}"
else
  for domain_dir in "${SKILLS_ROOT}"/*/; do
    domain=$(basename "${domain_dir}")
    [[ "${domain}" == .* ]] && continue
    [[ -d "${domain_dir}" ]] && install_domain "${domain}" "" "${TARGET}"
  done
fi

echo "Done."
