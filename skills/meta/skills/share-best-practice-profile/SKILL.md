---
name: share-best-practice-profile
description: Use when the user wants to publish or share a practice profile — e.g., "share my team profile", "publish this profile", "open-source my backend-defaults profile".
source: Grimoire project conventions; GitHub Gist and repository conventions
tags: [profile, sharing, open-source, ecosystem, contributor]
related: [review-best-practice-profile, write-best-practice-profile, apply-best-practice-profile]
---

# Share Best Practice Profile

Publish a validated practice profile so others can install and use it.

## Why This Is Best Practice

**Adopted by:** The `eslint-config-*` ecosystem (3000+ published configs on npm) and VS Code extension packs demonstrate the value of community-shared configuration bundles. A single `eslint-config-airbnb` install replaced thousands of manual ESLint rule configurations. Grimoire profiles follow the same pattern — one shareable file replaces manual skill selection for everyone who shares the same context.
**Impact:** A shared profile eliminates onboarding friction for new team members and external contributors who need to activate the same practices. A profile shared publicly becomes a reusable community asset — the same leverage as an open-source library.
**Why best:** Sharing without review leads to broken profiles that fail silently for others. The gate (`review-best-practice-profile` first) ensures what is shared actually works.

Sources: npm `eslint-config-*` ecosystem; VS Code extension pack format; Grimoire `docs/profiles.md`

## Steps

### 1. Confirm the profile is reviewed

If `review-best-practice-profile` has not been run on this profile, run it now and resolve any `FAIL` results before continuing. Warnings are acceptable.

---

### 2. Choose sharing format

```
Share as:
  [g] GitHub Gist    — single file, easy to update, direct install URL
  [r] GitHub repo    — full repo, recommended for maintained profiles
  [c] Clipboard      — copy TOML to clipboard for private sharing
```

**Gist** — best for personal or small team profiles.
**Repo** — best for maintained, versioned, community profiles. Use naming convention: `grimoire-profile-<name>`.

---

### 3. Prepare for sharing

Add a `description` field if absent — required for others to understand the profile:

```toml
name = "my-team"
description = "Backend team defaults — DDD, SOLID, clean architecture boundaries."
```

For repo shares, generate a minimal README:

```markdown
# grimoire-profile-my-team

Practice profile for grimoire — backend team defaults.

## Install

\`\`\`bash
curl -fsSL https://raw.githubusercontent.com/<user>/grimoire-profile-my-team/main/my-team.toml \
  -o .grimoire/profiles/my-team.toml
\`\`\`

## Activate

\`\`\`toml
# .grimoire/settings.toml
profiles = ["my-team"]
\`\`\`

## Skills included

- apply-solid-principles
- apply-domain-driven-design
- apply-low-coupling
```

---

### 4. Publish

**Gist:**
```bash
gh gist create .grimoire/profiles/my-team.toml --public --desc "grimoire profile: my-team"
```

**Repo:**
```bash
gh repo create grimoire-profile-my-team --public
git init && git add . && git commit -m "feat: initial profile"
git remote add origin https://github.com/<user>/grimoire-profile-my-team.git
git push -u origin main
```

---

### 5. Confirm with install instructions

```
✓ Published: https://gist.github.com/<user>/<id>

Others install with:
  curl -fsSL https://gist.githubusercontent.com/<user>/<id>/raw/my-team.toml \
    -o ~/.grimoire/profiles/my-team.toml

Then activate:
  profiles = ["my-team"]   # in .grimoire/settings.toml
```

## Common Mistakes

**Sharing before reviewing.** Always run `review-best-practice-profile` first. A broken profile shared publicly is harder to fix than one caught before publish.

**Missing description.** Others can't evaluate a profile without knowing what it's for. Always include a `description` field before sharing.
