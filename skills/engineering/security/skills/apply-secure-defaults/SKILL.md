---
name: apply-secure-defaults
description: Use when configuring a new application, framework, or service — ensuring security is enabled out of the box by setting restrictive defaults for cookies, headers, database connections, and framework security features.
source: 'OWASP Top 10 Proactive Controls C05 C09 (owasp.org/www-project-top-ten/); OWASP Application Security Verification Standard 14.1; CIS Software Supply Chain Security Guide; Django/Rails security configuration documentation'
tags: [security, owasp, secure-defaults, configuration, hardening, framework-security, developer]
---

# Apply Secure Defaults

Configure frameworks and infrastructure with security enabled by default — HttpOnly/Secure/SameSite cookies, HTTPS enforcement, debug disabled, strong cipher suites, and framework security middleware — preventing entire vulnerability classes through configuration rather than per-feature implementation.

## Why This Is Best Practice

**Adopted by:** OWASP Top 10 Proactive Controls C05 (Secure By Default Configurations) and C09 (Leverage Security Frameworks and Libraries). Django, Rails, Spring Security, and ASP.NET Core all ship with security defaults that must be explicitly disabled — enabling them is zero-effort compared to implementing the equivalent from scratch. OWASP ASVS 14.1 (Build and Deploy) defines the configuration security requirements. Google's BeyondProd and Netflix's "Paved Road" security model both use security-by-default infrastructure as the primary mechanism for consistent security posture across thousands of services.
**Impact:** The 2021 Log4Shell vulnerability (CVSS 10.0, estimated $6.2B remediation cost) was exploited partly because Java application servers had JNDI lookups enabled by default — disabled by default would have prevented exploitation on most targets without a single line of code change. Django's default `DEBUG = False` in production prevents stack traces from being served to users; teams that deploy with `DEBUG = True` expose full database queries, environment variables, and code to anyone who triggers an error. A CIS Benchmark analysis (2022) found that 65% of remediated cloud misconfigurations were default settings that were never reviewed.
**Why best:** Per-feature security implementation (add CSRF protection here, add XSS escaping there) creates gaps — developers forget or don't know every security requirement. Secure defaults mean a developer who knows nothing about security still gets CSRF protection, secure cookies, and HTTPS enforcement automatically. The alternative (security off by default, enabled per-feature) results in security that's only as good as the developer's security knowledge.

Sources: OWASP Proactive Controls C05, C09; OWASP ASVS 14.1; Django security documentation; Rails security guide; Netflix Paved Road security model

## Steps

1. **Django: verify all security middleware and settings are enabled**:

   ```python
   # settings.py — production security checklist

   # Run Django's built-in check: python manage.py check --deploy
   # All items below must be verified

   DEBUG = False  # CRITICAL — never True in production

   ALLOWED_HOSTS = ["app.company.com"]  # explicit allowlist, not ["*"]

   # HTTPS enforcement
   SECURE_SSL_REDIRECT = True
   SECURE_HSTS_SECONDS = 31536000  # 1 year
   SECURE_HSTS_INCLUDE_SUBDOMAINS = True
   SECURE_HSTS_PRELOAD = True

   # Cookie security
   SESSION_COOKIE_SECURE = True      # only over HTTPS
   SESSION_COOKIE_HTTPONLY = True     # no JS access
   SESSION_COOKIE_SAMESITE = "Lax"   # CSRF protection
   SESSION_COOKIE_AGE = 1800         # 30 min idle timeout
   CSRF_COOKIE_SECURE = True
   CSRF_COOKIE_HTTPONLY = True
   CSRF_COOKIE_SAMESITE = "Lax"

   # Security middleware (order matters)
   MIDDLEWARE = [
       "django.middleware.security.SecurityMiddleware",  # must be first
       "django.middleware.clickjacking.XFrameOptionsMiddleware",
       "django.middleware.csrf.CsrfViewMiddleware",
       # ... other middleware
   ]

   SECURE_CONTENT_TYPE_NOSNIFF = True
   X_FRAME_OPTIONS = "DENY"
   SECURE_REFERRER_POLICY = "strict-origin-when-cross-origin"
   ```

2. **Express.js/Node.js: apply Helmet and secure defaults**:

   ```javascript
   import express from 'express';
   import helmet from 'helmet';
   import session from 'express-session';

   const app = express();

   // Helmet sets 11 security headers in one call
   app.use(helmet({
     contentSecurityPolicy: {
       directives: {
         defaultSrc: ["'self'"],
         scriptSrc: ["'self'"],  // no unsafe-inline
         objectSrc: ["'none'"],
         upgradeInsecureRequests: [],
       },
     },
     hsts: { maxAge: 31536000, includeSubDomains: true, preload: true },
   }));

   // Session security
   app.use(session({
     secret: process.env.SESSION_SECRET,  // from secrets manager
     name: '__Host-session',  // __Host- prefix requires Secure flag
     resave: false,
     saveUninitialized: false,
     cookie: {
       secure: true,
       httpOnly: true,
       sameSite: 'lax',
       maxAge: 30 * 60 * 1000,  // 30 minutes
     },
   }));

   // Disable Express fingerprinting
   app.disable('x-powered-by');

   // Never serve raw errors in production
   app.use((err, req, res, next) => {
     console.error(err);  // log internally
     res.status(500).json({ error: 'Internal server error' });  // generic externally
   });
   ```

3. **Spring Boot: security configuration baseline**:

   ```java
   @Configuration
   @EnableWebSecurity
   public class SecurityConfig {

       @Bean
       public SecurityFilterChain filterChain(HttpSecurity http) throws Exception {
           http
               // HTTPS enforcement
               .requiresChannel(channel -> channel.anyRequest().requiresSecure())
               // HSTS
               .headers(headers -> headers
                   .httpStrictTransportSecurity(hsts -> hsts
                       .maxAgeInSeconds(31536000)
                       .includeSubDomains(true)
                       .preload(true))
                   .contentTypeOptions(Customizer.withDefaults())
                   .frameOptions(frame -> frame.deny())
               )
               // Session management
               .sessionManagement(session -> session
                   .sessionCreationPolicy(SessionCreationPolicy.IF_REQUIRED)
                   .maximumSessions(1)
                   .sessionFixation(SessionManagementConfigurer.SessionFixationConfigurer::newSession)
               )
               // CSRF enabled by default — verify not disabled
               .csrf(Customizer.withDefaults());
           return http.build();
       }
   }
   ```

   ```properties
   # application-production.properties
   server.ssl.enabled=true
   server.error.include-stacktrace=never  # no stack traces in responses
   server.error.include-exception=false
   management.endpoints.web.exposure.include=health,info  # not all actuator endpoints
   spring.jpa.show-sql=false  # no SQL in logs
   ```

4. **Database connection security defaults**:

   ```python
   # PostgreSQL (psycopg3) — secure connection defaults
   import psycopg

   conn = psycopg.connect(
       conninfo=f"postgresql://{user}:{password}@{host}/{dbname}",
       sslmode="require",           # enforce TLS
       connect_timeout=10,          # fail fast, not hang
       options="-c statement_timeout=30000"  # 30s max query time
   )

   # SQLAlchemy — pool security settings
   from sqlalchemy import create_engine

   engine = create_engine(
       DATABASE_URL,
       pool_pre_ping=True,          # verify connection before use
       pool_size=10,
       max_overflow=20,
       pool_timeout=30,
       connect_args={
           "sslmode": "require",
           "connect_timeout": 10,
           "options": "-c statement_timeout=30000"
       }
   )
   ```

5. **Infrastructure security defaults checklist**:

   ```bash
   # Verify security defaults via Django check
   python manage.py check --deploy --fail-level WARNING

   # Node.js: audit npm packages for known issues with secure defaults
   npm audit --audit-level=high

   # Scan for security misconfigurations with bandit (Python)
   bandit -r src/ -ll  # only HIGH and MEDIUM severity

   # Check TLS configuration
   testssl.sh --severity HIGH --color 0 https://app.company.com
   ```

   ```markdown
   Security defaults checklist for new services:
   □ DEBUG=False / equivalent disabled in production
   □ HTTPS enforced (redirect + HSTS)
   □ Session cookies: Secure + HttpOnly + SameSite=Lax
   □ CSRF protection enabled (not disabled)
   □ Security headers: X-Frame-Options, X-Content-Type-Options, Referrer-Policy
   □ Generic error messages (no stack traces)
   □ Database TLS enforced
   □ Logging does not include passwords/tokens/PII
   □ Dependency versions pinned in lockfile
   □ No secrets in environment variables or config files
   ```

## Rules

- Run `python manage.py check --deploy` (Django) or equivalent framework security check before every production deployment — frameworks detect their own misconfigurations.
- Security middleware must run before application code — placing `SecurityMiddleware` after custom middleware means requests bypass security headers.
- `DEBUG = True` in production is a critical misconfiguration — Django, Rails, and Flask all expose credentials and internal state in debug mode.
- Default deny for HTTP methods — explicitly allow GET/POST/PUT/DELETE and reject everything else (HEAD, TRACE, CONNECT, OPTIONS unless needed).

## Common Mistakes

- **Disabling CSRF protection for "convenience"** — CSRF protection is disabled per-view as a shortcut in many tutorials; each disabled view is a vulnerability.
- **`SameSite=None` cookies without `Secure`** — `SameSite=None` requires `Secure` per the spec; without it, the cookie is rejected by modern browsers and the CSRF protection is lost.
- **Framework security docs read once, not maintained** — security configurations accumulate technical debt; review against current framework documentation at least annually.
- **Permissive CSP as a "starting point"** — `Content-Security-Policy: default-src *` is worse than no CSP because it gives false assurance; start strict and loosen only as needed.
