---
name: apply-asymmetric-matching-strategy
description: Use when competing with limited resources against a stronger opponent across multiple simultaneous contests — deliberately sacrificing the unwinnable match to concentrate strength where it produces a winning majority
source: Sima Qian "Shiji" (史记, Records of the Grand Historian, ~94 BC) — account of Sun Bin's advice to Tian Ji; Von Neumann & Morgenstern "Theory of Games and Economic Behavior" (1944) — assignment problem foundations; Kuhn "The Hungarian Method for the Assignment Problem" (Naval Research Logistics, 1955)
tags: [game-theory, resource-allocation, competition, asymmetric-strategy, matching, tian-ji, sun-bin, assignment-problem]
verified: true
---

# Apply Asymmetric Matching Strategy

Rank your assets and your opponent's assets, sacrifice the slot you will certainly lose, and redeploy those resources to win the contests that determine the overall outcome.

## Why This Is Best Practice

**Origin:** During the Warring States period (~4th century BC), the nobleman Tian Ji (田忌) repeatedly lost horse-racing bets against King Wei of Qi despite having competitive horses, because he matched horses of equal class — and lost narrowly each time. The strategist Sun Bin (孙膑) advised him to deliberately enter his worst horse against the king's best (accepting that loss), then match his best horse against the king's second-best, and his second-best against the king's weakest. Tian Ji won 2 of 3 races and won the bet. The king's lineup did not change; only the assignment did.

**Adopted by:** The underlying mathematics is formalised as the assignment problem in operations research, solved optimally by the Hungarian algorithm (Kuhn, 1955). It is used in logistics, military resource allocation, competitive sports scheduling (seeding), and auction strategy. Military doctrine (Sun Tzu's *Art of War*, which Sun Bin descended from) explicitly addresses avoiding strength-against-strength engagement. Consulting firms (McKinsey, BCG) apply the same logic to portfolio prioritisation: concede markets where you cannot win, concentrate where you can.

**Impact:** The strategy converts a position of overall disadvantage into a winning majority by exploiting information asymmetry and sequencing freedom. It is counter-intuitive: most people match strength-for-strength (which leads to symmetric losses) or spread resources evenly (which leads to symmetric mediocrity). The asymmetric match wins more often than either default and does so with the same total resources.

**Why best:** The strategy works because winning a contest by one unit counts the same as winning it by ten. A narrow win against the opponent's second-best is worth as much as an impossible win against their best. Concentrating to guarantee wins where possible, and deliberately accepting the certain loss, is a higher-expected-value allocation than symmetric matching.

Sources: Sima Qian, *Shiji* (史记) vol. 65 — Sun Tzu and Sun Bin biographies; Von Neumann & Morgenstern, *Theory of Games and Economic Behavior* (1944); Kuhn, *Naval Research Logistics Quarterly* (1955) — Hungarian algorithm; Sun Tzu, *The Art of War* — chapter 6 (Weak Points and Strong)

## Steps

### Step 1: Define the contest structure
Clarify: How many simultaneous contests are there? What determines the overall winner (majority, total score, threshold)? Can you choose assignment freely, or is assignment partially constrained?

Asymmetric matching applies when:
- There are ≥ 3 simultaneous contests
- You have assignment freedom (you choose which of your assets faces which of theirs)
- Winning more contests than losing is sufficient (not required to win every contest)
- You have information about the opponent's assets before locking in assignments

If you cannot observe the opponent's lineup, or have no assignment freedom, skip to Step 6 (adaptations).

### Step 2: Rank all assets — yours and your opponent's
List every asset you control and every asset your opponent controls. Rank both lists from strongest to weakest on the single dimension that determines contest outcomes.

Be honest and objective. Overrating your own assets leads to misallocation.

| Rank | Your assets | Opponent's assets |
|------|-------------|------------------|
| 1 (best) | A | X |
| 2 | B | Y |
| 3 (worst) | C | Z |

If assets are multidimensional, identify the dimension that matters most for each specific contest type and rank on that dimension.

### Step 3: Identify certain losses
Find the contests where your asset is clearly outmatched regardless of assignment — the opponent's strongest asset beats your entire lineup, or you have one slot that will lose no matter what you do.

In Tian Ji's case: his worst horse loses to the king's best horse regardless of assignment. That is the designated sacrifice slot.

**Principle:** You cannot eliminate the certain loss; you can only choose *which* contest it falls in and ensure it is only one.

### Step 4: Assign the sacrifice — put your weakest against their strongest
Deliberately enter your weakest asset against the opponent's strongest. Accept this loss explicitly. This is not a mistake — it is the strategy.

This frees your remaining (stronger) assets to face the opponent's remaining (weaker) lineup at a systematic advantage.

### Step 5: Match remaining assets to maximise wins
With the sacrifice slot assigned, allocate your remaining assets against the opponent's remaining assets to maximise wins in those contests.

**General rule:** After the sacrifice, you have N−1 of your assets facing N−1 of their assets, and your relative strength has improved. Match your next-best against their weakest, and your second-best against their second-weakest — win both.

**In a 3-contest scenario:**
- Your worst vs. their best → lose (sacrifice)
- Your best vs. their second-best → win
- Your second-best vs. their worst → win
- Result: 2 wins, 1 loss — overall winner

**In a 5-contest scenario (majority = 3 wins needed):**
- Sacrifice 2 slots (your two worst against their two best)
- Win the remaining 3 with your three best against their three weakest
- Result: 3 wins, 2 losses — overall winner

### Step 6: Execute and hold the assignment
Commit to the assignment before the contest begins. Do not second-guess the sacrifice slot mid-competition. Changing the assignment under pressure typically reverts to strength-for-strength matching and loses the advantage.

If the opponent can observe and react to your assignment before locking in their own, move simultaneously (hidden assignment) or use a randomised strategy to prevent them from counter-matching.

## Rules

- Never match strength-for-strength by default. Symmetric matching against a superior opponent produces symmetric losses. Asymmetric matching can convert an overall disadvantage into a majority win.
- The sacrifice is intentional and must be chosen, not discovered. If you don't explicitly assign the sacrifice, it assigns itself to your second-best contest by accident.
- The strategy requires information about the opponent's lineup. If their lineup is unknown, the optimisation degrades to probability-based assignment — still better than random, but not as precise.
- Do not over-sacrifice. Sacrificing more than N−1 slots (in an N-contest, majority-wins structure) surrenders the advantage. One sacrifice in a 3-contest scenario is optimal; two sacrifices loses the majority.
- If the opponent has the same strategic awareness, they will try to force your best asset into their best slot. Counter by withholding your assignment or revealing it simultaneously.

## Examples

**Competitive sales territory allocation (3 territories, 3 salespeople):**
Your team: one top performer (A), one mid (B), one weak (C). Competitor's team: two strong reps, one weak. Competitor plans their best rep in the highest-value territory. Assign C to that territory (accept that loss), A to the second-highest-value territory (beat their mid rep), B to the lowest (beat their weak rep). Win 2 of 3 territories, majority of revenue.

**Product feature prioritisation against a better-resourced competitor:**
Your startup has engineering capacity for 3 features this quarter; the incumbent is building the same 3 features with 3× your team. Don't spread effort equally — concede the feature where they have insurmountable advantage (sacrifice). Concentrate full capacity on 2 features where speed-to-market or domain insight gives you an edge. Ship those first, own those segments.

**Job candidate placement across 3 open roles:**
You have 3 candidates and 3 competing firms offering the same candidates placements. Firm has already stated strong preference for their top candidate in role 1. Don't contest role 1 — redirect your strongest candidate to role 2 (where the competing firm has a weaker candidate), and your second-strongest to role 3. Win placements in 2 of 3 roles.

**Auction strategy across multiple simultaneous lots:**
You have a fixed budget competing across 5 lots against a better-funded bidder. Identify the 2 lots they most need (they will bid to win at any price — unwinnable for you). Bid minimally there (sacrifice), concentrate your budget on the 3 remaining lots where their interest is lower. Win 3 of 5.

**Military resource allocation (historical analogy):**
Sun Bin later applied the same logic at the Battle of Guiling (353 BC): instead of matching Wei army strength-for-strength at the contested city, he attacked the Wei capital (their weakness) while the Wei army was stretched. The "sacrifice" was the contested city; the win was the capital.

## Common Mistakes

**Refusing to sacrifice — trying to win everything:** Spreading resources to remain competitive in the unwinnable slot means weakening the winnable slots. You lose the certain loss anyway, and now you lose the winnable contests too.

**Misidentifying the sacrifice slot:** Sacrificing a winnable contest because it seems hard, rather than the genuinely unwinnable one, produces 1 win, 2 losses instead of 2 wins, 1 loss. Rank rigorously before assigning.

**Symmetric matching instinct:** The default human response is to match equal against equal. This is wrong when the opponent is better overall — you lose narrowly across the board. Resist the instinct.

**Revealing the assignment too early:** If the opponent sees your assignment and can react, they will move their strongest asset to contest yours. Reveal assignment simultaneously or use sealed bids.

**Applying where contest structure doesn't fit:** The strategy requires assignment freedom and a majority-wins or threshold structure. In winner-takes-all (single contest), resource averaging (cumulative score), or no-information scenarios, the same logic does not apply cleanly — adapt or use a different framework.
