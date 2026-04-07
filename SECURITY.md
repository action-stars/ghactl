# Security Policy

## Supported Versions

Only the latest release receives security fixes.

| Version | Supported          |
| ------- | ------------------ |
| Latest  | :white_check_mark: |
| Older   | :x:                |

## Reporting a Vulnerability

Report vulnerabilities through [GitHub Security Advisories](https://github.com/action-stars/ghactl/security/advisories/new). Do not open a public issue.

### What to include

- Description of the vulnerability and its impact
- Steps to reproduce or a proof of concept
- Suggested fixes, if you have any

### What to expect

- **Acknowledgement** within 3 business days.
- **Assessment** of severity. We may ask follow-up questions.
- **Fix and disclosure** coordinated with you before a public release.

## Scope

Security issues include:

- Remote code execution
- Privilege escalation
- Path traversal or file overwrite
- Credential or secret leakage
- Dependency vulnerabilities with a known exploit

Not in scope:

- Local denial of service through resource exhaustion
- Issues requiring physical access to the runner
- Dependency vulnerabilities without an exploit path in `ghactl`
