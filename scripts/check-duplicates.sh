#!/usr/bin/env bash
# Detects near-duplicate SKILL.md files added in a PR.
# Fails if similarity >= 0.7 unless duplicate-reviewed: true is set in frontmatter.
set -euo pipefail

# Detect new SKILL.md files vs origin/main (CI) or HEAD~1 (local)
# Allow NEW_SKILLS to be pre-set via environment (useful for testing)
if [ -z "${NEW_SKILLS:-}" ]; then
  if git rev-parse origin/main >/dev/null 2>&1; then
    BASE="origin/main"
  else
    BASE="HEAD~1"
  fi
  # --diff-filter=A detects Added files only; Copied (C) files also match.
  # Renamed (R) files are not detected — a renamed+edited duplicate will slip through.
  NEW_SKILLS=$(git diff --name-only --diff-filter=AC "$BASE" HEAD -- '*/SKILL.md' 2>/dev/null || true)
fi

if [ -z "$NEW_SKILLS" ]; then
  echo "No new SKILL.md files found. Skipping duplicate check."
  exit 0
fi

COUNT=$(echo "$NEW_SKILLS" | grep -c .)
echo "Checking $COUNT new SKILL.md file(s) for near-duplicates..."

cd "$(git rev-parse --show-toplevel)"
export NEW_SKILLS

node -e "
const fs = require('fs');
const path = require('path');

const THRESHOLD = 0.7;

function parseFrontmatter(content) {
  const match = content.match(/^---\n([\s\S]*?)\n---/);
  if (!match) return {};
  const fm = {};
  const lines = match[1].split('\n');
  let currentKey = null;
  let inArray = false;
  let arrayValues = [];

  for (const line of lines) {
    if (line.startsWith('- ') && inArray) {
      arrayValues.push(line.slice(2).trim().replace(/['\"\[\]]/g, ''));
      continue;
    }
    const kv = line.match(/^([\w-]+)\s*:\s*(.*)/);
    if (kv) {
      if (inArray) { fm[currentKey] = arrayValues; arrayValues = []; inArray = false; }
      const [, key, val] = kv;
      const trimVal = val.trim();
      if (trimVal === '' || trimVal === '[]') {
        inArray = true; arrayValues = []; currentKey = key;
      } else if (trimVal.startsWith('[')) {
        fm[key] = trimVal.replace(/[\[\]]/g, '').split(',').map(s => s.trim().replace(/['\"]/g, ''));
      } else {
        fm[key] = trimVal.replace(/^['\"]|['\"]$/g, '');
        currentKey = key; inArray = false;
      }
    }
  }
  if (inArray) fm[currentKey] = arrayValues;
  return fm;
}

function tagOverlap(tags1, tags2) {
  if (!tags1 || !tags2 || tags1.length === 0 || tags2.length === 0) return 0;
  const set1 = new Set(tags1.map(t => t.toLowerCase()));
  const set2 = new Set(tags2.map(t => t.toLowerCase()));
  const intersection = [...set1].filter(t => set2.has(t)).length;
  const union = new Set([...set1, ...set2]).size;
  return union === 0 ? 0 : intersection / union;
}

const STOP_WORDS = new Set(['use','when','a','an','the','is','are','for','to','in','of','and','or','with','that','this','as','by','on','at','from','be','been','being','have','has','had','do','does','did','will','would','could','should','may','might','it','its','not','but','if','you','your']);

function descOverlap(desc1, desc2) {
  if (!desc1 || !desc2 || typeof desc1 !== 'string' || typeof desc2 !== 'string') return 0;
  const tokenize = s => s.toLowerCase().split(/\W+/).filter(w => w.length > 2 && !STOP_WORDS.has(w));
  const words1 = new Set(tokenize(desc1));
  const words2 = new Set(tokenize(desc2));
  const intersection = [...words1].filter(w => words2.has(w)).length;
  const union = new Set([...words1, ...words2]).size;
  return union === 0 ? 0 : intersection / union;
}

function getDomain(filePath) {
  const parts = filePath.replace(/\\\\/g, '/').split('/');
  const idx = parts.indexOf('skills');
  return idx >= 0 ? (parts[idx + 1] || '') : '';
}

function computeScore(newPath, newFm, existingPath, existingFm) {
  const t = tagOverlap(newFm.tags, existingFm.tags);
  const d = descOverlap(newFm.description, existingFm.description);
  const dom = getDomain(newPath) === getDomain(existingPath) ? 1.0 : 0.0;
  return ((t * 2) + (d * 3) + (dom * 1)) / 6;
}

function findAllSkills(dir) {
  const results = [];
  if (!fs.existsSync(dir)) return results;
  for (const entry of fs.readdirSync(dir, { withFileTypes: true })) {
    const fullPath = path.join(dir, entry.name);
    if (entry.isDirectory()) results.push(...findAllSkills(fullPath));
    else if (entry.name === 'SKILL.md') results.push(fullPath);
  }
  return results;
}

const newSkillPaths = process.env.NEW_SKILLS.trim().split('\n').filter(Boolean);
const allSkills = findAllSkills('skills');
let hasFailure = false;

for (const newSkillPath of newSkillPaths) {
  if (!fs.existsSync(newSkillPath)) {
    console.log('Skipping (not found on disk): ' + newSkillPath);
    continue;
  }

  const newContent = fs.readFileSync(newSkillPath, 'utf8');
  const newFm = parseFrontmatter(newContent);
  const isReviewed = String(newFm['duplicate-reviewed']).toLowerCase() === 'true';

  const scores = [];
  for (const existingPath of allSkills) {
    if (path.resolve(existingPath) === path.resolve(newSkillPath)) continue;
    const existingContent = fs.readFileSync(existingPath, 'utf8');
    const existingFm = parseFrontmatter(existingContent);
    const score = computeScore(newSkillPath, newFm, existingPath, existingFm);
    if (score >= 0.4) scores.push({ path: existingPath, name: existingFm.name || existingPath, score });
  }

  scores.sort((a, b) => b.score - a.score);
  const top3 = scores.slice(0, 3);
  const topScore = top3.length > 0 ? top3[0].score : 0;

  if (topScore >= THRESHOLD) {
    const label = isReviewed
      ? '⚠️  WARNING (duplicate-reviewed bypasses failure): ' + newSkillPath
      : '❌ NEAR-DUPLICATE DETECTED: ' + newSkillPath;
    console.log(label);
    for (const m of top3) {
      const icon = m.score >= THRESHOLD ? '🔴' : '🟡';
      console.log('   ' + icon + ' ' + m.name + ' (' + m.path + ') score=' + m.score.toFixed(2));
    }
    if (!isReviewed) {
      console.log('   → Extend the existing skill, OR add duplicate-reviewed: true to frontmatter with justification in the PR body.');
      hasFailure = true;
    }
  } else {
    console.log('✅ ' + newSkillPath + ' — no near-duplicates (top score: ' + topScore.toFixed(2) + ')');
  }
}

process.exit(hasFailure ? 1 : 0);
"
