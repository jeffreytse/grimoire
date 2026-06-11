---
name: prevent-path-traversal
description: Use when building file upload, file download, file serving, or any feature where user input influences a filesystem path — including filename parameters, directory parameters, or file include logic.
source: 'OWASP Path Traversal Cheat Sheet; OWASP File Upload Cheat Sheet (owasp.org/www-project-cheat-sheets); OWASP Top 10 2021 A01; CWE-22; CWE-434'
tags: [security, owasp, path-traversal, file-upload, directory-traversal, developer, input-validation]
---

# Prevent Path Traversal

Canonicalize and validate all filesystem paths constructed from user input — verifying they resolve within an intended base directory — to block directory traversal attacks that read or overwrite arbitrary server files.

## Why This Is Best Practice

**Adopted by:** OWASP Top 10 2021 A01 (Broken Access Control) includes path traversal. CWE-22 is one of the most common CVE categories — 2,000+ CVEs listed in NVD. Apache HTTP Server (CVE-2021-41773), Pulse Secure (CVE-2019-11510), and Citrix (CVE-2019-19781) all had critical path traversal vulnerabilities. NIST NVD rates path traversal CVEs up to CVSS 10.0 (maximum criticality).
**Impact:** CVE-2021-41773 (Apache HTTP Server path traversal) allowed unauthenticated remote code execution via a crafted URL — patched in hours after public disclosure, but thousands of servers were exploited before patching. Pulse Secure CVE-2019-11510 enabled reading VPN credentials without authentication, affecting governments and enterprises globally. Path traversal on file upload enables webshell upload → full server compromise.
**Why best:** Stripping `../` sequences from user input is the legacy approach — it fails against encoded variants (`%2e%2e%2f`, `%252e`, `..%2f`), null bytes, and OS-specific separators. Canonicalization via `realpath()` (or equivalent) resolves all encoding variants to an absolute path, and the prefix check (`startswith(base_dir)`) is applied to the canonical result.

Sources: OWASP Path Traversal Cheat Sheet; CVE-2021-41773 analysis; CWE-22; OWASP File Upload Cheat Sheet

## Steps

1. **Canonicalize paths and verify they stay within the base directory**:

   ```python
   import os

   BASE_DIR = '/var/app/uploads'

   def safe_path(user_filename):
       # Remove any path components from user input
       filename = os.path.basename(user_filename)
       # Construct the candidate path
       candidate = os.path.join(BASE_DIR, filename)
       # Canonicalize — resolves .., symlinks, encoded chars
       canonical = os.path.realpath(candidate)
       # Verify it's within the base directory
       if not canonical.startswith(os.path.realpath(BASE_DIR) + os.sep):
           raise SecurityError(f"Path traversal detected: {user_filename}")
       return canonical
   ```

2. **Reject filenames with path separators** — strip or reject inputs containing `/`, `\`, `..`, null bytes, or non-printable characters:

   ```python
   import re

   def sanitize_filename(name):
       # Allow only alphanumeric, dash, underscore, dot
       if not re.match(r'^[a-zA-Z0-9._-]+$', name):
           raise ValueError("Invalid filename characters")
       if name.startswith('.'):
           raise ValueError("Hidden files not allowed")
       if len(name) > 255:
           raise ValueError("Filename too long")
       return name
   ```

3. **For file downloads, use indirect references** — store files by UUID or hash, never expose real filenames to users:

   ```python
   import uuid

   # On upload: store with UUID
   file_id = str(uuid.uuid4())
   stored_path = os.path.join(BASE_DIR, file_id)
   shutil.copy(uploaded_file, stored_path)
   db.store(file_id=file_id, original_name=filename, owner=user_id)

   # On download: look up UUID, never use user input as path
   @app.route('/download/<file_id>')
   def download(file_id):
       record = db.get(file_id=file_id, owner=current_user.id)
       if not record:
           abort(404)
       return send_file(os.path.join(BASE_DIR, file_id),
                        download_name=record.original_name)
   ```

4. **Validate file content type on upload** — check MIME type from file magic bytes, not the `Content-Type` header or file extension:

   ```python
   import magic  # python-magic library

   ALLOWED_MIME_TYPES = {'image/jpeg', 'image/png', 'image/gif', 'application/pdf'}

   def validate_file_type(file_bytes):
       detected = magic.from_buffer(file_bytes[:2048], mime=True)
       if detected not in ALLOWED_MIME_TYPES:
           raise ValueError(f"File type not allowed: {detected}")
   ```

5. **Store uploaded files outside the web root** — files in the web root are directly accessible; files outside it must be explicitly served:

   ```
   /var/app/
     web/                 ← web root, served by Nginx
       index.html
     uploads/             ← outside web root, not directly accessible
       uuid1.jpg
   ```

   Never store uploads in `static/`, `public/`, or any directory served by the web server.

6. **Restrict file extensions for uploads** — allowlist permitted extensions, reject all others:

   ```python
   ALLOWED_EXTENSIONS = {'.jpg', '.jpeg', '.png', '.gif', '.pdf'}

   def allowed_extension(filename):
       _, ext = os.path.splitext(filename.lower())
       return ext in ALLOWED_EXTENSIONS
   ```

   Defense-in-depth: rename files on storage (step 3's UUID approach eliminates extension attacks entirely).

## Rules

- Always use `realpath()` (or equivalent) for canonicalization — string manipulation of paths is insufficient.
- Check the canonical path starts with `base_dir + os.sep`, not just `base_dir` — without `os.sep`, `../uploads_malicious` would pass a `startswith('/uploads')` check.
- On Windows, path separators are both `\` and `/`, and drive-letter paths (`C:\`) must be handled — prefer UUID-based storage over path validation on Windows.
- Symlinks inside the base directory can point outside it — `realpath()` resolves them; ensure the resolved target is also within bounds.

## Common Mistakes

- **Stripping `../` with regex** — `....//` after stripping `../` becomes `../`; recursive stripping fails with `%2e%2e/` encoding.
- **Using `os.path.join` without `realpath` check** — `os.path.join('/uploads', '../etc/passwd')` resolves to `/etc/passwd` on Unix.
- **Allowing `.php`, `.py`, `.sh` extensions in uploads** — if the file is served from the web root, it executes as server-side code.
- **Trusting Content-Type header for file type validation** — attackers set it to `image/jpeg` while uploading a PHP webshell.
