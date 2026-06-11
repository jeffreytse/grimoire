---
name: apply-docker-security
description: Use when writing Dockerfiles, configuring container runtimes, or deploying containerized workloads — to harden containers against privilege escalation, image vulnerabilities, and container escape.
source: 'OWASP Docker Security Cheat Sheet (owasp.org/www-project-cheat-sheets); CIS Docker Benchmark; NIST SP 800-190 (Application Container Security); CWE-250'
tags: [security, owasp, docker, containers, devops, developer]
---

# Apply Docker Security

Harden Docker containers by running as non-root, using read-only filesystems, scanning images for vulnerabilities, and applying seccomp/AppArmor profiles — preventing privilege escalation and container escape.

## Why This Is Best Practice

**Adopted by:** OWASP Docker Security Cheat Sheet and CIS Docker Benchmark (Center for Internet Security) are the two authoritative references. NIST SP 800-190 (Application Container Security Guide, 2017) provides federal guidance. Google Cloud, AWS ECS, and Azure Container Apps all enforce non-root and read-only filesystem policies in their hardened runtime configurations. Kubernetes Pod Security Standards (baseline/restricted) mandate non-root and non-privileged containers.
**Impact:** Container escape vulnerabilities (CVE-2019-5736 runc, CVE-2020-15257 containerd) allow processes inside containers to gain host-level access — most exploitable when the container runs as root. The 2021 Codecov supply chain attack used compromised Docker images in CI pipelines to steal credentials from thousands of organizations. Trivy and Snyk Container scans find high/critical CVEs in 60%+ of production images according to Snyk's State of Open Source Security 2023.
**Why best:** Running containers with the `--privileged` flag or as root with volume mounts to `/` is equivalent to running on the host with sudo. Rootless containers, read-only filesystems, and capability dropping eliminate entire classes of container escape and privilege escalation attacks with minimal operational overhead.

Sources: OWASP Docker Security Cheat Sheet; CIS Docker Benchmark v1.6; NIST SP 800-190; CWE-250

## Steps

1. **Run as a non-root user — never use `USER root` in production images**:

   ```dockerfile
   FROM node:20-alpine

   # Create a non-root user
   RUN addgroup -S appgroup && adduser -S appuser -G appgroup

   WORKDIR /app
   COPY --chown=appuser:appgroup . .
   RUN npm ci --production

   # Drop to non-root before CMD
   USER appuser

   CMD ["node", "server.js"]
   ```

   For existing images that default to root: override with `--user 1000:1000` at runtime.

2. **Use minimal base images** — reduce attack surface by minimizing installed packages:

   ```dockerfile
   # Alpine: ~5MB, minimal attack surface
   FROM python:3.12-alpine

   # Distroless: no shell, no package manager, minimal CVE surface
   FROM gcr.io/distroless/python3:nonroot

   # Multi-stage: build dependencies stay in builder, only runtime artifacts in final
   FROM python:3.12-alpine AS builder
   COPY requirements.txt .
   RUN pip install --prefix=/install -r requirements.txt

   FROM gcr.io/distroless/python3:nonroot
   COPY --from=builder /install /usr/local
   COPY app.py .
   CMD ["app.py"]
   ```

3. **Make the filesystem read-only** — prevent runtime file modification:

   ```bash
   # Docker CLI
   docker run --read-only --tmpfs /tmp --tmpfs /run myapp

   # Docker Compose
   services:
     app:
       read_only: true
       tmpfs:
         - /tmp
         - /var/run
   ```

   ```yaml
   # Kubernetes pod spec
   securityContext:
     readOnlyRootFilesystem: true
   ```

4. **Drop all capabilities and add only what's needed**:

   ```bash
   docker run --cap-drop ALL --cap-add NET_BIND_SERVICE myapp
   ```

   ```yaml
   # Kubernetes
   securityContext:
     capabilities:
       drop: ["ALL"]
       add: ["NET_BIND_SERVICE"]  # only if binding port <1024
   ```

   Most apps need zero Linux capabilities — run on port ≥1024 to avoid even `NET_BIND_SERVICE`.

5. **Scan images for vulnerabilities in CI**:

   ```bash
   # Trivy (recommended — free, comprehensive)
   trivy image --exit-code 1 --severity HIGH,CRITICAL myapp:latest

   # Grype
   grype myapp:latest --fail-on high

   # In CI (GitHub Actions):
   - uses: aquasecurity/trivy-action@master
     with:
       image-ref: myapp:${{ github.sha }}
       exit-code: '1'
       severity: 'HIGH,CRITICAL'
   ```

6. **Never pass secrets as environment variables or build args**:

   ```dockerfile
   # BAD — secrets visible in image history
   ARG DATABASE_URL
   ENV DATABASE_URL=${DATABASE_URL}

   # GOOD — use runtime secret mounts
   # Docker BuildKit secret mount (available at build time, not stored in image):
   RUN --mount=type=secret,id=npmrc,target=/root/.npmrc npm install
   ```

   At runtime: use Docker secrets, Kubernetes Secrets (ideally from a Secrets Manager), or environment injection from a vault.

7. **Apply seccomp and AppArmor profiles** — restrict syscall surface:

   ```bash
   # Use Docker's default seccomp profile (denies 44 dangerous syscalls)
   docker run --security-opt seccomp=default.json myapp

   # Kubernetes: enforce restricted pod security standard
   # Add label to namespace:
   kubectl label namespace production pod-security.kubernetes.io/enforce=restricted
   ```

## Rules

- Never use `--privileged` in production — it gives the container all Linux capabilities and disables seccomp/AppArmor.
- Never mount the Docker socket (`/var/run/docker.sock`) inside a container — it gives full host access.
- Pin base image tags to digests in production: `FROM node:20-alpine@sha256:...` prevents supply chain attacks via tag mutation.
- Set `no-new-privileges: true` to prevent setuid binaries from escalating inside the container.

## Common Mistakes

- **`FROM ubuntu:latest` as base image** — includes hundreds of unnecessary packages, each a potential CVE.
- **Storing secrets in `.env` files copied into the image** — visible in `docker history` and image layers.
- **Not setting resource limits** — a single runaway container can starve all others on the host.
- **Running as root because "it's just a container"** — container boundaries are not security boundaries without hardening.
