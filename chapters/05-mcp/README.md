# 5. Connect to host tasks through MCP, then inspect the boundary

**Goal:** let the team read the live host Bean and append its result. So far the
launcher has copied a task snapshot; the sandbox cannot update that host task.
MCP supplies a tool interface, and the SBX gateway connects the sandbox client to
our narrowly scoped host adapter.

Read [SBX's MCP gateway documentation](https://docs.docker.com/ai/sandboxes/mcp-gateway/).
Our adapter is a native **host stdio process**. It needs no public server, host
Docker engine, Jira/Linear account or OAuth server. The app and its database still
run inside SBX. The host Bean directory is not mounted into the sandbox.

## 1. Install and check the supplied MCP server

```bash
# HOST
cd "$WORKSHOP"
export FACTORY_CONTROL_DIR="$CONTROL"
scripts/install-mcp.sh --from-local
"$CONTROL/bin/beans-mcp" --check \
  --beans-bin "$CONTROL/bin/beans" \
  --beans-config "$CONTROL/beans/.beans.yml" \
  --beans-data "$CONTROL/beans/.beans"
```

Chapter 00 downloaded the native asset into `dist/` for your host platform.
`--from-local` installs that checksum-verified download. If it is missing, rerun
`./scripts/get-materials.sh` from the repository root. No Go build is needed.
The self-check should report a readable workshop backlog and its tasks.

To see the imperative equivalent of an environment declaration, register a
read-only server once:

```bash
# HOST
sbx mcp add wad-beans-intro --command "$CONTROL/bin/beans-mcp" \
  --args="--beans-bin=$CONTROL/bin/beans" \
  --args="--beans-config=$CONTROL/beans/.beans.yml" \
  --args="--beans-data=$CONTROL/beans/.beans" \
  --dir "$CONTROL"
sbx mcp ls
```

This registers a host command; it does not copy that executable into SBX. To attach
such a registration to an existing sandbox, the command is
`sbx mcp load wad-beans-intro --sandbox SANDBOX_NAME`. We will use the environment
file for our actual team so registration and attachment are reproducible together.
The `wad-beans-intro` registration is read-only and can be removed after inspecting
it: `sbx mcp rm wad-beans-intro`.

## 2. Add the real server to the existing recipe

```bash
# HOST
mkdir -p "$WORKSHOP/chapters/my-05"
cd "$WORKSHOP/chapters/my-05"
cat ../my-04/sbxenv.yaml > sbxenv.yaml
# Append just the MCP section; retain the kits you already composed.
sed -n '/^mcp:/,$p' ../05-mcp/sbxenv.yaml >> sbxenv.yaml
for f in chapter.env PROMPT.md launch; do cat "../05-mcp/$f" > "$f"; done
cat ../my-04/team.tsv > team.tsv
chmod +x launch
cat sbxenv.yaml
sbx env plan sbxenv.yaml --env-arg name=wad-ch-05 \
  --env-arg port=3106 --env-arg "control_dir=$CONTROL"
```

Read the added section before approving it. `command` names the host executable;
`args` fixes the Bean binary, config and data directory. The daemon needs absolute
host paths because it does not inherit your shell setup. The name is scoped to this
run. `--enable-presenter-note-tool` is a historical flag name: here it enables the
attendee exercise's append-only notes too. The disposable backlog marker from
chapter 02 is also required; either alone is insufficient.

**Access changes here.** The host subprocess itself runs with host-user permissions;
`--dir` is not a security sandbox. We trust this adapter's implementation to constrain
its tools to the configured backlog. It exposes task reads and append-only notes,
not arbitrary filesystem access, shell commands or task completion. The gateway is
the connection point; adding a gateway does not make every possible host server safe.

## 3. Launch the team and start the real feature

```bash
# HOST — terminal A, in chapters/my-05
WORKSHOP_APP_REPO="$WARMUP" WORKSHOP_APP_REF=HEAD ./launch wad-ch-05
```

In terminal B:

```bash
# HOST
curl -fsS http://127.0.0.1:3106/healthz
sbx exec wad-ch-05 bash -lc 'cat ~/work/task-id; test ! -e ~/work/task.json'
sbx exec wad-ch-05 claude mcp list
sbx exec wad-ch-05 herdr agent list
```

Expected before assignment: the task ID is `wad-102` and no task snapshot exists.
The gateway client is configured. Neither observation yet proves a successful tool
call. Wait for the three agents' initialization acknowledgments, then:

```bash
# HOST
sbx exec wad-ch-05 /home/agent/work/bin/assign
sbx exec wad-ch-05 herdr agent read developer --lines 40
```

Send the assignment once. The Pi coordinator asks the gateway-capable Claude
**developer** to fetch `get_task` and save `~/work/task.json`. The developer implements
assignment/resolution, sends a review to QA, handles feedback, refreshes the app and
calls `add_task_note`. The coordinator can use the bridge through that role; we do
not assume every installed harness is automatically configured as an MCP client.

Watch actual tool calls and inspect the durable communication when needed:

```bash
# HOST
sbx exec -it wad-ch-05 bash
```

```bash
# SANDBOX
export FACTORY_DIR="$HOME/work/factory"
handoff inbox human --json
herdr agent read coordinator --lines 40
herdr agent read qa --lines 40
```

Exit this inspection shell when done. The launch terminal continues holding the
sandbox open. Let the main feature run while discussing the controls below.

## 4. Exercise a real network boundary

An agent may legitimately need an additional download. The observed example is
Playwright's Chromium download from `cdn.playwright.dev`. Browser testing is not
our learning objective; handling the access request is.

From a sandbox inspection shell, try:

```bash
# SANDBOX
cd ~/work/app
npx playwright install chromium
```

If the current policy blocks the CDN, inspect the exact error. If your account
already allows it, say so; do not report a denial you did not see. You can interrupt
the download once the access behavior is clear.

The human can grant access for just this sandbox from a **host** terminal:

```bash
# HOST — deliberate operator action after inspecting the denial
sbx policy allow network --sandbox wad-ch-05 cdn.playwright.dev
```

Retry the download if you want to verify access. A redirect can require another
hostname: inspect that error before considering another scoped rule. Do not solve
it with unrestricted network access. This exercise does not require running browser
tests or finishing the browser download.

Alternatively, declare the known requirement for future environments in a small
mixin and add it to the next recipe's `kits` list:

```yaml
schemaVersion: "2"
kind: mixin
name: workshop-browser-access
permissions:
  network:
    allow: [cdn.playwright.dev]
```

Save this as `chapters/my-browser-access/spec.yaml`, validate it with
`sbx kit validate`, then reference `../my-browser-access` in a **future** recipe.
That is a network declaration, not a promise that all browser dependencies are
installed. Editing a recipe does not change the running team's environment.

Organization AI governance can add centrally managed controls. The presenter may
show a named tool allowed/denied on an enrolled account. That is separate from this
network rule and is not required for the ungoverned attendee path. See
[SBX security documentation](https://docs.docker.com/ai/sandboxes/security/).

## 5. Observe write-back and save the real result

After the team reports that review is complete:

```bash
# HOST
"$CONTROL/bin/beans" --config "$CONTROL/beans/.beans.yml" \
  --beans-path "$CONTROL/beans/.beans" show wad-102
sbx exec wad-ch-05 git -C /home/agent/work/app status --short
sbx exec wad-ch-05 git -C /home/agent/work/app log -1 --oneline
```

Open <http://127.0.0.1:3106> and try the assignment/resolution behavior described in
the Bean. Look for a result note containing the commit, checks and review findings.
The note is the agents' report; it is not independent host verification. If no note
appears, inspect the developer's MCP call and reported error before repeating work.

When changes are committed and the team is idle:

```bash
# HOST — use a new destination if this one already exists
sbx cp wad-ch-05:/home/agent/work/app "$WORKSHOP/.local/feature-app"
```

**Result:** a real task traveled from host Beans through MCP, was implemented and
reviewed inside SBX, and gained a host result note. No host acceptance system or
automatic task closure is required. Preserve the result, end the hold and stop
`wad-ch-05` before launching the next sandbox on a memory-constrained machine.

Next: [a human decision through SSH](../06-human/README.md).
