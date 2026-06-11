---
name: design-grpc-service
description: Use when building a gRPC service — to define the proto contract first, apply field numbering and naming conventions, model errors with google.rpc.Status, handle streaming patterns, and maintain backward compatibility across versions.
source: 'Google API Design Guide (aip.dev); Protocol Buffers Language Guide (protobuf.dev); gRPC documentation (grpc.io); Effective Proto (buf.build/docs/best-practices); Google "API Improvement Proposals" AIP-0140, AIP-0158, AIP-0216'
tags: [grpc, protobuf, api-design, contract-first, rpc, microservices, streaming]
verified: true
---

# Design gRPC Service

Write the `.proto` file as the authoritative contract before any implementation — apply consistent naming, safe field numbering, typed error codes, and versioning rules that preserve backward compatibility.

## Why This Is Best Practice

**Adopted by:** Google uses Protocol Buffers and gRPC for virtually all internal service communication. Lyft, Netflix, Dropbox, and Cloudflare use gRPC for internal microservices. The Buf toolchain (buf.build) brings CI-enforced linting and breaking-change detection, adopted by thousands of teams. Google AIP defines naming, pagination, and error conventions for all Google APIs over gRPC.

**Impact:** Protobuf's binary encoding and HTTP/2 multiplexing produce 5–10× smaller payloads than equivalent JSON/REST. One `.proto` generates correct stubs in Go, Java, Python, C++, Rust, and 10+ languages simultaneously. A well-designed proto is reused for years; a poorly designed one accumulates workarounds in every consumer.

**Why best:** Writing the proto first forces explicit decisions about field types, enums, error codes, and streaming before any implementation begins. Protobuf's compatibility rules make safe evolution possible — but only if the initial design follows them.

Sources: Google AIP (aip.dev); Protocol Buffers Language Guide; Buf best practices; gRPC docs

## Steps

### Step 1: Establish proto file structure and package naming

```protobuf
syntax = "proto3";

package example.inventory.v1;

option go_package = "github.com/example/api/inventory/v1;inventoryv1";
option java_multiple_files = true;
option java_package = "com.example.inventory.v1";

import "google/protobuf/timestamp.proto";
import "google/protobuf/field_mask.proto";
import "google/rpc/error_details.proto";
```

### Step 2: Name services and RPCs following AIP conventions

```protobuf
service InventoryService {
  rpc GetProduct(GetProductRequest) returns (Product);
  rpc ListProducts(ListProductsRequest) returns (ListProductsResponse);
  rpc CreateProduct(CreateProductRequest) returns (Product);
  rpc UpdateProduct(UpdateProductRequest) returns (Product);
  rpc DeleteProduct(DeleteProductRequest) returns (google.protobuf.Empty);
  rpc ArchiveProduct(ArchiveProductRequest) returns (ArchiveProductResponse);
}
```

Naming rules: RPC names `PascalCase` verb-first; field names `snake_case`; enum values `UPPER_SNAKE_CASE` prefixed with enum name to avoid collision.

### Step 3: Design request and response messages

```protobuf
message Product {
  string id = 1;                             // field 1–15: 1 byte on wire — use for frequent fields
  string name = 2;
  int64 price_cents = 3;                     // money as integer — no float imprecision
  string currency_code = 4;
  ProductStatus status = 5;
  google.protobuf.Timestamp created_at = 6;  // always Timestamp, never string/int for time
  google.protobuf.Timestamp updated_at = 7;
}

enum ProductStatus {
  PRODUCT_STATUS_UNSPECIFIED = 0;  // proto3 default is 0 — never assign semantic meaning to 0
  PRODUCT_STATUS_ACTIVE = 1;
  PRODUCT_STATUS_DRAFT = 2;
  PRODUCT_STATUS_ARCHIVED = 3;
}

message UpdateProductRequest {
  Product product = 1;
  google.protobuf.FieldMask update_mask = 2;  // distinguishes absent from zero — required for PATCH
}
```

### Step 4: Model list responses with cursor pagination (AIP-0158)

```protobuf
message ListProductsRequest {
  int32 page_size = 2;
  string page_token = 3;
}

message ListProductsResponse {
  repeated Product products = 1;
  string next_page_token = 2;  // empty when no more pages
}
```

### Step 5: Use `google.rpc.Status` for rich error details

```go
// Validation error with field-level details
st, _ := status.New(codes.InvalidArgument, "invalid product").
    WithDetails(&errdetails.BadRequest{
        FieldViolations: []*errdetails.BadRequest_FieldViolation{
            {Field: "price_cents", Description: "must be a positive integer"},
        },
    })
return nil, st.Err()
```

gRPC → HTTP status mapping: `INVALID_ARGUMENT`→400, `UNAUTHENTICATED`→401, `PERMISSION_DENIED`→403, `NOT_FOUND`→404, `ALREADY_EXISTS`→409, `RESOURCE_EXHAUSTED`→429, `INTERNAL`→500, `UNAVAILABLE`→503.

### Step 6: Define streaming RPCs for appropriate patterns

```protobuf
service DataService {
  rpc WatchProducts(WatchProductsRequest) returns (stream ProductEvent);   // server streaming: live updates
  rpc ImportProducts(stream ImportProductRequest) returns (ImportProductsResponse); // client streaming: bulk import
  rpc SyncInventory(stream SyncRequest) returns (stream SyncResponse);    // bidi: real-time collaboration
}
```

Add `string resume_token` to streaming request/response so clients reconnect without missing events.

### Step 7: Reserve field numbers when removing fields

```protobuf
message Product {
  string id = 1;
  string name = 2;
  int64 price_cents = 4;
  reserved 3;              // never reuse — causes silent data corruption
  reserved "legacy_price";
}
```

### Step 8: Version the package, not the message

Increment `v1` → `v2` for breaking changes. Serve both versions from the same binary during migration; consumers migrate at their own pace.

Breaking changes: removing/renaming a field, changing field type or number, removing an enum value, singular↔repeated.

### Step 9: Enforce with Buf in CI

```bash
buf lint                                          # naming/style violations
buf breaking --against '.git#branch=main'         # breaking change detection
buf generate                                      # regenerate stubs
```

## Rules

- Field numbers 1–15 are one byte on the wire — reserve for the most frequently populated fields.
- Enum value 0 must always be `{ENUM_NAME}_UNSPECIFIED` — proto3 treats absent fields as 0.
- Never delete or reuse a field number — always `reserved` removed numbers and names.
- Use `google.protobuf.Timestamp` for all timestamps.
- Use `google.protobuf.FieldMask` for update operations to distinguish absent from zero.
- All proto changes go through `buf lint` and `buf breaking` in CI before merge.

## Common Mistakes

**Using `string` for money:** Introduces parsing ambiguity. Use `int64 price_cents` or `google.type.Money`.

**Reusing a removed field number:** Silent data corruption when old and new code coexist. Always `reserved` it.

**Semantic meaning on enum 0:** `STATUS_ACTIVE = 0` silently activates every message with an unset status. Use `UNSPECIFIED = 0`.

**Not using FieldMask for updates:** Without it, a partial update zeros all omitted fields — proto3 cannot distinguish absent from zero.

## Examples

```bash
# Generate stubs, implement, test with grpcurl
buf generate
grpcurl -plaintext localhost:9090 list
grpcurl -plaintext -d '{"id": "prod_01HXYZ"}' \
  localhost:9090 example.inventory.v1.InventoryService/GetProduct
```
