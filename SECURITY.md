# Security policy

## Reporting a vulnerability

Please report vulnerabilities privately through GitHub:
**[Report a vulnerability](https://github.com/kinglet-dev/core/security/advisories/new)**
(Security tab → "Report a vulnerability").

Please don't open a public issue for security problems.

Include what you found, how to reproduce it (a sample SVG helps), and the
impact you expect.

## What core trusts

Nothing in the SVG it reads. The reader parses one input within limits on
size, nesting depth, element count and path complexity, refuses DOCTYPEs
that define entities, never follows references inside the file, and uses no
network. Inputs over a limit are refused with an error naming the limit.

The machine-readable contact is published at
https://kinglet.dev/.well-known/security.txt (RFC 9116).
