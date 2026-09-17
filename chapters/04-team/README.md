# 4. Turn assistants into a team

We can run assistants with shared guidance. Now we need them to divide work and
communicate. [Herdr](https://github.com/herdrdev/herdr) manages their terminal
sessions. Small file messages carry assignments and replies; a notification gives
the recipient a turn to read them.

We will see one message travel through the team before trusting it with a feature.

## 1. Add session management

Open `chapters/kits/herdr/spec.yaml`. Like the Pi kit, it pins a tool version,
allows its download hosts and installs the executable. It supplies a session
manager; the role instructions decide what those sessions should do.

Add this entry to the existing `kits` list in `factory/sbxenv.yaml`:

```yaml
  - source: ../chapters/kits/herdr
```

You now combine three capabilities: ACR installs guidance, Pi provides another
assistant, and Herdr manages the team's sessions.

[Docker Agent](https://docs.docker.com/ai/docker-agent/) is another way to define
agent teams, models, instructions and tools. It could also run inside SBX. This
workshop uses Herdr with coding assistants, following the author's own setup;
the environment and access-control concepts apply to either approach.

## 2. Give the assistants different responsibilities

Create `factory/team.tsv` in your editor. Each line has four fields, separated by spaces or tabs:

```text
coordinator	pi	anthropic	claude-sonnet-5
developer	claude	anthropic	claude-sonnet-5
qa	claude	anthropic	claude-sonnet-5
```

Each row means **role, assistant, provider, model**. Use the working Pi model from
chapter 03. The coordinator routes work and questions; the developer changes code;
QA checks the proposed change against the requirements.

If you only have Claude subscription access, use this content instead:

```text
coordinator	claude	anthropic	claude-sonnet-5
developer	claude	anthropic	claude-sonnet-5
qa	claude	anthropic	claude-sonnet-5
```

The sandbox lets you test these choices before relying on them: change a role's
assistant or model, recreate the environment, and give it the same small handoff
exercise. Kits keep the installation repeatable; the role table keeps the comparison
separate from the communication protocol.

These are separate conversations with separate responsibilities. Using different
providers can add another review perspective; it does not change how roles exchange
work. The [mixed-provider demonstration](../03-pi/MIXED-MODELS.md) explains the
additional credential and assistant configuration for Pi/Gemini and Codex.

Read the three short role briefs in `chapters/support/roles/`. Notice that only the
developer writes application code, and QA reviews a specific commit. A model name
alone does not establish those responsibilities. Keep Claude as the developer for
now; it will use the MCP gateway in the next chapter.

Keep `MODE=manual`, `SESSION=shell` and `USE_ACR=1` in `factory/chapter.env`.
In the SANDBOX tab, at the workshop root:

```bash
./scripts/launch-factory.sh wad-ch-04
```

## 3. Start the sessions and understand what that means

In the SANDBOX shell:

```bash
start-team
```

This supplied helper reads `team.tsv` and performs three Herdr operations:

| Operation | Effect |
|---|---|
| `herdr server` | Starts the session service inside this sandbox. |
| `herdr agent start` | Launches the assistant selected for a role. |
| `herdr agent prompt` | Gives that assistant its role brief. |

Open `chapters/support/bin/start-team` on the host if you want the complete working
example. Workspace IDs, environment variables and terminal creation are plumbing
around these three operations.

Inside SBX, inspect the sessions once:

```bash
herdr agent list
crew logs coordinator
```

`list` should show coordinator, developer and QA. `crew logs` shows an assistant's
terminal, useful for understanding its introduction or diagnosing a login problem.
The `crew` helper is our small human interface to this team's messages; it is not
an SBX command. `crew help` lists its operations.

## 4. See a file message and its wakeup

We will ask a small question without changing the app. In the SANDBOX shell:

```bash
crew ask "Ask developer to list the test scripts in package.json, ask QA to check that list, and send me the combined answer. Do not change the application."
crew watch
```

`ask` waits for role acknowledgments, stores your message, and notifies the
coordinator. `watch` displays readable messages as the team passes work between
roles. Expect a developer answer, a QA check and a combined reply to you. Press
Ctrl-C to return to the shell; the agents continue running.

If the conversation stays quiet, we want to distinguish an assistant still
working from one waiting for login or failing a model request. Press Ctrl-C to
leave the message view, then inspect the coordinator in the SANDBOX shell:

```bash
crew status
crew logs coordinator
```

`status` shows the task stage and messages addressed to you. `logs` shows the
assistant's current terminal output; it does not open an interactive conversation
with that assistant. If the coordinator has handed work to the developer, inspect
`crew logs developer` next; use `crew logs qa` for the review. Look for an active
request, an authentication prompt or an API error. If you need help, show the
instructor that output and your role table. Once the role can continue, use
`crew watch` to follow replies. A quiet message view alone does not tell us whether
an assistant has stopped.

What happened underneath? These are the two operations the helper combines,
shown as a **reference**, not another request to send:

```text
handoff send --to coordinator --from human --kind question --body "Your request"
crew-notify coordinator
```

A file is durable: the next role can read the exact request. But writing a file
alone does not give a running assistant a new turn. Herdr's notification does that.
This is also why explicit messages matter when a harness's busy/idle status is
unreliable. You can inspect `~/work/factory/messages/` inside SBX to see the stored
conversation. Agents use `crew send ROLE "message"` to combine storage and notification too.
Future exercises use `crew` so neither you nor a model needs to remember a separate
wakeup after every message.

## 5. Make startup repeatable

In HOST, set `MODE=team` in `factory/chapter.env`. A future launch now runs
`start-team` for you. Starting sessions and assigning a task remain separate:
`crew submit` will give them the selected task when you are ready.

Type `exit` in the SANDBOX shell. In HOST:

```bash
sbx env rm factory/sbxenv.yaml --env-arg name=wad-ch-04
```

We have a team that can discuss work. Next we give it a controlled connection to
the host backlog so it can fetch requirements and write a result note.

Next: [connect host tools through MCP](../05-mcp/README.md).
