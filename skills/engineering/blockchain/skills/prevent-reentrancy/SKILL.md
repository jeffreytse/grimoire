---
name: prevent-reentrancy
description: Use when writing Solidity smart contracts that transfer ETH or call external contracts — applying the Checks-Effects-Interactions pattern and ReentrancyGuard to prevent reentrancy attacks.
source: 'OWASP Smart Contract Top 10 SC05 (owasp.org/www-project-smart-contract-top-10/); SWC-107 (Smart Contract Weakness Classification); OpenZeppelin ReentrancyGuard documentation; Ethereum Foundation security documentation'
tags: [security, owasp, solidity, smart-contracts, reentrancy, ethereum, blockchain, developer]
---

# Prevent Reentrancy

Apply Checks-Effects-Interactions (CEI) pattern and OpenZeppelin's ReentrancyGuard to prevent reentrancy attacks — the attack vector behind the $60M DAO hack and dozens of subsequent DeFi exploits.

## Why This Is Best Practice

**Adopted by:** OWASP Smart Contract Top 10 SC05 (Reentrancy Attacks). SWC-107 (Smart Contract Weakness Classification Registry) is the authoritative vulnerability taxonomy. OpenZeppelin's Contracts library (used in 70%+ of DeFi protocols including Uniswap, Aave, Compound) provides `ReentrancyGuard` as the standard defense. The Ethereum Foundation's Smart Contract Security Best Practices mandate CEI pattern for all external calls.
**Impact:** The 2016 DAO hack exploited reentrancy to drain 3.6M ETH (~$60M at the time, ~$11B at 2021 peak prices), triggering the Ethereum hard fork. Since 2020: Cream Finance lost $18.8M (2021), Fei Protocol lost $80M (2022), and Euler Finance lost $197M (2023) — all reentrancy variants. DeFi reentrancy losses exceed $1B according to Rekt.news (2023). A single missing `nonReentrant` modifier causes complete fund drainage in protocols with external calls.
**Why best:** CEI pattern (validate, then update state, then call external) ensures that when a malicious contract re-enters the function, the state is already updated and the re-entry path fails the checks. The alternative (optimistic pattern: call first, update state after) allows repeated withdrawals before balance is decremented — the exact DAO vulnerability. ReentrancyGuard adds a mutex that catches cross-function reentrancy which CEI alone doesn't prevent.

Sources: OWASP Smart Contract Top 10 SC05; SWC-107; OpenZeppelin ReentrancyGuard source; Ethereum Foundation Smart Contract Best Practices

## Steps

1. **Apply Checks-Effects-Interactions pattern — update state before external calls**:

   ```solidity
   // SPDX-License-Identifier: MIT
   pragma solidity ^0.8.19;

   // BAD — vulnerable to reentrancy
   contract VulnerableVault {
       mapping(address => uint256) public balances;

       function withdraw() external {
           uint256 amount = balances[msg.sender];
           require(amount > 0, "No balance");

           // INTERACTION before EFFECT — reentrancy window here
           (bool success, ) = msg.sender.call{value: amount}("");
           require(success, "Transfer failed");

           balances[msg.sender] = 0;  // too late — state updated after external call
       }
   }

   // GOOD — CEI pattern
   contract SecureVault {
       mapping(address => uint256) public balances;

       function withdraw() external {
           // CHECK
           uint256 amount = balances[msg.sender];
           require(amount > 0, "No balance");

           // EFFECT — update state before external call
           balances[msg.sender] = 0;

           // INTERACTION — external call last
           (bool success, ) = msg.sender.call{value: amount}("");
           require(success, "Transfer failed");
       }
   }
   ```

2. **Use OpenZeppelin ReentrancyGuard for all external-calling functions**:

   ```solidity
   // SPDX-License-Identifier: MIT
   pragma solidity ^0.8.19;

   import "@openzeppelin/contracts/utils/ReentrancyGuard.sol";

   contract DeFiProtocol is ReentrancyGuard {
       mapping(address => uint256) public deposits;
       IERC20 public token;

       function deposit(uint256 amount) external nonReentrant {
           token.transferFrom(msg.sender, address(this), amount);
           deposits[msg.sender] += amount;
       }

       function withdraw(uint256 amount) external nonReentrant {
           require(deposits[msg.sender] >= amount, "Insufficient balance");
           deposits[msg.sender] -= amount;  // EFFECT before INTERACTION
           token.transfer(msg.sender, amount);
       }
   }
   ```

3. **Protect cross-function reentrancy — same mutex across related functions**:

   ```solidity
   // Cross-function reentrancy: attacker calls withdrawETH → enters depositETH during callback
   contract CrossFunctionSafe is ReentrancyGuard {
       mapping(address => uint256) public ethBalance;
       mapping(address => uint256) public tokenBalance;

       // Both functions protected by same mutex
       function withdrawETH(uint256 amount) external nonReentrant {
           require(ethBalance[msg.sender] >= amount);
           ethBalance[msg.sender] -= amount;
           (bool ok, ) = msg.sender.call{value: amount}("");
           require(ok);
       }

       function depositToken(uint256 amount) external nonReentrant {
           // If called during withdrawETH callback, mutex blocks this
           tokenBalance[msg.sender] += amount;
       }
   }
   ```

4. **Use pull-over-push payment pattern for ETH distribution**:

   ```solidity
   // GOOD — pull pattern: recipient claims funds instead of contract pushing
   contract PullPayment {
       mapping(address => uint256) private _pendingWithdrawals;

       function _asyncTransfer(address recipient, uint256 amount) internal {
           _pendingWithdrawals[recipient] += amount;
       }

       function withdrawPayments() external nonReentrant {
           uint256 payment = _pendingWithdrawals[msg.sender];
           require(payment > 0, "No pending payment");
           _pendingWithdrawals[msg.sender] = 0;
           (bool ok, ) = msg.sender.call{value: payment}("");
           require(ok, "Transfer failed");
       }
   }
   ```

5. **Test reentrancy attacks in your test suite**:

   ```solidity
   // Test helper: malicious contract that re-enters on receive
   contract ReentrancyAttacker {
       IVault public target;
       uint256 public attackCount;

       receive() external payable {
           if (attackCount < 5 && address(target).balance >= 1 ether) {
               attackCount++;
               target.withdraw();  // re-enter
           }
       }

       function attack() external payable {
           target.deposit{value: msg.value}();
           target.withdraw();
       }
   }
   ```

   ```javascript
   // Hardhat test
   it("should prevent reentrancy", async function () {
     const attacker = await ReentrancyAttacker.deploy(vault.address);
     await attacker.attack({ value: ethers.parseEther("1") });
     expect(await attacker.attackCount()).to.equal(0);  // reentrancy blocked
   });
   ```

## Rules

- CEI order is: Checks (require/revert), then Effects (state changes), then Interactions (external calls/transfers) — never call external contracts before updating state.
- `transfer()` and `send()` have a 2300 gas stipend that prevents most reentrancy — but this is not reliable post-EIP-1884 and should not replace CEI or ReentrancyGuard.
- Read-only reentrancy (attacker reads state during callback to manipulate price oracle) is not prevented by ReentrancyGuard — also apply view function guards on price-sensitive reads.
- Every function that interacts with an external contract or transfers ETH must have `nonReentrant`.

## Common Mistakes

- **CEI without ReentrancyGuard on multi-function contracts** — CEI prevents same-function reentrancy; cross-function reentrancy requires the mutex.
- **Using `transfer()` and believing it's safe** — the 2300 gas stipend is not guaranteed in all contexts (proxies, `receive()` hooks with state writes).
- **Not testing with a malicious contract** — unit tests with EOA (externally owned accounts) don't trigger reentrancy; test with a contract that implements `receive()` with re-entry logic.
- **Forgetting reentrancy in ERC777/ERC1155 token hooks** — `_beforeTokenTransfer` and `tokensReceived` hooks execute external code during transfers, creating reentrancy windows even without explicit `.call()`.
