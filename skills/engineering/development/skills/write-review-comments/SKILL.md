---
name: write-review-comments
description: Use when writing code review comments — to phrase feedback in a way that is kind, clear, and actionable; distinguishes blocking from optional suggestions; and helps the author improve without feeling attacked.
source: 'Google Engineering Practices "How to Write Code Review Comments" (google.github.io/eng-practices/review/reviewer/comments.html); Google Engineering Practices "Handling Pushback in Code Reviews" (google.github.io/eng-practices/review/reviewer/pushback.html)'
tags: [code-review, feedback, communication, collaboration, developer-experience]
verified: true
---

# Write Review Comments

Write code review comments that are kind, explain the reasoning behind each request, distinguish blocking issues from optional suggestions, and help the author improve rather than just comply.

## Why This Is Best Practice

**Adopted by:** Google's eng-practices guide is the canonical reference for code review at Google and is widely adopted across the industry. The practices described have evolved from millions of code reviews across tens of thousands of engineers. Many companies (Shopify, Stripe, Airbnb) reference it directly in their engineering handbooks.

**Impact:** Poorly phrased review comments cause more friction than slow reviews. When comments feel like personal attacks, authors become defensive, arguments replace progress, and team relationships deteriorate. Well-phrased comments produce the same code quality outcomes with far less friction — and make both reviewer and author more effective.

**Why best:** The goal of a code review comment is not to express the reviewer's opinion — it is to improve the code. Comments that fail to explain their reasoning, use dismissive language, or don't distinguish "must fix" from "nice to have" create ambiguity and waste time. The techniques below eliminate the most common causes of review friction while maintaining high standards.

Sources: Google Engineering Practices; Winters, Manshreck & Wright *Software Engineering at Google* (2020)

## Steps

### Step 1: Be kind — critique the code, not the author

Frame every comment around the code, not the person who wrote it. Avoid language that implies incompetence or laziness.

```
# BAD — attacks the author
"Why would you do it this way? This is clearly wrong."
"This is a terrible approach."
"Did you even test this?"

# GOOD — addresses the code
"This approach will break when the list is empty — can we add a guard here?"
"I'd suggest extracting this logic into a helper function for readability."
"This might not handle the case where user is null — worth adding a test?"
```

You can maintain high standards and still be polite. Most author frustration in code review is about comment tone, not the rigor of the standards.

### Step 2: Explain your reasoning

Never leave a bare instruction without a reason. The author needs to understand *why* to learn, to evaluate whether the comment applies in edge cases, and to avoid making the same mistake again.

```
# BAD — no reasoning
"Change this to a map."
"Don't use a global here."
"Extract this into a function."

# GOOD — reasoning included
"A map lookup here is O(1) vs O(n) for the current linear scan — worth changing
 given this runs in the request hot path."
"A global makes this hard to test in isolation; dependency-inject the config instead."
"This logic appears in three places — extracting to a function prevents future drift."
```

If you can't explain why, reconsider whether the comment is warranted.

### Step 3: Use `Nit:` for non-blocking suggestions

Prefix non-blocking, optional suggestions with `Nit:` (short for "nitpick"). This immediately signals to the author that the comment is a suggestion they can take or leave — it does not block approval.

```
Nit: This variable name could be more descriptive — `retryCount` instead of `n`?

Nit: Extra blank line here.

Nit: Consider using `enumerate()` instead of tracking the index manually.
```

Without `Nit:`, every comment looks like a blocker. Authors spend time debating minor style issues that the reviewer didn't actually care about strongly. `Nit:` eliminates that ambiguity.

Use sparingly — a PR with 20 nits is exhausting even if none are blocking.

### Step 4: Ask questions instead of giving orders (when appropriate)

For complex design decisions where the author may have context you don't, ask rather than assert. This invites explanation and may reveal that the author is correct.

```
# Order — appropriate when the issue is clear-cut
"This will cause a race condition — the mutex must be held across both reads."

# Question — appropriate when you're uncertain or the author may have context
"Would it be simpler to use a read-through cache here? I'm not sure what the
 write frequency looks like on this data."

# Question — invites the author to explain a non-obvious choice
"Is there a reason this is using a raw SQL query rather than the ORM?
 Happy to keep it if there's a performance reason, just want to understand."
```

Don't ask questions as a passive-aggressive way of pointing out problems. "Did you consider testing this?" is not a question — it's a veiled accusation. Ask genuine questions; give direct feedback for clear issues.

### Step 5: Mark blocking vs. non-blocking explicitly

If your tool doesn't distinguish blocking from non-blocking comments automatically, state it:

```
# Blocking — must be resolved before approval
"[blocking] This leaks the database connection on error. Close it in a defer."

# Non-blocking — nice to have, doesn't hold up the PR
"[non-blocking] This might be easier to read with an early return."

# FYI — no action needed, just information
"[fyi] There's a related ticket open for this: JIRA-1234."
```

Ambiguous "blocking" status is one of the most common causes of review stalls — the author doesn't know what to act on and the reviewer thinks the author is ignoring them.

### Step 6: When the author pushes back — who is right?

When an author disagrees with a comment, take their argument seriously:

1. **Consider their perspective first.** They are closer to the code. Does their argument make sense? Does it improve or at least maintain code health? If they're right, say so and drop the comment.

2. **If you still believe your suggestion is better**, explain *why* more thoroughly. Acknowledge their point and add the additional context that makes your case. Repeat as needed — sometimes it takes a few rounds.

3. **Hold firm on code health, not on style preferences.** If the disagreement is about a style preference with no code health impact, defer to the author. If it genuinely improves maintainability, testability, or correctness, continue advocating politely.

4. **Never yield simply to end the conflict.** Approving a CL with a known code health problem to avoid friction degrades the codebase in small steps. Stay polite but clear.

```
# Productive pushback response
"I see your point — the current approach works for the happy path. My concern is
 that when X condition occurs (which has happened before in JIRA-456), the lock
 isn't released. I think the try/finally approach I suggested handles both cases.
 Can we discuss?"
```

## Rules

- Every blocking comment must explain its reasoning. A comment with no explanation is a request without context — authors can't learn from it or evaluate edge cases.
- Use `Nit:` for anything that doesn't block approval. Reserve unmarked comments for real blockers.
- Never yield to pushback to end an argument. Either the author is right (update your position) or they aren't (hold your position, more clearly).
- Critique code, not the person. Never imply incompetence; always imply the fix is easy and the author is capable of doing it.

## Common Mistakes

**Leaving ambiguous blocking status:** Every comment looks like a blocker by default. Explicitly mark non-blocking items or use `Nit:`. Ambiguity wastes both parties' time.

**Giving orders without reasons:** "Change X to Y" is not a comment — it's a command. Commands don't teach, don't invite discussion, and can't be evaluated for correctness without context.

**Using questions as passive-aggressive criticism:** "Did you test this?" or "Is this really the best approach?" are not genuine questions — they're disapproval dressed as curiosity. Give direct feedback instead.

**Yielding on code health to avoid conflict:** Approving substandard code to prevent an uncomfortable conversation is the primary mechanism by which codebases decay. Be kind, be clear, hold the standard.
