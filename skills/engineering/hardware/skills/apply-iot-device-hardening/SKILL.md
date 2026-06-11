---
name: apply-iot-device-hardening
description: Use when preparing IoT devices for production deployment — removing default credentials, disabling debug interfaces, encrypting stored data, and implementing physical tamper detection.
source: 'OWASP IoT Top 10 I9 I10 (owasp.org/www-project-internet-of-things/); ETSI EN 303 645; NIST SP 800-193; UK PSTI Act (2024 Product Security and Telecommunications Infrastructure Act)'
tags: [security, owasp, iot, device-hardening, credentials, tamper-detection, embedded, hardware, developer]
---

# Apply IoT Device Hardening

Remove default credentials, disable JTAG/UART debug interfaces, encrypt data at rest, and implement physical tamper detection — preventing both remote and physical compromise of deployed IoT devices.

## Why This Is Best Practice

**Adopted by:** OWASP IoT Top 10 I9 (Insecure Default Settings) and I10 (Lack of Physical Hardening). ETSI EN 303 645 Section 5.1 mandates unique per-device credentials and prohibits default passwords. The UK Product Security and Telecommunications Infrastructure (PSTI) Act 2024 bans universal default passwords for consumer IoT — manufacturers face fines up to £10M. NIST SP 800-193 requires firmware resilience against physical access. Medical IoT (FDA guidance 2023), industrial IoT (IEC 62443), and automotive IoT (UNECE WP.29) all mandate these controls.
**Impact:** Mirai botnet's primary infection vector was default credentials (admin/admin, root/1234) on Telnet and HTTP — 600,000 devices compromised within weeks of release. Shodan 2023 found 2.3 million IoT devices still using factory default credentials. Physical hardware attacks allow credential extraction from unprotected UART/JTAG interfaces in under 5 minutes with a $30 FTDI adapter and publicly available firmware extraction tools. NIST estimates 30% of IoT breaches involve physical access to the device.
**Why best:** Remote hardening (network services, TLS) protects against internet-facing attacks; physical hardening adds a layer against supply chain attacks, insider threats, and theft of deployed devices. Default credentials are unique in that they're publicly known to attackers before the device is even deployed — per-device unique credentials eliminate this entire attack class at zero ongoing cost.

Sources: OWASP IoT Top 10 I9, I10; ETSI EN 303 645 sections 5.1 and 5.2; UK PSTI Act 2024; FDA "Cybersecurity in Medical Devices" guidance (2023)

## Steps

1. **Generate unique per-device credentials at manufacturing time**:

   ```python
   # Manufacturing provisioning script — runs once per device
   import secrets
   import hashlib

   def provision_device(device_serial: str, device_mac: str) -> dict:
       # Generate cryptographically random password (not derived from serial/MAC)
       admin_password = secrets.token_urlsafe(16)

       # Generate unique device certificate
       private_key = generate_ec_key()
       cert = sign_device_certificate(
           private_key=private_key,
           subject=f"CN={device_serial},O=Company,C=US",
           ca_key=MANUFACTURING_CA_KEY
       )

       # Store: password hash in device, plaintext in HSM-backed provisioning DB
       device_config = {
           "admin_password_hash": bcrypt.hash(admin_password),
           "device_cert": cert,
           "device_key": private_key,
           "serial": device_serial,
       }

       # Label on physical device: serial + initial password (for setup only)
       # Force password change on first login
       return device_config

   # NEVER: password = hashlib.md5(device_serial.encode()).hexdigest()
   # NEVER: same default password across entire product line
   ```

2. **Disable JTAG and UART debug interfaces in production**:

   ```c
   /* STM32 — disable JTAG/SWD via option bytes at end of production test */
   void disable_debug_interfaces(void) {
       FLASH_OBProgramInitTypeDef ob_init = {0};
       HAL_FLASHEx_OBGetConfig(&ob_init);

       /* Disable SWD interface — write-once OTP operation */
       ob_init.OptionType = OPTIONBYTE_USER;
       ob_init.USERType = OB_USER_TZEN;
       ob_init.USERConfig = OB_TZEN_ENABLE;  /* Enable TrustZone */

       /* Set RDP (Read-Out Protection) Level 2 — no debug, no readback */
       ob_init.OptionType |= OPTIONBYTE_RDP;
       ob_init.RDPLevel = OB_RDP_LEVEL_2;  /* PERMANENT — cannot be reversed */

       HAL_FLASH_Unlock();
       HAL_FLASH_OB_Unlock();
       HAL_FLASHEx_OBProgram(&ob_init);
       HAL_FLASH_OB_Lock();
       HAL_FLASH_Lock();
   }
   ```

   ```bash
   # For embedded Linux: disable UART console in production
   # Remove console= from kernel command line in bootloader config
   # /etc/inittab or systemd getty: comment out ttyS0 entries

   # Disable JTAG in Device Tree
   # &jtag { status = "disabled"; };
   ```

3. **Encrypt sensitive data at rest using hardware security**:

   ```c
   /* Use hardware AES engine + device-unique key (UID) for storage encryption */
   #include "aes.h"

   /* Derive storage key from device UID — unique per device, never leaves hardware */
   void derive_storage_key(uint8_t key_out[32]) {
       /* STM32 unique device ID: 96-bit value burned at manufacturing */
       uint32_t uid[3] = {
           HAL_GetUIDw0(), HAL_GetUIDw1(), HAL_GetUIDw2()
       };
       /* HKDF with device UID as IKM and fixed info */
       hkdf_sha256(key_out, 32,
                   (uint8_t *)uid, sizeof(uid),
                   (uint8_t *)"storage-key-v1", 15,
                   NULL, 0);
   }

   void encrypt_credentials(const char *plaintext, uint8_t *ciphertext_out) {
       uint8_t key[32];
       derive_storage_key(key);
       uint8_t nonce[12];
       crypto_get_random(nonce, sizeof(nonce));
       aes_gcm_encrypt(key, nonce, (uint8_t *)plaintext, strlen(plaintext),
                       ciphertext_out);
       memset(key, 0, sizeof(key));  /* clear key from stack */
   }
   ```

4. **Implement tamper detection and response**:

   ```c
   /* Hardware tamper pin connected to case-open switch */
   void tamper_detection_init(void) {
       /* Configure tamper pin as interrupt */
       GPIO_InitTypeDef gpio = {
           .Pin = TAMPER_PIN,
           .Mode = GPIO_MODE_IT_RISING,  /* rising = case opened */
           .Pull = GPIO_PULLDOWN,
       };
       HAL_GPIO_Init(TAMPER_PORT, &gpio);
       HAL_NVIC_EnableIRQ(TAMPER_IRQn);
   }

   void TAMPER_IRQHandler(void) {
       if (__HAL_GPIO_EXTI_GET_FLAG(TAMPER_PIN)) {
           /* Log tamper event to persistent counter */
           increment_tamper_counter();
           log_security_event("physical_tamper_detected");

           /* Policy: clear sensitive keys on tamper */
           secure_erase_credentials();

           /* Require re-provisioning — device will not boot normally */
           set_boot_flag(BOOT_FLAG_TAMPERED);
           HAL_NVIC_SystemReset();
       }
   }
   ```

5. **Force credential change on first login**:

   ```python
   # Embedded web server — enforce password change workflow
   def handle_login(username: str, password: str, session) -> Response:
       if not verify_credentials(username, password):
           return Response(401, "Invalid credentials")

       user = get_user(username)

       # Force password change if using factory default
       if user.must_change_password:
           session["pending_password_change"] = True
           return redirect("/change-password")

       session["authenticated"] = True
       return redirect("/dashboard")

   def handle_change_password(new_password: str, session) -> Response:
       if len(new_password) < 12:
           return Response(400, "Password must be at least 12 characters")
       if not has_complexity(new_password):
           return Response(400, "Password must include letters, numbers, and symbols")

       update_password(session["user"], bcrypt.hash(new_password))
       session["pending_password_change"] = False
       return Response(200, "Password updated")
   ```

## Rules

- RDP Level 2 on STM32 (and equivalent on other MCUs) is permanent and irreversible — verify the production firmware is correct before enabling it.
- Per-device credentials must be generated with a CSPRNG — never derived from predictable values (MAC address, serial number, hash of serial number).
- Tamper detection response must clear keys before reset — a tampered device that reboots normally with keys intact defeats the tamper protection.
- Debug access restrictions must be in the release build configuration, verified by the manufacturing test procedure.

## Common Mistakes

- **Same credential on all devices of the same model** — a single device's credential extraction compromises all devices of that model; per-device credentials are mandatory.
- **UART debug console left enabled in release build** — `CONFIG_CONSOLE_UART=y` in Kconfig; audit release build configs before tapeout/production flash.
- **Storing credentials in plaintext flash** — flash memory is readable with a logic analyzer or by desoldering; encrypt all stored credentials with a hardware-derived key.
- **Tamper detection only in application, not bootloader** — an attacker can replace the application; tamper detection must be in the first-stage bootloader or hardware-enforced.
