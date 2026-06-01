# Domain Safety Guide

Skills in health, law, finance, and psychology domains can influence decisions with real consequences. grimoire encodes best practices, not personalized advice — that distinction must be explicit in every skill in these domains.

Missing a required disclaimer is a **REJECT finding** in `review-skill`. This is not optional.

---

## Why this matters

A skill that says "take 500mg of ibuprofen" without qualification is giving medical advice. A skill that says "sign this clause without an attorney" is giving legal advice. grimoire skills must stay on the safe side of that line by:

1. Encoding the general practice — what professionals do, in aggregate
2. Citing the evidence tier so readers know the strength of the claim
3. Never prescribing specific doses, diagnoses, or decisions for individual situations
4. Including the required footer on every skill in the domain

---

## Health / Medicine

### Evidence tier tagging

Every factual claim in a health skill must carry an evidence tier tag. Place the tag in square brackets immediately after the claim.

| Tag | Meaning |
|-----|---------|
| `[RCT]` | Randomized controlled trial evidence |
| `[SR]` | Systematic review or meta-analysis |
| `[Consensus]` | Expert consensus or clinical guideline |
| `[Practical]` | Widely used practice without strong RCT support |

```markdown
# Good
Progressive overload produces greater strength gains than fixed-load training [SR, Ralston et al. 2017].

# Good — lower evidence tier but explicitly tagged
Many coaches use 3-second eccentric tempo for hypertrophy [Practical].

# Bad — no evidence tier, claim appears authoritative
Progressive overload is the best way to build strength.
```

`[RCT]` and `[SR]` are strongest. `[Consensus]` and `[Practical]` are acceptable but must be labeled. Unlabeled factual claims in health skills are a REJECT finding.

### What NOT to include

- Specific medication names, doses, or schedules
- Diagnoses for specific conditions or symptoms
- Personalized treatment protocols ("if you have X, do Y")
- Claims that override clinical judgment

### Credible sources

| Source type | Examples |
|-------------|---------|
| Clinical guidelines | WHO Clinical Guidelines, CDC recommendations |
| Medical institutions | NIH, Mayo Clinic, Cleveland Clinic |
| Professional bodies | ACSM Position Stands, NSCA CSCS content, ADA guidelines |
| Systematic reviews | Cochrane Reviews, JAMA, NEJM systematic reviews |
| Sports science | ACSM/NSCA position stands for exercise and nutrition |

Avoid: individual practitioner blogs, supplement company research, anecdotal case studies.

### Required footer

Every health / medicine skill must end with:

```markdown
> For personal health decisions, consult a qualified healthcare provider.
```

This is exact required text. Do not paraphrase.

---

## Law

### Jurisdiction requirements

Every law skill must state its jurisdiction scope explicitly in the first paragraph or in a dedicated "Scope" line:

```markdown
**Scope:** US law (federal and Delaware state). General principles may differ in other jurisdictions.
```

```markdown
**Scope:** EU General Data Protection Regulation (GDPR). Non-EU practitioners should verify local equivalents.
```

```markdown
**Scope:** General common law contract principles. Specific jurisdictions vary.
```

Omitting jurisdiction is a REJECT finding. "International" or "global" is not a valid jurisdiction — law is jurisdiction-specific.

### What NOT to include

- Advice on specific legal situations ("based on what you described, you should...")
- Interpretation of specific contract clauses in a reader's context
- Recommendation to sign, reject, or litigate in a specific situation

### Credible sources

| Source type | Examples |
|-------------|---------|
| Bar association standards | ABA Model Rules of Professional Conduct |
| Practice standards | BigLaw firm practice guides (Latham, Skadden, Kirkland) |
| Model agreements | ISDA Master Agreement, NVCA model term sheets, YC SAFE |
| Regulatory bodies | SEC guidance, FTC regulations, GDPR official text |
| Legal academic | Restatements of Law (American Law Institute) |

Avoid: individual attorney blogs, LegalZoom articles, general "legal tips" content.

### Required footer

Every law skill must end with:

```markdown
> This is educational content, not legal advice. Consult qualified legal counsel for your situation.
```

This is exact required text. Do not paraphrase.

---

## Finance / Investing

### Performance disclaimers

Where a skill references historical returns, performance data, or investment outcomes, include:

```markdown
> Past performance does not guarantee future results.
```

This is in addition to the required footer — it goes adjacent to the specific claim, not only at the end.

### What NOT to include

- Specific ticker symbols or investment recommendations
- Asset allocation advice tailored to individual circumstances
- Tax advice specific to a jurisdiction or individual situation
- Claims that a specific investment strategy will produce specific returns

### Credible sources

| Source type | Examples |
|-------------|---------|
| Professional standards | CFA Institute Standards of Practice |
| Investment firms | Bridgewater research, Goldman Sachs investment frameworks |
| Academic finance | Fama-French factor research, Shiller CAPE data |
| Regulatory bodies | SEC guidelines, FINRA rules, FCA (UK) standards |
| Central banks | Federal Reserve research, BIS working papers |

Avoid: personal finance influencer content, proprietary fund marketing materials, cryptocurrency project whitepapers.

### Required footer

Every finance / investing skill must end with:

```markdown
> This is educational content, not financial advice. Consult a licensed financial advisor for personal decisions.
```

This is exact required text. Do not paraphrase.

---

## Psychology / Mental Health

### What NOT to include

- Specific medication recommendations or dosages
- Diagnostic criteria applied to describe a reader's likely condition
- Treatment protocols prescribed for specific people ("if you feel X, do Y for Z weeks")
- Claims that a psychological technique will resolve a specific clinical condition

Encoding general psychological principles (e.g., cognitive reframing as a general technique) is acceptable. Directing a reader to apply clinical interventions to themselves is not.

### Credible sources

| Source type | Examples |
|-------------|---------|
| Clinical guidelines | APA Clinical Practice Guidelines |
| Evidence-based therapy | Beck Institute CBT research, DBT evidence base (Linehan) |
| Research bodies | NIMH, NIDA, WHO Mental Health Gap Action Programme |
| Academic journals | JAMA Psychiatry, Psychological Medicine, Journal of Abnormal Psychology |

Avoid: self-help books without clinical backing, wellness blogs, pop psychology content.

### Required footer

Every psychology / mental health skill must end with:

```markdown
> For mental health concerns, consult a qualified mental health professional.
```

This is exact required text. Do not paraphrase.

---

## Finding credible sources for regulated domains

| Domain | Start here |
|--------|-----------|
| Health / fitness | ACSM, NSCA, Cochrane Library |
| Clinical medicine | WHO, NIH PubMed, Mayo Clinic |
| Law (US) | ABA, Restatements, model agreements (NVCA, YC, ISDA) |
| Finance | CFA Institute, Federal Reserve, BIS |
| Investing | Fama-French research library, Shiller data, NBER working papers |
| Psychology | APA, Beck Institute, NIMH |
| Mental health | WHO MhGAP, APA Practice Guidelines, NICE (UK) |

**If you can't find a guideline from a major institution for a regulated domain claim, the claim may not be majority-adopted at the top tier.** Reconsider whether the practice qualifies before writing the skill. A practice that credentialed professionals don't follow at majority is not a grimoire skill.

---

## Quick checklist for regulated domain skills

Before running `review-skill`:

- [ ] Health: all factual claims carry `[RCT]`, `[SR]`, `[Consensus]`, or `[Practical]` tags
- [ ] Health: no medication doses, diagnoses, or personalized treatment protocols
- [ ] Law: jurisdiction stated explicitly in first paragraph
- [ ] Law: no advice on specific situations
- [ ] Finance: performance disclaimer on any historical return claim
- [ ] Finance: no specific ticker symbols or personalized allocation advice
- [ ] Psychology: no diagnostic language applied to the reader
- [ ] Psychology: no clinical treatment protocols prescribed
- [ ] All: required footer present (exact text, not paraphrased)
- [ ] All: source is from the credible source list above (not a blog or influencer)
