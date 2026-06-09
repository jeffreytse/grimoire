---
name: apply-protected-variations
description: Use when a point of variation or instability is identified in a design — especially when a third-party dependency, algorithm, or external interface is likely to change and other components should not be affected.
source: Craig Larman, "Applying UML and Patterns", 3rd ed., Prentice Hall, 2004
tags: [grasp, protected-variations, design-principles, oop, stability, developer, encapsulation, open-closed]
related: [apply-polymorphism, apply-indirection, apply-low-coupling, apply-solid-principles, apply-strangler-fig-pattern]
---

# Apply Protected Variations

Identify the points of instability in a design; wrap them in a stable interface so that variations at those points do not ripple through the system.

## Why This Is Best Practice

**Adopted by:** GRASP (Larman, 2004); the most general expression of a principle behind SOLID's Open/Closed Principle (closed to modification via stable interface), Dependency Inversion (depend on abstractions), and encapsulation (hide implementation behind a stable API). Protected Variations is the *why* behind these patterns — they exist because instability at one point must not force change everywhere.
**Impact:** Alistair Cockburn's Hexagonal Architecture (2005) and Robert Martin's Clean Architecture (2017) are architectural applications of Protected Variations at scale — the stable domain is wrapped in adapter interfaces that protect it from unstable external systems (UI, database, third-party APIs). The Strangler Fig pattern (for legacy system replacement) is Protected Variations in time: wrap the unstable legacy code behind a stable interface before replacing it.
**Why best:** Every system has two kinds of code: stable (domain rules, core logic) and volatile (third-party APIs, external formats, infrastructure providers). Without Protected Variations, volatile code is directly imported into stable code — every change to the volatile part forces changes in the stable part. Protected Variations draws a line and builds a stable interface at that line. This is the principle that makes software tolerant of change.

Sources: Larman, "Applying UML and Patterns" (Prentice Hall, 2004); Cockburn, "Hexagonal Architecture" (2005); Martin, "Clean Architecture" (Prentice Hall, 2017)

## Steps

### 1. Identify points of instability

Before designing, explicitly enumerate what is likely to change. Sources of instability:

- **Third-party dependencies** — their APIs change with version upgrades
- **External systems** — payment providers, shipping APIs, auth providers can be swapped
- **Storage technology** — relational vs. document vs. event-sourced
- **Business rules** — discount logic, pricing, eligibility rules change with the business
- **Output formats** — JSON today, Protobuf next quarter, CSV for a batch job

Mark these in design documents. Each is a candidate for a protective interface.

---

### 2. Define a stable interface at the boundary

The stable interface expresses what the rest of the system needs, not what the volatile implementation provides. Design the interface from the consumer's perspective.

```python
# Volatile: payment providers change (Stripe → Adyen, or adding PayPal)
# Stable interface: what our system needs from any payment provider

class PaymentGateway(ABC):          # stable — won't change when provider changes
    @abstractmethod
    def charge(self, amount: Money, source: PaymentSource) -> ChargeResult: ...

    @abstractmethod
    def refund(self, charge_id: str, amount: Money) -> RefundResult: ...

# Volatile implementations hide behind the stable interface
class StripeGateway(PaymentGateway): ...
class AdyenGateway(PaymentGateway): ...
class PayPalGateway(PaymentGateway): ...
```

---

### 3. Use polymorphism for behavioral variation

When the point of instability is algorithm or behavior selection (which algorithm to use, which provider to call), use Polymorphism to protect against variation.

```python
# Volatile: tax calculation rules vary by country, change with tax laws
class TaxCalculator(ABC):
    @abstractmethod
    def calculate(self, order: Order) -> Money: ...

class USTaxCalculator(TaxCalculator):   # US rules
    def calculate(self, order): ...

class EUVATCalculator(TaxCalculator):   # EU VAT rules
    def calculate(self, order): ...

class OrderService:
    def __init__(self, tax_calculator: TaxCalculator):   # stable — unaffected by tax changes
        self.tax_calculator = tax_calculator
```

---

### 4. Use the Adapter pattern for external interface instability

When the volatile point is an external API whose format or protocol differs from what the system needs, use an Adapter to translate.

```python
# External shipping API has a messy interface that changes with API versions
class FedExShippingAPI:
    def create_shipment_v3(self, shipper_acct, weight_oz, dest_zip, service_type): ...

# Stable interface our system uses
class ShippingProvider(ABC):
    @abstractmethod
    def ship(self, order: Order, destination: Address) -> TrackingNumber: ...

# Adapter translates — only this class changes when FedEx API changes
class FedExAdapter(ShippingProvider):
    def ship(self, order, destination):
        return self.fedex_api.create_shipment_v3(
            shipper_acct=FEDEX_ACCOUNT,
            weight_oz=order.weight_oz,
            dest_zip=destination.zip,
            service_type="GROUND"
        )
```

---

### 5. Apply at architectural level: Hexagonal / Ports-and-Adapters

At the application architecture level, Protected Variations means the domain (core) is surrounded by ports (stable interfaces) and adapters (volatile implementations). The domain never imports infrastructure.

```
[Web Handler] → [Port: OrderUseCase] → [Domain: OrderService, Order]
[REST API]    ↗                        ↘ [Port: OrderRepository]
                                           ↘ [MySQL Adapter | InMemory Adapter]
```

The domain is stable. All variation (web vs. API, MySQL vs. MongoDB, Stripe vs. PayPal) lives outside it, behind ports.

## Common Mistakes

**Protecting against everything.** Not every class needs an interface. Wrapping stable internal utility classes in interfaces adds indirection without protection. Only protect against genuine variation points.

**Stable interface that leaks volatile details.** An interface method `def charge_stripe_v3(self, card_token: str)` is not stable — it leaks Stripe and version details. Stable interfaces use domain vocabulary, not provider vocabulary.

**Premature protection.** Designing elaborate interfaces for variation that never happens is YAGNI. Apply Protected Variations when there is a known variation point (current or highly probable), not as a general precaution.

## When NOT to Use

- **Internal implementation details** — algorithm internals within a single, focused class do not need a protective interface; the class boundary is sufficient
- **Truly stable third-party APIs** — the Go standard library, Python `datetime`, Java `java.lang.*` — wrapping these in interfaces adds no protection because they don't change in breaking ways
- **Simple scripts and one-off tools** — variation points are not relevant when the code runs once and is discarded
- **When variation is known to be final** — if a third-party contract is legally binding and the interface is versioned/frozen, the variation point no longer exists
