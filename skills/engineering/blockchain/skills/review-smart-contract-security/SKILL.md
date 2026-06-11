---
name: review-smart-contract-security
description: Use when auditing a Solidity smart contract or reviewing a protocol before deployment — systematically checking all 10 OWASP Smart Contract Top 10 vulnerability classes with Slither/Foundry test procedures.
source: 'OWASP Smart Contract Top 10 (owasp.org/www-project-smart-contract-top-10/); SWC Registry (swcregistry.io); Slither documentation; Trail of Bits audit methodology'
tags: [security, owasp, solidity, smart-contracts, audit, ethereum, blockchain]
---

# Review Smart Contract Security

Audit smart contracts against the OWASP Smart Contract Top 10 using Slither static analysis and Foundry invariant testing — covering SC01 through SC10 with specific test procedures, detection commands, and remediation steps.

## Why This Is Best Practice

**Adopted by:** OWASP Smart Contract Top 10 is the authoritative vulnerability taxonomy. Trail of Bits, OpenZeppelin Security, and Consensys Diligence — the three most prominent smart contract audit firms — all use structured checklists covering SC01–SC10 equivalent categories. Slither (Trail of Bits) is the standard static analysis tool, used in Ethereum Foundation security tooling. MakerDAO, Compound, and Uniswap conduct external audits plus internal pre-deployment reviews against these classes before each major release.
**Impact:** Smart contract audit findings consistently cluster in these 10 categories. Immunefi's "DeFi Bug Bounty Report" (2023) found that 73% of critical findings were in categories covered by SC01–SC10. Protocols that undergo structured audits have 4× lower incident rates than unaudited protocols (DeFi Safety score analysis, 2022). The $1B+ in DeFi losses attributed to smart contract exploits in 2022–2023 would have been significantly reduced by systematic SC01–SC10 review.
**Why best:** Manual code review without a structured checklist misses vulnerability classes — a reviewer focused on reentrancy may overlook integer overflow or oracle manipulation. The SC01–SC10 framework provides completeness; Slither and Foundry provide automated detection for the mechanical findings, freeing human reviewers to focus on business logic vulnerabilities that tools miss.

Sources: OWASP Smart Contract Top 10; SWC Registry; Trail of Bits Slither documentation; Immunefi Bug Bounty Report (2023)

## Steps

### Pre-Audit Setup

```bash
# Install analysis tools
pip install slither-analyzer
forge install foundry-rs/forge-std  # Foundry standard library

# Run Slither on the contract directory
slither . --checklist --markdown-root "contracts/"

# Run Slither detectors individually
slither . --detect reentrancy-eth,reentrancy-no-eth
slither . --detect unprotected-upgrade,suicidal
slither . --detect oracle-manipulation,msg-value-loop
```

### SC01 — Access Control

```bash
# Slither: find unprotected functions
slither . --detect unprotected-upgrade,suicidal,controlled-delegatecall
```

**Check:**
- [ ] All admin/privileged functions have `onlyOwner` or `onlyRole` modifier
- [ ] `initialize()` functions use `initializer` modifier (proxy contracts)
- [ ] `DEFAULT_ADMIN_ROLE` assigned to multi-sig, not EOA
- [ ] Two-step ownership transfer (`Ownable2Step`) for all Ownable contracts

### SC02 — Integer Overflow/Underflow (pre-Solidity 0.8.x)

```bash
slither . --detect overflow-before-cast,tautology
```

**Check:**
- [ ] Solidity version ≥ 0.8.0 (built-in overflow protection) OR using SafeMath
- [ ] No explicit `unchecked {}` blocks without documented justification
- [ ] Type casting from larger to smaller types checked for truncation

### SC03 — Timestamp Dependence

```bash
slither . --detect weak-prng,timestamp
```

**Check:**
- [ ] No `block.timestamp` used for randomness or lottery outcomes
- [ ] Time-locked operations allow ±15 second tolerance for timestamp manipulation
- [ ] No `block.number` as a proxy for time in cross-chain deployments

### SC04 — Reentrancy

```bash
slither . --detect reentrancy-eth,reentrancy-no-eth,reentrancy-benign,reentrancy-events
```

**Check:**
- [ ] CEI pattern in all functions with external calls
- [ ] `ReentrancyGuard.nonReentrant` on all public/external functions making external calls
- [ ] ERC777/ERC1155 token hooks analyzed for reentrancy paths

### SC05 — Unprotected Ether Withdrawal (SWC-105)

```bash
slither . --detect suicidal,locked-ether,arbitrary-send-eth
```

**Check:**
- [ ] `selfdestruct` not callable by unauthorized addresses
- [ ] ETH transfer functions restricted to authorized recipients
- [ ] Contracts that hold ETH can withdraw it (no locked ETH)

### SC06 — Oracle Manipulation

```bash
# Manual check — Slither doesn't fully detect oracle patterns
grep -r "slot0()" contracts/  # spot price reading
grep -r "getReserves()" contracts/  # Uniswap V2 spot price
```

**Check:**
- [ ] No spot price reads from AMM pools (`pool.slot0()`, `pair.getReserves()`)
- [ ] Chainlink feeds have staleness check (`updatedAt + MAX_STALENESS > block.timestamp`)
- [ ] Price deviation circuit breaker for liquidation functions
- [ ] TWAP period ≥ 30 minutes for any price used in financial calculations

### SC07 — Logic Errors (Business Logic Bugs)

```bash
# Foundry invariant testing
forge test --match-test invariant
```

Write invariant tests:
```solidity
// Foundry invariant: total supply never exceeds max
function invariant_totalSupplyLeqMax() external view {
    assertLe(token.totalSupply(), MAX_SUPPLY);
}

// Invariant: vault solvency — assets >= liabilities
function invariant_vaultSolvent() external view {
    assertGe(vault.totalAssets(), vault.totalLiabilities());
}
```

### SC08 — Insecure Randomness

```bash
slither . --detect weak-prng
grep -r "block\.timestamp\|block\.number\|blockhash\|block\.prevrandao" contracts/
```

**Check:**
- [ ] No block variables used as entropy for outcomes with financial value
- [ ] Chainlink VRF used for all on-chain randomness
- [ ] Commit-reveal scheme implemented for user-supplied randomness
- [ ] VRF callback not executable by unauthorized addresses

### SC09 — Gas Limit and Denial of Service

```bash
slither . --detect costly-loop,msg-value-loop
```

**Check:**
- [ ] No unbounded loops over user-supplied or growing arrays
- [ ] Pull-over-push payment pattern for fund distribution
- [ ] External calls in loops (e.g., sending ETH to each participant) refactored to pull
- [ ] `DoS with (Unexpected) Revert` — one failing recipient doesn't block all withdrawals

### SC10 — Front-Running

```bash
grep -rn "amountOutMinimum: 0\|minAmountOut: 0" contracts/
grep -rn "deadline.*block\.timestamp" contracts/
```

**Check:**
- [ ] Swap calls have non-zero `amountOutMinimum`
- [ ] Deadline parameter is user-specified, not `block.timestamp`
- [ ] Sensitive operations use commit-reveal or private mempool submission
- [ ] Mint/auction functions randomize trait assignment after reveal

### Final Checklist

```
Audit Summary:
□ SC01 Access Control — all privileged functions protected
□ SC02 Integer Safety — 0.8.x or SafeMath, no unsafe unchecked
□ SC03 Timestamp — no block vars for randomness
□ SC04 Reentrancy — CEI + nonReentrant on all external calls
□ SC05 Ether Withdrawal — no arbitrary send, no locked ETH
□ SC06 Oracle — TWAP/Chainlink with staleness checks
□ SC07 Logic — invariant tests written and passing
□ SC08 Randomness — Chainlink VRF for all random outcomes
□ SC09 Gas/DoS — no unbounded loops, pull payment pattern
□ SC10 Front-running — slippage protection, deadline params
□ Slither 0 high/medium findings
□ External audit completed (for TVL > $100k)
```

## Rules

- Run Slither before every PR merge — CI integration: `slither . --fail-on high`.
- All high/critical Slither findings must be fixed or documented with a false-positive justification.
- Invariant tests must cover core financial invariants (solvency, supply caps, fee bounds).
- Any protocol with TVL > $100k should have an external professional audit before mainnet deployment.

## Common Mistakes

- **Only auditing the main contract, not libraries and interfaces** — reentrancy and oracle bugs often exist in inherited contracts or libraries.
- **Treating Slither PASS as audit complete** — Slither catches mechanical patterns but misses business logic bugs (wrong formula, wrong price direction) that require human review.
- **Testing only on local fork, not on mainnet fork** — Foundry `vm.createFork(mainnet_rpc)` tests real oracle conditions and real liquidity depths.
- **Not re-auditing after contract changes** — a minor change to fee calculation or access control can introduce new vulnerability paths; re-run the full checklist after any significant change.
