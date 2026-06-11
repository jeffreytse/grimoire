---
name: apply-smart-contract-access-control
description: Use when writing Solidity contracts with admin functions, upgradeable proxies, or multi-role permissions — implementing Ownable, role-based access control, and two-step ownership transfer to prevent unauthorized function execution.
source: 'OWASP Smart Contract Top 10 SC01 SC06 (owasp.org/www-project-smart-contract-top-10/); SWC-105 SWC-106; OpenZeppelin AccessControl documentation; Ethereum Foundation security documentation'
tags: [security, owasp, solidity, smart-contracts, access-control, rbac, ethereum, blockchain, developer]
---

# Apply Smart Contract Access Control

Restrict privileged contract functions with OpenZeppelin's Ownable and AccessControl — using two-step ownership transfer, role-based permissions, and timelocks — preventing unauthorized fund withdrawal, parameter manipulation, and upgrade hijacking.

## Why This Is Best Practice

**Adopted by:** OWASP Smart Contract Top 10 SC01 (Access Control) and SC06 (Vulnerable Access Control). SWC-105 (Unprotected Ether Withdrawal) and SWC-106 (Unprotected SELFDESTRUCT Instruction) are the canonical CWEs. OpenZeppelin's Contracts (used by Uniswap, Compound, Aave) provides `Ownable2Step` and `AccessControl` as the standard implementations. Ethereum Foundation security documentation mandates access control for all privileged state changes.
**Impact:** SWC-105 is the most commonly exploited smart contract vulnerability class. The 2022 Nomad Bridge hack ($190M) exploited an access control bug that allowed any address to call an initialization function. The 2021 Poly Network hack ($611M) exploited an access control flaw in the `EthCrossChainManager` contract. The 2022 Wormhole hack ($320M) involved an unguarded `complete_wrapped` function. Missing `onlyOwner` on a withdrawal function typically results in complete fund loss.
**Why best:** Solidity's default function visibility is `public` — all functions are callable by any address unless explicitly restricted. Adding `require(msg.sender == owner)` manually is error-prone and doesn't handle role hierarchies. OpenZeppelin's `AccessControl` provides audited, battle-tested RBAC with events for all role grants/revocations — the audit trail is essential for forensic analysis after incidents.

Sources: OWASP Smart Contract Top 10 SC01, SC06; SWC-105, SWC-106; OpenZeppelin Contracts documentation; Ethereum Smart Contract Best Practices

## Steps

1. **Use `Ownable2Step` — two-step ownership transfer prevents accidental loss**:

   ```solidity
   // SPDX-License-Identifier: MIT
   pragma solidity ^0.8.19;

   import "@openzeppelin/contracts/access/Ownable2Step.sol";

   contract Protocol is Ownable2Step {
       uint256 public fee;

       constructor(address initialOwner) Ownable(initialOwner) {}

       // Restricted to owner only
       function setFee(uint256 newFee) external onlyOwner {
           require(newFee <= 100, "Fee cannot exceed 100 basis points");
           fee = newFee;
       }

       function emergencyWithdraw() external onlyOwner {
           (bool ok, ) = owner().call{value: address(this).balance}("");
           require(ok, "Transfer failed");
       }
   }

   // Ownable2Step: transferOwnership sets pending owner, new owner must call acceptOwnership()
   // Prevents transferring to wrong address (which would permanently lock admin functions)
   ```

2. **Use `AccessControl` for multi-role systems**:

   ```solidity
   // SPDX-License-Identifier: MIT
   pragma solidity ^0.8.19;

   import "@openzeppelin/contracts/access/AccessControl.sol";

   contract DeFiVault is AccessControl {
       bytes32 public constant ADMIN_ROLE = keccak256("ADMIN_ROLE");
       bytes32 public constant OPERATOR_ROLE = keccak256("OPERATOR_ROLE");
       bytes32 public constant PAUSER_ROLE = keccak256("PAUSER_ROLE");

       constructor(address admin) {
           _grantRole(DEFAULT_ADMIN_ROLE, admin);  // can grant/revoke other roles
           _grantRole(ADMIN_ROLE, admin);
       }

       function setFeeRecipient(address recipient) external onlyRole(ADMIN_ROLE) {
           // only ADMIN_ROLE can change fee recipient
       }

       function processWithdrawal(address user) external onlyRole(OPERATOR_ROLE) {
           // operators can process but cannot change protocol parameters
       }

       function pause() external onlyRole(PAUSER_ROLE) {
           // separate pause role for emergency response
       }
   }
   ```

3. **Protect initialization functions — prevent re-initialization**:

   ```solidity
   import "@openzeppelin/contracts/proxy/utils/Initializable.sol";

   contract UpgradeableProtocol is Initializable, AccessControl {
       uint256 public fee;

       // initializer modifier ensures this can only be called once
       function initialize(address admin, uint256 initialFee) external initializer {
           _grantRole(DEFAULT_ADMIN_ROLE, admin);
           fee = initialFee;
       }

       // For upgradeable contracts: add reinitializer version guard
       function initializeV2(uint256 newParam) external reinitializer(2) {
           // only runs on upgrade to V2, not callable again after
       }
   }
   ```

4. **Use a Timelock for high-impact admin operations**:

   ```solidity
   import "@openzeppelin/contracts/governance/TimelockController.sol";

   // Deployer: create timelock with 48-hour delay for all critical ops
   address[] memory proposers = new address[](1);
   address[] memory executors = new address[](1);
   proposers[0] = multiSigAddress;
   executors[0] = multiSigAddress;

   TimelockController timelock = new TimelockController(
       48 hours,     // minDelay — operations wait 48 hours before execution
       proposers,
       executors,
       address(0)   // no admin — timelock is self-administered
   );

   // Then set the timelock as the owner of your protocol contracts
   protocol.transferOwnership(address(timelock));
   ```

5. **Emit events on all role changes for auditability**:

   ```solidity
   // OpenZeppelin AccessControl emits RoleGranted/RoleRevoked automatically
   // For custom access control, always emit:
   event OwnershipTransferred(address indexed previousOwner, address indexed newOwner);
   event RoleGranted(bytes32 indexed role, address indexed account);
   event RoleRevoked(bytes32 indexed role, address indexed account);

   // Monitor these events off-chain for unauthorized role changes
   ```

6. **Use multi-sig for protocol admin keys**:

   ```
   For production protocols:
   - Deploy Gnosis Safe (multi-sig) as the contract owner
   - Require M-of-N signers (e.g., 3-of-5) for admin transactions
   - Never use a single EOA as the sole owner of a protocol with TVL > $100k
   - Hardware wallets for all multi-sig signers
   ```

## Rules

- Never leave `initialize()` functions callable after deployment — use `initializer` modifier or `_disableInitializers()` in the constructor of implementation contracts.
- `DEFAULT_ADMIN_ROLE` in OpenZeppelin `AccessControl` can grant any role — restrict it to a multi-sig, not a developer EOA.
- Test that all privileged functions revert when called by an unauthorized address — add a negative test for each `onlyRole`/`onlyOwner` function.
- Timelock delays should match the severity of the operation: parameter changes (24–48h), contract upgrades (72h+), emergency pause (0h with separate pause role).

## Common Mistakes

- **`Ownable` single-step transfer** — `transferOwnership(wrongAddress)` permanently loses admin control; use `Ownable2Step`.
- **Constructor sets owner to `msg.sender` which is a deploy script** — constructor ownership goes to the deployer account; immediately transfer to a multi-sig post-deployment.
- **`onlyOwner` on functions that should be callable by users** — access control bugs go both ways: too permissive (no guard) and too restrictive (wrong guard).
- **Not testing with a non-owner address** — Foundry: use `vm.prank(attacker)` to call owner-only functions and verify revert.
