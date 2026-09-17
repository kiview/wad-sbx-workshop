# 5. Give the team controlled access to host tools

Our agents can work together, but their task is a snapshot. We want them to read
current requirements from Beans and write a result note back to it. The backlog
stays on the host; mounting the application does not expose those files.

[MCP](https://modelcontextprotocol.io/docs/getting-started/intro) lets an assistant
discover and call a server's named tools. [SBX's MCP gateway](https://docs.docker.com/ai/sandboxes/mcp-gateway/)
connects the assistant inside the sandbox to our small host server:

```text
Claude developer → SBX MCP gateway → host Beans MCP server → workshop backlog
```

## 1. Understand the tool boundary

The workshop includes a prebuilt MCP adapter for Beans. It exposes listing tasks,
reading a task and, when enabled, appending a note. It does not expose arbitrary
shell commands or task deletion. We can make the integration useful without giving
the agent our entire host filesystem or connecting a production issue tracker.

Install the downloaded adapter in HOST:

```bash
./scripts/install-mcp.sh --from-local
```

This installs the native executable downloaded during setup. The gateway starts
it as a host process and exchanges messages over stdin/stdout (**stdio**). No host
Docker engine, public web service or separate OAuth login is needed.

Open `scripts/beans-mcp` in your editor. This short entry point selects the same
workshop backlog as `scripts/beans`. By default it is read-only; `--allow-notes`
enables the result-note tool. The underlying adapter also requires the disposable
backlog marker created in chapter 02.

## 2. Declare the connection alongside the sandbox

Add this section to `factory/sbxenv.yaml`:

```yaml
mcp:
  servers:
    - name: "${{ env.args.name }}-beans"
      command: "${{ env.fileDir }}/../scripts/beans-mcp"
      args: [--allow-notes]
```

`command` is the **host** program the gateway will start. `${{ env.fileDir }}` is
where this environment file lives, so the path points to the workshop script.
`args` permits notes on our practice backlog. The name combines the sandbox name
with `-beans`, giving each environment its own registration.

These fields describe both registration and attachment. You can also register and
load servers with `sbx mcp add` and `sbx mcp load`; recording them here makes the
connection part of every new factory environment.

This is a trust decision. The MCP adapter is a host program with your host user's
permissions. We trust its implementation to limit its tools to the selected
backlog. The gateway supplies the connection; the adapter supplies that narrow
interface. A protocol does not automatically make a powerful tool safe.

## 3. Tell the team when to use its new tools

Edit `factory/chapter.env` to contain:

```text
TASK=wad-102
MODE=mcp
USE_ACR=1
SESSION=shell
```

`MODE=mcp` starts the team and sends only a task ID, so the agent must fetch the
current requirements. `SESSION=shell` gives us a place to explore the tools before
submitting work.

Replace `factory/PROMPT.md` with:

```markdown
Ask the developer to retrieve the task named in ~/work/task-id through the Beans
MCP get_task tool and save it to ~/work/task.json. Do not invent its requirements.
The developer should inspect the mounted project's documentation, install its
dependencies and start any required services inside SBX, then implement the task.
Ask QA to review the exact commit against the task and coding policy.
After review, refresh the web app on 0.0.0.0:8080 and leave it running.
Ask the developer to append a Beans result note with the change, commit and actual
check outcomes. Leave the task open. Send the human a summary and how to try it.
```

The environment provides access. This prompt explains **when and why** the team
should use it. Claude is our gateway-connected developer; Pi can ask that role for
the requirements through the communication system we just built.

In HOST, preview the connection:

```bash
sbx env plan factory/sbxenv.yaml --env-arg name=wad-ch-05
```

Look for the host program, note permission and application mount. In SANDBOX,
starting at the workshop root:

```bash
./scripts/launch-factory.sh wad-ch-05
```

## 4. Explore the actual MCP tools

In the SANDBOX shell, open a personal Claude conversation:

```bash
claude
```

This is your exploration session, separate from the team's developer. Type `/mcp`
to inspect the gateway connection, then ask:

> Use the Beans MCP tools to list the workshop tasks and read wad-102. Explain the
> requested feature. Do not implement it or write a note yet.

Watch the tool calls. You should see the assignment-and-resolution requirements
stored in the host backlog. Ask what information is required to resolve an incident.
This proves the complete connection, rather than just a server registration.

Type `/exit` to return to the SANDBOX shell. Now give the job to the team:

```bash
crew submit
crew watch
```

Submission waits for the roles to be ready. Expect the developer to fetch the task,
implement it, and pass it through QA. If a role needs attention, `crew logs developer`
(or `qa` or `coordinator`) shows its terminal. Ctrl-C leaves the message view;
it does not stop the agents. `crew status` shows the stage and replies addressed to you.

## 5. Let the human handle a request for more access

After the team finishes its feature, try a normal development need: downloading a
browser for future UI checks. In the SANDBOX shell:

```bash
crew ask "Ask the developer to try installing Playwright Chromium for this project. If a download is blocked, tell me the exact destination and why it is needed. Do not change host policy."
crew watch
```

Read the reported failure. A common blocked destination is `cdn.playwright.dev`.
If that is the destination your agent reports and you decide to allow it, use HOST:

```bash
sbx policy allow network --sandbox wad-ch-05 cdn.playwright.dev
```

This grants network access only for this sandbox and host. Return to the SANDBOX
shell (Ctrl-C leaves `crew watch`) and ask it to retry:

```bash
crew ask "The reported download host is now allowed. Ask the developer to retry and report the outcome."
crew watch
```

Redirects may introduce another destination; inspect the actual request before
allowing it. If the original download already succeeds, the current policy permits
it—there is no failure to fix. The lesson is that the agent reports a need and the
host operator decides whether to grant it.

For a requirement you want every future worker to have, create
`factory/browser-access/spec.yaml` in your editor:

```yaml
schemaVersion: "2"
kind: mixin
name: workshop-browser-access
permissions:
  network:
    allow: [cdn.playwright.dev]
```

Use the destinations you actually needed. In HOST:

```bash
sbx kit validate factory/browser-access
```

Add this entry to `factory/sbxenv.yaml`'s existing `kits` list:

```yaml
  - source: ./browser-access
```

The kit records the decision for the next sandbox creation; it does not modify the
current one. On the next chapter's plan, look for this extra kit. Central AI
governance can apply additional organization rules; the presenter will demonstrate
that separately. [Network policy and tool governance](https://docs.docker.com/ai/sandboxes/security/)
are distinct controls.

## 6. See the factory's result

Open **<http://127.0.0.1:3102>**. Try assigning and resolving an incident. In HOST:

```bash
./scripts/beans show wad-102
git -C sample-app log -1 --oneline
```

The task should contain the team's result note. The commit and changed code are
already in `sample-app/`. Compare the actual app behavior with the requirements;
a note saying “done” is not a substitute for trying the result.

When the team is finished, exit the SANDBOX shell. In HOST:

```bash
sbx env rm factory/sbxenv.yaml --env-arg name=wad-ch-05
sbx mcp rm wad-ch-05-beans
```

The first command removes the sandbox. The second removes its named host MCP
registration; this SBX release keeps registrations after sandbox removal. Keep the application as it
is; the next task builds on your assignment-and-resolution feature.

Next: [answer a human question through SSH](../06-human/README.md).
