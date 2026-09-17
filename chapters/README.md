# From one agent to a software factory

Start with an assistant you can let act. Finish with a team you can give work.
A sample application supplies the UI, API and PostgreSQL workload as we build the factory.
In chapter 07, you will bring another project and your own task.
You will assemble the environment and tools that let the team do that work.

Start with [chapter 00](00-setup/README.md), then follow each chapter's README
from the top. Each starts from a problem in the previous setup, introduces a capability, and
lets you use it before making it automatic. Read the explanation around each
command and take time to explore the result. Host and sandbox commands are labelled.

The third-party tools and workflow are the author's choices for this workshop,
not Docker endorsements. See [about the tool choices](../README.md#about-the-tool-choices).

[Presenter-only cloud and mounts](08-presenter/README.md)

| Directory | The new capability | What you can point to afterward |
|---|---|---|
| [00-setup](00-setup) | Install and authenticate | Standalone SBX, tools, dedicated host backlog |
| [01-agent](01-agent) | Let one agent act in SBX | Isolated files, Docker inside SBX, board in your browser |
| [02-launcher](02-launcher) | Launch work from Beans + sbxenv | Task snapshot and app delivered by a small launcher |
| [02.5-acr](02.5-acr) | First kit: distribute guidance | Four coding rules and a review skill materialized by ACR |
| [03-pi](03-pi) | Add another assistant with a kit | Pi reads the same policy; role/provider choices are explicit |
| [04-team](04-team) | Add Herdr with another kit | Coordinator → developer → QA, using files plus wakeups |
| [05-mcp](05-mcp) | Bridge to the host through MCP | Team reads a demo task from the host and appends its result note |
| [06-human](06-human) | Intervene through SSH | A recorded human decision resumes the same team |
| [07-factory](07-factory) | Use the factory on your own project | Your repository, task and project instructions driving the same team |

```mermaid
flowchart LR
  subgraph Host
    B[Beans backlog] --- M[Small Beans MCP server]
    L[Small launcher] --> S[SBX environment]
  end
  M <--> G[SBX MCP gateway]
  subgraph Sandbox
    S --> H[Herdr sessions]
    H --> P[Pi coordinator]
    P <-->|files + wakeups| D[Developer]
    P <-->|files + wakeups| Q[QA]
    A[ACR rules + skill] --> D
    A --> Q
    D --> App[Board + Postgres containers]
  end
  G <--> D
  Human[Human via SSH] --> P
```

## Run a completed reference chapter (catch-up)

The main walkthrough uses `./scripts/launch-factory.sh` with the configuration
you build in `factory/`. If you want to try a chapter's completed setup instead,
its reference launcher supplies that setup and an app checkpoint. For example:

```bash
# HOST — from the workshop repository
./chapters/03-pi/launch wad-ch-03
```

Stop the previous sandbox first to free port 3102. Leave this terminal open, then
join the sandbox at [chapter 03, step 3](03-pi/README.md#3-enter-the-sandbox-and-start-pi-yourself).
This runs the supplied configuration; it does not update the files in your own
`factory/` directory. The task tracker remains on the host.

Chapter 1 mounts `sample-app/`: its changes are already on the host. Later factory
chapters use isolated source snapshots. To carry one of those results forward,
retrieve the sandbox's app with
`sbx cp NAME:/home/agent/work/app ./saved-app`, then launch the next chapter with
`./chapters/04-team/launch new-name ./saved-app`.
This transfers committed source only. Commit inside SBX before copying. Don't execute
agent-generated application code on the host.

## The pieces you will use

`chapters/support/launch` starts the sandbox and supplies its task and source code.
`chapters/support/bin` contains the helpers for starting the app, managing Herdr sessions
and passing messages. `chapters/support/roles` describes the coordinator, developer and QA
responsibilities. You will inspect these pieces as you introduce them, then reuse
them for the next task.

## Keep one SBX session open

SBX v0.45.0-rc2 auto-stops a sandbox 30 seconds after its last client session
disconnects, even if background processes are running inside. `launch` therefore
ends with one foreground `sbx exec NAME sleep infinity`. Leave that host terminal
open; use another terminal for the exercises. Ctrl-C releases the session hold. A stopped sandbox preserves files but loses live
agent processes. Retrieve work before stopping; use a fresh sandbox name when repeating an exercise.
