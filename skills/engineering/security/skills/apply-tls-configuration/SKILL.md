---
name: apply-tls-configuration
description: Use when configuring HTTPS on a web server, API gateway, or load balancer — covers cipher suites, TLS version enforcement, certificate management, and HSTS.
source: 'OWASP Transport Layer Security Cheat Sheet (owasp.org/www-project-cheat-sheets); Mozilla SSL Configuration Generator; NIST SP 800-52 Rev 2; CWE-326'
tags: [security, owasp, tls, https, ssl, certificates, devops, developer]
---

# Apply TLS Configuration

Configure TLS 1.2+ with modern cipher suites, valid certificates, and HSTS — eliminating downgrade attacks, weak encryption, and certificate validation failures.

## Why This Is Best Practice

**Adopted by:** NIST SP 800-52 Rev 2 mandates TLS 1.2 minimum for US federal systems. PCI DSS v4.0 Requirement 4.2 requires TLS 1.2+ and explicitly prohibits SSL/TLS 1.0/1.1. Google Chrome marks all HTTP as "Not Secure" since 2018. Let's Encrypt has issued over 3 billion certificates, removing cost as a barrier. Mozilla's Server Side TLS guidelines are the de facto industry reference used by AWS, Cloudflare, and Fastly.
**Impact:** POODLE (SSLv3), BEAST (TLS 1.0), SWEET32 (3DES), LUCKY13, and DROWN are all cryptographic attacks against weak TLS versions and ciphers — each enabled by misconfigured servers. Qualys SSL Labs A+ rating requires TLS 1.2+, forward secrecy, and HSTS — servers scoring below B are vulnerable to known-exploitable attacks. Forward secrecy (ECDHE) ensures past sessions cannot be decrypted even if the private key is later compromised.
**Why best:** Custom cipher configuration over accepting server defaults — because OpenSSL and platform defaults often include legacy ciphers for backward compatibility. Mozilla's "Intermediate" profile balances security and broad client compatibility; "Modern" maximizes security for internal services.

Sources: OWASP TLS Cheat Sheet; Mozilla SSL Config Generator; NIST SP 800-52 Rev 2; Qualys SSL Labs research; CWE-326

## Steps

1. **Enforce TLS 1.2 minimum; prefer TLS 1.3**:

   **Nginx:**
   ```nginx
   ssl_protocols TLSv1.2 TLSv1.3;
   ```

   **Apache:**
   ```apache
   SSLProtocol -all +TLSv1.2 +TLSv1.3
   ```

   **Node.js:**
   ```javascript
   const tls = require('tls');
   // Node defaults to TLS 1.2+ since v12 — verify with:
   console.log(tls.DEFAULT_MIN_VERSION); // should be TLSv1.2
   ```

2. **Configure forward-secret cipher suites (Mozilla Intermediate profile)**:

   **Nginx:**
   ```nginx
   ssl_ciphers ECDHE-ECDSA-AES128-GCM-SHA256:ECDHE-RSA-AES128-GCM-SHA256:ECDHE-ECDSA-AES256-GCM-SHA384:ECDHE-RSA-AES256-GCM-SHA384:ECDHE-ECDSA-CHACHA20-POLY1305:ECDHE-RSA-CHACHA20-POLY1305:DHE-RSA-AES128-GCM-SHA256:DHE-RSA-AES256-GCM-SHA384;
   ssl_prefer_server_ciphers off;
   ```

   All ciphers use ECDHE/DHE (forward secrecy) and AEAD (authenticated encryption). No RC4, 3DES, or CBC-mode ciphers.

3. **Obtain and renew certificates automatically** — use Let's Encrypt (free, 90-day, auto-renewable):

   ```bash
   # Certbot (most common)
   certbot --nginx -d example.com -d www.example.com
   # Auto-renewal via cron (certbot installs this automatically)
   ```

   For AWS: use ACM (AWS Certificate Manager) — free certificates, auto-renewed, no manual rotation.
   For Kubernetes: use cert-manager with Let's Encrypt ClusterIssuer.

4. **Enable HSTS with preload**:

   ```nginx
   add_header Strict-Transport-Security "max-age=63072000; includeSubDomains; preload" always;
   ```

   Submit to HSTS preload list at `hstspreload.org` after testing. See `apply-security-headers` for full context.

5. **Enable OCSP stapling** — server proves certificate validity without client round-trip:

   ```nginx
   ssl_stapling on;
   ssl_stapling_verify on;
   ssl_trusted_certificate /path/to/chain.pem;
   resolver 8.8.8.8 1.1.1.1 valid=300s;
   ```

6. **Set DH parameters for DHE ciphers** (2048-bit minimum):

   ```bash
   openssl dhparam -out /etc/nginx/dhparam.pem 2048
   ```

   ```nginx
   ssl_dhparam /etc/nginx/dhparam.pem;
   ```

7. **Test with Qualys SSL Labs** — run `ssllabs.com/ssltest/` against your domain. Target A+ grade. Address all flagged findings before production launch.

8. **Validate certificate chain** — ensure full chain (leaf + intermediates) is served. Missing intermediates cause failures on some clients:

   ```bash
   openssl s_client -connect example.com:443 -showcerts < /dev/null
   # Verify chain shows: leaf → intermediate → root
   ```

## Rules

- Never use self-signed certificates in production — clients will either reject them or users will train themselves to click through warnings.
- Certificate pinning (HPKP) is deprecated — it caused widespread outages when misused; do not implement.
- `ssl_prefer_server_ciphers off` is correct for TLS 1.3 (client chooses) — leave it off.
- Wildcard certificates (`*.example.com`) cover one subdomain level only — they don't cover `sub.sub.example.com`.

## Common Mistakes

- **Using default OpenSSL ciphers** — includes RC4, 3DES, and export ciphers on older OpenSSL versions.
- **Not renewing certificates before expiry** — automate renewal; calendar reminders are insufficient at scale.
- **Serving HTTP and HTTPS on the same domain without a redirect** — configure HTTP to redirect 301 to HTTPS permanently.
- **Forgetting to include the intermediate certificate** — chain failures are the #1 TLS deployment error.
