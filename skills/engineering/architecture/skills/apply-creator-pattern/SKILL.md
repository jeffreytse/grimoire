---
name: apply-creator-pattern
description: Use when deciding which class should be responsible for creating instances of another class — especially when object creation is scattered across the codebase, or when constructors are called in classes that have no other relationship to the created object.
source: Craig Larman, "Applying UML and Patterns", 3rd ed., Prentice Hall, 2004
tags: [grasp, creator, design-principles, oop, responsibility-assignment, developer, object-creation, low-coupling]
related: [apply-information-expert, apply-low-coupling, apply-solid-principles, apply-pure-fabrication]
---

# Apply Creator Pattern

Assign the responsibility for creating an instance of class A to the class B that aggregates, contains, closely uses, or has the initialization data for A.

## Why This Is Best Practice

**Adopted by:** GRASP (Larman, 2004); the pattern is the principled basis behind factory methods in the Gang of Four "Design Patterns" and behind dependency injection container wiring in Spring, .NET DI, and Guice. Any enterprise OOP framework with a dependency injection model is implementing Creator systematically.
**Impact:** Arbitrary object creation — constructors called wherever they are needed — scatters instantiation logic across the codebase. When the constructor signature changes (new parameter, different type), every call site must be updated. Creator localizes instantiation to the class with the strongest relationship to the created object, reducing the number of classes that need updating when construction changes.
**Why best:** Creation is a form of coupling: every class that calls `new X()` is coupled to X's constructor. Creator minimizes this coupling by assigning creation to the class that is already coupled to X for other reasons. If `Order` already aggregates `LineItem` objects, it is already coupled to `LineItem` — `Order` creating `LineItem` adds no new coupling. A `ReportGenerator` creating `LineItem` would add coupling with no justification.

Sources: Larman, "Applying UML and Patterns" (Prentice Hall, 2004); Gamma et al., "Design Patterns" (Addison-Wesley, 1994); Seemann, "Dependency Injection Principles, Practices, and Patterns" (Manning, 2019)

## Steps

### 1. Apply the Creator criteria

Assign B the responsibility to create A if one or more of the following is true (more = stronger candidate):

| Criterion | Example |
|---|---|
| B **aggregates** A | `Order` aggregates `LineItem` → `Order` creates `LineItem` |
| B **contains** A (composition) | `Document` contains `Paragraph` → `Document` creates `Paragraph` |
| B **closely uses** A | `OrderProcessor` uses `Receipt` only → `OrderProcessor` creates `Receipt` |
| B **has initialization data** for A | `OrderController` has all data to create `Order` → `OrderController` creates `Order` |

If multiple classes meet criteria, choose the one with the most criteria.

```python
# Bad — ReportGenerator creates LineItem with no relationship to it
class ReportGenerator:
    def summarize(self, order_data):
        items = [LineItem(d["name"], d["price"], d["qty"]) for d in order_data]
        # ReportGenerator has no other reason to know about LineItem

# Good — Order aggregates LineItems; Order creates them
class Order:
    def add_item(self, name, price, qty):
        self.line_items.append(LineItem(name, price, qty))
        # Order already owns LineItems; creation adds no new coupling
```

---

### 2. Use factory methods when creation logic is complex

When creating an object requires non-trivial logic (validation, lookups, conditional construction), encapsulate creation in a factory method on the Creator class.

```python
class Order:
    @classmethod
    def from_cart(cls, cart, customer):
        order = cls(customer_id=customer.id)
        for cart_item in cart.items:
            if cart_item.product.in_stock:            # non-trivial: stock check
                order.add_item(cart_item.product, cart_item.quantity)
        order.apply_loyalty_discount(customer)        # non-trivial: discount
        return order
```

---

### 3. Centralize repeated construction in one place

When the same object is constructed in multiple places with the same or similar parameters, it belongs in a single factory method or factory class owned by the Creator.

```python
# Bad — Address constructed in 3 places
class OrderController:
    def create(self, data):
        addr = Address(data["street"], data["city"], data["zip"])

class CustomerController:
    def update(self, data):
        addr = Address(data["street"], data["city"], data["zip"])

class ShippingService:
    def ship(self, data):
        addr = Address(data["street"], data["city"], data["zip"])

# Good — centralized factory; only one place to change when Address evolves
class Address:
    @classmethod
    def from_form_data(cls, data):
        return cls(data["street"], data["city"], data["zip"])
```

---

### 4. Separate creation from use (SRP on construction)

A class that creates an object and also uses it heavily is both a Creator and a User — two responsibilities. When creation logic becomes complex, extract a Factory or Builder.

```python
# Bad — OrderService both creates and processes orders (creation logic is complex)
class OrderService:
    def process(self, request):
        order = Order(
            user=self.user_repo.find(request.user_id),
            items=[self.product_repo.find(pid) for pid in request.product_ids],
            coupon=self.coupon_repo.find(request.coupon_code) if request.coupon_code else None,
        )
        return self._process_order(order)

# Good — factory handles complex creation; service handles processing
class OrderFactory:
    def from_request(self, request) -> Order: ...

class OrderService:
    def __init__(self, factory: OrderFactory): ...
    def process(self, request):
        order = self.factory.from_request(request)
        return self._process_order(order)
```

---

### 5. Override Creator with Dependency Injection for testability

When the created class has external dependencies (database, HTTP), creating it directly inside the Creator introduces coupling to infrastructure. Use dependency injection: the Creator receives the factory or the created object instead of constructing it.

```python
# Bad — OrderProcessor constructs PaymentGateway; test requires real gateway
class OrderProcessor:
    def process(self, order):
        gateway = StripeGateway(api_key=os.environ["STRIPE_KEY"])
        return gateway.charge(order.total())

# Good — gateway injected; test injects a fake
class OrderProcessor:
    def __init__(self, gateway: PaymentGateway):
        self.gateway = gateway

    def process(self, order):
        return self.gateway.charge(order.total())
```

## Common Mistakes

**Creating objects in classes with no relationship to them.** If class B has no aggregation, containment, usage, or data relationship with A, B should not create A. Find or create the class that does have this relationship.

**Over-centralizing in a single factory.** A `GlobalFactory` that creates all objects in the system is not applying Creator — it is creating a coupling hub. Each type should be created by its most closely related class.

**Confusing Creator with Dependency Injection.** Creator determines *which class* is responsible for creation. DI determines *how* dependencies are supplied. Both can apply to the same object: Creator says "Order creates LineItem"; DI says "OrderService receives OrderFactory via injection."

## When NOT to Use

- **When creation requires capabilities the Creator class should not have** — if creating X requires database access and the Creator is a domain object, use a repository or factory service instead (Pure Fabrication)
- **When the created object is a value object** — value objects (`Money`, `Address`, `DateRange`) are typically created at the call site with their initialization data; applying Creator here adds indirection without benefit
- **When a DI container manages creation** — in applications wired by a DI container (Spring, .NET DI, Guice), the container is the Creator for application-layer objects; don't fight it
