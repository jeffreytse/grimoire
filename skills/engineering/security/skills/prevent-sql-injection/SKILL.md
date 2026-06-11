---
name: prevent-sql-injection
description: Use when writing code that constructs database queries using any data that originates outside the application — user input, API parameters, headers, cookies, or environment variables.
source: 'OWASP SQL Injection Prevention Cheat Sheet (owasp.org/www-project-cheat-sheets); OWASP Top 10 2021 A03; CWE-89; NIST SP 800-53 SI-10'
tags: [security, owasp, sql-injection, database, developer, input-validation]
---

# Prevent SQL Injection

Eliminate SQL injection by never interpolating untrusted data into SQL strings — use parameterized queries or ORMs exclusively.

## Why This Is Best Practice

**Adopted by:** Mandated by PCI DSS v4.0 Requirement 6.2.4, NIST SP 800-53 SI-10, and OWASP Top 10 2021 (A03:Injection). Google, Microsoft, and Amazon all require parameterized queries in secure coding standards. Every major ORM (Hibernate, SQLAlchemy, ActiveRecord, Prisma) defaults to parameterization.
**Impact:** SQL injection is the #1 cause of data breach in web applications — responsible for 65% of database-related breaches (Verizon DBIR 2023). Parameterized queries eliminate the entire vulnerability class, not individual instances. NIST NVD lists thousands of CVEs annually where string concatenation was the root cause.
**Why best:** String escaping (e.g., `mysql_real_escape_string`) is the alternative — it fails under character encoding edge cases, stored procedure contexts, and developer error. Parameterized queries structurally separate code from data, making injection impossible regardless of input content.

Sources: OWASP SQL Injection Prevention Cheat Sheet; Verizon DBIR 2023; CWE-89; PCI DSS v4.0

## Steps

1. **Use parameterized queries (prepared statements) for every query** — pass user data as bound parameters, never concatenate it into the SQL string.

   ```python
   # BAD — injectable
   cursor.execute(f"SELECT * FROM users WHERE id = {user_id}")

   # GOOD — parameterized
   cursor.execute("SELECT * FROM users WHERE id = %s", (user_id,))
   ```

   ```java
   // BAD
   stmt.executeQuery("SELECT * FROM users WHERE name = '" + name + "'");

   // GOOD
   PreparedStatement ps = conn.prepareStatement("SELECT * FROM users WHERE name = ?");
   ps.setString(1, name);
   ```

2. **Use an ORM for standard CRUD** — ORMs parameterize automatically. Only drop to raw SQL when necessary.

   ```python
   # Django ORM — safe by default
   User.objects.filter(username=username)

   # Raw SQL in ORM — still safe with params
   User.objects.raw("SELECT * FROM users WHERE id = %s", [user_id])
   ```

3. **For dynamic query structure (table names, column names, ORDER BY)** — whitelist identifiers explicitly; never accept them from user input. Column and table names cannot be bound parameters.

   ```python
   ALLOWED_SORT_COLUMNS = {"name", "created_at", "email"}
   if sort_col not in ALLOWED_SORT_COLUMNS:
       raise ValueError("Invalid sort column")
   cursor.execute(f"SELECT * FROM users ORDER BY {sort_col}")
   ```

4. **Apply least-privilege database accounts** — the app's DB user should have only SELECT/INSERT/UPDATE/DELETE on specific tables. No DROP, no `xp_cmdshell`, no access to system tables. Defense-in-depth: limits damage if injection does occur.

5. **For stored procedures** — use parameterized calls, not concatenated SQL inside the procedure body. Stored procedures are not inherently safe if they build dynamic SQL internally.

6. **Enable query logging and anomaly detection** — log all queries in staging; alert on patterns like `UNION SELECT`, `--`, `1=1` in production to catch exploitation attempts.

## Rules

- Never use string formatting, concatenation, or f-strings to build SQL with external data.
- `LIKE` clauses need both parameterization AND wildcard escaping: `LIKE ?` with value `%term%` is safe; `LIKE '%{term}%'` is not.
- Second-order injection: data stored in the DB can later be used in queries — treat retrieved data as untrusted if it will be re-queried.

## Common Mistakes

- **Using an ORM for SELECT but raw SQL for batch inserts or reports** — the one unparameterized query is the one that gets exploited.
- **Treating numeric inputs as safe** — `int(user_id)` then string concat is safer but still wrong; use params for type safety and consistency.
- **Assuming stored procedures are safe** — they are only safe if they themselves use parameterized internal queries.
