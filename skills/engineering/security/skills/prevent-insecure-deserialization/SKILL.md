---
name: prevent-insecure-deserialization
description: Use when deserializing data from untrusted sources — user-submitted cookies, API request bodies, message queue payloads, or any format that reconstructs objects (pickle, Java serialization, YAML, PHP serialize).
source: 'OWASP Deserialization Cheat Sheet (owasp.org/www-project-cheat-sheets); OWASP Top 10 2021 A08; CWE-502; NIST NVD deserialization CVE history'
tags: [security, owasp, deserialization, pickle, java-serialization, rce, developer]
---

# Prevent Insecure Deserialization

Never deserialize untrusted data using native serialization formats — use data-only formats (JSON, Protobuf) with schema validation, or enforce strict type allowlists when native serialization is unavoidable.

## Why This Is Best Practice

**Adopted by:** OWASP Top 10 2021 A08 (Software and Data Integrity Failures) covers insecure deserialization. Oracle, Apache, Red Hat, and Cisco have all patched critical RCE vulnerabilities caused by Java deserialization (CVE-2015-4852, CVE-2016-0792, etc). Node.js `node-serialize`, Python `pickle`, PHP `unserialize`, and Ruby `Marshal.load` are all documented RCE vectors when given untrusted input.
**Impact:** Apache Commons Collections deserialization RCE (2015) affected WebLogic, JBoss, Jenkins, and dozens of Java application servers — enabling remote code execution with zero authentication. Python `pickle.loads(user_input)` is equivalent to `eval(user_input)` — it executes arbitrary code. GitHub's 2017 Enterprise Server RCE was caused by Ruby Marshal deserialization. Insecure deserialization is rated CVSS 9.8 (Critical) in most CVEs.
**Why best:** Denylisting dangerous classes (Java serialization filters, etc.) is the alternative — it requires knowing every gadget chain in advance and fails when new chains are discovered. Using data-only formats (JSON) eliminates the vulnerability class: JSON can represent data structures but cannot encode executable code.

Sources: OWASP Deserialization Cheat Sheet; Apache Commons Collections CVE-2015-4852; CWE-502; Java Serialization "AppSecCali 2015 — Marshalling Pickles" (Frohoff & Lawrence)

## Steps

1. **Prefer data-only serialization formats** — JSON, XML (with XXE disabled), MessagePack, Protobuf, or Avro. These encode data, not executable objects.

   ```python
   # BAD — pickle executes arbitrary code on load
   import pickle
   obj = pickle.loads(user_input)

   # GOOD — JSON only represents data
   import json
   data = json.loads(user_input)  # safe; validate schema after
   ```

2. **If native serialization is required, validate/authenticate the data before deserializing**:

   Sign the serialized blob with an HMAC and verify the signature before deserializing:

   ```python
   import hmac, hashlib, pickle

   SECRET = b'server-secret-key-min-32-bytes-here'

   def serialize_signed(obj):
       payload = pickle.dumps(obj)
       sig = hmac.new(SECRET, payload, hashlib.sha256).hexdigest()
       return sig + ':' + payload.hex()

   def deserialize_verified(signed_data):
       sig, hex_payload = signed_data.split(':', 1)
       payload = bytes.fromhex(hex_payload)
       expected = hmac.new(SECRET, payload, hashlib.sha256).hexdigest()
       if not hmac.compare_digest(sig, expected):
           raise ValueError("Signature mismatch — data tampered")
       return pickle.loads(payload)
   ```

3. **Java — implement serialization filters (JEP 290)**:

   ```java
   // Java 9+ serialization filter — allowlist only expected classes
   ObjectInputFilter filter = ObjectInputFilter.Config.createFilter(
       "com.example.SafeClass;!*"  // allow only SafeClass, reject all others
   );
   ObjectInputStream ois = new ObjectInputStream(inputStream);
   ois.setObjectInputFilter(filter);
   Object obj = ois.readObject();
   ```

   Alternatively, use safer alternatives: Jackson (with `@JsonTypeInfo` restrictions), Kryo with registered-classes-only mode, or FST with class registration.

4. **PHP — avoid `unserialize()` with untrusted input**:

   ```php
   // BAD
   $obj = unserialize($_COOKIE['data']);

   // GOOD — use JSON
   $data = json_decode($_COOKIE['data'], true);

   // If unserialize is unavoidable, use allowed_classes:
   $obj = unserialize($data, ['allowed_classes' => ['SafeClass']]);
   ```

5. **Ruby — avoid `Marshal.load` with untrusted data**:

   ```ruby
   # BAD
   obj = Marshal.load(user_supplied_string)

   # GOOD — JSON
   require 'json'
   data = JSON.parse(user_supplied_string)
   ```

6. **Validate schema after deserialization** — even with safe formats, verify the structure matches expectations before using the data:

   ```python
   import json
   from jsonschema import validate

   schema = {"type": "object", "properties": {"user_id": {"type": "integer"}}, "required": ["user_id"]}
   data = json.loads(user_input)
   validate(instance=data, schema=schema)
   ```

7. **Log and alert on deserialization failures** — unexpected deserialization errors often indicate exploit attempts:

   ```python
   try:
       obj = deserialize(input_data)
   except (ValueError, TypeError, pickle.UnpicklingError) as e:
       logger.warning("Deserialization failure from %s: %s", request.remote_addr, e)
       raise BadRequest("Invalid data format")
   ```

## Rules

- `pickle.loads`, `yaml.load` (without `Loader=yaml.SafeLoader`), `Marshal.load`, and Java's `ObjectInputStream` are unsafe on untrusted input — no exceptions.
- YAML `yaml.load()` without SafeLoader allows arbitrary Python object construction — always use `yaml.safe_load()`.
- Signing serialized data prevents tampering but does not prevent attacks if the signing key is compromised — defense-in-depth (type allowlists) is still needed.
- Deserialization gadget chains evolve — new chains are discovered in existing libraries regularly; allowlisting is more durable than denylisting.

## Common Mistakes

- **Using `yaml.load(input)` instead of `yaml.safe_load(input)`** — YAML's default loader supports Python object tags and executes arbitrary constructors.
- **Base64-encoding the serialized blob and assuming it's safe** — encoding != encryption != authentication. Attackers can decode, modify, and re-encode.
- **Java `ObjectInputStream` without filters** — the default behavior accepts all classes, enabling gadget-chain RCE.
- **Trusting signed JWTs for complex object reconstruction** — JWT payload should contain simple claims, not serialized objects.
