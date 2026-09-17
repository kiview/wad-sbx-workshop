# 2. Make starting a task repeatable

Running Claude by hand lets us see what happens. Now we want to start another task
without repeating the setup. We will first put sandbox settings in an `sbxenv.yaml`
environment file and run it with SBX. Then we will add a task tracker, Beans, for our demo tasks and a small
host launcher that applies the environment file, supplies a task and delivers a private copy
of the app. You will see the SBX command before using the wrapper.

## 1. Run the smallest SBX environment file

An [environment file](https://docs.docker.com/ai/sandboxes/configuration/environment-files/)
records the choices you would otherwise type into `sbx run` or `sbx create`.
From your workshop repository, create a practice directory:

```bash
# HOST
mkdir -p factory
```

Create `factory/sbxenv.yaml` with this content:

```yaml
schemaVersion: "1"
name: wad-env-first
agent: claude
sandboxOptions:
  skills: off
  cpus: 4
  memory: 8g
```

| Line | What it does |
|---|---|
| `schemaVersion` | Selects the environment-file format. |
| `name` | Gives this sandbox a name we can use later. |
| `agent` | Chooses the built-in Claude environment. |
| `skills: off` | Keeps this example independent of skills installed on your host. |
| `cpus`, `memory` | Assigns resources to the sandbox. |

Preview it, then create it:

```bash
# HOST
sbx env plan factory/sbxenv.yaml
sbx env create factory/sbxenv.yaml --auto-approve
sbx run --name wad-env-first
```

The first line shows the proposed environment. The second creates it; `--auto-approve`
applies the plan without another confirmation. The third opens Claude in it.
You have now started an agent using an environment file instead of the flags from chapter 1.

Inside Claude, try `!docker version`. Then exit Claude and stop this practice sandbox:

```bash
# HOST
sbx stop wad-env-first
```

This environment file has no workspace mount. That is deliberate for the factory: each job
will receive its own copy of the app. Your chapter-1 edits remain in the host's
`sample-app/` directory.

## 2. Give the host a task list

We use [Beans](https://github.com/hmans/beans), a small file-based task tracker.
The task is a written contract an agent can read, rather than a prompt we keep retyping.

```bash
# HOST — from the workshop repository
./scripts/install-beans.sh
./scripts/backlog-init.sh --disposable
```

The first script downloads the pinned Beans executable. The second creates our
practice backlog with four tasks. `--disposable` identifies it as workshop data;
later we will let the agents append result notes through MCP.

Read the demo task from the first coding exercise:

```bash
# HOST
"$(pwd)/.local/chapters/bin/beans" --config "$(pwd)/.local/chapters/beans/.beans.yml" \
  --beans-path "$(pwd)/.local/chapters/beans/.beans" show wad-101
```

`--config` selects our Beans configuration; `--beans-path` selects its task directory;
`show wad-101` reads one task. It is the same task you completed in chapter 1.

## 3. Build the app's environment file

We want each launch to create a sandbox, make the sample app available in the
browser, and give the agent a specific coding task. Three files describe those
choices:

| File | What it tells the system to do |
|---|---|
| `factory/sbxenv.yaml` | Tells SBX which agent environment to create and which application port to publish. |
| `factory/chapter.env` | Tells our launcher which demo task to load and which startup steps to run. |
| `factory/PROMPT.md` | Gives a coding agent instructions for approaching that task. The team will read this when we introduce task assignment. |

The supplied starting configuration names the sandbox through an argument,
publishes the app on port 3102, and selects the first demo task. Put those settings
in place:

```bash
# HOST — from the workshop repository
cp chapters/02-launcher/sbxenv.yaml factory/sbxenv.yaml
cp chapters/02-launcher/chapter.env factory/chapter.env
cp chapters/02-launcher/PROMPT.md factory/PROMPT.md
```

Open `factory/sbxenv.yaml` and find `args`, `sandboxOptions` and `ports`. These
control sandbox creation; the task and prompt files are read by our workshop
helpers. We will keep extending this same environment file.

Compared with the tiny example, this environment file makes the sandbox name an argument and
publishes the application's port 8080 at `http://127.0.0.1:3102` on your laptop.
The port stays the same throughout the factory exercises. Stop each chapter's
sandbox before starting the next one so that port is available.

Preview the environment file:

```bash
sbx env plan factory/sbxenv.yaml --env-arg name=wad-ch-02
```

`factory/sbxenv.yaml` selects your environment file. `--env-arg name=wad-ch-02`
supplies the name declared in its `args` section. Find the port mapping in the plan;
it comes from the file, so you do not need to pass it on the command line.

## 4. Add the small amount of host glue

Open `chapters/support/launch`. Its job is to:

1. Read a Bean and export the committed app source.
2. Create the environment through `sbx env create`, using the environment file's arguments.
3. Send the app and task into SBX and run the supplied app-startup helper there.

For reference, this is the SBX command the script executes for this job. **Read it;
the launch step below runs it for you.** It changes `plan` to `create` and adds
`--auto-approve` to apply the environment file without another confirmation:

```text
sbx env create factory/sbxenv.yaml --env-arg name=wad-ch-02 --auto-approve
```

That is the host harness: start an environment and give it the source and task.
`factory/chapter.env` picks the task and startup behavior.
`scripts/launch-factory.sh` calls this shared helper with your `factory/` directory.
For this chapter, `TASK=wad-101` selects the warm-up.

Now let the launcher create and prepare the task environment:

```bash
# HOST — terminal A, from the workshop repository
./scripts/launch-factory.sh wad-ch-02
```

The argument names the sandbox. The launcher takes your committed app from
`sample-app/` and reads the configuration in `factory/`. The script runs the same SBX creation command with
`name=wad-ch-02`, then supplies the sample app and demo task and starts the app.
It exports committed files, so ask the chapter-1 agent to commit first if needed.

Leave this terminal open: its final SBX connection keeps the sandbox awake.

## 5. Open the agent and inspect what it received

In a second host terminal:

```bash
# HOST
sbx run --name wad-ch-02
```

This opens Claude in the environment we created. Tell it:

> Work in /home/agent/work/app. Read /home/agent/work/task.json. Explain the task
> and check that our warm-up change is already present. Do not implement it again.

Open **<http://127.0.0.1:3102>**. You should see your app, started by the supplied
setup helper. We now have a repeatable environment plus delivery of a task and code.
A fresh worker still does not know our coding conventions; that is the next chapter.

Exit Claude when done, end terminal A's waiting command with Ctrl-C, then stop:

```bash
# HOST
sbx stop wad-ch-02
```

We only inspected the app here. Keep using the same `sample-app/` source for the next
chapters; no copying is needed.

Next: [kits and ACR](../02.5-acr/README.md).
