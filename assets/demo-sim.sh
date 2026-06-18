#!/usr/bin/env bash
# Simulated Claude Code TUI demo — analyze-best-practice-problem → plan → 4 skills → outcome

ORANGE='\033[38;5;174m'
AMBER='\033[38;5;172m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BOLD='\033[1m'
DIM='\033[2m'
NC='\033[0m'

DLINE='────────────────────────────────────────────────────────────────────────────────'
QUESTION="We're launching our SaaS in 48 hours. I'm terrified something will break. What do we do?"
_WRAPPED=0

pause() { sleep "$1"; }

type_cmd() {
  local cmd="$1" i char
  printf "${BOLD}❯${NC} "
  for (( i=0; i<${#cmd}; i++ )); do
    char="${cmd:$i:1}"; printf "%s" "$char"; sleep 0.055
  done
}

type_at_bottom() {
  local msg="$1"
  local term_h; term_h=$(tput lines)
  local r=$(( term_h - 2 ))
  printf "\0337"  # save cursor (content position)
  # Redraw full box then type into input row
  printf "\033[${r};1H${DIM}${DLINE}${NC}"
  printf "\033[$(( r + 2 ));1H${DIM}${DLINE}${NC}"
  printf "\033[$(( r + 1 ));1H\033[2K${BOLD}❯${NC} "
  local i char
  for (( i=0; i<${#msg}; i++ )); do
    char="${msg:$i:1}"; printf "%s" "$char"; sleep 0.055
  done
  pause 0.4
  # Clear bottom input box (simulate send)
  printf "\033[$(( r + 1 ));1H\033[2K${DIM}❯${NC} "
  pause 0.2
  # Restore cursor, print message in history pane
  printf "\0338"
  printf "${DIM}❯ %s${NC}\n\n" "$msg"
}

repin_box() {
  printf "\0337"
  local term_h; term_h=$(tput lines)
  local r=$(( term_h - 2 ))
  printf "\033[${r};1H${DIM}${DLINE}${NC}"
  printf "\033[$(( r + 1 ));1H${DIM}❯${NC} "
  printf "\033[$(( r + 2 ));1H${DIM}${DLINE}${NC}"
  printf "\0338"
}

header() {
  printf "\033[r"
  printf " ${ORANGE}▐▛███▜▌${NC}   ${BOLD}Claude Code${NC} ${DIM}v2.1.137${NC}\n"
  printf "${ORANGE}▝▜█████▛▘${NC}  ${DIM}Sonnet 4.6 · grimoire${NC}\n"
  printf "  ${ORANGE}▘▘ ▝▝${NC}    ${DIM}~/Projects/grimoire${NC}\n"
  echo ""
}

think() {
  local word="$1" secs="$2"
  printf "${AMBER}✳${NC} ${word}… ${DIM}(thinking)${NC}\n"
  pause "$secs"
  printf "${DIM}✻ Churned for ${secs}s${NC}\n"
  pause 0.3
}

prompt_at_bottom() {
  local cmd="$1"
  local term_h; term_h=$(tput lines)
  local max_line=$(( ${#DLINE} - 3 ))  # chars per row after "❯ " prefix

  # Start with 3-row box (1 input line)
  local r=$(( term_h - 2 ))
  printf "\033[${r};1H${DIM}${DLINE}${NC}"
  printf "\033[$(( r + 2 ));1H${DIM}${DLINE}${NC}"
  printf "\033[$(( r + 1 ));1H\033[2K${BOLD}❯${NC} "

  local i char
  if [[ ${#cmd} -le $max_line ]]; then
    _WRAPPED=0
    for (( i=0; i<${#cmd}; i++ )); do
      char="${cmd:$i:1}"; printf "%s" "$char"; sleep 0.055
    done
  else
    _WRAPPED=1
    local chunk="${cmd:0:$max_line}"
    local wp=${#chunk}
    while [[ $wp -gt 0 && "${cmd:$(( wp - 1 )):1}" != " " ]]; do (( wp-- )); done
    [[ $wp -eq 0 ]] && wp=$max_line
    local line1="${cmd:0:$wp}"; line1="${line1% }"
    local line2="${cmd:$wp}";   line2="${line2# }"

    # type line1 on row r+1
    for (( i=0; i<${#line1}; i++ )); do
      char="${line1:$i:1}"; printf "%s" "$char"; sleep 0.055
    done
    pause 0.1

    # expand UP: move line1 from r+1 → r, shift top DLINE to r-1, type line2 on r+1
    printf "\033[${r};1H\033[2K${BOLD}❯${NC} ${line1}"   # redraw line1 on row r
    printf "\033[$(( r - 1 ));1H${DIM}${DLINE}${NC}"     # new top DLINE at r-1
    printf "\033[$(( r + 1 ));1H\033[2K  "               # clear r+1, indent — cursor lands here
    # bottom DLINE at r+2 stays — within terminal bounds

    # type line2 on r+1 (cursor already positioned)
    for (( i=0; i<${#line2}; i++ )); do
      char="${line2:$i:1}"; printf "%s" "$char"; sleep 0.055
    done
  fi
}

send_animation() {
  local term_h; term_h=$(tput lines)
  local r=$(( term_h - 2 ))
  if [[ $_WRAPPED -eq 1 ]]; then
    printf "\033[${r};1H\033[2K${DIM}❯${NC} "   # clear line1
    printf "\033[$(( r + 1 ));1H\033[2K"         # clear line2
    _WRAPPED=0
  else
    printf "\033[$(( r + 1 ));1H\033[2K${DIM}❯${NC} "
  fi
}

pin_prompt() {
  local cmd="$1"
  local term_h; term_h=$(tput lines)
  local r=$(( term_h - 2 ))
  printf "\033[${r};1H${DIM}${DLINE}${NC}"
  printf "\033[$(( r + 1 ));1H${BOLD}❯${NC} ${cmd}"
  printf "\033[$(( r + 2 ));1H${DIM}${DLINE}${NC}"
  printf "\033[5;1H"
}

pin_empty() {
  local term_h; term_h=$(tput lines)
  local r=$(( term_h - 2 ))
  printf "\033[${r};1H${DIM}${DLINE}${NC}"
  printf "\033[$(( r + 1 ));1H${DIM}❯${NC} "
  printf "\033[$(( r + 2 ));1H${DIM}${DLINE}${NC}"
  printf "\033[5;$(( r - 1 ))r"
  printf "\033[5;1H"
}

history_line() {
  local msg="$1"
  local max=$(( ${#DLINE} - 3 ))
  if [[ ${#msg} -le $max ]]; then
    printf "${DIM}❯ %s${NC}\n\n" "$msg"
  else
    local chunk="${msg:0:$max}"
    local wp=${#chunk}
    while [[ $wp -gt 0 && "${msg:$(( wp - 1 )):1}" != " " ]]; do (( wp-- )); done
    [[ $wp -eq 0 ]] && wp=$max
    local line1="${msg:0:$wp}"; line1="${line1% }"
    local line2="${msg:$wp}";   line2="${line2# }"
    printf "${DIM}❯ %s\n  %s${NC}\n\n" "$line1" "$line2"
  fi
}

analyze() {
  local ans1="First time at this scale. Maybe 500 concurrent users at launch."

  clear; header
  pin_empty
  history_line "$QUESTION"
  pause 0.4
  think "Routing" 2
  echo ""

  printf "${AMBER}⏺${NC} ${BOLD}analyze-best-practice-problem${NC}\n\n"; pause 0.5
  printf "  Have you shipped at this scale before, and what's\n"
  printf "  your expected peak concurrent user count?\n\n"
  pause 1.8
  type_at_bottom "$ans1"

  think "Scoping" 1.5
  echo ""
  printf "${AMBER}⏺${NC} ${BOLD}Problem statement${NC}\n\n"; pause 0.4
  printf "  ${AMBER}⎿${NC}  Situation   First production launch, 48h runway, ~500 concurrent users\n"; pause 0.35
  printf "  ${AMBER}⎿${NC}  Goal        Launch with confidence — know what breaks before users do\n";  pause 0.35
  printf "  ${AMBER}⎿${NC}  Root cause  No pre-launch protocol in place\n"
  pause 2.2
}

plan_route() {
  echo ""; pause 0.4
  think "Matching skills" 2
  echo ""
  printf "${AMBER}⏺${NC} ${BOLD}plan-best-practice-solution${NC}  ${DIM}— 4 skills detected${NC}\n\n"; pause 0.5
  printf "  ${AMBER}⎿${NC}  Step 1  ${BOLD}apply-premortem${NC}           ${DIM}business/strategy${NC}        ${GREEN}✓${NC}\n"; pause 0.4
  printf "  ${AMBER}⎿${NC}  Step 2  ${BOLD}design-slo${NC}                ${DIM}engineering/reliability${NC}  ${GREEN}✓${NC}\n"; pause 0.4
  printf "  ${AMBER}⎿${NC}  Step 3  ${BOLD}plan-incident-response${NC}    ${DIM}engineering/devops${NC}       ${GREEN}✓${NC}\n"; pause 0.4
  printf "  ${AMBER}⎿${NC}  Step 4  ${BOLD}run-game-day${NC}              ${DIM}engineering/reliability${NC}  ${GREEN}✓${NC}\n"; pause 0.5
  echo ""
  printf "${AMBER}⏺${NC} ${BOLD}Plan confirmed${NC}\n\n"
  printf "  ${AMBER}⎿${NC}  apply-premortem  ${DIM}→${NC}  design-slo  ${DIM}→${NC}  plan-incident-response  ${DIM}→${NC}  run-game-day\n"
  echo ""; pause 0.5
  printf "  Apply step 1 now?\n"
  pause 1.5
  type_at_bottom "yes"
  pause 0.4
}

skill_premortem() {
  echo ""; pause 0.4
  think "Running apply-premortem" 2
  echo ""
  printf "${AMBER}⏺${NC} ${BOLD}apply-premortem${NC}  ${DIM}business/strategy${NC}\n\n"; pause 0.4
  printf "  Imagine it's 72 hours post-launch and it failed. What happened?\n\n"; pause 0.6
  printf "  ${AMBER}⎿${NC}  ${RED}HIGH${NC}    DB connection exhaustion under concurrent signups\n"; pause 0.35
  printf "  ${AMBER}⎿${NC}  ${RED}HIGH${NC}    Payment webhook timeouts causing silent failures\n"; pause 0.35
  printf "  ${AMBER}⎿${NC}  ${YELLOW}MEDIUM${NC}  Email rate limits hit during onboarding burst\n"; pause 0.35
  printf "  ${AMBER}⎿${NC}  ${DIM}LOW${NC}     CDN cache poisoning serving stale auth tokens\n"
  pause 0.6; echo ""
  printf "  Fix the HIGH items before launch. Not after.\n"
  printf "  ${DIM}Source: Klein, Sources of Power (1998) — prospective hindsight${NC}\n"
  pause 0.7; echo ""
  printf "  Step 1 done. Continue to step 2 ${DIM}(design-slo)${NC}?\n"
  pause 1.5
  type_at_bottom "yes"
  pause 0.4
}

skill_slo() {
  echo ""; pause 0.4
  think "Running design-slo" 2
  echo ""
  printf "${AMBER}⏺${NC} ${BOLD}design-slo${NC}  ${DIM}engineering/reliability${NC}\n\n"; pause 0.4
  printf "  Define \"good enough\" before launch — or you won't know when you've failed.\n\n"; pause 0.6
  printf "  ${AMBER}⎿${NC}  API p99 latency  < 500ms        Error budget: 43 min/month\n"; pause 0.35
  printf "  ${AMBER}⎿${NC}  Signup flow      99.5%% success   Error budget: 3.6 hr degraded\n"; pause 0.35
  printf "  ${AMBER}⎿${NC}  Payment flow     99.9%% success   Error budget: 43 min degraded\n"
  pause 0.6; echo ""
  printf "  ${YELLOW}⚠${NC}  Without this: \"the site feels slow\" is not actionable at 2am.\n"
  printf "  ${DIM}Source: Google SRE Book (2016) ch.4${NC}\n"
  pause 0.7; echo ""
  printf "  Step 2 done. Continue to step 3 ${DIM}(plan-incident-response)${NC}?\n"
  pause 1.5
  type_at_bottom "yes"
  pause 0.4
}

skill_incident() {
  echo ""; pause 0.4
  think "Running plan-incident-response" 2
  echo ""
  printf "${AMBER}⏺${NC} ${BOLD}plan-incident-response${NC}  ${DIM}engineering/devops${NC}\n\n"; pause 0.4
  printf "  Define roles before the incident — not during it.\n\n"; pause 0.6
  printf "  ${AMBER}⎿${NC}  Incident Lead  [you]            Coordinates, owns status page\n"; pause 0.35
  printf "  ${AMBER}⎿${NC}  Tech Lead      [senior eng]     Diagnoses and fixes\n"; pause 0.35
  printf "  ${AMBER}⎿${NC}  Comms          [founder/PM]     Emails customers, updates investors\n"
  pause 0.6; echo ""
  printf "  ${YELLOW}⚠${NC}  Without this: everyone stares at the same terminal.\n"
  printf "  ${DIM}Source: PagerDuty Incident Response Guide (2020)${NC}\n"
  pause 0.7; echo ""
  printf "  Step 3 done. Continue to step 4 ${DIM}(run-game-day)${NC}?\n"
  pause 1.5
  type_at_bottom "yes"
  pause 0.4
}

skill_gameday() {
  echo ""; pause 0.4
  think "Running run-game-day" 3
  echo ""
  printf "${AMBER}⏺${NC} ${BOLD}run-game-day${NC}  ${DIM}engineering/reliability${NC}\n\n"; pause 0.4
  printf "  Break it before your users do.\n\n"; pause 0.6
  printf "  ${AMBER}⎿${NC}  T-48h  Kill the database — measure recovery time\n"; pause 0.35
  printf "  ${AMBER}⎿${NC}  T-48h  Flood API to 10× expected load — observe behavior\n"; pause 0.35
  printf "  ${AMBER}⎿${NC}  T-24h  Simulate payment webhook timeout — verify order state\n"
  pause 0.6; echo ""
  printf "  ${GREEN}✓${NC}  If game day breaks staging, it saved your launch.\n"
  printf "  ${DIM}Source: Netflix Chaos Engineering (2011)${NC}\n"
  pause 2.5
}

final_outcome() {
  repin_box
  echo ""
  printf "${DIM}${DLINE}${NC}\n\n"
  printf "${AMBER}⏺${NC} ${BOLD}Pre-launch protocol complete${NC}\n\n"; pause 0.4
  printf "  ${GREEN}⎿  Risks mapped${NC}         DB pool + webhook — fix today, not post-launch\n";   pause 0.4
  printf "  ${GREEN}⎿  SLOs defined${NC}         you'll know within 5 min if launch is failing\n";    pause 0.4
  printf "  ${GREEN}⎿  Roles assigned${NC}       no one freezes at 2am — everyone has a job\n";      pause 0.4
  printf "  ${GREEN}⎿  Game day scheduled${NC}   break it in staging today, not in prod tomorrow\n"
  pause 0.8; echo ""
  printf "  You don't need a bigger team. You need the same protocol ${GREEN}${BOLD}Google uses${NC}.  ${AMBER}✓${NC}\n"
  repin_box
  pause 4.0
}

clear
header
pause 1.2
prompt_at_bottom "$QUESTION"
pause 1.5
send_animation
pause 0.3
analyze
plan_route
skill_premortem
skill_slo
skill_incident
skill_gameday
final_outcome
sleep 9999
