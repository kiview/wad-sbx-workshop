# Third-party notices — beans-mcp

`beans-mcp` is a statically linked Go binary. The modules below are compiled into it.
Run `go list -m all` in `mcp/beans` to regenerate this list, and
`go mod download -x <module>` to fetch a module's own licence text.

| Module | Version | Licence |
|---|---|---|
| github.com/modelcontextprotocol/go-sdk | v1.8.0 | MIT |
| github.com/google/jsonschema-go | v0.4.3 | BSD-3-Clause |
| github.com/segmentio/encoding | v0.5.4 | MIT |
| github.com/segmentio/asm | v1.1.3 | MIT |
| github.com/yosida95/uritemplate/v3 | v3.0.2 | BSD-3-Clause |
| github.com/golang-jwt/jwt/v5 | v5.3.1 | MIT |
| golang.org/x/oauth2 | v0.35.0 | BSD-3-Clause |
| golang.org/x/sync | v0.20.0 | BSD-3-Clause |
| golang.org/x/sys | v0.41.0 | BSD-3-Clause |
| golang.org/x/time | v0.15.0 | BSD-3-Clause |
| cloud.google.com/go/compute/metadata | v0.3.0 | Apache-2.0 |
| Go standard library | go1.26.0 | BSD-3-Clause |

## Separately pinned host dependency, not redistributed

The Beans CLI (`github.com/hmans/beans`, v0.4.2) is **not** bundled in this binary or in
any release asset of this repository. `scripts/install-beans.sh` downloads the upstream
release and verifies its published checksum. Redistribution was deliberately not taken
on, so its licence and install behaviour remain the upstream project's.
