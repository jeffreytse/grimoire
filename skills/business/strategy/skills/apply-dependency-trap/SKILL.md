---
name: apply-dependency-trap
description: Use when customer or partner adoption at scale requires removing upfront barriers — to offer genuine value freely, let switching costs accumulate naturally, then shift to monetisation once the dependency is established
source: Thirty-Six Stratagems #16 欲擒故縱 "In Order to Capture, First Let It Go" (三十六計, Sawyer trans. 1994); Verstappen "The 36 Strategies of Ancient China" (2007); Parker et al. "Platform Revolution" (2016); Shapiro & Varian "Information Rules" (1999)
tags: [thirty-six-stratagems, dependency, freemium, platform, switching-costs, lock-in, let-go-to-capture]
verified: true
---

# Apply Dependency Trap

Offer genuine, unconstrained value during an adoption phase that accumulates switching costs naturally — then shift to monetisation once the dependency makes alternatives costly to pursue.

## Why This Is Best Practice

**Origin:** Stratagem #16 of the Thirty-Six Stratagems: "In order to capture, first let it go" (欲擒故縱). In the original military context: when an enemy is fleeing, pursue them too hard and they will fight with the desperation of a cornered force; let them go and they relax, then capture them on your terms. In business: present the customer with a free and unconstrained adoption path; let dependency accumulate; capture them when switching becomes costlier than staying.

**Adopted by:** The dependency trap is the mechanism behind the most successful platform and software growth models of the past two decades. Slack offered unlimited free messaging until teams were deeply dependent on channels, integrations, and message history — then introduced pricing for large teams. Dropbox offered 2GB free, let users build years of file dependencies, then monetised through storage tiers. HashiCorp open-sourced Terraform and Vault, built a global infrastructure automation community, then introduced enterprise features that the most dependent users — enterprises — needed. Shapiro and Varian's *Information Rules* (1999) formalised the economics: lock-in is created by switching costs, and switching costs accumulate through use, data, integrations, and trained habits.

**Impact:** Premature monetisation — presenting pricing before the customer has experienced value — is the most common early failure mode in platform and SaaS businesses. The customer's cost-benefit calculation before dependency is unfavourable: they bear the switching cost from their current solution without yet experiencing the dependency value that will make future switching painful. The dependency trap delays monetisation until the switching cost calculation runs the other way.

**Why best:** The dependency trap is not deceptive when executed correctly — the value offered during the adoption phase must be genuine, and the monetisation model should be predictable to customers at the time of adoption. The competitive advantage is time: starting monetisation before dependency accumulates loses customers to competitors who let the dependency build. Starting after dependency is established monetises a committed base.

Sources: Thirty-Six Stratagems #16 (Sawyer trans. 1994); Shapiro & Varian, *Information Rules* (1999) — switching costs and lock-in; Parker, Van Alstyne & Choudary, *Platform Revolution* (2016)

## Steps

### Step 1: Define the capture outcome — what dependency do you want to establish?

Before designing the free phase, identify the endpoint:

| Dependency type | How it accumulates | Switching cost |
|----------------|-------------------|---------------|
| Data | Customer data stored in the platform | Migration, data loss, re-entry |
| Integrations | Product connected to other tools via APIs | Re-integration effort, potential downtime |
| Trained workflows | Team habits built around the product | Re-training, productivity loss during transition |
| Network effects | Value comes from connections to others on the platform | Loss of connections, inability to reach others |
| Ecosystem | Third-party apps, extensions, community built around the platform | Loss of ecosystem-specific capabilities |
| Embedded decisions | Historical decisions encoded in the system | Loss of institutional memory, auditability |

The capture outcome should be one or more dependency types that are genuinely costly to reverse once established.

### Step 2: Design the "let go" phase — what you offer freely, and why it builds dependency

The free phase must pass two tests:
1. **Genuinely valuable without the dependency:** If the free offering is not valuable on its own, it does not attract adoption.
2. **Creates accumulating switching costs during use:** If use of the free offering does not build dependency, there is no trap — only charity.

Design the free phase to maximise both:

| Element | Free phase design |
|---------|------------------|
| Pricing | Free, or priced below the customer's cost of attention (trial, freemium, open-source) |
| Onboarding | Remove all friction from adoption; ask for nothing that increases the cost of starting |
| Feature scope | Include the features that most powerfully build dependency during free use; gate only what is valuable specifically to committed buyers |
| Data | Store customer data, integrations, and decisions in the platform during free use — this is the dependency accumulation mechanism |

### Step 3: Remove friction from the adoption phase — let the dependency accumulate without pressure

During the free phase:
- Do not push premature monetisation conversation; let the product create the dependency
- Do not require commitments that increase the psychological cost of starting
- Provide genuine customer success support to maximise adoption depth
- Monitor usage signals that indicate dependency accumulation: integration count, data volume, active users, workflow touchpoints

Premature monetisation disrupts the dependency accumulation process. Sales pressure during the free phase signals the end of the free offering and activates the customer's switching cost calculation before dependency is sufficient.

### Step 4: Identify the dependency threshold — when switching cost exceeds replacement cost

The transition to monetisation is viable when the customer's cost of switching to an alternative exceeds the cost of paying for the current platform. Signals that the threshold has been crossed:

- Integration count: each integration with another tool raises switching cost by an estimated amount
- Data volume: the customer's historical data has accumulating value that would be lost on migration
- User base: team habits and training investment are sunk; re-training cost is calculable
- Business process: key decisions are encoded in the platform and would need to be re-executed

Do not wait for the customer to calculate this themselves. Transition to monetisation before they consciously notice the dependency, when the dependency is strong enough that the cost of leaving exceeds the cost of paying.

### Step 5: Shift to monetisation in a way that preserves the relationship

The monetisation transition must not feel like a trap sprung — even if the dependency makes staying rational, a customer who feels deceived will leave on principle. Execute the transition:
- Signal the monetisation model before the free phase begins — transparency about the eventual paid structure is not disqualifying if the free value is real
- Make the pricing transition gradual — grandfather existing free usage into a discounted tier; increase pricing over time rather than immediately
- Provide genuine additional value in paid tiers that is not available free — the paid tier is an upgrade, not a removal of what was free
- Treat the customer success function as a retention mechanism — renewal decisions are made when dependencies are reassessed

## Rules

- The free offering must be genuinely valuable. A free offering that is intentionally deficient to create frustration is not a dependency trap — it is bait-and-switch, which produces churn and reputational damage. The free value must be sufficient for a customer to build genuine dependency.
- Be transparent about the monetisation model. Customers who discover a freemium trap they were not informed of become adversaries. Disclose the eventual pricing model; most customers will adopt anyway if the free value is genuine.
- Monitor dependency accumulation — not vanity metrics. The measure of the free phase is not registered users or downloads; it is integration count, data volume, and workflow touchpoints. Monitor what creates switching cost, not what creates growth numbers.
- Do not monetise before the dependency threshold. Premature monetisation with insufficient dependency produces churn. The dependency calculation must run in your favour at the time of monetisation.
- Design the paid tier to be worth paying for. If the only reason to pay is dependency — not additional value — the paid tier is rent-seeking. Add genuine value to the paid tier so that the monetisation is defensible to the customer and to their board.

## Examples

**Slack: free messaging → paying for message history:**
Slack offered unlimited team communication free for an unlimited time. During free use, teams built channels, connected integrations (GitHub, Jira, Salesforce, Google Drive), and encoded 6–18 months of institutional knowledge in message history. When Slack introduced message history limits for free plans, the switching cost was the loss of archived conversations, decisions, and integration configurations. Teams that had accumulated significant history faced the calculation: pay Slack or lose access to institutional memory. The dependency threshold had been crossed; most teams paid. Slack's 2019 IPO revealed that free-to-paid conversion rates exceeded 30% at teams above 10 members — among the highest in SaaS.

**HashiCorp: open-source Terraform → enterprise licensing:**
HashiCorp open-sourced Terraform in 2014, enabling any organisation to adopt infrastructure-as-code at zero cost. Over 8 years, Terraform became the industry-standard tool for cloud infrastructure management: enterprise engineering teams built years of Terraform configurations, modules, and workflows. By 2022, HashiCorp's tools were embedded in the CI/CD pipelines of tens of thousands of companies. The enterprise tier — offering collaboration features, audit logging, and governance — addressed specifically the needs of organisations whose dependency was so deep that switching to an alternative (Pulumi, CDK) would require rewriting years of infrastructure code. HashiCorp's 2023 licensing change (BSL) was controversial precisely because the dependency was deep enough that most enterprise customers had no practical alternative.

**Developer tools embedded in CI/CD pipelines:**
A developer productivity startup offers a free code review tool that integrates into GitHub Actions, GitLab CI, and Jenkins. During the free phase, engineering teams connect the tool to their repositories, configure rules, and build review workflows around its output. After 12 months, the average customer has 200+ rules configured, 3 years of review history, and 8 integrated tools. The switching cost is re-configuration, re-training, and loss of historical data. The startup introduces paid tiers for teams above 10 developers. Churn at the monetisation transition is under 10% because the switching cost exceeds the subscription cost for all but the smallest teams.

## Common Mistakes

**Free offering that does not build dependency:** Offering a free product that is useful but generates no integration, data, or workflow commitment produces adopters who will leave without paying when the free period ends. Design the free offering to specifically maximise dependency accumulation.

**Premature monetisation:** Introducing pricing before the dependency threshold is crossed produces churn. The customer's switching cost calculation has not yet run in your favour; they can leave without significant cost. Wait for dependency accumulation signals before transitioning.

**Opacity about the monetisation model:** Customers who discover the monetisation model as a surprise — particularly if the free offering was presented as permanent — feel deceived. Transparency about the eventual paid model does not prevent adoption if the free value is genuine.

**Paid tier that only restricts free features:** If the paid tier is designed only to restore features removed from the free offering, it is rent-seeking. Add genuine value to paid tiers — collaboration features, security controls, audit capabilities, SLA guarantees — so the monetisation is a real upgrade, not a hostage payment.

**Ignoring ethical obligations at dependency threshold:** The dependency trap, fully executed, gives you significant pricing power. Using that power to extract unreasonable rents from customers with no alternative — rather than pricing reasonably for genuine value — creates existential backlash. The regulatory, reputational, and competitive risk of extractive pricing at dependency is significant.
