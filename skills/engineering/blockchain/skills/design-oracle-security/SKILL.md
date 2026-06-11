---
name: design-oracle-security
description: Use when writing DeFi smart contracts that rely on price feeds, randomness, or external data — preventing oracle manipulation via TWAP, multiple oracle sources, and circuit breakers.
source: 'OWASP Smart Contract Top 10 SC02 (owasp.org/www-project-smart-contract-top-10/); SWC-116; Uniswap V3 TWAP documentation; Chainlink Data Feeds documentation'
tags: [security, owasp, solidity, smart-contracts, oracle, defi, price-manipulation, ethereum, blockchain, developer]
---

# Design Oracle Security

Use TWAP prices and multiple oracle sources with staleness checks and circuit breakers — preventing flash loan-based oracle manipulation that has drained over $1B from DeFi protocols.

## Why This Is Best Practice

**Adopted by:** OWASP Smart Contract Top 10 SC02 (Oracle Manipulation). Uniswap V3's TWAP oracle is the on-chain reference implementation used by dozens of protocols. Chainlink Data Feeds (used by Aave, Compound, Synthetix) is the leading off-chain oracle with cryptographic attestation. The DeFi Security Alliance's oracle best practices mandate multiple sources for any price used in liquidations or collateral valuation.
**Impact:** Oracle manipulation is the leading DeFi attack vector by total value lost. Mango Markets lost $117M (2022) via MANGO token price manipulation. Euler Finance lost $197M (2023) partially via oracle manipulation. Cream Finance's $130M loss (2021) and the Harvest Finance $34M loss (2020) were both price oracle attacks. CoinGecko's "DeFi Security Incident Review" (2023) attributes $1.3B+ in losses to oracle vulnerabilities. Spot price oracles (single AMM pool, single block) are manipulable via flash loans with zero capital risk.
**Why best:** Spot price oracles read from a single block — flash loans allow an attacker to borrow millions, manipulate the price in that block, exploit the oracle, and repay — all in one transaction with no capital. TWAP (Time-Weighted Average Price) averages over many blocks, making manipulation expensive and visible. Multiple oracle sources require an attacker to simultaneously manipulate independent systems.

Sources: OWASP Smart Contract Top 10 SC02; SWC-116; Uniswap V3 TWAP oracle documentation; Chainlink security overview; Euler Finance post-mortem (2023)

## Steps

1. **Use Uniswap V3 TWAP instead of spot price for on-chain oracles**:

   ```solidity
   // SPDX-License-Identifier: MIT
   pragma solidity ^0.8.19;

   import "@uniswap/v3-core/contracts/interfaces/IUniswapV3Pool.sol";
   import "@uniswap/v3-periphery/contracts/libraries/OracleLibrary.sol";

   contract TWAPOracle {
       IUniswapV3Pool public immutable pool;
       uint32 public constant TWAP_PERIOD = 1800;  // 30-minute TWAP

       constructor(address _pool) {
           pool = IUniswapV3Pool(_pool);
       }

       function getTWAPPrice() public view returns (uint256 price) {
           uint32[] memory secondsAgos = new uint32[](2);
           secondsAgos[0] = TWAP_PERIOD;
           secondsAgos[1] = 0;

           (int56[] memory tickCumulatives, ) = pool.observe(secondsAgos);

           int56 tickCumulativeDiff = tickCumulatives[1] - tickCumulatives[0];
           int24 arithmeticMeanTick = int24(tickCumulativeDiff / int56(int32(TWAP_PERIOD)));

           uint160 sqrtPriceX96 = TickMath.getSqrtRatioAtTick(arithmeticMeanTick);
           price = FullMath.mulDiv(sqrtPriceX96, sqrtPriceX96, FixedPoint96.Q96);
       }
   }
   ```

2. **Use Chainlink Data Feeds with staleness checks**:

   ```solidity
   import "@chainlink/contracts/src/v0.8/interfaces/AggregatorV3Interface.sol";

   contract ChainlinkPriceConsumer {
       AggregatorV3Interface public immutable priceFeed;
       uint256 public constant MAX_STALENESS = 3600;  // 1 hour max age

       constructor(address _priceFeed) {
           priceFeed = AggregatorV3Interface(_priceFeed);
       }

       function getPrice() public view returns (int256) {
           (
               uint80 roundId,
               int256 price,
               ,
               uint256 updatedAt,
               uint80 answeredInRound
           ) = priceFeed.latestRoundData();

           // Staleness check — stale prices are as dangerous as manipulated ones
           require(block.timestamp - updatedAt <= MAX_STALENESS, "Price feed stale");
           // Completeness check
           require(answeredInRound >= roundId, "Incomplete round");
           require(price > 0, "Invalid price");

           return price;
       }
   }
   ```

3. **Use multiple oracle sources — median of 3 prevents single-source manipulation**:

   ```solidity
   contract MultiOracleConsumer {
       AggregatorV3Interface[3] public priceFeeds;

       constructor(address[3] memory feeds) {
           for (uint i = 0; i < 3; i++) {
               priceFeeds[i] = AggregatorV3Interface(feeds[i]);
           }
       }

       function getMedianPrice() public view returns (int256) {
           int256[3] memory prices;
           for (uint i = 0; i < 3; i++) {
               (, prices[i], , uint256 updatedAt, ) = priceFeeds[i].latestRoundData();
               require(block.timestamp - updatedAt <= 3600, "Feed stale");
           }
           return _median3(prices[0], prices[1], prices[2]);
       }

       function _median3(int256 a, int256 b, int256 c) internal pure returns (int256) {
           if ((a >= b && a <= c) || (a >= c && a <= b)) return a;
           if ((b >= a && b <= c) || (b >= c && b <= a)) return b;
           return c;
       }
   }
   ```

4. **Implement circuit breakers — reject prices that deviate too far from TWAP**:

   ```solidity
   contract OracleWithCircuitBreaker {
       int256 public constant MAX_PRICE_DEVIATION_BPS = 500;  // 5% max deviation

       function validatePrice(int256 spotPrice, int256 twapPrice) internal pure {
           int256 deviation = abs(spotPrice - twapPrice) * 10000 / twapPrice;
           require(deviation <= MAX_PRICE_DEVIATION_BPS,
               "Price deviates >5% from TWAP — possible manipulation");
       }

       function abs(int256 x) internal pure returns (int256) {
           return x >= 0 ? x : -x;
       }
   }
   ```

5. **Use price averaging over multiple blocks for liquidation prices**:

   ```solidity
   contract SafeLiquidation {
       mapping(address => uint256[]) private priceHistory;
       uint256 private constant PRICE_SAMPLES = 5;

       function recordPrice(address asset, uint256 price) external {
           priceHistory[asset].push(price);
           if (priceHistory[asset].length > PRICE_SAMPLES) {
               // Remove oldest
               for (uint i = 0; i < priceHistory[asset].length - 1; i++) {
                   priceHistory[asset][i] = priceHistory[asset][i + 1];
               }
               priceHistory[asset].pop();
           }
       }

       function getLiquidationPrice(address asset) public view returns (uint256) {
           uint256[] memory history = priceHistory[asset];
           require(history.length == PRICE_SAMPLES, "Insufficient price history");
           uint256 sum;
           for (uint i = 0; i < history.length; i++) sum += history[i];
           return sum / PRICE_SAMPLES;
       }
   }
   ```

## Rules

- Never use `pool.slot0().sqrtPriceX96` (spot price from a single AMM pool) for any financial calculation — it's manipulable with a flash loan.
- TWAP period must be long enough that manipulation is economically unviable: 30 minutes minimum, 1 hour for protocols with large TVL.
- Always check Chainlink `updatedAt` timestamp — a stale feed returning old prices is as dangerous as a manipulated feed.
- Oracle failures (reverts, stale data) should pause protocol operations, not fall back to an unchecked default price.

## Common Mistakes

- **`block.timestamp` as a price oracle input** — block timestamps are manipulable by miners/validators within ~15 seconds; use cumulative price accumulators.
- **Single-source oracle for collateral pricing** — if the single source is compromised or goes down, the entire protocol's liquidation system fails.
- **Not testing oracle manipulation in fork tests** — Foundry's `vm.createFork()` allows testing against mainnet state; simulate flash loan price manipulation to verify TWAP resistance.
- **Chainlink circuit breaker not accounting for decimals** — Chainlink returns prices with 8 decimals; failing to normalize before comparison causes incorrect deviation calculations.
