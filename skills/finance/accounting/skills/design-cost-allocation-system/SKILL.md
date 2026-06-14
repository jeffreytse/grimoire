---
name: design-cost-allocation-system
description: Use when assigning overhead or shared costs to products, services, departments, or customers — e.g., "how do I allocate overhead?", "activity-based costing vs. traditional?", "which products are actually profitable?", "cost per customer segment?"
source: Kaplan & Cooper "Cost and Effect" (1998, activity-based costing); CIMA management accounting standards; FASB overhead allocation guidance; Harvard Business School case studies on ABC implementation
tags: [finance, accounting, cost-allocation, activity-based-costing, overhead, management-accounting, profitability]
verified: true
---

# Design Cost Allocation System

Build a structured approach for allocating indirect costs to cost objects (products, services, departments, customers) to reveal true profitability and inform pricing decisions.

## Why This Is Best Practice

**Adopted by:** Activity-Based Costing (ABC) was developed at Harvard Business School (Kaplan & Cooper, 1988) and adopted by major manufacturers including Caterpillar, HP, and Boeing to understand true product profitability. The CIMA (Chartered Institute of Management Accountants) includes cost allocation as a core competency. Every cost accounting system — from QuickBooks to SAP — is built around cost allocation principles.
**Impact:** Kaplan & Cooper documented that traditional volume-based allocation systematically over-costs high-volume simple products and under-costs low-volume complex products — causing companies to reprice wrongly, over-invest in unprofitable products, and under-invest in genuinely profitable ones. Companies that switched to ABC-based allocation in documented case studies (Siemens, Kanthal, Hewlett-Packard) revealed that 20–40% of products were losing money when true costs were allocated.
**Why best:** Unallocated or badly allocated overhead makes financial statements look like cost centers while operations management is flying blind on true profitability. A well-designed cost allocation system converts a P&L into a decision tool: which products to grow, which to discontinue, which customers to serve at what price.

## Steps

1. **Identify cost objects** — What do you need to know the cost of? Products, SKUs, service lines, customers, channels, departments. Be specific; "the business" is not a cost object.

2. **Separate direct and indirect costs** — Direct: material, direct labor — traceable with reasonable effort to the cost object. Indirect (overhead): utilities, rent, depreciation, management salaries, IT — must be allocated.

3. **Choose the allocation method:**
   - **Traditional (volume-based)**: allocate overhead as a % of direct labor hours, machine hours, or revenue. Simple; accurate only when all products consume overhead in proportion to volume. Fails in diverse product/service environments.
   - **Activity-Based Costing (ABC)**: identify cost drivers (what actually causes overhead to be incurred?), create activity pools, allocate based on actual consumption. More complex; dramatically more accurate for diverse products.
   - **Direct method**: allocate service department costs directly to production departments (no inter-department allocation). Simple; loses accuracy when service departments heavily serve each other.
   - **Step-down method**: allocate service department costs in sequence, recognizing partial inter-department services. Better than direct; less accurate than reciprocal.

4. **Implement ABC (if applicable):**
   a. Identify activities: order processing, machine setup, quality inspection, customer support, delivery.
   b. Assign costs to activity pools: salary for customer support reps → customer support pool.
   c. Identify cost drivers for each pool: order processing → number of orders; setup → number of setups.
   d. Calculate cost per driver unit: customer support pool $500k ÷ 5,000 support tickets = $100/ticket.
   e. Assign costs to products/customers: Product A required 200 tickets → $20,000 customer support cost allocated.

5. **Validate the model** — Total allocated cost must equal total overhead incurred (no over/under allocation at year-end). If using standard rates: reconcile actual vs. standard overhead quarterly; adjust for significant variances.

6. **Calculate product/customer profitability** — Revenue − direct costs − allocated overhead = true product margin. Rank products/customers by true profitability. Identify the 20% generating 120% of profits and the 20% losing money (the "whale curve" is common).

7. **Use for pricing and mix decisions** — Products below their full cost are subsidized by profitable ones. Pricing below full cost is only defensible if (a) it's strategic (new customer acquisition) or (b) fixed costs are truly unavoidable and contribution margin is positive.

## Rules

- An allocation is always arbitrary to some degree — the goal is "good enough to make better decisions," not accounting perfection.
- Review allocation bases annually — a company that grew its e-commerce channel from 5% to 40% of revenue has fundamentally different cost drivers than when the allocation system was designed.
- Allocated costs should not be used for short-run decisions — in the short run, fixed overhead is sunk; contribution margin (price − variable cost) is the relevant metric.
- Activity-based costing is worth implementing when: overhead > 30% of total cost, product diversity is high, and pricing decisions are consequential.

## Examples

**Manufacturer with two products — traditional vs. ABC:**
Traditional allocation (% of direct labor): Product A (high-volume, simple) gets 70% of overhead; Product B (low-volume, complex) gets 30%.
ABC reveals: Product B requires 3× the machine setups, 5× the inspection hours, 2× the purchasing transactions.
ABC reallocation: Product B actually consumes 55% of overhead; Product A: 45%.
Result: Product B is priced 20% below true cost; Product A is overpriced vs. competition. Pricing corrections follow.

## Common Mistakes

- **Allocating overhead as a flat % of revenue** — High-revenue products absorb more overhead even if they're operationally simple; this systematically misprices complex low-revenue products.
- **Over-engineering the ABC system** — 50 activity pools and 30 cost drivers creates a maintenance burden that exceeds its value. 10–15 well-chosen activity pools handles 90% of the accuracy gain.
- **Using fully allocated cost for make-vs-buy decisions** — If overhead is fixed regardless of the decision, only variable cost is relevant. Fully allocated cost used for short-run decisions leads to outsourcing profitable activities.

---

> **Finance disclaimer:** This skill encodes professional best practices for educational purposes. It is not financial advice. Consult a licensed financial advisor before making investment decisions.
