---
name: review-iot-security
description: Use when auditing an IoT product or deployment — systematically checking all 10 OWASP IoT Top 10 vulnerability classes with test procedures, firmware analysis commands, and network scanning.
source: 'OWASP IoT Top 10 2018 (owasp.org/www-project-internet-of-things/); ETSI EN 303 645; NIST SP 800-193; Shodan IoT security research; Binwalk documentation'
tags: [security, owasp, iot, audit, firmware, network, embedded, hardware]
---

# Review IoT Security

Audit IoT devices against the OWASP IoT Top 10 (2018) using network scanning, firmware extraction, and traffic analysis — covering I1 through I10 with specific test commands and remediation guidance.

## Why This Is Best Practice

**Adopted by:** OWASP IoT Top 10 (2018) is the authoritative IoT vulnerability taxonomy. ETSI EN 303 645 (2020) is the European IoT security standard, now mandatory for consumer IoT in the EU. NIST SP 800-193 provides US federal guidance. IoT security assessors at Rapid7, Pen Test Partners (who publicized IoT vulnerabilities at scale including smart locks, baby monitors, and industrial controllers), and FirmWire all use structured I1–I10 equivalent checklists. Pwnie Express' IoT security assessments use network scanning + firmware analysis as the two primary audit phases.
**Impact:** Pen Test Partners' "IoT Hall of Shame" (2016–2023) documents 50+ consumer IoT products with critical vulnerabilities across all I1–I10 categories. Shodan regularly finds millions of exposed IoT devices in each category. The FCC's 2023 "Cyber Trust Mark" for IoT devices requires I1–I10 coverage as a condition of certification. UK's PSTI Act 2024 makes I1 (default passwords), I3 (insecure interfaces), and I9 (insecure software) legally mandatory — violations carry fines of £10M or 4% of global revenue.
**Why best:** IoT audits without a structured checklist miss vulnerability classes — a reviewer who assesses network services may overlook physical interfaces (UART/JTAG) or privacy data collection. I1–I10 provides completeness for both network-facing and physical attack surfaces, which is unique to IoT vs. web/mobile application security reviews.

Sources: OWASP IoT Top 10 (2018); ETSI EN 303 645; Pen Test Partners IoT security research; Shodan State of IoT Security (2023)

## Steps

### Pre-Audit Reconnaissance

```bash
# Network discovery — find IoT devices on network
nmap -sV -O --script=banner 192.168.1.0/24

# Shodan CLI — check if device is exposed to internet
shodan host <DEVICE_IP>

# Extract firmware from vendor's OTA update URL or device flash
binwalk -e firmware.bin    # extract filesystem
strings firmware.bin | grep -E "password|secret|api_key|token|admin"
```

### I1 — Weak, Guessable, or Hardcoded Passwords

```bash
# Check for default credentials
hydra -L /usr/share/wordlists/routers.txt -P /usr/share/wordlists/routers.txt \
  192.168.1.1 http-form-post "/login:user=^USER^&pass=^PASS^:Login failed"

# Check firmware for hardcoded credentials
binwalk -e firmware.bin
grep -r "password\|passwd\|secret\|token" _firmware.bin.extracted/
strings _firmware.bin.extracted/squashfs-root/etc/passwd
```

**Findings to look for:** Universal default password, same password for all devices, credentials in firmware strings.

**Fix:** Unique per-device credentials at manufacture, forced change on first login, no factory default shared across devices.

### I2 — Insecure Network Services

```bash
# Port scan device
nmap -sV -p- --script=vuln 192.168.1.X

# Check for Telnet (23), unencrypted HTTP (80), SNMP v1/v2 (161), FTP (21)
nmap -sV -p 21,23,80,161 192.168.1.X

# Check Telnet banner for software version
nc 192.168.1.X 23
```

**Findings to look for:** Telnet open, HTTP without auth, SNMP v1/v2 with community string "public"/"private".

**Fix:** Disable Telnet, enforce HTTPS, upgrade to SNMP v3 with authentication.

### I3 — Insecure Ecosystem Interfaces

```bash
# Intercept mobile app API calls (set up Burp Suite proxy on phone)
# Check for: HTTP not HTTPS, missing auth headers, broken object-level auth

# Test cloud API for auth bypass
curl -H "X-Device-ID: victim_device_id" https://api.vendor.com/v1/device/status

# Check cloud dashboard for IDOR
# Change device_id in URL to another user's device
```

**Findings to look for:** Unauthenticated REST API endpoints, missing authorization on cloud dashboard, IDOR in device management APIs.

### I4 — Lack of Secure Update Mechanism

```bash
# Capture OTA update traffic
tcpdump -i eth0 host 192.168.1.X -w ota_capture.pcap
wireshark ota_capture.pcap

# Check: HTTP vs HTTPS, signature verification
# Attempt downgrade: serve old firmware version to device
```

**Findings to look for:** OTA over HTTP, no signature verification, no rollback protection, missing integrity check.

**Fix:** HTTPS OTA with certificate pinning, ECDSA signature verification in bootloader, monotonic counter for rollback prevention.

### I5 — Use of Insecure or Outdated Components

```bash
# Extract and check firmware components
binwalk -e firmware.bin
# Check versions in extracted filesystem
cat _firmware.bin.extracted/squashfs-root/etc/openwrt_release 2>/dev/null
find _firmware.bin.extracted/ -name "*.so" -exec strings {} \; | grep -E "OpenSSL [0-9]"

# Check for CVEs in identified components
trivy fs _firmware.bin.extracted/squashfs-root/
```

**Findings to look for:** OpenSSL < 3.0, Linux kernel < 5.15, BusyBox < 1.35 with known CVEs.

### I6 — Insufficient Privacy Protection

```bash
# Capture device traffic and check for PII
tcpdump -i eth0 host 192.168.1.X -w traffic.pcap
tshark -r traffic.pcap -T fields -e http.request.uri -e http.file_data | grep -E "email|name|location|mac"

# Check cloud storage for unencrypted PII
# Review privacy policy against actual data collected
```

**Findings to look for:** Location data sent to third-party analytics, device MAC/serial sent without user consent, no data deletion mechanism.

### I7 — Insecure Data Transfer and Storage

```bash
# Check for unencrypted local storage
binwalk -e firmware.bin
find _firmware.bin.extracted/ -name "*.db" -o -name "*.sqlite" | xargs strings | \
  grep -E "password|key|secret|token"

# Check MQTT traffic
mosquitto_sub -h 192.168.1.X -t '#' -v  # unauth MQTT
```

**Findings to look for:** SQLite databases with plaintext credentials, unencrypted MQTT on port 1883, config files with API keys.

### I8 — Lack of Device Management

```bash
# Check if devices can be identified and managed
# Can you revoke one device without affecting others?
# Is there certificate rotation mechanism?
# Is there a decommissioning workflow?

# Check for shared credentials across fleet
# If one device's cert matches another's cert — shared credential
```

**Findings to look for:** Shared certificates/keys across device fleet, no revocation mechanism, no decommissioning procedure.

### I9 — Insecure Default Settings

```bash
# Check for unnecessary enabled services
nmap -sV 192.168.1.X  # compare open ports to vendor documentation

# Check for unnecessary file shares
smbclient -L 192.168.1.X -N 2>/dev/null
showmount -e 192.168.1.X 2>/dev/null

# Check debug mode in web interface
curl http://192.168.1.X/debug
curl http://192.168.1.X/phpinfo.php
```

**Findings to look for:** SMB shares open by default, debug endpoints accessible, verbose error messages with stack traces.

### I10 — Lack of Physical Hardening

```bash
# Physical assessment checklist:
# [ ] UART/JTAG headers present on PCB (test with FTDI adapter)
# [ ] SPI/NAND flash chip accessible (test with flashrom)
# [ ] Flash dump possible without authentication
# [ ] OTP fuses burned (debug disabled in production)
```

```bash
# Test UART access
screen /dev/ttyUSB0 115200
# Receiving boot log = UART accessible

# Test flash extraction
flashrom -p ft2232_spi:type=2232H,port=A -r flash_dump.bin
binwalk -e flash_dump.bin
```

**Findings to look for:** Unauthenticated root shell on UART, JTAG accessible post-production, flash readable without authentication.

### Audit Summary Checklist

```
□ I1  Passwords — per-device unique, no hardcoded in firmware
□ I2  Network services — no Telnet/FTP/unauthenticated HTTP
□ I3  Ecosystem interfaces — authenticated APIs, no IDOR
□ I4  Update mechanism — signed firmware, HTTPS OTA, rollback protection
□ I5  Components — no components with known HIGH/CRITICAL CVEs
□ I6  Privacy — data minimization, consent, deletion mechanism
□ I7  Data storage — encrypted at rest, TLS in transit
□ I8  Device management — per-device certs, revocation, decommission
□ I9  Default settings — only required services enabled
□ I10 Physical — debug interfaces disabled in production firmware
```

## Rules

- Run this audit on every new hardware revision — a new PCB revision may re-enable JTAG headers removed in a prior version.
- Network scan (I1–I3) and firmware analysis (I4–I7, I9) are independent — do both, as firmware may show vulnerabilities not visible from the network.
- Physical assessment (I10) requires physical access to the device — lab testing one unit covers the design; manufacturing QA must verify each production unit has debug interfaces disabled.

## Common Mistakes

- **Auditing development firmware instead of production firmware** — debug features enabled in development are often forgotten in production builds; audit the firmware image actually shipped to customers.
- **Accepting vendor documentation of "no UART" without testing** — test physical interfaces with an oscilloscope or logic analyzer; undocumented UART ports appear in many consumer IoT products.
- **Not testing with a production account** — administrative backdoors may only appear on non-admin accounts; test with the same account type as end users.
