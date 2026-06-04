#!/usr/bin/env bash
# Simulated Claude Code TUI demo — no actual commands executed

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
    char="${cmd:$i:1}"
    printf "%s" "$char"
    sleep 0.055
  done
}

header() {
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

# User types question at bottom with animation (Phase 1: idle view)
prompt_at_bottom() {
  local cmd="$1"
  local term_h
  term_h=$(tput lines)
  local r=$(( term_h - 2 ))
  printf "\033[${r};1H${DIM}${DLINE}${NC}"
  printf "\033[$(( r + 2 ));1H${DIM}${DLINE}${NC}"
  printf "\033[$(( r + 1 ));1H"
  type_cmd "$cmd"
}

# Pin question at bottom instantly, return cursor to row 5 (Phase 2: response view)
pin_prompt() {
  local cmd="$1"
  local term_h
  term_h=$(tput lines)
  local r=$(( term_h - 2 ))
  printf "\033[${r};1H${DIM}${DLINE}${NC}"
  printf "\033[$(( r + 1 ));1H${BOLD}❯${NC} ${cmd}"
  printf "\033[$(( r + 2 ));1H${DIM}${DLINE}${NC}"
  printf "\033[5;1H"
}

intro() {
  # User types the natural language problem at bottom
  prompt_at_bottom "$QUESTION"
  pause 1.5

  # Claude identifies which grimoire skills to apply
  clear
  header
  pin_prompt "$QUESTION"
  echo ""
  pause 0.5
  think "Analyzing your situation" 3
  echo ""
  printf "${AMBER}⏺${NC} ${BOLD}Breaking this down with grimoire best practices${NC}\n"
  echo ""
  pause 0.4
  printf "  ${AMBER}⎿${NC}  ${DIM}Step 1${NC}  apply-first-principles    ${DIM}— question the core assumption${NC}\n";   pause 0.4
  printf "  ${AMBER}⎿${NC}  ${DIM}Step 2${NC}  calculate-fire-number     ${DIM}— your real financial target${NC}\n";     pause 0.4
  printf "  ${AMBER}⎿${NC}  ${DIM}Step 3${NC}  design-pricing-strategy   ${DIM}— what to charge as a consultant${NC}\n"; pause 0.4
  printf "  ${AMBER}⎿${NC}  ${DIM}Step 4${NC}  write-value-proposition   ${DIM}— pitch to land the first client${NC}\n"
  pause 2.5
}

skill1() {
  clear
  header
  pin_prompt "$QUESTION"
  echo ""
  pause 0.4
  think "Running apply-first-principles" 2
  echo ""
  printf "${AMBER}⏺${NC} ${BOLD}First Principles${NC}  ${DIM}Aristotle · Descartes · Musk (TED 2013)${NC}\n"
  echo ""
  pause 0.4
  printf "  Assumption  ${DIM}\"I need a coding job for income security\"${NC}\n"
  echo ""
  pause 0.35
  printf "  ${AMBER}⎿${NC}  Is income tied to employment?  ${GREEN}No${NC} ${DIM}— retainers give same predictability${NC}\n"; pause 0.35
  printf "  ${AMBER}⎿${NC}  Is my value in writing syntax?  ${GREEN}No${NC} ${DIM}— judgment, domain depth, trust${NC}\n";    pause 0.35
  printf "  ${AMBER}⎿${NC}  Can AI replace those?           ${GREEN}No${NC} ${DIM}— not for years${NC}\n"
  pause 0.5
  echo ""
  printf "  Rebuilt  ${BOLD}AI removed the junior work.  You are now 3× more valuable.${NC}\n"
  pause 2.2
}

skill2() {
  clear
  header
  pin_prompt "$QUESTION"
  echo ""
  pause 0.4
  think "Running calculate-fire-number" 2
  echo ""
  printf "${AMBER}⏺${NC} ${BOLD}FIRE Number${NC}  ${DIM}Trinity Study · Cooley, Hubbard, Walz 1998${NC}\n"
  echo ""
  pause 0.4
  printf "  ${AMBER}⎿${NC}  Monthly burn     \$5,400  ${DIM}(mortgage + living)${NC}\n";                  pause 0.35
  printf "  ${AMBER}⎿${NC}  FIRE number      ${BOLD}\$1,620,000${NC}  ${DIM}(4.0%% safe withdrawal rate)${NC}\n"; pause 0.35
  printf "  ${AMBER}⎿${NC}  Current savings  \$67,000  ${DIM}(4.1%% of target)${NC}\n";                   pause 0.35
  printf "  ${AMBER}⎿${NC}  Bare runway      ${YELLOW}12.4 months${NC}  ${DIM}then mortgage at risk${NC}\n"
  pause 0.5
  echo ""
  printf "  At \$12k/mo consulting → infinite runway + ${GREEN}${BOLD}FIRE at 51${NC}\n"
  pause 2.2
}

skill3() {
  clear
  header
  pin_prompt "$QUESTION"
  echo ""
  pause 0.4
  think "Running design-pricing-strategy" 2
  echo ""
  printf "${AMBER}⏺${NC} ${BOLD}Pricing Strategy${NC}  ${DIM}Simon-Kucher & Partners · value-based model${NC}\n"
  echo ""
  pause 0.4
  printf "  ${AMBER}⎿${NC}  Hourly contractor  \$90/hr   ${RED}← AI will compress this to zero${NC}\n"; pause 0.35
  printf "  ${AMBER}⎿${NC}  AI-augmented rate  \$150/hr  ${DIM}(still anchors you to time)${NC}\n";     pause 0.35
  printf "  ${AMBER}⎿${NC}  Outcome retainer   ${BOLD}\$12,000/mo${NC}  ${GREEN}← recommended${NC}\n"
  pause 0.5
  echo ""
  printf "  Logic   deliver in 2wks what 3 devs took 2mo — charge for that\n"; pause 0.35
  printf "  Never   hourly on long engagements — AI makes hours irrelevant\n"
  pause 0.5
  echo ""
  printf "  ${BOLD}\$12k/mo${NC} = \$144k/yr  ${DIM}· same salary, no boss, own the IP${NC}\n"
  pause 2.2
}

skill4() {
  clear
  header
  pin_prompt "$QUESTION"
  echo ""
  pause 0.4
  think "Running write-value-proposition" 2
  echo ""
  printf "${AMBER}⏺${NC} ${BOLD}Value Proposition${NC}  ${DIM}Geoffrey Moore \"Crossing the Chasm\" (1991)${NC}\n"
  echo ""
  pause 0.4
  printf "  ${AMBER}⎿${NC}  For     Series A startups that can't afford a 3-person dev team\n"; pause 0.35
  printf "  ${AMBER}⎿${NC}  Who     need senior engineering + AI-accelerated delivery\n";       pause 0.35
  printf "  ${AMBER}⎿${NC}  Unlike  dev shops, offshore teams, or a 6-month hire\n";           pause 0.35
  printf "  ${AMBER}⎿${NC}  Ours    ${BOLD}15yr expertise + AI → 1 engineer, 3× output${NC}\n"
  pause 0.5
  echo ""
  printf "  Proof  ${DIM}\"Shipped in 6 weeks what their team quoted 6 months\"${NC}\n"; pause 0.35
  printf "  CTA    ${DIM}\"30-min discovery call → proposal in 48hrs\"${NC}\n"
  pause 0.6
  echo ""
  printf "${AMBER}⏺${NC}  First \$12,000 retainer in 60 days.  FIRE at 51.  ${BOLD}AI is the superpower.${NC}  ${AMBER}✓${NC}\n"
  pause 3.5
}

clear
header
pause 1.2
intro
skill1
skill2
skill3
skill4
