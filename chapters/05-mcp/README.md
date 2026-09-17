# 5. Give the team a tool for the host backlog

The team can work from a task snapshot. Now we want it to read the current task and
write its result back to Beans on the host. We will connect a small Beans MCP server
through [SBX's MCP gateway](https://docs.docker.com/ai/sandboxes/mcp-gateway/),
inspect its tools with Claude, then give the team an assignment. Because this
crosses the sandbox boundary, we will also examine which operations we allow and
how to change access deliberately.

The connection will be:

```text
Claude inside SBX → SBX MCP gateway → Beans MCP server on host → host backlog
```

## 1. Understand the server we are giving access to

The server is a prebuilt native program downloaded in chapter 00. Install it now:

```bash
# HOST — from the workshop repository
./scripts/install-mcp.sh --from-local
```

This checks and installs the downloaded executable into our workshop's local tools
directory. It does not start a public web server. The gateway will run the program
on the host and communicate with it over stdin/stdout, commonly called **stdio**.
No host Docker engine or separate OAuth service is involved.

There are three useful operations: list tasks, read one task, and append a task
note. A server chooses which operations it exposes; MCP is their interface, not a
reason to give an agent arbitrary host access.

## 2. Register a server directly with SBX

First see how a host command becomes an MCP server registration:

```bash
# HOST
sbx mcp add wad-beans-intro --command "$CONTROL/bin/beans-mcp" \
  --args="--beans-bin=$CONTROL/bin/beans" \
  --args="--beans-config=$CONTROL/beans/.beans.yml" \
  --args="--beans-data=$CONTROL/beans/.beans" \
  --dir "$CONTROL"
```

| Part | What it tells the gateway |
|---|---|
| `wad-beans-intro` | The registration's name. |
| `--command` | The native host program to start. |
| `--beans-bin` | Which Beans executable the server should use. |
| `--beans-config`, `--beans-data` | Which workshop backlog it should read. |
| `--dir` | The host process's working directory. This is not a filesystem restriction. |

Each `--args` passes an argument to the server. Absolute paths matter because the
gateway does not run inside your current interactive shell.

```bash
sbx mcp ls
```

Find your named registration. This step registered a **read-only** tool source; it
did not mount host files into a sandbox. To attach a registered server to an existing
sandbox you would use `sbx mcp load NAME --sandbox SANDBOX`. For our factory, we will
record registration and attachment together in the recipe instead.

Remove this practice registration before moving on:

```bash
sbx mcp rm wad-beans-intro
```

## 3. Put the connection in the recipe

Prepare the chapter directory on the host:

```bash
mkdir -p "$WORKSHOP/chapters/my-05"
cd "$WORKSHOP/chapters/my-05"
cat ../my-04/sbxenv.yaml > sbxenv.yaml
cat ../my-04/team.tsv > team.tsv
for file in chapter.env PROMPT.md launch; do cat "../05-mcp/$file" > "$file"; done
chmod +x launch
```

You are keeping your kits and model choices. The supplied chapter settings select
`wad-102`, the assignment/resolution feature, and `MODE=mcp`. That mode delivers
only a task ID: the agents must obtain the task through the new tool.

Open the supplied `chapters/05-mcp/sbxenv.yaml`. Copy its **`mcp:` section** to the
end of your `chapters/my-05/sbxenv.yaml`. Read it before saving:

- `command` is the same host program you registered manually.
- `args` selects the same backlog.
- The server name uses this sandbox's name so separate runs have separate registrations.
- `--enable-presenter-note-tool` enables appending notes. The name is historical;
  we use it for the learner exercise too.

### Decide what the agent should be allowed to change

This is the point where we give the sandbox authority over something on the host.
The adapter limits its tools to the configured backlog. A write requires both the
note-tool flag and the disposable marker created in chapter 02. It cannot close a
task or run an arbitrary command through its tool interface.

The adapter itself is a host process with host-user permissions, so its implementation
is part of what we trust. The gateway connects it; the gateway does not automatically
make a broadly privileged server narrow. For this exercise we are choosing a small,
inspectable set of operations on practice data.

Preview the updated environment and then launch it:

```bash
# HOST — terminal A, in chapters/my-05
sbx env plan sbxenv.yaml --env-arg name=wad-ch-05 \
  --env-arg port=3106 --env-arg "control_dir=$CONTROL"
./launch wad-ch-05 "$WARMUP"
```

Look for the MCP server in the plan. The launcher then calls `sbx env create`,
prepares the app and starts the team, as before. Leave this terminal open.

## 4. Explore the tools in Claude before assigning work

In a second host terminal:

```bash
sbx run --name wad-ch-05
```

This opens a conversation for **you** to explore the environment. It is separate
from the waiting Herdr developer session. Type `/mcp` to inspect Claude's MCP
connections, then return to the conversation and ask:

> Use the Beans MCP tools to list the workshop tasks and read wad-102. Explain the
> requested feature to me. Do not implement it or change the backlog yet.

Watch which tools Claude calls. You should be discussing the host's actual task,
not a task invented from the application source. Read the requirements together:
what must assignment do, and what information is required to resolve an incident?

Exit this exploration conversation when you understand the task. We will now give
it to the team we built in chapter 04.

## 5. Let the team use the same bridge

In terminal B, after exiting Claude to the host, enter a shell in the sandbox:

```bash
# HOST
sbx exec -it wad-ch-05 bash
```

Inside that shell:

```bash
# SANDBOX
cd ~/work
export FACTORY_DIR="$HOME/work/factory"
cat task-id
cat bin/assign
```

`task-id` identifies the job. `assign` is the two-step handoff from chapter 04,
wrapped around this chapter's task: leave a message for the coordinator, then wake
its session. Once all roles have finished their introductions, run it once:

```bash
./bin/assign
```

The coordinator routes the job. The Claude developer can reach the gateway, so it
reads the Bean and supplies the task to the team. Pi does not need its own MCP
client in this setup: it can ask the gateway-capable role to obtain the information.

Follow the developer and reviewer from the same shell:

```bash
herdr agent read developer --lines 40
herdr agent read qa --lines 40
```

These commands show their conversations. Read the implementation plan and review
feedback. The agents may need several turns; there is no need to keep sending the
assignment. When the work is reviewed, the developer refreshes the app and appends
a result note through MCP.

## 6. Handle a legitimate request for more access

While the feature runs, consider another normal development need: downloading a
browser. Our observed example is Chromium's download from `cdn.playwright.dev`.
From your sandbox shell:

```bash
cd ~/work/app
npx playwright install chromium
```

`npx` runs the app's Playwright tool; `install chromium` asks it to fetch a browser.
We are interested in the access boundary, not in running browser tests. Read the
error if the download is blocked. If your policy already permits it, simply observe
that difference; you can interrupt the download once you have seen the behavior.

The agent cannot fix a host access rule from this shell. Keep the sandbox shell
open in terminal B. Open **terminal C on the host**, enter the workshop repository,
and decide whether this sandbox should reach that particular download host:

```bash
# HOST — terminal C, from the workshop repository
source ./scripts/workshop-env.sh
sbx policy allow network --sandbox wad-ch-05 cdn.playwright.dev
```

`allow network` grants network access; `--sandbox` limits the rule to this one
sandbox; the final argument names the destination. Return to the sandbox shell in
terminal B and retry the download if you want to see the effect. A redirect may name another destination; inspect
that request before adding another rule.

For a dependency every future worker needs, a kit can record the requirement:

```yaml
schemaVersion: "2"
kind: mixin
name: workshop-browser-access
permissions:
  network:
    allow: [cdn.playwright.dev]
```

Save that as `chapters/my-browser-access/spec.yaml` on the host if you want to keep
it, validate it with `sbx kit validate`, and add `../my-browser-access` to a future
recipe's kits. You have moved a manual operator decision into a declared environment
requirement. Editing a recipe does not change the running sandbox.

An enrolled organization can add central AI governance. The presenter may show it
at the end. Network access and named MCP-tool governance are different controls;
you do not need organization access to complete this local exercise. See
[SBX security](https://docs.docker.com/ai/sandboxes/security/).

## 7. See the result on both sides

When the team reports completion, open **<http://127.0.0.1:3106>** and try assigning
and resolving an incident. Does the behavior match the task you read earlier?

In host terminal C, where you sourced the workshop variables, read the same Bean again:

```bash
"$CONTROL/bin/beans" --config "$CONTROL/beans/.beans.yml" \
  --beans-path "$CONTROL/beans/.beans" show wad-102
```

This is the same task-reading command from chapter 02. This time it should include
the team's result note: what changed, the commit and the review/check outcomes.
The agents have written back through the bridge. Their note remains their report;
we have not added a separate host acceptance service.

Back in terminal B, inspect the working tree and latest commit:

```bash
# SANDBOX — terminal B
cd ~/work/app
git status --short
git log -1 --oneline
```

Wait for the team to finish and commit its work. `git status` shows any remaining
uncommitted changes; `git log` identifies the saved result. Type `exit` to leave
the sandbox shell, then save the app from the host:

```bash
# HOST — terminal B
sbx cp wad-ch-05:/home/agent/work/app "$WORKSHOP/.local/feature-app"
```

`sbx cp` copies from the named sandbox path to a new host directory. Unlike chapter
1's mounted app, this factory job worked on a private source snapshot, so we now
retrieve it. If the destination exists from an earlier run, choose a new name.

End the launcher with Ctrl-C and stop `wad-ch-05` after saving your result. We have
connected task, team and result. Next we handle a question the team cannot answer.

Next: [join the team through SSH](../06-human/README.md).
