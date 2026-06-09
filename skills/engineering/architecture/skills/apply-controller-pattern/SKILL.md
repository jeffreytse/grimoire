---
name: apply-controller-pattern
description: Use when deciding which non-UI class should handle a system operation or external event — especially when business logic leaks into UI handlers, or when the same system operation is triggered from multiple interfaces (web, CLI, API).
source: Craig Larman, "Applying UML and Patterns", 3rd ed., Prentice Hall, 2004
tags: [grasp, controller, design-principles, oop, responsibility-assignment, developer, separation-of-concerns, testability]
related: [apply-information-expert, apply-high-cohesion, apply-low-coupling, apply-solid-principles]
---

# Apply Controller Pattern

Assign the responsibility for receiving and handling system events to a non-UI class that represents either the overall system or a use-case scenario.

## Why This Is Best Practice

**Adopted by:** GRASP (Larman, 2004); the foundation of MVC's C (Controller), MVP's Presenter, and MVVM's ViewModel. Explicitly required in layered architectures (Clean Architecture — Uncle Bob; Hexagonal Architecture — Alistair Cockburn). Any application built with Spring MVC, Django views, Rails controllers, or ASP.NET is implementing GRASP Controller.
**Impact:** When business logic is placed directly in UI event handlers (button click → database write), it cannot be tested without a UI, cannot be invoked from a second interface (API, CLI), and cannot be reused. Separating system operations into Controller classes — independent of any UI framework — means the same use case can be invoked from web, mobile, API, and test without duplication. Google's Application Architecture Guide (Android) and Apple's App Architecture Guide both mandate a Controller/ViewModel separation for exactly this reason.
**Why best:** The alternative — UI event handlers that contain business logic — produces fat views, untestable code, and logic that cannot be reused across interfaces. GRASP Controller is the design decision that makes layered architecture work: it is the class that sits at the boundary between UI and domain, receives external inputs, and delegates to domain objects.

Sources: Larman, "Applying UML and Patterns" (Prentice Hall, 2004); Martin, "Clean Architecture" (Prentice Hall, 2017); Cockburn, "Hexagonal Architecture" (2005)

## Steps

### 1. Choose the Controller type

Two valid options — choose based on system complexity:

**Facade Controller:** One class represents the entire system. Appropriate for simple systems with few use cases.
```python
class SystemController:
    def place_order(self, order_data): ...
    def cancel_order(self, order_id): ...
    def register_user(self, user_data): ...
```

**Use-Case Controller:** One class per use case or user session. Appropriate when the system has many operations and a Facade Controller would grow too large.
```python
class PlaceOrderController:
    def execute(self, order_data): ...

class CancelOrderController:
    def execute(self, order_id): ...
```

Default to Use-Case Controllers when a Facade Controller exceeds ~10 operations.

---

### 2. Keep UI logic out of the Controller

The Controller handles the system operation — it does not know how the result will be displayed.

```python
# Bad — Controller renders HTML; tied to web UI
class OrderController:
    def place_order(self, request):
        order = Order(request.POST)
        order.save()
        return render_template("order_success.html", order=order)  # UI in Controller

# Good — Controller returns a result; caller (view) decides how to present it
class PlaceOrderController:
    def execute(self, order_data) -> OrderResult:
        order = self.order_factory.create(order_data)
        self.order_repo.save(order)
        return OrderResult(order_id=order.id, total=order.total())
```

---

### 3. Keep business logic out of the Controller

The Controller coordinates — it delegates to domain objects for actual business logic.

```python
# Bad — Controller contains business rules
class PlaceOrderController:
    def execute(self, order_data):
        if order_data["total"] > 10000:                    # business rule in Controller
            raise Exception("Order exceeds maximum value")
        if not self.inventory.all_in_stock(order_data):    # domain logic in Controller
            raise Exception("Items out of stock")
        # Controller should not know these rules

# Good — domain objects enforce their own rules; Controller only coordinates
class PlaceOrderController:
    def execute(self, order_data):
        order = self.order_factory.create(order_data)      # creation delegated
        order.validate()                                   # rules in domain object
        self.inventory.reserve(order)                      # domain service handles stock
        self.order_repo.save(order)
        return OrderResult.from_order(order)
```

---

### 4. Handle one use case per Controller (Use-Case Controller pattern)

A Use-Case Controller has one primary responsibility: coordinating the steps of a single use case. It should be testable by calling its `execute` method directly, with no UI setup.

```python
# Test — Controller is independent of HTTP, HTML, and web framework
def test_place_order_charges_payment():
    controller = PlaceOrderController(
        order_factory=FakeOrderFactory(),
        inventory=FakeInventory(in_stock=True),
        payment=FakePaymentGateway(),
        order_repo=FakeOrderRepository(),
    )
    result = controller.execute({"items": [...], "payment": {...}})
    assert result.payment_status == "charged"
```

---

### 5. Expose Controllers as the entry point for all interfaces

The same Controller class should handle the same use case regardless of which interface invokes it — web, CLI, message queue, or test.

```python
# Web handler — delegates to Controller
@app.post("/orders")
def web_place_order(request):
    result = PlaceOrderController(...).execute(request.json)
    return jsonify(result)

# CLI — delegates to same Controller
def cli_place_order(args):
    result = PlaceOrderController(...).execute(vars(args))
    print(result)

# Test — calls Controller directly, no HTTP
def test_place_order():
    result = PlaceOrderController(...).execute(test_data)
    assert result.order_id is not None
```

## Common Mistakes

**Fat Controller.** A Controller that contains business rules, data transformation, validation, and persistence logic is an anemic domain model in reverse — all logic has moved from domain objects to Controllers. Controllers should coordinate; domain objects should compute.

**One Controller per HTTP route.** Mapping one Controller class to one URL endpoint is an MVC framework convention, not the GRASP pattern. A GRASP Controller maps to a use case, which may be exposed via multiple routes (REST, GraphQL, RPC).

**Returning framework objects from Controllers.** A Controller that returns `HttpResponse`, `flask.Response`, or `express.Response` is coupled to the web framework — it cannot be called from CLI or tests without the framework. Return plain domain objects or result types.

## When NOT to Use

- **CRUD-only applications with no business logic** — a Controller that just calls `repo.save(data)` adds a layer with no value; in simple CRUD, the framework handler and the repository are sufficient
- **Event-driven architectures** — message consumers and event handlers serve the same role as Controllers; apply the same principles but don't add a separate Controller layer on top of a well-designed event handler
- **Microservices with a single operation** — a service that does exactly one thing is already a Controller by design; there is no UI layer to separate from
