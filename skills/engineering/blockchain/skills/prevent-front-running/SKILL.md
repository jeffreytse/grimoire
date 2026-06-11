---
name: prevent-front-running
description: Use when writing DeFi smart contracts with user-submitted transactions — implementing commit-reveal, slippage tolerance, and MEV-resistant patterns to prevent front-running and sandwich attacks.
source: 'OWASP Smart Contract Top 10 SC08 (owasp.org/www-project-smart-contract-top-10/); Ethereum Foundation MEV documentation; Flashbots research; Uniswap slippage documentation'
tags: [security, owasp, solidity, smart-contracts, mev, front-running, defi, ethereum, blockchain, developer]
---

# Prevent Front Running

Implement commit-reveal, strict slippage bounds, and MEV-resistant submission via Flashbots to prevent front-running, sandwich attacks, and transaction reordering that extracts value from users.

## Why This Is Best Practice

**Adopted by:** OWASP Smart Contract Top 10 SC08 (Front-Running). Flashbots (ethereum.org/en/developers/docs/mev) is the canonical solution for MEV protection, used by Ethereum validators and adopted by Uniswap, Curve, and 1inch for protected routing. Uniswap V3's `amountOutMinimum` parameter is the standard slippage protection pattern. The Ethereum Foundation's MEV documentation mandates commit-reveal for any outcome dependent on transaction order.
**Impact:** Flashbots' MEV-Explore dashboard estimates cumulative MEV extracted from Ethereum exceeds $1.5B (2020–2024). Sandwich attacks are the most common user-facing MEV: an attacker sees a large swap in the mempool, buys before it (raising the price), lets the victim's swap execute at a worse price, then sells — consistently extracting value from victims with no capital risk. JaredFromSubway.eth (a prominent MEV bot) extracted $40M in 2023 from sandwich attacks alone. The 2022 Solana NFT mint exploits used front-running bots to win "fair" lotteries.
**Why best:** Transaction ordering in public mempools is not first-come-first-served — miners/validators can reorder, insert, or censor transactions within a block. Slippage bounds limit the price a swap can accept, making sandwiches unprofitable. Commit-reveal hides transaction intent until the block is sealed. Private mempool submission (Flashbots, MEV Blocker) routes transactions directly to validators, bypassing the public mempool entirely.

Sources: OWASP Smart Contract Top 10 SC08; Flashbots MEV documentation; Ethereum Foundation MEV guide; Uniswap V3 router documentation

## Steps

1. **Enforce strict slippage tolerance in AMM swaps**:

   ```solidity
   // SPDX-License-Identifier: MIT
   pragma solidity ^0.8.19;

   import "@uniswap/v3-periphery/contracts/interfaces/ISwapRouter.sol";

   contract SlippageProtectedSwap {
       ISwapRouter public immutable router;

       constructor(address _router) {
           router = ISwapRouter(_router);
       }

       function swapWithProtection(
           address tokenIn,
           address tokenOut,
           uint256 amountIn,
           uint256 minAmountOut,  // caller must specify minimum — never 0
           uint24 fee
       ) external returns (uint256 amountOut) {
           require(minAmountOut > 0, "Slippage protection required");

           ISwapRouter.ExactInputSingleParams memory params = ISwapRouter.ExactInputSingleParams({
               tokenIn: tokenIn,
               tokenOut: tokenOut,
               fee: fee,
               recipient: msg.sender,
               deadline: block.timestamp + 20 minutes,  // also require deadline
               amountIn: amountIn,
               amountOutMinimum: minAmountOut,  // reverts if price moves beyond tolerance
               sqrtPriceLimitX96: 0
           });

           amountOut = router.exactInputSingle(params);
       }
   }
   ```

2. **Implement commit-reveal for order-sensitive operations**:

   ```solidity
   contract CommitRevealAuction {
       struct Commitment {
           bytes32 hash;
           uint256 blockNumber;
       }

       mapping(address => Commitment) public commitments;
       uint256 private constant COMMIT_PERIOD = 20;   // blocks
       uint256 private constant REVEAL_DEADLINE = 40; // blocks after commit

       // Phase 1: Submit hash(bid, salt) — intent hidden from mempool watchers
       function commit(bytes32 commitment) external {
           commitments[msg.sender] = Commitment({
               hash: commitment,
               blockNumber: block.number
           });
       }

       // Phase 2: Reveal actual bid after commit period
       function reveal(uint256 bid, bytes32 salt) external {
           Commitment memory c = commitments[msg.sender];
           require(block.number >= c.blockNumber + COMMIT_PERIOD, "Commit period not ended");
           require(block.number <= c.blockNumber + REVEAL_DEADLINE, "Reveal deadline passed");
           require(
               keccak256(abi.encodePacked(bid, salt, msg.sender)) == c.hash,
               "Invalid reveal"
           );
           // Process bid — all bids revealed before winner determined
           _processBid(msg.sender, bid);
           delete commitments[msg.sender];
       }
   }
   ```

3. **Add deadline parameter to time-sensitive transactions**:

   ```solidity
   function executeOperation(
       address asset,
       uint256 amount,
       uint256 deadline  // always require explicit deadline
   ) external {
       require(block.timestamp <= deadline, "Transaction expired");
       // Prevents delayed mining of a transaction when market conditions change
   }
   ```

4. **Use Flashbots for private transaction submission (off-chain)**:

   ```javascript
   // ethers.js + Flashbots SDK
   import { FlashbotsBundleProvider } from "@flashbots/ethers-provider-bundle";
   import { ethers } from "ethers";

   const provider = new ethers.providers.JsonRpcProvider(RPC_URL);
   const flashbotsProvider = await FlashbotsBundleProvider.create(
     provider,
     ethers.Wallet.createRandom(),  // reputation signer
     "https://relay.flashbots.net"
   );

   const signedBundle = await flashbotsProvider.signBundle([{
     signer: wallet,
     transaction: {
       chainId: 1,
       to: CONTRACT_ADDRESS,
       data: swapCalldata,
       gasLimit: 200000,
       maxFeePerGas: ethers.utils.parseUnits("50", "gwei"),
       maxPriorityFeePerGas: ethers.utils.parseUnits("2", "gwei"),
     }
   }]);

   // Submit directly to validators — bypasses public mempool
   const bundleSubmission = await flashbotsProvider.sendRawBundle(
     signedBundle,
     targetBlockNumber
   );
   ```

5. **For NFT mints: random assignment after reveal**:

   ```solidity
   // Prevent sniper bots from selecting rare traits at mint time
   contract AntiSnipeMint {
       bool public revealed = false;
       uint256 public revealBlock;
       uint256 public revealSeed;

       function mint(uint256 quantity) external payable {
           // Mint token IDs without revealing traits
           _mintTokens(msg.sender, quantity);
       }

       function reveal() external {
           require(!revealed && block.number > revealBlock);
           // Use future block hash — not known at mint time
           revealSeed = uint256(blockhash(revealBlock));
           revealed = true;
       }

       function getTraits(uint256 tokenId) external view returns (uint256) {
           require(revealed, "Not yet revealed");
           return uint256(keccak256(abi.encodePacked(revealSeed, tokenId))) % NUM_TRAITS;
       }
   }
   ```

## Rules

- Never use `amountOutMinimum: 0` in swap calls — this accepts any price including one manipulated by a sandwicher.
- Deadline parameters must be caller-supplied and checked against `block.timestamp` — a transaction sitting in the mempool for hours may execute under completely different market conditions.
- Commit-reveal salts must be random and kept secret until reveal — a guessable or reused salt defeats the commitment scheme.
- For large trades, consider using aggregators with MEV protection (CoW Protocol, 1inch Fusion) rather than routing directly through DEXes.

## Common Mistakes

- **`deadline: block.timestamp`** — using the current block's timestamp as the deadline means the transaction never expires, defeating the purpose.
- **Small slippage only for large swaps** — a 0.5% slippage tolerance on a $1M swap is still $5,000 of extractable value; consider splitting large swaps.
- **Commit-reveal without salt** — `keccak256(bid)` is brute-forceable for small bid spaces; always combine with a large random salt.
- **Not testing with a fork that simulates MEV** — Foundry supports mainnet fork testing; simulate front-running by having a bot contract observe and front-run your test transaction.
