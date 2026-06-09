---
name: apply-facade-pattern
description: Use when you want to provide a simple, unified interface to a complex subsystem — reducing the number of objects clients must interact with and lowering the coupling between client code and subsystem internals.
source: Gamma, Helm, Johnson, Vlissides, "Design Patterns: Elements of Reusable Object-Oriented Software" (1994) pp. 185–193; AWS SDK high-level clients (S3 Transfer Manager); Spring Boot auto-configuration; Python's `smtplib` wrapper in `email.mime`
tags: [design-patterns, structural, facade, oop, developer, simplification, subsystem-decoupling]
related: [apply-adapter-pattern, apply-law-of-demeter, apply-solid-principles]
---

# Apply Facade Pattern

Provide a unified, simplified interface to a complex subsystem, hiding its internal complexity from clients.

## Why This Is Best Practice

**Adopted by:** AWS SDK's high-level clients (S3 Transfer Manager wraps 15+ low-level
S3 API calls behind `upload_file()` — used by millions of applications), Spring Boot
auto-configuration (a facade over Spring's XML-heavy configuration subsystem), Django's
`send_mail()` (a facade over `smtplib`, `MIMEText`, `SMTP`, and `EmailMessage`),
and every OS API (POSIX `open()` is a facade over filesystem drivers, VFS, inode
resolution, and buffer caches).
**Impact:** GoF documents that facades reduce the number of objects clients interact with
and simplify the learning curve for subsystems. AWS Transfer Manager's single
`upload_file()` call replaced a 5-step multipart upload sequence — adoption of S3 for
large files increased dramatically after its introduction.
**Why best:** Direct use of a complex subsystem forces every client to understand all
subsystem interactions, dependencies, and initialization order. A facade centralizes
that knowledge, reducing client code to the essential operations and isolating subsystem
changes to the facade.

Sources: Gamma et al. (1994) pp. 185–193; AWS S3 Transfer Manager documentation;
Django `django.core.mail` source

## Steps

### Step 1: Identify the subsystem and the operations clients actually need

```
Subsystem: video encoding pipeline
Classes: VideoDecoder, FrameExtractor, AudioExtractor, Encoder, Muxer, FileWriter

Client needs: transcode(input_path, output_path, format)
```

Clients rarely need all subsystem operations — only the orchestrated result.

### Step 2: Write the facade — orchestrate subsystem calls, expose one method

```python
class VideoTranscoder:
    def __init__(self):
        self._decoder = VideoDecoder()
        self._audio = AudioExtractor()
        self._encoder = Encoder()
        self._muxer = Muxer()
        self._writer = FileWriter()

    def transcode(self, input_path: str, output_path: str, fmt: str) -> None:
        video_frames = self._decoder.decode(input_path)
        audio_track = self._audio.extract(input_path)
        encoded_video = self._encoder.encode(video_frames, fmt)
        muxed = self._muxer.mux(encoded_video, audio_track)
        self._writer.write(muxed, output_path)
```

### Step 3: Keep subsystem classes accessible for advanced use

Facade hides complexity but does not seal the subsystem. Clients that need fine-grained
control can still use `VideoDecoder`, `Encoder`, etc. directly. The facade is a
convenience, not a lock.

### Step 4: Don't add business logic to the facade

```python
# Wrong — business logic in the facade
def transcode(self, input_path, output_path, fmt):
    if not self._billing.is_paid(user_id):
        raise BillingError("subscription required")
    ...

# Right — facade only orchestrates subsystem calls
# Billing belongs in a service layer above the facade
```

### Step 5: Use one facade per coherent subsystem area

```python
# Good — one facade per subsystem
class VideoFacade: ...
class AudioFacade: ...
class StorageFacade: ...

# Wrong — one mega-facade for everything
class MediaSystemFacade:  # does encoding, billing, storage, notifications
    ...
```

## When NOT to Use

- **When clients legitimately need fine-grained subsystem control** — a facade that exposes too many methods is not a facade, it's just a class. Expose only the operations that cover 90% of use cases.
- **When the subsystem is simple** — a 2-class subsystem doesn't need a facade. Introduce it when clients are interacting with 5+ classes to accomplish routine tasks.

## Common Mistakes

**Facade that delegates to itself recursively.** A facade whose methods call other facade methods creates a tangle. Each facade method should call subsystem classes, not other facade methods.

**Hiding all subsystem access.** If the facade is the only way in, advanced use cases are impossible. Keep subsystem classes public; the facade is a convenience layer, not an access control mechanism.

**Treating the facade as a God Object.** A facade with 40 methods is just a big class. If the facade grows to cover unrelated subsystems, split it.
