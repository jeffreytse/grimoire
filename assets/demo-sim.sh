#!/usr/bin/env bash
# Simulated Claude Code TUI demo — analyze-problem → plan → 3 skills → outcome

ORANGE='\033[38;5;174m'
AMBER='\033[38;5;172m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BOLD='\033[1m'
DIM='\033[2m'
NC='\033[0m'

DLINE='────────────────────────────────────────────────────────────────────────────────'
QUESTION="I'm 42, AI just took my job, I have a mortgage. What do I do?"

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
  local r=$(( term_h - 2 ))
  printf "\033[${r};1H${DIM}${DLINE}${NC}"
  printf "\033[$(( r + 2 ));1H${DIM}${DLINE}${NC}"
  printf "\033[$(( r + 1 ));1H"
  type_cmd "$cmd"
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
  printf "${DIM}❯ %s${NC}\n\n" "$1"
}

analyze() {
  local ans1="Both. I need income fast but I also can't lose the house."
  local ans2="Still current — I have maybe 3 months of savings."

  clear; header
  pin_empty
  history_line "$QUESTION"
  pause 0.4
  think "Routing" 2
  echo ""

  # Q1 — accumulates on screen, user types at bottom
  printf "${AMBER}⏺${NC} ${BOLD}analyze-problem${NC}\n\n"; pause 0.5
  printf "  What outcome are you trying to achieve — stabilize finances,\n"
  printf "  find new income, or both at once?\n\n"
  pause 1.8
  type_at_bottom "$ans1"

  # Q2 — continues below Q1 + user message in history
  printf "${AMBER}⏺${NC} ${BOLD}analyze-problem${NC}\n\n"; pause 0.4
  printf "  Is the mortgage currently in arrears, or are you still\n"
  printf "  current with payments?\n\n"
  pause 1.8
  type_at_bottom "$ans2"

  # Scoping — continues below Q2 + user message in history
  think "Scoping" 1.5
  echo ""
  printf "${AMBER}⏺${NC} ${BOLD}Problem statement${NC}\n\n"; pause 0.4
  printf "  ${AMBER}⎿${NC}  Situation   Sudden job loss at 42, 3-month cash runway, active mortgage\n"; pause 0.35
  printf "  ${AMBER}⎿${NC}  Goal        Secure income and protect housing simultaneously\n";            pause 0.35
  printf "  ${AMBER}⎿${NC}  Root cause  External (AI displacement), not performance\n"
  pause 2.2
}

plan_route() {
  echo ""; pause 0.4
  think "Matching skills" 2
  echo ""
  printf "${AMBER}⏺${NC} ${BOLD}plan-best-practice-solution${NC}  ${DIM}— 3 domains detected${NC}\n\n"; pause 0.5
  printf "  ${AMBER}⎿${NC}  Step 1  ${BOLD}design-budget${NC}            ${DIM}finance/personal-finance${NC}  ${GREEN}✓${NC}\n"; pause 0.4
  printf "  ${AMBER}⎿${NC}  Step 2  ${BOLD}design-debt-payoff-plan${NC}  ${DIM}finance/personal-finance${NC}  ${GREEN}✓${NC}\n"; pause 0.4
  printf "  ${AMBER}⎿${NC}  Step 3  Career path — ${YELLOW}multiple practices apply${NC}\n"
  pause 0.5; echo ""
  repin_box
  printf "\0337"
  printf "  ${AMBER}?${NC} ${BOLD}Which practice for step 3?${NC}\n\n"
  printf "  ${AMBER}❯${NC} ${AMBER}★${NC} ${BOLD}run-scenario-planning${NC}    ${DIM}map career paths by speed to income  ← recommended${NC}\n"
  printf "        ${DIM}design-go-to-market      build a freelance/consulting pipeline${NC}\n"
  pause 2.2
  printf "\0338\033[J"
  printf "  ${GREEN}✔${NC} Step 3  ${BOLD}run-scenario-planning${NC}  ${GREEN}← selected${NC}\n"
  printf "\0337"
  local _th; _th=$(tput lines)
  printf "\033[$(( _th - 2 ));1H${DIM}${DLINE}${NC}"
  printf "\033[$(( _th - 1 ));1H${DIM}❯${NC} "
  printf "\033[${_th};1H${DIM}${DLINE}${NC}"
  printf "\0338"
  pause 0.6; echo ""
  printf "${AMBER}⏺${NC} ${BOLD}Plan confirmed${NC}\n\n"
  printf "  ${AMBER}⎿${NC}  design-budget  ${DIM}→${NC}  design-debt-payoff-plan  ${DIM}→${NC}  run-scenario-planning\n"
  echo ""; pause 0.5
  printf "  Apply step 1 now?\n"
  pause 1.5
  type_at_bottom "yes"
  pause 0.4
}

skill_budget() {
  echo ""; pause 0.4
  think "Running design-budget" 3
  echo ""
  printf "${AMBER}⏺${NC} ${BOLD}design-budget${NC}  ${DIM}finance/personal-finance${NC}\n\n"; pause 0.4
  printf "  ${AMBER}⎿${NC}  Mortgage     \$2,100  ${DIM}·${NC}  Utilities  \$280  ${DIM}·${NC}  Insurance  \$420\n"; pause 0.35
  printf "  ${AMBER}⎿${NC}  Fixed total  ${BOLD}\$2,895/mo${NC}\n"
  pause 0.5; echo ""
  printf "  ${AMBER}⎿${NC}  Cut dining   \$600 ${DIM}→${NC} \$150  ${DIM}·${NC}  Streaming \$120 ${DIM}→${NC} \$40  ${DIM}·${NC}  Gym ${DIM}→${NC} pause\n"; pause 0.35
  printf "  ${AMBER}⎿${NC}  Saves        \$610/mo\n"
  pause 0.5; echo ""
  printf "  ${AMBER}⎿${NC}  Savings  \$18,000  ${DIM}·${NC}  Burn  \$2,285/mo  ${DIM}·${NC}  Runway  ${YELLOW}7.9 months${NC}\n"
  pause 0.5; echo ""
  printf "  ${YELLOW}⚠${NC}  Target 12-month runway. Gap: \$9,370 — forbearance closes this in step 2.\n"
  pause 0.7; echo ""
  printf "  Step 1 done. Continue to step 2 ${DIM}(design-debt-payoff-plan)${NC}?\n"
  pause 1.5
  type_at_bottom "yes"
  pause 0.4
}

skill_debt() {
  echo ""; pause 0.4
  think "Running design-debt-payoff-plan" 2
  echo ""
  printf "${AMBER}⏺${NC} ${BOLD}design-debt-payoff-plan${NC}  ${DIM}finance/personal-finance${NC}\n\n"; pause 0.4
  printf "  ${AMBER}⎿${NC}  1. Forbearance  pause 3–12 months, no credit penalty  ${GREEN}← do this now${NC}\n"; pause 0.35
  printf "  ${AMBER}⎿${NC}  2. Refinance    lower rate if equity > 20%%  ${DIM}— 30–45 days${NC}\n";              pause 0.35
  printf "  ${AMBER}⎿${NC}  3. HELOC        equity to credit line  ${DIM}— last resort only${NC}\n"
  pause 0.5; echo ""
  printf "  Call lender today. Cite involuntary job loss.\n"
  printf "  Federal law requires forbearance offer. Get it in writing.\n"
  pause 0.5; echo ""
  printf "  ${AMBER}⎿${NC}  Runway after forbearance  ${GREEN}7.9 + 5.5 = 13.4 months${NC}\n"
  pause 0.7; echo ""
  printf "  Step 2 done. Continue to step 3 ${DIM}(run-scenario-planning)${NC}?\n"
  pause 1.5
  type_at_bottom "yes"
  pause 0.4
}

skill_career() {
  echo ""; pause 0.4
  think "Running run-scenario-planning" 3
  echo ""
  printf "${AMBER}⏺${NC} ${BOLD}run-scenario-planning${NC}  ${DIM}business/strategy${NC}\n\n"; pause 0.4
  printf "  ${AMBER}⎿${NC}  A. Freelance on current skills    ${GREEN}60–90 days    ← fastest${NC}\n"
  printf "        2–3 retainers at \$8–12k/mo  ${DIM}·${NC}  warm outreach this week\n"; pause 0.5; echo ""
  printf "  ${AMBER}⎿${NC}  B. Same field, new employer       ${YELLOW}90–120 days${NC}\n"
  printf "        Risk: AI follows you  ${DIM}·${NC}  hedge toward judgment/strategy work\n"; pause 0.5; echo ""
  printf "  ${AMBER}⎿${NC}  C. Adjacent field pivot           ${DIM}6–12 months${NC}\n"
  printf "        Domain expertise in AI-augmented fields\n"; pause 0.5; echo ""
  printf "  ${AMBER}⎿${NC}  D. Reskill for AI-resistant role  ${DIM}12–24 months${NC}\n"
  printf "        Only viable with 13+ months secured\n"
  pause 0.6; echo ""
  printf "  Recommended: A + B in parallel. First retainer closes the gap fastest.\n"
  pause 2.5
}

final_outcome() {
  repin_box
  echo ""
  printf "${DIM}${DLINE}${NC}\n\n"
  printf "${AMBER}⏺${NC} ${BOLD}Final outcome${NC}\n\n"; pause 0.4
  printf "  ${GREEN}⎿  Budget locked${NC}    \$2,285/mo burn  ${DIM}·${NC}  \$610/mo freed\n";                 pause 0.4
  printf "  ${GREEN}⎿  Housing secured${NC}  forbearance call scheduled  ${DIM}—${NC}  6 months protected\n"; pause 0.4
  printf "  ${GREEN}⎿  Runway${NC}           ${YELLOW}7.9${NC}  ${DIM}→${NC}  ${GREEN}13.4 months${NC}\n";        pause 0.4
  printf "  ${GREEN}⎿  Career path${NC}      freelance outreach starts today  ${DIM}·${NC}  job search parallel\n"
  pause 0.8; echo ""
  printf "  You have ${GREEN}${BOLD}13 months${NC} and a plan. First action: 5 warm messages today.  ${AMBER}✓${NC}\n"
  repin_box
  pause 4.0
}

clear
header
pause 1.2
prompt_at_bottom "$QUESTION"
pause 1.5
# Send animation: clear bottom input box
term_h=$(tput lines)
printf "\033[$(( term_h - 1 ));1H\033[2K${DIM}❯${NC} "
pause 0.3
analyze
plan_route
skill_budget
skill_debt
skill_career
final_outcome
sleep 9999
