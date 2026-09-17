# 7. Give the factory another job

We have connected repeatable environments, coding guidance, agent roles, backlog
tools and a way to ask a human for help. Now we will reuse that setup for another
job: adding a service filter to the board. You will follow the team, try its change
in the browser and bring the source home. This is open-lab time; if your earlier
job is still running, follow that conversation and try its result first.

## 1. Look at what you are taking away

Think back through the things you configured:

| You wanted… | You added… |
|---|---|
| An agent able to run commands and containers | SBX |
| The same environment for another task | sbxenv and a small launcher |
| Shared coding guidance | ACR kit and a policy package |
| A choice of assistant and model | Pi and explicit role configuration |
| Agents to pass work to each other | Herdr sessions and file handoffs |
| Access to a host task without a host-directory mount | Beans MCP through the gateway |
| A way to answer a missing requirement | SSH and a recorded decision |

Your task note, changed app and saved source show what this particular run produced.
Your recipes and role instructions are the reusable part.

## 2. Change the job, keep the environment

On the host:

```bash
mkdir -p "$WORKSHOP/chapters/my-07"
cd "$WORKSHOP/chapters/my-07"
cat ../my-06/sbxenv.yaml > sbxenv.yaml
cat ../my-06/team.tsv > team.tsv
for file in chapter.env PROMPT.md launch; do cat "../07-factory/$file" > "$file"; done
chmod +x launch
```

The recipe and team are copied unchanged. Open `chapter.env`: `TASK=wad-104` now
asks for a service filter, and the browser port is 3108. There is no new kit to learn.

Continue from the chapter-06 application you saved:

```bash
# HOST — terminal A, in chapters/my-07
./launch wad-ch-07 "$WORKSHOP/.local/reopen-app"
```

The arguments choose a fresh sandbox name and your saved app, including the reopen
behavior you just decided. If you stopped after chapter 05, use `"$WORKSHOP/.local/feature-app"`
instead. If you did not save either result, run `./launch wad-ch-07` instead: it
starts from the supplied assignment/resolution checkpoint. Choose the starting
point that actually exists from your exercise.

## 3. Hand over the next task

In your second host terminal, enter the sandbox:

```bash
sbx exec -it wad-ch-07 bash
```

Inside, look at the configured job and submit it after the roles are ready:

```bash
# SANDBOX shell
cd ~/work
cat task-id
./bin/assign
```

The task ID should be `wad-104`. `assign` uses the same handoff as before. Follow the
team at your own pace with `herdr agent read coordinator --lines 40`, then read the
developer and QA conversations. Notice what changes and what stays the same when
the task changes: the coordinator still routes, the developer still implements,
and QA still reviews.

When they finish, open <http://127.0.0.1:3108>. Try filtering by `checkout-api`. Does
the result count follow the visible incidents? Do combined filters behave as you
would expect? Use the app before reading the agents' conclusion.

Exit the sandbox shell and read the host task:

```bash
# HOST
"$CONTROL/bin/beans" --config "$CONTROL/beans/.beans.yml" \
  --beans-path "$CONTROL/beans/.beans" show wad-104
```

The configuration and data arguments point to the same host backlog. Only the task
ID changed. Read the result note and compare it with what you tried in the browser.

## 4. Keep the source and decide what to adapt

Once the team has committed its result, copy it out from your host terminal:

```bash
# HOST
sbx cp wad-ch-07:/home/agent/work/app "$WORKSHOP/.local/service-filter-app"
```

This preserves the app in a new host directory. Keep your `chapters/my-*` recipes
as well. End the launch terminal with Ctrl-C and run `sbx stop wad-ch-07` when done.
Stopping preserves files, but it ends the running processes; it is not a pause of
a live model conversation.

For your own project, start with one change: replace the sample app and its startup
helper, or write a task for a different feature. Decide which guidance belongs in
skills, which tools belong in kits, and which host operations deserve an MCP tool.


Next: [explore what else SBX can do](../08-presenter/README.md), with the presenter.
