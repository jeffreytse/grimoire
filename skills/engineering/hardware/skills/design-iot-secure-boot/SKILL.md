---
name: design-iot-secure-boot
description: Use when designing firmware update mechanisms for IoT devices — implementing code signing, secure boot chains, and rollback protection to prevent malicious firmware from running on deployed hardware.
source: 'OWASP IoT Top 10 I4 (owasp.org/www-project-internet-of-things/); NIST SP 800-193 (Platform Firmware Resiliency); IETF SUIT manifest standard (RFC 9019); Arm TrustZone secure boot documentation'
tags: [security, owasp, iot, firmware, secure-boot, code-signing, embedded, hardware, developer]
---

# Design IoT Secure Boot

Implement cryptographic code signing for firmware, verified boot chain from ROM to application, and rollback prevention — ensuring only authorized firmware runs on IoT devices even after physical access.

## Why This Is Best Practice

**Adopted by:** OWASP IoT Top 10 I4 (Lack of Secure Update Mechanism). NIST SP 800-193 (Platform Firmware Resiliency) is the US federal standard for firmware integrity. The IETF SUIT (Software Updates for Internet of Things) manifest standard (RFC 9019) defines the protocol. Arm TrustZone secure boot is implemented in the majority of modern IoT SoCs (STM32H5, NXP i.MX RT, Nordic nRF9161). AWS IoT Jobs, Azure Device Update, and Mbed TLS all implement SUIT-compatible update pipelines.
**Impact:** Mirai botnet (2016, 600,000+ devices compromised) and its variants exploited devices with no firmware verification — malicious firmware was flashed over default credentials and telnet. VPNFilter (2018, 500,000 routers compromised) survived factory resets because it persisted in flash memory beyond the OS partition. NIST SP 800-193 states that "the ability to detect and recover from corrupt or maliciously modified firmware is a foundational security control" — without it, physical access equals full device compromise.
**Why best:** OTA updates without signature verification allow any attacker who intercepts the update channel (or the update server itself) to push malicious firmware to the entire device fleet. Code signing ensures the device verifies the firmware author before executing it. Secure boot ties the trust chain to hardware (ROM key, OTP fuses) so no software-level attack can bypass it.

Sources: OWASP IoT Top 10 I4; NIST SP 800-193; IETF RFC 9019 (SUIT Manifest); Arm Platform Security Architecture (PSA) documentation

## Steps

1. **Sign firmware images before distribution**:

   ```bash
   # Generate signing key (keep private key in HSM for production)
   openssl ecparam -name prime256v1 -genkey -noout -out firmware-signing-key.pem
   openssl ec -in firmware-signing-key.pem -pubout -out firmware-signing-key-pub.pem

   # Sign firmware binary
   openssl dgst -sha256 -sign firmware-signing-key.pem \
     -out firmware-v1.2.3.bin.sig firmware-v1.2.3.bin

   # Verify (device does this)
   openssl dgst -sha256 -verify firmware-signing-key-pub.pem \
     -signature firmware-v1.2.3.bin.sig firmware-v1.2.3.bin
   ```

2. **Implement ECDSA signature verification in the bootloader**:

   ```c
   #include "mbedtls/ecdsa.h"
   #include "mbedtls/sha256.h"

   /* Public key burned into OTP fuses or stored in secure storage */
   static const uint8_t SIGNING_PUBLIC_KEY[] = {
       0x04, /* uncompressed point marker */
       /* 64-byte public key X and Y coordinates */
   };

   int verify_firmware(const uint8_t *firmware, size_t firmware_len,
                       const uint8_t *signature, size_t sig_len) {
       uint8_t hash[32];
       mbedtls_sha256(firmware, firmware_len, hash, 0);

       mbedtls_ecdsa_context ecdsa;
       mbedtls_ecdsa_init(&ecdsa);

       int ret = mbedtls_ecp_point_read_binary(
           &ecdsa.private_grp, &ecdsa.private_Q,
           SIGNING_PUBLIC_KEY, sizeof(SIGNING_PUBLIC_KEY)
       );
       if (ret != 0) return VERIFY_FAILED;

       ret = mbedtls_ecdsa_read_signature(&ecdsa, hash, sizeof(hash),
                                          signature, sig_len);
       mbedtls_ecdsa_free(&ecdsa);

       return (ret == 0) ? VERIFY_OK : VERIFY_FAILED;
   }
   ```

3. **Implement rollback protection using a monotonic counter**:

   ```c
   /* Store anti-rollback counter in OTP fuses — cannot be decremented */

   typedef struct {
       uint32_t version;       /* firmware semantic version */
       uint32_t anti_rollback; /* must be >= current OTP counter */
   } FirmwareHeader;

   int check_rollback(const FirmwareHeader *header) {
       uint32_t otp_counter = read_otp_anti_rollback_counter();

       /* New firmware must have counter >= current value */
       if (header->anti_rollback < otp_counter) {
           return ROLLBACK_DETECTED;  /* reject downgrade */
       }

       /* After successful boot, advance OTP counter to new firmware's value */
       if (header->anti_rollback > otp_counter) {
           /* Only write OTP after successful boot, not at verification time */
           schedule_otp_update(header->anti_rollback);
       }
       return ROLLBACK_OK;
   }
   ```

4. **Implement A/B partition update with revert capability**:

   ```c
   /* Dual-bank flash layout:
    * Bank A: active firmware (currently running)
    * Bank B: update target (download here before activating)
    */

   typedef enum { BANK_A = 0, BANK_B = 1 } FlashBank;

   UpdateResult apply_update(const uint8_t *firmware, size_t len,
                             const uint8_t *signature, size_t sig_len) {
       /* Step 1: verify before touching active partition */
       if (verify_firmware(firmware, len, signature, sig_len) != VERIFY_OK) {
           return UPDATE_SIGNATURE_INVALID;
       }

       FlashBank inactive = get_inactive_bank();

       /* Step 2: erase and write to inactive bank only */
       flash_erase_bank(inactive);
       flash_write(get_bank_address(inactive), firmware, len);

       /* Step 3: verify what was written */
       if (!verify_bank_contents(inactive, firmware, len)) {
           return UPDATE_WRITE_ERROR;
       }

       /* Step 4: set boot flag — reverts to A if B fails to boot 3 times */
       set_pending_boot_bank(inactive, max_boot_attempts=3);

       return UPDATE_PENDING_REBOOT;
   }
   ```

5. **OTA update client with transport security**:

   ```c
   /* Use TLS 1.2+ with certificate pinning for update server */
   #include "mbedtls/ssl.h"

   #define UPDATE_SERVER_CERT_FINGERPRINT \
       "\xAB\xCD\xEF..." /* SHA-256 of server certificate */

   int download_update(const char *url, uint8_t *buffer, size_t *len) {
       /* Pin server certificate — prevent MitM from serving malicious firmware */
       ssl_config.verify_mode = MBEDTLS_SSL_VERIFY_REQUIRED;
       ssl_config.ca_chain = &server_cert_chain;

       /* Verify fingerprint after TLS handshake */
       const mbedtls_x509_crt *peer_cert = mbedtls_ssl_get_peer_cert(&ssl);
       uint8_t fingerprint[32];
       mbedtls_sha256(peer_cert->raw.p, peer_cert->raw.len, fingerprint, 0);

       if (memcmp(fingerprint, UPDATE_SERVER_CERT_FINGERPRINT, 32) != 0) {
           return ERROR_CERT_MISMATCH;
       }
       /* ... download proceeds ... */
   }
   ```

## Rules

- The signing private key must never be stored on the device — it lives in an HSM or air-gapped signing server; the device only stores the public key.
- Never verify firmware signature in application code — the bootloader must verify before jumping to application; application-level verification is too late.
- Test rollback protection explicitly: flash an older firmware version and verify the device rejects it.
- Disable JTAG/SWD debug interfaces in production via OTP fuses — debuggers allow bypassing signature verification.

## Common Mistakes

- **Burning the same signing key into all devices** — a single key compromise affects the entire fleet; use a key hierarchy with per-batch keys where feasible.
- **Checking version number for rollback prevention** — version numbers are in flash and changeable; use OTP hardware counters that can only increment.
- **Downloading firmware over HTTP** — even with a valid signature, HTTP OTA allows detection of update patterns and MITM size information leakage; use HTTPS.
- **Not testing the revert path** — A/B update that can't revert leaves devices bricked if new firmware has a boot failure; test the failure scenario.
