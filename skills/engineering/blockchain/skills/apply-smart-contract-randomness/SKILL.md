---
name: apply-smart-contract-randomness
description: Use when writing Solidity contracts that need randomness — NFT mints, lotteries, games — replacing insecure block variable entropy with Chainlink VRF or commit-reveal to prevent miner manipulation.
source: 'OWASP Smart Contract Top 10 SC10 (owasp.org/www-project-smart-contract-top-10/); SWC-120 (Weak Sources of Randomness); Chainlink VRF documentation; Ethereum Foundation security documentation'
tags: [security, owasp, solidity, smart-contracts, randomness, vrf, nft, ethereum, blockchain, developer]
---

# Apply Smart Contract Randomness

Use Chainlink VRF (Verifiable Random Function) or commit-reveal schemes instead of block variables — preventing validators from manipulating outcomes in lotteries, NFT mints, and on-chain games.

## Why This Is Best Practice

**Adopted by:** OWASP Smart Contract Top 10 SC10 (Insufficient Randomness). SWC-120 (Weak Sources of Randomness from Chain Attributes) is the canonical weakness. Chainlink VRF is the industry standard, used by Axie Infinity, Apes NFT, and dozens of DeFi protocols for verifiable on-chain randomness. The Ethereum Foundation's security documentation explicitly warns against using `block.prevrandao`, `block.timestamp`, or `blockhash` for any outcome with financial value.
**Impact:** The SmartBillions lottery was exploited in 2017 — an attacker manipulated their own transaction to only execute when `block.blockhash` produced a winning result. The Fomo3D game (2018) was gamed by miners who manipulated block timestamps to time their entries. Axie Infinity's early NFT minting using `blockhash` was exploited before migration to Chainlink VRF. Using `keccak256(block.timestamp, msg.sender)` as randomness allows the miner who produces the block to selectively include or exclude transactions based on the outcome.
**Why best:** All EVM block variables (`block.timestamp`, `block.prevrandao`, `blockhash`, `block.difficulty`) are either known in advance by validators or manipulable with modest resources. Chainlink VRF generates randomness off-chain with a cryptographic proof that the result was not manipulated — the proof is verifiable on-chain before the random value is used. Commit-reveal provides a cheaper alternative where manipulation resistance comes from economic cost rather than cryptographic proof.

Sources: OWASP Smart Contract Top 10 SC10; SWC-120; Chainlink VRF v2.5 documentation; Ethereum Foundation Smart Contract Best Practices

## Steps

1. **Use Chainlink VRF v2.5 for verifiable on-chain randomness**:

   ```solidity
   // SPDX-License-Identifier: MIT
   pragma solidity ^0.8.19;

   import "@chainlink/contracts/src/v0.8/vrf/dev/VRFConsumerBaseV2Plus.sol";
   import "@chainlink/contracts/src/v0.8/vrf/dev/interfaces/IVRFCoordinatorV2Plus.sol";

   contract FairLottery is VRFConsumerBaseV2Plus {
       IVRFCoordinatorV2Plus public immutable vrfCoordinator;
       uint256 public immutable subscriptionId;
       bytes32 public immutable keyHash;

       mapping(uint256 => address) private requestToPlayer;
       mapping(address => uint256) public winnings;

       uint32 private constant CALLBACK_GAS_LIMIT = 200_000;
       uint16 private constant REQUEST_CONFIRMATIONS = 3;
       uint32 private constant NUM_WORDS = 1;

       constructor(address coordinator, uint256 subId, bytes32 _keyHash)
           VRFConsumerBaseV2Plus(coordinator)
       {
           vrfCoordinator = IVRFCoordinatorV2Plus(coordinator);
           subscriptionId = subId;
           keyHash = _keyHash;
       }

       function enterLottery() external payable returns (uint256 requestId) {
           require(msg.value == 0.01 ether, "Entry fee required");

           requestId = vrfCoordinator.requestRandomWords(
               VRFV2PlusClient.RandomWordsRequest({
                   keyHash: keyHash,
                   subId: subscriptionId,
                   requestConfirmations: REQUEST_CONFIRMATIONS,
                   callbackGasLimit: CALLBACK_GAS_LIMIT,
                   numWords: NUM_WORDS,
                   extraArgs: VRFV2PlusClient._argsToBytes(
                       VRFV2PlusClient.ExtraArgsV1({nativePayment: false})
                   )
               })
           );
           requestToPlayer[requestId] = msg.sender;
       }

       // Called by Chainlink with cryptographically proven random number
       function fulfillRandomWords(uint256 requestId, uint256[] calldata randomWords) internal override {
           address player = requestToPlayer[requestId];
           uint256 outcome = randomWords[0] % 100;  // 0–99
           if (outcome < 10) {  // 10% win chance
               winnings[player] += address(this).balance / 2;
           }
       }
   }
   ```

2. **Use commit-reveal for cheaper manipulation resistance**:

   ```solidity
   contract CommitRevealRandom {
       mapping(address => bytes32) public commitments;
       mapping(address => uint256) public revealBlock;
       uint256 private constant REVEAL_DELAY = 5;  // must reveal after 5 blocks

       // Phase 1: player commits hash of their secret
       function commit(bytes32 commitment) external {
           commitments[msg.sender] = commitment;
           revealBlock[msg.sender] = block.number + REVEAL_DELAY;
       }

       // Phase 2: player reveals secret, combined with block hash for randomness
       function reveal(uint256 secret) external returns (uint256 randomNumber) {
           require(block.number >= revealBlock[msg.sender], "Too early to reveal");
           require(block.number <= revealBlock[msg.sender] + 250, "Reveal window expired");

           bytes32 commitment = commitments[msg.sender];
           require(keccak256(abi.encodePacked(secret)) == commitment, "Invalid reveal");

           // Combine player secret with blockhash — neither party can manipulate alone
           bytes32 blockHash = blockhash(revealBlock[msg.sender]);
           randomNumber = uint256(keccak256(abi.encodePacked(secret, blockHash)));

           delete commitments[msg.sender];
       }
   }
   ```

3. **Never use these as randomness sources**:

   ```solidity
   // ALL of these are manipulable by validators or predictable:
   uint256 bad1 = uint256(blockhash(block.number));     // always 0 for current block
   uint256 bad2 = block.timestamp;                       // manipulable ±15 seconds
   uint256 bad3 = block.prevrandao;                     // biasable by validators
   uint256 bad4 = uint256(keccak256(abi.encodePacked(block.timestamp, msg.sender)));

   // A validator/miner who sees a losing outcome can simply drop the transaction
   // and re-include it in the next block they produce.
   ```

4. **For NFT mints: use VRF to assign traits after mint, not during**:

   ```solidity
   contract FairNFT {
       uint256[] public tokenIds;
       mapping(uint256 => uint256) public tokenTraits;  // assigned after VRF response

       function mint() external payable returns (uint256 tokenId) {
           tokenId = tokenIds.length;
           tokenIds.push(tokenId);
           // Don't assign traits yet — request VRF
           _requestVRFForToken(tokenId);
       }

       function _fulfillVRF(uint256 tokenId, uint256 randomWord) internal {
           // Assign traits after mint — attacker can't abort the tx after seeing traits
           tokenTraits[tokenId] = randomWord % NUM_TRAIT_COMBINATIONS;
       }
   }
   ```

## Rules

- `block.prevrandao` (formerly `block.difficulty`) is not secure randomness — validators can bias it by choosing when to propose a block.
- VRF requests are asynchronous — design the contract to handle the two-transaction flow (request → callback) rather than expecting immediate randomness.
- Never use the same VRF random word for multiple independent outcomes — use `keccak256(abi.encodePacked(randomWord, index))` to derive multiple values from one VRF response.
- Chainlink VRF requires a funded subscription — test that your contract handles VRF callback failure (e.g., out-of-LINK) gracefully.

## Common Mistakes

- **Using `blockhash(block.number - 1)`** — this is slightly better than current block but still known to the miner who just produced that block before the current miner.
- **Commit-reveal without reveal window expiry** — if the player never reveals, their slot is locked; enforce an expiry after which the commitment is cleared.
- **Requesting multiple random words and treating them as independent** — Chainlink VRF v2 returns correlated words derived from one seed; use `keccak256` derivation for independence.
- **Front-running VRF reveals** — the VRF callback transaction is visible in the mempool; design the contract so observing the random number before acting on it provides no advantage.
