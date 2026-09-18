# Beans MCP adapter

This small Go program connects the SBX MCP gateway to the workshop's host backlog.
It runs the pinned Beans CLI as a subprocess, with an explicit configuration and
data directory. No Go installation, host Docker engine or HTTP server is needed
to use the prebuilt binary.

```text
Agent in SBX → SBX MCP gateway → beans-mcp → Beans → workshop task files
```

The adapter provides `get_task` and `list_tasks`. The workshop also enables
`add_task_note`, which requires a backlog marked with `.workshop-disposable`.
It does not provide arbitrary commands, paths, task deletion or task completion.

## Install or upgrade

From the workshop repository, run:

```bash
./scripts/install-beans.sh
./scripts/install-mcp.sh
```

The installer downloads the pinned adapter release and verifies its SHA-256
checksum. Chapter 00 also downloads it through `get-materials.sh`; chapter 05
uses `install-mcp.sh --from-local` to install that verified download. Restart any
sandbox already using the adapter after an upgrade so its gateway starts the new
executable.

Release binaries cover Windows x64, macOS ARM64/x64 and Linux ARM64/x64. Windows
uses `beans.exe` and `beans-mcp.exe`; the workshop wrapper runs through Git Bash.
Adapter support is separate from SBX host support: an Intel macOS adapter build
does not make Intel Macs supported SBX hosts.

Version 0.1.1 fixes the Unix execute-bit check that rejected Windows executables.
The adapter still checks Unix permissions on macOS/Linux, and actually invokes
`beans version` on every platform before starting the MCP server. Windows child
processes receive Windows home/temp variables without inheriting model credentials
or `BEANS_PATH` from the host.

## Run directly

The workshop's `scripts/beans-mcp` selects its backlog for you. For another host
integration, launch the native binary with three explicit paths:

```bash
beans-mcp --beans-bin /absolute/path/to/beans \
  --beans-config /absolute/path/to/.beans.yml \
  --beans-data /absolute/path/to/.beans
```

Use Windows paths and `.exe` names when invoking the binary from PowerShell.
Each path is a separate argument, so paths containing spaces are supported by
this adapter. The sample application's separate path limitation is unchanged.

`--enable-presenter-note-tool` enables result notes on a disposable backlog.
`--check` reports configuration and a backlog read as JSON; `--version` reports
build provenance. MCP uses stdin/stdout, while diagnostics go to stderr.

## Maintain and test

The source and its existing protocol tests were brought here from the original
workshop implementation. This component retains its MIT license in `LICENSE`;
the surrounding workshop uses Apache-2.0. Dependencies are listed in
`THIRD-PARTY-NOTICES.md` and pinned in `go.mod`/`go.sum`.

Maintainers need Go 1.26 or later. From the workshop repository:

```bash
./scripts/install-beans.sh
./scripts/build-mcp.sh
cd mcp/beans
go test -count=1 ./...
go vet ./...
```

Tests use a real Beans executable and MCP client/server processes. They cover
reads, filtering, notes, invalid IDs, isolation from an unrelated backlog,
startup failures, paths containing spaces, Unicode note limits, revoking note
access, simultaneous notes, and shutdown. CI requires Beans so
missing dependencies cannot silently skip protocol coverage. CI also exercises
the installed `scripts/beans-mcp` entrypoint, including Git Bash on Windows.

The `Beans MCP` workflow runs natively on all five release OS/architecture pairs,
then cross-compiles static binaries and packages checksums and provenance. These
checks validate the adapter and host wrapper, not SBX virtualization or agent login.

This adapter targets small workshop backlogs. Each Beans call has a 10-second
timeout and a 1 MiB output cap. A larger result produces `backlog_output_too_large`,
even if the caller requests a short list; `list_tasks` limits the returned task
count after reading the CLI response.

For a release, update `BEANS_MCP_VERSION` in `scripts/versions.env`, run the checks,
and push a matching `beans-mcp-vVERSION` tag at the tested source commit. The tag
workflow repeats the native checks and publishes only when all succeed. The
application bundle remains pinned separately to `materials-v0.1.0`.
