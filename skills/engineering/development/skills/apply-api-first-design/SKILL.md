---
name: apply-api-first-design
description: Use when starting a new API, adding a major feature to an existing API, or aligning multiple teams on an interface — to define the contract in an OpenAPI spec before writing any implementation code.
source: 'OpenAPI Initiative (openapis.org); Stoplight "API Design Guide" (2020); Swagger "API-First Development" guide; Stripe Engineering blog — "Designing APIs for humans"; Google API Design Guide (aip.dev); Richardson & Ruby "RESTful Web Services" (O\'Reilly, 2007)'
tags: [api, openapi, contract-first, api-design, rest, developer-experience, design]
verified: true
---

# Apply API-First Design

Write the OpenAPI specification before writing implementation code — use it as the contract between teams, generate stubs and mocks from it, and keep it as the authoritative source of truth throughout the API's lifetime.

## Why This Is Best Practice

**Adopted by:** Stripe publishes its OpenAPI spec as the canonical source for all client SDKs and documentation — no client code is written without the spec existing first. Twilio, Shopify, GitHub, and AWS each use spec-first workflows. Zalando (800+ engineers) mandates API-first as a company-wide engineering standard. The OpenAPI Initiative (Linux Foundation) maintains the spec; tooling ecosystem includes Stoplight, Redocly, Speakeasy, and hundreds of code generators.

**Impact:** Teams that design APIs code-first frequently ship breaking changes because the interface was an afterthought rather than a decision. Spec-first forces interface decisions before sunk cost accumulates. It enables frontend and backend teams to work in parallel — frontend uses a mock server from the spec while backend implements against it. Stripe attributes its developer experience reputation in part to the discipline of spec-first design.

**Why best:** Code-first API design produces specs that describe what was built, not what should have been built. Specs derived from code are often incomplete, inconsistent across endpoints, and lack the examples and error schemas that make an API usable. Writing the spec first forces clarity on resource modeling, naming, error handling, and pagination before any code exists — changes are cheap at this stage and expensive after.

Sources: OpenAPI Initiative; Stoplight API Design Guide; Zalando RESTful API Guidelines; Stripe Engineering; Google AIP (aip.dev)

## Steps

### Step 1: Define the API contract in OpenAPI before writing code

Start with a minimal but complete OpenAPI 3.1 document. At minimum: `info`, `paths` with at least one operation, and `components/schemas` for request/response bodies.

```yaml
openapi: 3.1.0
info:
  title: Payments API
  version: 1.0.0
  description: Accepts payment intents and manages payment lifecycle.

paths:
  /payment-intents:
    post:
      operationId: createPaymentIntent
      summary: Create a payment intent
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/CreatePaymentIntentRequest'
      responses:
        '201':
          description: Payment intent created
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/PaymentIntent'
        '400':
          $ref: '#/components/responses/BadRequest'
        '422':
          $ref: '#/components/responses/UnprocessableEntity'

components:
  schemas:
    CreatePaymentIntentRequest:
      type: object
      required: [amount, currency]
      properties:
        amount:
          type: integer
          description: Amount in smallest currency unit (e.g., cents)
          example: 2000
        currency:
          type: string
          description: ISO 4217 currency code
          example: usd
    PaymentIntent:
      type: object
      required: [id, amount, currency, status, created_at]
      properties:
        id:
          type: string
          example: pi_3OhBxz2eZvKYlo2C
        amount:
          type: integer
        currency:
          type: string
        status:
          type: string
          enum: [requires_payment_method, requires_confirmation, succeeded, canceled]
        created_at:
          type: string
          format: date-time
  responses:
    BadRequest:
      description: Invalid request parameters
      content:
        application/json:
          schema:
            $ref: '#/components/schemas/Error'
    UnprocessableEntity:
      description: Request valid but semantically incorrect
      content:
        application/json:
          schema:
            $ref: '#/components/schemas/Error'
  schemas:
    Error:
      type: object
      required: [type, message]
      properties:
        type:
          type: string
          example: invalid_request_error
        message:
          type: string
          example: "amount must be a positive integer"
        param:
          type: string
          example: amount
```

### Step 2: Review the spec as a team before any code is written

Schedule a spec review with all stakeholders — API consumers (frontend, mobile, partner integrators), API providers (backend engineers), and product. Treat it like a design review, not a code review.

Questions to answer in the review:
- Are resource names and operation IDs consistent with existing APIs?
- Are all error cases modeled? (validation errors, auth errors, not-found, rate limits)
- Are required vs. optional fields correct?
- Does the response schema include everything consumers need, without oversharing?
- Is pagination modeled for list endpoints?
- Are examples realistic and complete?

Changes to the spec at this stage cost nothing. Changes after implementation cost the time to refactor code, regenerate clients, and communicate breaking changes.

### Step 3: Spin up a mock server from the spec

Generate a mock server so frontend/mobile teams can begin integration immediately — before the backend is implemented. The mock returns example responses defined in the spec.

```bash
# Prism (Stoplight) — most widely used
npm install -g @stoplight/prism-cli
prism mock openapi.yaml
# → mock server at http://localhost:4010

# Test it
curl -X POST http://localhost:4010/payment-intents \
  -H "Content-Type: application/json" \
  -d '{"amount": 2000, "currency": "usd"}'
# → returns example PaymentIntent from spec
```

```bash
# Alternatively: openapi-mock (Go, no Node required)
brew install muonsoft/openapi-mock/openapi-mock
openapi-mock serve --specification openapi.yaml
```

### Step 4: Generate server stubs from the spec

Use code generation to produce server-side scaffolding — route handlers, request/response types, and validation — directly from the spec. Implement business logic inside generated stubs; do not write routes or types by hand.

```bash
# OpenAPI Generator (supports 50+ languages/frameworks)
npx @openapitools/openapi-generator-cli generate \
  -i openapi.yaml \
  -g nodejs-express-server \
  -o ./server

# For Go
npx @openapitools/openapi-generator-cli generate \
  -i openapi.yaml \
  -g go-server \
  -o ./server

# For FastAPI (Python) — generates Pydantic models + route stubs
npx @openapitools/openapi-generator-cli generate \
  -i openapi.yaml \
  -g python-fastapi \
  -o ./server
```

Generated code provides:
- Route definitions matching spec paths exactly
- Request validation against schema (reject invalid requests before business logic)
- Response types that match the schema

### Step 5: Generate client SDKs from the spec

If you have consumers in multiple languages, generate client SDKs rather than writing them by hand. The spec ensures the SDK matches the server.

```bash
# TypeScript client (for frontend)
npx @openapitools/openapi-generator-cli generate \
  -i openapi.yaml \
  -g typescript-fetch \
  -o ./clients/typescript

# Python client
npx @openapitools/openapi-generator-cli generate \
  -i openapi.yaml \
  -g python \
  -o ./clients/python
```

Publish SDKs automatically: on every spec change, CI regenerates and publishes updated clients. Stripe's SDK generation pipeline works this way.

### Step 6: Keep the spec as the authoritative source of truth

Never update the implementation without updating the spec first. Treat spec-vs-implementation drift as a build failure.

```bash
# Validate spec is valid OpenAPI
npx @redocly/cli lint openapi.yaml

# In CI: check implementation matches spec (e.g., schemathesis for Python)
pip install schemathesis
schemathesis run openapi.yaml --url http://localhost:8000

# Spectral: enforce style rules on the spec itself
npx @stoplight/spectral-cli lint openapi.yaml --ruleset .spectral.yaml
```

Add spec linting to the PR check — a PR that changes an API endpoint must also update the spec, or CI fails.

## Rules

- The spec is reviewed and approved before implementation begins. Implementation follows the spec; it does not define the spec.
- Never auto-generate the spec from implementation code as the primary workflow — this produces docs-from-code, not contract-first design. Code generation flows from spec → code, not code → spec.
- Every endpoint must have at least one error response modeled. Undocumented errors are breaking changes in disguise.
- Breaking changes to the spec require a version increment or deprecation cycle — field removal, type change, required field addition, and enum value removal are all breaking.
- The spec lives in version control alongside the code. Spec changes go through the same PR review process as code changes.

## Common Mistakes

**Generating the spec from code and calling it API-first:** Tools like Springdoc or FastAPI's auto-generated OpenAPI are useful for documentation but are not API-first design. The spec reflects what was built, after the fact, with no prior design review.

**Designing only the happy path:** Specs that model only 200 responses leave consumers without guidance on how to handle errors, rate limits, and validation failures. Model all response codes your API can return.

**Skipping the team review:** Writing the spec alone and moving straight to implementation misses the primary benefit — forcing alignment before sunk cost accumulates. The review is the point.

**Treating generated code as the source of truth:** If the generator is re-run over modified generated code, manual changes are overwritten. Keep a clean separation: generated code in its own directory, business logic in a separate layer that imports from generated types.

## Examples

**Stripe's spec-first workflow:**
Stripe maintains a single `openapi.yaml` in a private repo. Every API change starts as a spec PR reviewed by the API design council. On merge, CI automatically regenerates the Ruby, Python, Node, Go, Java, and .NET client libraries and publishes them to package managers. No client code is written by hand.

**Parallel development with mock server:**
```
Week 1: Backend writes spec, team reviews, spec merged
Week 1: Frontend team runs `prism mock openapi.yaml`, begins UI integration
Week 2: Backend implements against spec
Week 3: Frontend swaps mock URL for real backend — integration works on first try
```

**CI enforcement:**
```yaml
# .github/workflows/api.yml
- name: Lint OpenAPI spec
  run: npx @redocly/cli lint openapi.yaml

- name: Check for breaking changes
  run: npx @openapitools/openapi-diff-cli compare main:openapi.yaml HEAD:openapi.yaml --fail-on-incompatible
```
