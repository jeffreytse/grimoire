---
name: review-deployment-checklist
description: Use when verifying a deployment meets all pre-release criteria before pushing to production.
source: Google SRE Book, "Site Reliability Engineering", 2016
tags: [deployment, reliability, devops, error-reduction]
emerging: true
---

# Review Deployment Checklist

Verify every pre-release criterion before a production deployment.

## Why This Is Best Practice

**Adopted by:** Google, Netflix, Stripe.
**Impact:** Google SRE checklists reduced production incidents by 30% (Google SRE Book, 2016).
**Why best:** Eliminates recall failures under time pressure.

Sources: Google SRE Book (2016)

## Steps

1. Run the full test suite and confirm all tests pass.
2. Verify rollback procedure is documented and tested.
3. Confirm monitoring alerts are active for the affected services.
