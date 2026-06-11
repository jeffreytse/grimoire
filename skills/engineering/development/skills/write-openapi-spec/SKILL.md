---
name: write-openapi-spec
description: Use when writing or reviewing an OpenAPI 3.x specification — to apply naming conventions, schema reuse, error modeling, pagination, security schemes, and examples that make the spec complete, consistent, and usable as a contract.
source: 'OpenAPI Specification 3.1.0 (spec.openapis.org); Zalando RESTful API Guidelines (opensource.zalando.com/restful-api-guidelines); Google API Design Guide (aip.dev); RFC 7807 "Problem Details for HTTP APIs"'
tags: [api, openapi, swagger, rest, api-design, documentation, contract-first, developer-experience]
verified: true
---

# Write OpenAPI Spec

Structure an OpenAPI 3.1 specification with consistent naming, reusable components, complete error schemas, realistic examples, and correct security declarations — so it serves as a usable contract rather than a documentation afterthought.

## Why This Is Best Practice

**Adopted by:** Zalando's RESTful API Guidelines are an industry-standard reference adopted by hundreds of teams. GitHub, Stripe, and Twilio publish OpenAPI specs as the canonical source for SDK generation and developer documentation. RFC 7807 (Problem Details for HTTP APIs) is implemented natively by Spring, FastAPI, and most modern frameworks.

**Impact:** An incomplete or inconsistent OpenAPI spec produces unusable generated clients, incorrect documentation, and consumer confusion. A high-quality spec generates correct SDKs, mocks, validation, and docs simultaneously — a poor spec generates noise overridden manually in every consumer.

**Why best:** OpenAPI's flexibility allows many ways to describe the same API. Without conventions, specs across large teams become inconsistent, hard to lint, and difficult to generate clients from. The conventions below represent consensus across Zalando, Google, and the OpenAPI tooling ecosystem.

Sources: OpenAPI 3.1.0 spec; Zalando RESTful API Guidelines; Google AIP; RFC 7807

## Steps

### Step 1: Set `info` and `servers`

```yaml
openapi: 3.1.0
info:
  title: Inventory API
  version: 2.1.0
  description: Manages product inventory across warehouses.
  contact:
    name: API Support
    email: api-support@example.com

servers:
  - url: https://api.example.com/v2
    description: Production
  - url: https://sandbox.api.example.com/v2
    description: Sandbox
```

### Step 2: Name paths and operationIds consistently

`operationId` is the most important field — it becomes function names in generated SDKs. Convention: `{verb}{Resource}` camelCase.

```yaml
paths:
  /products:
    get:
      operationId: listProducts     # collection → listX
    post:
      operationId: createProduct    # create → createX
  /products/{productId}:
    get:
      operationId: getProduct       # single resource → getX
    patch:
      operationId: updateProduct    # partial update → updateX
    delete:
      operationId: deleteProduct
  /products/{productId}/variants:
    get:
      operationId: listProductVariants  # nested → listParentChild
```

Path rules: plural nouns (`/products`), camelCase params (`{productId}`), hyphens for multi-word (`/payment-intents`).

### Step 3: Define all schemas as reusable `$ref` components

Never inline schemas used more than once. Define in `components/schemas`, reference everywhere.

```yaml
components:
  schemas:
    Product:
      type: object
      required: [id, name, price, currency, status, created_at]
      properties:
        id:
          type: string
          example: prod_01HXYZ
          readOnly: true
        name:
          type: string
          minLength: 1
          maxLength: 255
        price:
          type: integer
          description: Price in smallest currency unit (cents)
          example: 2999
        currency:
          type: string
          pattern: '^[A-Z]{3}$'
          example: USD
        status:
          type: string
          enum: [active, draft, archived]
        created_at:
          type: string
          format: date-time
          readOnly: true

    CreateProductRequest:
      type: object
      required: [name, price, currency]
      properties:
        name:
          type: string
        price:
          type: integer
          minimum: 1
        currency:
          type: string
          pattern: '^[A-Z]{3}$'
```

Separate `CreateXRequest` / `UpdateXRequest` from the response `X` schema — they're almost never identical.

### Step 4: Model errors using RFC 7807 Problem Details

```yaml
components:
  schemas:
    Error:
      type: object
      required: [type, title, status]
      properties:
        type:
          type: string
          format: uri
          example: https://api.example.com/errors/validation-error
        title:
          type: string
          example: Validation Error
        status:
          type: integer
          example: 422
        detail:
          type: string
          example: "price must be a positive integer"
        errors:
          type: array
          items:
            type: object
            required: [field, message]
            properties:
              field:
                type: string
              message:
                type: string

  responses:
    BadRequest:
      description: Invalid request parameters
      content:
        application/problem+json:
          schema:
            $ref: '#/components/schemas/Error'
    NotFound:
      description: Resource not found
      content:
        application/problem+json:
          schema:
            $ref: '#/components/schemas/Error'
    TooManyRequests:
      description: Rate limit exceeded
      headers:
        Retry-After:
          schema:
            type: integer
      content:
        application/problem+json:
          schema:
            $ref: '#/components/schemas/Error'
```

Define shared responses for 400, 401, 403, 404, 429, 422. Reference with `$ref` on every operation — never inline per-endpoint.

### Step 5: Model pagination for list endpoints

```yaml
components:
  schemas:
    ProductList:
      type: object
      required: [data, meta]
      properties:
        data:
          type: array
          items:
            $ref: '#/components/schemas/Product'
        meta:
          type: object
          required: [has_more]
          properties:
            has_more:
              type: boolean
            next_cursor:
              type: string
              example: cursor_eyJpZCI6InByb2RfMDFIWFlaIn0
```

Query params: `limit` (integer, default 20, max 100), `starting_after` (cursor string). Cursor-based is preferred; offset for small stable collections.

### Step 6: Declare security schemes

```yaml
components:
  securitySchemes:
    BearerAuth:
      type: http
      scheme: bearer
      bearerFormat: JWT

security:
  - BearerAuth: []  # global default

paths:
  /products:
    get:
      security: []  # override: public endpoint
```

### Step 7: Write realistic examples

```yaml
paths:
  /products:
    post:
      requestBody:
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/CreateProductRequest'
            examples:
              minimal:
                value: {name: Basic Widget, price: 999, currency: USD}
              full:
                value: {name: Premium Widget XL, price: 4999, currency: EUR}
      responses:
        '201':
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Product'
              example:
                id: prod_01HXYZ
                name: Basic Widget
                price: 999
                currency: USD
                status: active
                created_at: "2024-01-15T10:30:00Z"
```

### Step 8: Lint before committing

```bash
npx @redocly/cli lint openapi.yaml           # structural errors, broken $refs
npx @stoplight/spectral-cli lint openapi.yaml # style rules
npx @openapitools/openapi-diff openapi.yaml openapi.prev.yaml  # breaking changes
```

## Rules

- All schemas used more than once → `$ref` to `components/schemas`, never inlined twice.
- Every operation declares all response codes it can return; use shared `$ref` responses.
- `operationId` unique, `{verb}{Resource}` camelCase convention throughout.
- `readOnly: true` on server-generated fields (`id`, `created_at`) — generators exclude them from request schemas.
- Never `type: object` without `properties` — generates `any`/`dict` in clients.

## Common Mistakes

**`200` for everything:** Use `201` for POST, `204` for DELETE. `200` everywhere signals afterthought design.

**Missing `required` array:** Omitting it makes all fields optional in generated clients — forces null checks on data that's always present.

**Inlining schemas per-endpoint:** Inline schemas can't be reused and create silent inconsistencies when one copy is updated and others aren't.

**No examples:** Empty mock server responses and unusable docs. Every schema property should have `example`; every operation at least one request/response example.

**`additionalProperties: true` on responses:** Tells generators the response type is open-ended. Use `false` to signal the schema is complete.
