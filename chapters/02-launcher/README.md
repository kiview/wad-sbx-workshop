# 2. Make starting a task repeatable

Running Claude by hand lets us see what happens. Now we want to start another task
without repeating the setup. We will first put sandbox settings in an `sbxenv.yaml`
recipe and run it with SBX. Then we will add a task tracker, Beans, for our demo tasks and a small
host launcher that applies the recipe, supplies a task and delivers a private copy
of the app. You will see the SBX command before using the wrapper.

## 1. Run the smallest SBX recipe

An [environment file](https://docs.docker.com/ai/sandboxes/configuration/environment-files/)
records the choices you would otherwise type into `sbx run` or `sbx create`.
From your workshop repository, create a practice directory:

```bash
# HOST
mkdir -p chapters/my-first-env
```

Create `chapters/my-first-env/sbxenv.yaml` with this content:

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
sbx env plan chapters/my-first-env/sbxenv.yaml
sbx env create chapters/my-first-env/sbxenv.yaml --auto-approve
sbx run --name wad-env-first
```

The first line shows the proposed environment. The second creates it; `--auto-approve`
applies the plan without another confirmation. The third opens Claude in it.
This is the recipe equivalent of starting an agent directly in chapter 1.

Inside Claude, try `!docker version`. Then exit Claude and stop this practice sandbox:

```bash
# HOST
sbx stop wad-env-first
```

This recipe has no workspace mount. That is deliberate for the factory: each job
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

## 3. Build the app's environment recipe

We will keep this exercise's sandbox configuration in `chapters/my-02/`. This
lets you change the settings while keeping the supplied example available to
compare with. Create the directory and copy the starting files:

```bash
# HOST — from the workshop repository
mkdir -p chapters/my-02
for file in sbxenv.yaml chapter.env PROMPT.md launch; do
  cat "chapters/02-launcher/$file" > "chapters/my-02/$file"
done
chmod +x chapters/my-02/launch
```

`mkdir` creates your chapter directory. The loop copies the four reference files
into it while keeping your terminal at the workshop root. `chmod` makes the launcher runnable.
Open `chapters/my-02/sbxenv.yaml` in your editor. It is SBX's environment definition.
`chapter.env` configures our workshop helper, `PROMPT.md` gives the agent the chapter's
brief, and `launch` is the small entry point that connects those pieces.

Compared with the tiny example, this recipe adds **arguments** and a **port mapping**.
Arguments let the same file create differently named sandboxes. The port mapping
connects a host browser port to the app's port 8080 inside SBX.

Preview the full recipe before asking our helper to create it:

```bash
# HOST — from the workshop repository
sbx env plan chapters/my-02/sbxenv.yaml \
  --env-arg name=wad-ch-02 \
  --env-arg port=3102 \
  --env-arg "control_dir=$(pwd)/.local/chapters"
```

| Argument | Meaning |
|---|---|
| `chapters/my-02/sbxenv.yaml` | Read your chapter recipe. |
| `name=wad-ch-02` | Name the sandbox we will create. |
| `port=3102` | Publish its app on host port 3102. |
| `control_dir=$(pwd)/.local/chapters` | Supply our host workshop-data path; MCP will use it later. This does not mount it. |

Find those values in the plan. Planning has not started an app or occupied the port.
You already created the tiny environment yourself; next we will let the launcher
perform creation and app delivery together.

## 4. Add the small amount of host glue

Open `chapters/support/launch`. Its job is to:

1. Read a Bean and export the committed app source.
2. Create the environment through `sbx env create`, using the recipe's arguments.
3. Send the app and task into SBX and run the supplied app-startup helper there.

For reference, this is the SBX command the script executes for this job. **Read it;
the launch step below runs it for you.** It changes `plan` to `create` and adds
`--auto-approve` to apply the recipe without another confirmation:

```text
sbx env create chapters/my-02/sbxenv.yaml --env-arg name=wad-ch-02 --env-arg port=3102 --env-arg "control_dir=$(pwd)/.local/chapters" --auto-approve
```

That is the host harness: start an environment and give it the source and task.
`chapter.env` picks the task and default port; `launch` is a short entry point to
this shared script. For this chapter, `TASK=wad-101` selects the warm-up.

Now let the launcher create and prepare the task environment:

```bash
# HOST — terminal A, from the workshop repository
./chapters/my-02/launch wad-ch-02 "$(pwd)/sample-app"
```

The first argument names the sandbox. The second selects your committed app from
chapter 1 (or its completed shortcut). The script runs the same SBX creation command with
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
