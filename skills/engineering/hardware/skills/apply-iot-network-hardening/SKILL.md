---
name: apply-iot-network-hardening
description: Use when designing IoT device network interfaces — disabling unnecessary services, enforcing TLS for all communication, segmenting devices on isolated VLANs, and securing management APIs.
source: 'OWASP IoT Top 10 I2 I3 (owasp.org/www-project-internet-of-things/); NIST SP 800-82 (Industrial Control Systems Security); ETSI EN 303 645 (IoT Cybersecurity); Shodan IoT security research'
tags: [security, owasp, iot, network, tls, vlan, firewall, embedded, hardware, developer]
---

# Apply IoT Network Hardening

Disable unnecessary network services, enforce mutual TLS for device-cloud communication, segment IoT devices on isolated VLANs, and restrict management interfaces to authorized addresses — preventing network-based device compromise.

## Why This Is Best Practice

**Adopted by:** OWASP IoT Top 10 I2 (Insecure Network Services) and I3 (Insecure Ecosystem Interfaces). ETSI EN 303 645 (European IoT cybersecurity standard, 2020) mandates disabling unused services and TLS for all communications. NIST SP 800-82 guides industrial IoT network segmentation. Shodan's 2023 report found 1.5 million IoT devices exposing Telnet (port 23), 500,000 exposing unencrypted MQTT (port 1883), and 2 million exposing unauthenticated HTTP management interfaces.
**Impact:** Mirai botnet (2016) infected 600,000 devices primarily by scanning for open Telnet and using default credentials — the services had no legitimate user need but were enabled by default. The Shodan IoT exposure report estimates 40% of internet-connected IoT devices expose at least one unauthenticated management interface. CISA's 2022 IoT advisory found that network segmentation failures are the primary vector allowing IoT device compromise to pivot to IT networks, causing incidents including the 2021 Oldsmar water treatment plant hack.
**Why best:** Open network services on IoT devices are attack surface that cannot be patched without disabling them — each Telnet port and unauthenticated HTTP endpoint is a permanent vulnerability unless removed from the firmware. Network segmentation limits blast radius: a compromised IoT device on its own VLAN cannot reach internal servers or workstations, containing the incident.

Sources: OWASP IoT Top 10 I2, I3; ETSI EN 303 645 section 4; Shodan State of IoT Security (2023); CISA "Securing the Internet of Things" advisory (2022)

## Steps

1. **Disable all unnecessary network services in firmware**:

   ```bash
   # On embedded Linux (OpenWRT/Yocto) — disable services at build time

   # Remove or disable Telnet entirely
   systemctl disable telnetd 2>/dev/null
   opkg remove telnet

   # Disable unauthenticated FTP
   systemctl disable vsftpd

   # Verify only intended ports are open
   netstat -tlnp | grep LISTEN
   # Should show only: SSH (22) on management interface, HTTPS (443) or MQTT-TLS (8883)
   ```

   ```c
   /* In embedded RTOS (FreeRTOS/Zephyr) — only start services explicitly required */
   void network_init(void) {
       /* DO NOT start: Telnet server, unauthenticated HTTP, SNMP v1/v2 */
       mqtt_tls_client_start();   /* MQTT over TLS only */
       https_server_start();      /* HTTPS management only */
       /* All other services: compile out via Kconfig */
   }
   ```

2. **Enforce TLS 1.2+ for all network communication**:

   ```c
   /* MQTT over TLS (Zephyr RTOS + Mbed TLS) */
   #include "mqtt_client.h"
   #include "mbedtls/ssl.h"

   struct mqtt_tls_config tls_cfg = {
       .peer_verify = TLS_PEER_VERIFY_REQUIRED,
       .cipher_list = "TLS-ECDHE-RSA-WITH-AES-128-GCM-SHA256:"
                      "TLS-ECDHE-RSA-WITH-AES-256-GCM-SHA384",
       .ca_cert = broker_ca_cert,
       .ca_cert_len = sizeof(broker_ca_cert),
       /* Mutual TLS: device presents its own certificate */
       .client_cert = device_cert,
       .private_key = device_private_key,
   };

   struct mqtt_client_config client_cfg = {
       .broker = "iot.company.com",
       .port = 8883,  /* MQTT over TLS — not 1883 (plaintext) */
       .tls = &tls_cfg,
   };
   ```

3. **Segment IoT devices on an isolated VLAN**:

   ```
   Network architecture for IoT deployment:

   Internet
       │
   Router/Firewall
       ├── VLAN 1 (Corporate): 192.168.1.0/24
       │   ├── Workstations
       │   └── Servers
       ├── VLAN 10 (IoT): 10.0.10.0/24
       │   ├── IoT Device A (10.0.10.1)
       │   ├── IoT Device B (10.0.10.2)
       │   └── IoT Hub (10.0.10.100)
       └── VLAN 20 (DMZ): 192.168.20.0/24

   Firewall rules:
   - VLAN 10 → Internet: ALLOW HTTPS (443), MQTT-TLS (8883) only
   - VLAN 10 → VLAN 1: DENY ALL (IoT cannot reach corporate network)
   - VLAN 1 → VLAN 10: ALLOW SSH to IoT Hub only (for management)
   ```

   ```bash
   # pfSense/OPNsense — deny IoT to LAN rule
   # Interface: IoT_VLAN (10.0.10.0/24)
   # Action: Block
   # Source: IoT_VLAN net
   # Destination: LAN net
   # Protocol: any
   ```

4. **Restrict management interface to specific addresses**:

   ```nginx
   # Nginx — restrict admin UI to management subnet only
   server {
       listen 443 ssl;
       location /admin {
           allow 10.0.1.0/24;   # management VLAN only
           deny all;
           proxy_pass http://localhost:8080;
       }
   }
   ```

   ```c
   /* In-firmware ACL for embedded HTTP server */
   bool is_management_ip_allowed(uint32_t client_ip) {
       /* Only allow management subnet: 10.0.1.0/24 */
       return (client_ip & 0xFFFFFF00) == 0x0A000100;
   }

   void http_server_request_handler(struct http_request *req) {
       if (is_admin_path(req->path) &&
           !is_management_ip_allowed(req->client_ip)) {
           send_response(req, HTTP_403_FORBIDDEN);
           return;
       }
   }
   ```

5. **Implement rate limiting and brute-force protection**:

   ```c
   /* Login attempt rate limiting */
   #define MAX_LOGIN_ATTEMPTS 5
   #define LOCKOUT_SECONDS 300

   typedef struct {
       uint32_t ip;
       uint8_t attempts;
       uint32_t lockout_until;
   } LoginAttemptRecord;

   static LoginAttemptRecord attempt_table[MAX_TRACKED_IPS];

   bool check_login_rate_limit(uint32_t client_ip) {
       LoginAttemptRecord *record = find_or_create_record(client_ip);

       if (get_unix_time() < record->lockout_until) {
           return false;  /* locked out */
       }

       record->attempts++;
       if (record->attempts >= MAX_LOGIN_ATTEMPTS) {
           record->lockout_until = get_unix_time() + LOCKOUT_SECONDS;
           log_security_event("brute_force_lockout", client_ip);
           return false;
       }
       return true;
   }
   ```

## Rules

- Telnet must be disabled in production firmware — SSH with key authentication is the minimum acceptable for remote shell access.
- MQTT must use TLS (port 8883) — plaintext MQTT (port 1883) transmits device data and commands in the clear.
- Factory reset must not re-enable disabled services — the secure default must be the only default.
- Device certificates for mutual TLS must be unique per device — shared certificates mean compromising one device compromises all.

## Common Mistakes

- **Enabling Telnet "for debugging" and forgetting to disable in production build** — use Kconfig or build flags to exclude debug services from release builds at compile time.
- **VLAN segmentation without firewall rules** — VLANs alone separate broadcast domains but don't enforce traffic policy; add firewall rules to block inter-VLAN routing.
- **Self-signed certificates without certificate pinning** — without pinning, any attacker with a CA can issue a valid certificate for your IoT broker; pin the specific CA or certificate.
- **Exposing management interface on the same interface as data traffic** — use separate network interfaces for management and data, or bind management to localhost only with SSH tunneling.
