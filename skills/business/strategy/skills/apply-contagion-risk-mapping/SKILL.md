---
name: apply-contagion-risk-mapping
description: Use when assessing organizational risk exposure — to map how cascade failures propagate through implicit network proximity, shared infrastructure, and correlated dependencies, because organizations that assess only direct relationships systematically underestimate their exposure to others' failures
source: "左传 僖公五年 (655 BC) — 宫之奇谏假道: 虢，虞之表也。虢亡，虞必从之。唇亡则齿寒 (Guo is the outer protection of Yu. If Guo falls, Yu must follow. When lips are gone, the teeth feel cold); Basel III systemic risk framework (BIS, 2010) — too-interconnected-to-fail; Brunnermeier & Pedersen 'CoVaR' systemic risk measure (2011) — adopted by Federal Reserve, ECB; Sheffi 'The Resilient Enterprise' (2005) — supply chain contagion; ISO 31000:2018; tech industry blast-radius analysis"
tags: [strategy, risk-management, resilience, supply-chain, systemic-risk, dependencies]
verified: true
---

# Apply Contagion Risk Mapping

Map your exposure to cascade failures from entities in your implicit network — not just explicit strategic partners — because failure propagates through proximity and shared infrastructure, not only through deliberate relationships.

## Why This Is Best Practice

左传 僖公五年 (655 BC) — 宫之奇谏假道:

> 虢，虞之表也。虢亡，虞必从之。唇亡则齿寒。

"Guo is the outer protection of Yu. If Guo falls, Yu must follow. When lips are gone, the teeth feel cold."

**Why best:** In 655 BC, the Duke of Jin requested passage through the state of Yu to attack the neighboring state of Guo. The minister Gong Zhiqi warned the Yu ruler against granting the passage: Yu and Guo were structurally interdependent — Guo provided outer protection to Yu without any formal alliance treaty requiring it. If Guo fell, Yu would be exposed and would inevitably follow. The Yu ruler granted the passage, Jin defeated Guo, and on the return march, Jin absorbed Yu as well. The classical commentators use this case as the canonical illustration of cascade failure through proximity: the threat did not come from a direct enemy relationship, but from the collapse of an adjacent entity that Yu had not recognized as load-bearing for its own survival. 唇亡齿寒 ("when lips are gone, the teeth feel cold") became standard vocabulary in Chinese strategic and risk analysis for implicit interdependence that only becomes visible when it fails.

**Basel III Systemic Risk Framework (BIS, 2010):** Enacted after the 2008 financial crisis, which demonstrated that individual firms can be solvent while the system collapses through contagion. Basel III introduced the concept of "systemically important financial institutions" (SIFIs) — entities whose failure creates cascade risk for others not directly connected to them contractually. The key finding: traditional risk models assessed bilateral exposure (what do I owe to, or own of, this counterparty?) but missed network-level exposure (what happens to me if an entity I depend on indirectly fails, or if shared funding sources dry up?). Basel III requires banks to hold additional capital buffers proportional to their systemic interconnectedness, not just their bilateral credit exposure. Standard across 28 BCBS member jurisdictions covering ~90% of global banking assets.

**Brunnermeier & Pedersen — CoVaR (2011, adopted by Federal Reserve and ECB):** CoVaR ("Conditional Value at Risk") measures the risk to the financial system conditional on a specific institution being in distress — capturing the spillover effect that bilateral VaR models miss. The key innovation: your systemic risk contribution is not the risk you hold but the risk you transmit. An institution with low standalone VaR but high CoVaR imposes large costs on others through shared exposures. The Federal Reserve incorporates CoVaR in systemic risk monitoring; the ECB uses it in its Financial Stability Review. Standard in financial stability research; referenced in IMF Global Financial Stability Reports.

**Yossi Sheffi — "The Resilient Enterprise" (Harvard/MIT, 2005):** Sheffi's systematic study of supply chain disruptions — the 2000 Ericsson/Nokia fire (single Phillips semiconductor plant destroyed; Ericsson lost $400M in sales, Nokia was unaffected because it had mapped single-source dependencies), the 2001 Taiwan earthquake (disrupted disk drive supply globally), the 2011 Japan earthquake (automotive supply chains worldwide halted because Tier-3 suppliers in Tohoku were the single source for specific components that no buyer had mapped) — establishes that the most damaging supply chain failures are from entities at Tier-2 and Tier-3 that buyers had never analyzed, often because there was no direct contractual relationship. Sheffi's framework: map the full supply network, not just direct suppliers; identify single points of failure in the implicit network; scenario-test cascade scenarios. Standard in supply chain resilience programs; taught at MIT CTL.

**Tech industry blast-radius analysis:** Widely adopted in cloud infrastructure engineering (AWS, Google, Netflix) to assess the scope of failures before they occur: if component X fails, what else fails with it? The discipline emerged from real incidents — a single AWS availability zone outage taking down applications that had not implemented cross-zone redundancy; a shared library or authentication service failing and cascading across products that shared it without that dependency being explicitly tracked. Standard practice in site reliability engineering; covered in the Google SRE book (Beyer et al., 2016).

**Why distinct from `design-alliance-strategy`:** design-alliance-strategy addresses building explicit coalitions to achieve competitive objectives — it is about constructing deliberate relationships. apply-contagion-risk-mapping addresses the risk exposure created by implicit network proximity without any deliberate relationship — the 宫之奇 problem, where the threat comes not from an enemy but from the structural collapse of an adjacent entity you never formally allied with.

**Why distinct from `apply-dependency-trap`:** apply-dependency-trap addresses deliberately creating lock-in for customers. apply-contagion-risk-mapping addresses the risk imposed on you by your implicit dependencies — the inverse: mapping where you are vulnerable to others' failures, not where you are creating vulnerability in others.

**Why distinct from `design-incident-response`:** design-incident-response designs the operational process for responding to failures after they occur. apply-contagion-risk-mapping is pre-event: mapping implicit exposure before failures materialize, to identify and reduce cascade paths.

**Adopted by:** Basel III systemic risk framework (BIS, 2010) — adopted by 28 BCBS member jurisdictions covering ~90% of global banking assets; Federal Reserve and ECB (CoVaR methodology); AWS, Google, and Netflix (blast-radius analysis as standard SRE practice); ISO 31000:2018 risk management standard; MIT CTL supply chain resilience programs (Sheffi framework).

**Impact:** Organizations that mapped Tier-2 and Tier-3 supply chain dependencies avoided the catastrophic losses others suffered in documented contagion events — Nokia's pre-mapped single-source dependencies allowed it to avoid the $400M revenue loss Ericsson sustained from the 2000 Phillips semiconductor plant fire; banks implementing Basel III systemic risk buffers maintained solvency through stress events that would have caused cascade failures under pre-2010 bilateral-only risk models.

## Steps

1. **Map the implicit network, not just the explicit one.** Begin by listing all direct relationships: suppliers, customers, infrastructure providers, financial counterparties, regulatory dependencies. Then extend to the second-order network: who do your direct partners depend on that you also depend on (correlated single sources)? What shared infrastructure or services do you and your key customers both rely on? What geographic concentrations exist? The failure point in most contagion scenarios is not in the direct relationship list — it is in the implicit shared dependencies that only become visible after failure.

2. **Classify dependencies by replaceability and failure impact.** For each node in the mapped network, assess two dimensions:
   - **Criticality:** if this entity fails, what function or capability do you lose? Rate as: loss of revenue only, operational disruption, existential threat.
   - **Replaceability:** how quickly could you source an alternative? Rate as: days (commodity), weeks (specialized), months+ (sole source or deeply embedded).
   
   Nodes that are high criticality and low replaceability are your structural vulnerabilities — the equivalent of the Guo-Yu relationship. These require active mitigation regardless of how stable the entity currently appears.

3. **Identify correlated failure scenarios.** Some dependencies are individually manageable but fail simultaneously under shared stress scenarios. Map correlated exposures:
   - **Geographic concentration:** multiple suppliers, customers, or data centers in the same region — at risk from the same natural disaster, regulatory event, or infrastructure failure
   - **Single upstream source:** multiple of your critical suppliers or partners sourcing from the same Tier-2 or Tier-3 supplier — a failure at that level cascades to all of them simultaneously
   - **Shared financial counterparties:** multiple business relationships that all depend on the same bank, insurance provider, or credit facility — at risk if that financial entity faces stress
   - **Platform or infrastructure concentration:** multiple functions or products depending on the same cloud provider, authentication service, or data feed
   
   The 2008 financial crisis, the Ericsson/Nokia fire, and the 2011 Japan earthquake supply chain disruptions all involved correlated failures that looked like independent risks until they materialized simultaneously.

4. **Scenario-test your most critical dependencies.** For each high-criticality, low-replaceability node, run a structured failure scenario: "This entity fails tomorrow. What happens to us in the first 24 hours? First 30 days? What is the path to restoration, and at what cost?" The scenario test surfaces hidden dependencies (you discover functions that you believed were independent but share the failing node) and quantifies impact in terms that justify mitigation investment. Use actual operational data, not best-case assumptions.

5. **Assess contagion paths, not just direct exposure.** Beyond direct dependencies, map the transmission mechanisms by which others' failures reach you indirectly:
   - **Reputation contagion:** if a supplier or partner is publicly implicated in misconduct, does that contaminate your own brand or regulatory standing?
   - **Market contagion:** if a major competitor fails, does that create a credit event or pricing disruption in your shared markets?
   - **Regulatory contagion:** if a peer organization in your sector violates a regulation, does that trigger sector-wide regulatory scrutiny that creates compliance burden for you?
   These paths are often missed because they involve no direct contractual relationship — precisely the 唇亡齿寒 dynamic.

6. **Implement structural interventions for identified vulnerabilities.** Awareness of exposure is not risk management. For each identified structural vulnerability, choose an intervention matched to the risk level:
   - **Dual-source critical supplies:** qualify a second supplier for single-source dependencies before failure occurs, not after
   - **Geographic distribution:** split critical operations, inventory, or infrastructure across regions that cannot fail simultaneously
   - **Inventory buffers:** hold strategic inventory of critical components at the tier where your Tier-1 supplier is sole-sourced
   - **Financial redundancy:** maintain credit facilities with multiple counterparties; avoid concentration of operating cash at a single bank
   - **Contractual cascade protections:** include step-in rights, business continuity requirements, or notification obligations in contracts with entities whose failure would be existential for you
   
   The appropriate intervention is sized to the criticality and replaceability scores from Step 2, not to worst-case scenarios for every dependency.

7. **Maintain and update the map.** Network topology changes continuously: suppliers consolidate, infrastructure providers merge, geographic concentrations shift, new shared dependencies emerge from product changes. Contagion risk maps built once and never updated reflect historical rather than current exposure. Review the map at minimum annually, and update immediately when a critical dependency changes — new supplier, new infrastructure provider, major acquisition by a key partner.

## Rules

- The most dangerous dependency is the one you don't know you have. Implicit dependencies — shared infrastructure, geographic concentration, upstream single sources — are invisible until they fail. Explicit effort to map them is required because they will not surface themselves.
- High stability of a critical node is not a reason to defer mitigation. The Yu ruler believed Guo was stable. The Ericsson purchasing team believed the Phillips plant was stable. Stability assessments of single points of failure are almost always made without visibility into the node's own vulnerabilities. Replaceability mitigation must precede failure, not follow it.
- Correlated risk is not the sum of individual risks. Two suppliers that each have 2% annual failure probability but source from the same Tier-2 supplier have a correlated failure rate that far exceeds either individual estimate. Assess correlated scenarios, not just independent risk.
- Contagion does not require a contractual relationship. The 唇亡齿寒 principle: failure propagates through proximity, shared infrastructure, and market structure — not only through bilateral contracts. Risk assessment limited to entities with which you have explicit contractual relationships will systematically miss the most dangerous exposure paths.

## Examples

**Supply chain:** A medical device manufacturer sources a critical component from a single European supplier. That supplier sources its raw material from two Asian suppliers who both source from the same mine in a single country. The manufacturer has never mapped beyond Tier-1. A labor dispute at the mine halts production. Both Asian suppliers stop shipping. The European supplier halts production. The medical device manufacturer's production stops within 6 weeks. Contagion risk mapping would have identified the Tier-3 geographic concentration and triggered strategic inventory buffering for a component that is 6-week lead time on good days.

**Financial services:** A regional bank holds what it believes to be a diversified commercial real estate loan portfolio across 40 borrowers in different property types. Contagion risk mapping reveals: 32 of the 40 borrowers have operating leases as their primary cash flow source; 28 of those operate in the same metropolitan area; the metropolitan area's largest employer (a single tech company) is at risk. A single employer relocation cascades: vacancy rates spike, operating cash flows collapse, 28 loans simultaneously enter default. Portfolio diversity by borrower name masked geographic and tenant concentration.

**Technology:** A SaaS company builds its product on a major cloud provider and uses that same provider's authentication, storage, and networking services. A regional availability zone outage disables all four simultaneously. Contagion risk mapping would have identified that services nominally in different categories share the same availability zone dependency and triggered cross-zone redundancy design before the outage, not after.

**Organizational:** A consulting firm has three major clients that collectively represent 65% of revenue. Contagion risk mapping reveals: all three clients are in the same industry vertical; all three have the same 2–3 private equity owners who sit on their boards; all three operate on the same fiscal year cycle with synchronized budget decisions. An industry-wide shock (regulatory change, commodity price shift, ESG-driven investor pullback) hits all three simultaneously, resulting in a coordinated contract pause. The firm's revenue concentration by client name masked the correlation in the client portfolio.

## Common Mistakes

- **Mapping only the first tier.** Tier-1 supplier lists, direct customer concentration, named counterparties — these are the beginning of the map, not the map. The most common catastrophic failures originate at Tier-2 or Tier-3, in geographic concentration, or in shared infrastructure that no tier analysis revealed.
- **Treating low historical failure rate as low current risk.** The 宫之奇 insight: the fact that Guo had never been attacked before was not a risk assessment, it was an absence of data. Low historical failure rate of a critical sole-source dependency is not a substitute for replaceability analysis.
- **Assessing dependencies independently rather than in correlation.** Individual dependency assessments miss correlated failure scenarios. The question is not "what is the probability this supplier fails?" but "what is the probability of simultaneous failure across the cluster of dependencies that share a common upstream source?"
- **Confusing mitigation awareness with mitigation.** Identifying a structural vulnerability without implementing a structural intervention leaves the exposure unchanged. Risk maps must produce action — dual sourcing, geographic distribution, inventory buffers, contractual protections — not just documentation.
