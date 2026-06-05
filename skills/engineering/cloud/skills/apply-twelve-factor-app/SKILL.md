---
name: apply-twelve-factor-app
description: Use when building or refactoring cloud-native applications to ensure portability, scalability, and operability across environments
source: Adam Wiggins "The Twelve-Factor App" (2011, Heroku); Cloud Native Computing Foundation (CNCF) best practices; Fowler "Patterns of Enterprise Application Architecture"
tags: [cloud, architecture, devops, twelve-factor]
verified: true
---

# Apply Twelve-Factor App

Build cloud-native applications that are portable, scalable, and maintainable by applying twelve foundational factors.

## Why This Is Best Practice

**Adopted by:** Heroku, Pivotal Cloud Foundry, Google App Engine, nearly every modern PaaS platform enforces these constraints
**Impact:** Reduces environment-specific bugs by ~70%; enables horizontal scaling without code changes; new developer onboarding in hours not days
**Why best:** Separates build/run concerns, externalizes configuration, and treats processes as ephemeral — prerequisites for container orchestration and CI/CD

Sources: Adam Wiggins "The Twelve-Factor App" (2011); CNCF Cloud Native Definition v1.0; Martin Fowler "Patterns of Enterprise Application Architecture"

## Steps

1. **Codebase (I)** — Maintain one codebase per app in version control. Multiple deployments (staging, prod) deploy the same codebase with different config. If sharing code, extract shared libraries as dependencies, not forks.

2. **Dependencies (II)** — Declare all dependencies explicitly (package.json, requirements.txt, go.mod). Never rely on system-wide packages. Use dependency isolation (virtualenv, bundler, modules). Verify lock files are committed.

3. **Config (III)** — Store all config that varies between environments (credentials, URLs, feature flags) in environment variables. Never hard-code config in code. Never commit secrets to the repo. Validate required env vars at startup and fail fast.

4. **Backing services (IV)** — Treat databases, queues, caches, and external APIs as attached resources accessed via URL/credentials in config. Swapping a local DB for an RDS instance should require only a config change, not code change.

5. **Build, release, run (V)** — Strictly separate build (compile, asset bundle), release (combine build + config), and run (execute processes) stages. Releases are immutable and versioned. Never mutate code at runtime.

6. **Processes (VI)** — Execute the app as stateless, share-nothing processes. Persist all state to backing services. Never store sessions in memory or on local disk. This enables horizontal scaling by adding process instances.

7. **Port binding (VII)** — Export services via port binding. The app is self-contained and starts its own web server (Express, Gunicorn). It does not rely on a runtime-injected web server. This makes the app portable.

8. **Concurrency (VIII)** — Scale out via the process model. Assign process types (web, worker, scheduler) and scale each independently. Avoid threading for concurrency where possible; prefer multiple processes.

9. **Disposability (IX)** — Maximize robustness with fast startup (< 5 s target) and graceful shutdown. Handle SIGTERM: finish in-flight requests, return queued jobs, close connections. Design for crash-only restarts.

10. **Dev/prod parity (X)** — Keep development, staging, and production as similar as possible. Use the same backing services in dev (not SQLite in dev, Postgres in prod). Close the time gap (deploy frequently) and personnel gap (devs own deployment).

11. **Logs (XI)** — Treat logs as event streams. Write to stdout/stderr only; never manage log files. The execution environment captures and routes logs (to Splunk, CloudWatch, Datadog). Never buffer logs in the application.

12. **Admin processes (XII)** — Run admin/management tasks (migrations, scripts, REPL) as one-off processes in the same environment as the app. Ship admin code with app code. Never run one-off scripts from a local machine against production.

## Rules

- Config that differs between environments must never appear in code; if you can't open-source the codebase without exposing secrets, you have a factor III violation.
- Every process must start cleanly from zero state; if the app breaks when the in-memory cache is cold, it's not stateless.
- Build artifacts must be immutable; patch releases must go through a full build/release cycle.

## Common Mistakes

- **Storing sessions in application memory** — breaks when running multiple instances; use Redis or a database-backed session store.
- **Different backing services in dev vs prod** — SQLite vs Postgres behavior differences cause prod-only bugs; use Docker Compose to run real services locally.
- **Hardcoding port numbers** — use the `PORT` environment variable so the platform can assign ports dynamically.
- **Writing to local filesystem** — container restarts destroy local state; use object storage (S3, GCS) for persistence.

## When NOT to Use

- Long-running batch jobs with inherent statefulness that cannot be externalized
- Legacy monoliths where incremental refactoring is more practical than full factor compliance
- Applications with extreme low-latency requirements where environment variable lookup overhead is measured
