---
name: design-iot-device-management
description: Use when building IoT fleet management infrastructure — designing secure device provisioning, certificate lifecycle, remote management, and decommissioning processes.
source: 'OWASP IoT Top 10 I8 (owasp.org/www-project-internet-of-things/); AWS IoT Core device provisioning documentation; Azure IoT Hub DPS documentation; NIST SP 800-183 (Networks of Things)'
tags: [security, owasp, iot, device-management, provisioning, certificates, fleet-management, hardware, developer]
---

# Design IoT Device Management

Implement zero-touch provisioning with X.509 device certificates, automated certificate rotation, and secure decommissioning — ensuring every device in the fleet has a unique verifiable identity and revocable access.

## Why This Is Best Practice

**Adopted by:** OWASP IoT Top 10 I8 (Lack of Device Management). AWS IoT Core and Azure IoT Hub both implement X.509 certificate-based device identity as their primary authentication mechanism. Google Cloud IoT Core (before deprecation) mandated JWT RS256 with per-device key pairs. NIST SP 800-183 (Networks of Things) defines the reference architecture for IoT device identity and lifecycle management. Apple's Automated Device Enrollment (ADE) demonstrates the enterprise standard for zero-touch provisioning at scale.
**Impact:** Rapid7's 2022 IoT security report found that 35% of enterprise IoT security incidents involved decommissioned devices that retained valid credentials and continued to communicate with back-end systems. The 2020 SolarWinds attack pivoted via management plane access — IoT management channels with persistent long-lived credentials present the same risk. Fleet management without revocation means a single compromised device retains access indefinitely — certificate-based auth with revocation lists limits blast radius.
**Why best:** Shared API keys for device authentication (the alternative) require rotating the key on all devices simultaneously when any device is compromised — operationally infeasible at 10,000+ device scale. Per-device X.509 certificates allow revoking a single device's access immediately via CRL or OCSP without affecting other devices. Automated rotation means certificates expire before an attacker can use an extracted certificate indefinitely.

Sources: OWASP IoT Top 10 I8; AWS IoT Core fleet provisioning documentation; NIST SP 800-183; Rapid7 IoT Security Report (2022)

## Steps

1. **Zero-touch provisioning with X.509 certificates**:

   ```python
   # Manufacturing provisioning service — called once per device at factory
   import boto3
   import json
   from datetime import datetime, timedelta

   iot_client = boto3.client("iot")

   def provision_device(serial_number: str, device_mac: str) -> dict:
       # Create unique thing in IoT registry
       thing_name = f"device-{serial_number}"
       iot_client.create_thing(
           thingName=thing_name,
           attributePayload={
               "attributes": {
                   "serial": serial_number,
                   "mac": device_mac,
                   "provisioned_at": datetime.utcnow().isoformat(),
               }
           }
       )

       # Generate device certificate (or use CSR from device)
       cert_response = iot_client.create_keys_and_certificate(setAsActive=True)

       # Attach policy — device can only publish/subscribe to its own topic
       iot_client.attach_policy(
           policyName="device-policy",
           target=cert_response["certificateArn"]
       )
       iot_client.attach_thing_principal(
           thingName=thing_name,
           principal=cert_response["certificateArn"]
       )

       return {
           "certificate_pem": cert_response["certificatePem"],
           "private_key": cert_response["keyPair"]["PrivateKey"],  # burn to device, destroy here
           "certificate_id": cert_response["certificateId"],
       }
   ```

2. **Device topic policy — restrict each device to its own namespace**:

   ```json
   {
     "Version": "2012-10-17",
     "Statement": [
       {
         "Effect": "Allow",
         "Action": ["iot:Publish"],
         "Resource": "arn:aws:iot:*:*:topic/devices/${iot:Connection.Thing.ThingName}/*"
       },
       {
         "Effect": "Allow",
         "Action": ["iot:Subscribe", "iot:Receive"],
         "Resource": [
           "arn:aws:iot:*:*:topicfilter/devices/${iot:Connection.Thing.ThingName}/commands",
           "arn:aws:iot:*:*:topic/devices/${iot:Connection.Thing.ThingName}/commands"
         ]
       },
       {
         "Effect": "Deny",
         "Action": "iot:*",
         "Resource": "arn:aws:iot:*:*:topic/devices/+/admin/*"
       }
     ]
   }
   ```

3. **Automated certificate rotation before expiry**:

   ```python
   # Lambda function — runs weekly, rotates certs expiring in < 30 days
   import boto3
   from datetime import datetime, timedelta

   iot = boto3.client("iot")

   def rotate_expiring_certificates(event, context):
       paginator = iot.get_paginator("list_certificates")
       for page in paginator.paginate(pageSize=100):
           for cert in page["certificates"]:
               cert_detail = iot.describe_certificate(certificateId=cert["certificateId"])
               expiry = cert_detail["certificateDescription"]["validity"]["notAfter"]

               days_until_expiry = (expiry - datetime.utcnow()).days
               if days_until_expiry < 30:
                   _rotate_certificate(cert["certificateId"])

   def _rotate_certificate(old_cert_id: str):
       # Get the thing attached to this certificate
       principals = iot.list_principal_things(
           principal=f"arn:aws:iot:...:cert/{old_cert_id}"
       )
       thing_name = principals["things"][0]

       # Issue new certificate
       new_cert = iot.create_keys_and_certificate(setAsActive=True)

       # Push new cert to device via shadow or job
       iot.update_thing_shadow(
           thingName=thing_name,
           payload=json.dumps({
               "state": {"desired": {"new_certificate": new_cert["certificatePem"]}}
           })
       )
       # Deactivate old cert after device confirms new cert active
       # (handled by device's cert rotation state machine)
   ```

4. **Secure decommissioning procedure**:

   ```python
   def decommission_device(serial_number: str, reason: str):
       thing_name = f"device-{serial_number}"

       # Step 1: revoke all certificates — device can no longer connect
       principals = iot.list_thing_principals(thingName=thing_name)
       for cert_arn in principals["principals"]:
           cert_id = cert_arn.split("/")[-1]
           iot.update_certificate(
               certificateId=cert_id,
               newStatus="REVOKED"
           )
           iot.detach_thing_principal(thingName=thing_name, principal=cert_arn)
           iot.delete_certificate(certificateId=cert_id, forceDelete=True)

       # Step 2: send decommission command to device (if online)
       try:
           iot.publish(
               topic=f"devices/{thing_name}/commands",
               qos=1,
               payload=json.dumps({"command": "factory_reset", "reason": reason})
           )
       except Exception:
           pass  # device may be offline — certificate revocation is sufficient

       # Step 3: archive device telemetry, then delete from registry
       archive_device_data(thing_name)
       iot.delete_thing(thingName=thing_name)

       # Step 4: audit log
       log_decommission_event(serial_number, reason, timestamp=datetime.utcnow())
   ```

5. **Health monitoring and anomaly detection**:

   ```python
   # AWS IoT Defender or custom Lambda — detect anomalous device behavior
   def check_device_health(thing_name: str) -> dict:
       violations = []

       # Check: device connecting from unexpected IP geolocation
       connections = get_recent_connections(thing_name)
       if connection_from_unexpected_geo(connections):
           violations.append("unexpected_geolocation")

       # Check: device publishing at unusually high rate (possible compromise/bot)
       msg_rate = get_message_rate(thing_name, window_minutes=5)
       if msg_rate > EXPECTED_MAX_RATE[get_device_type(thing_name)]:
           violations.append("abnormal_message_rate")

       # Check: device last seen timestamp (offline for unexpected duration)
       last_seen = get_last_connection_time(thing_name)
       if (datetime.utcnow() - last_seen).days > 7:
           violations.append("unexpected_offline")

       if violations:
           alert_security_team(thing_name, violations)

       return {"thing_name": thing_name, "violations": violations}
   ```

## Rules

- Device private keys must be generated on the device (or in an HSM) and must never be transmitted in plaintext — CSR-based provisioning is preferred over server-side key generation.
- Certificate revocation must be checked by the broker on every connection — OCSP or CRL lookups must be current (< 24 hours old).
- Decommissioned devices must be wiped if physically recovered — factory reset command or physical destruction prevents credential reuse.
- Each device must have its own certificate — shared certificates mean revoking one device revokes all devices sharing that certificate.

## Common Mistakes

- **Reusing certificates across devices** — a common shortcut for cost reduction; makes individual device revocation impossible and allows lateral movement between devices.
- **Not testing certificate rotation** — certificate rotation is complex; test that devices reconnect successfully after rotation before deploying to production fleet.
- **Provisioning with production credentials in factory** — use a bootstrap certificate with limited permissions for initial provisioning; swap to a full device certificate after device identity is established.
- **No certificate expiry monitoring** — certificates silently expire, causing devices to go offline; monitor expiry with 90-day, 30-day, and 7-day alerts.
