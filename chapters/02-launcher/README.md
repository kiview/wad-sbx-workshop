# 2. Make starting a task repeatable

**Goal:** describe the environment in a file, then connect a host task to it.
In chapter 1 you started Claude and explained how to run the app. For the next job,
we want to create the same environment and give the agent a written task.

We will build this in two steps: **an SBX recipe**, then **a small script that
combines that recipe with our app and task tracker**.

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
`$WARMUP` directory.

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

Read the warm-up task:

```bash
# HOST
"$CONTROL/bin/beans" --config "$CONTROL/beans/.beans.yml" \
  --beans-path "$CONTROL/beans/.beans" show wad-101
```

`--config` selects our Beans configuration; `--beans-path` selects its task directory;
`show wad-101` reads one task. It is the same task you completed in chapter 1.

## 3. Build the app's environment recipe

Create your chapter directory and copy the supplied starting files into it:

```bash
# HOST — from the workshop repository
mkdir -p chapters/my-02
cd chapters/my-02
for file in sbxenv.yaml chapter.env PROMPT.md launch; do
  cat "../02-launcher/$file" > "$file"
done
chmod +x launch
```

`mkdir` creates your working directory; `cd` enters it. The loop copies the four
reference files without changing the originals. `chmod` makes the launcher runnable.
Open `sbxenv.yaml` in your editor.

Compared with the tiny example, this recipe adds **arguments** and a **port mapping**.
Arguments let the same file create differently named sandboxes. The port mapping
connects a host browser port to the app's port 8080 inside SBX.

**This is the actual SBX command that creates our workshop environment. Run it:**

```bash
# HOST — in chapters/my-02
sbx env create sbxenv.yaml \
  --env-arg name=wad-recipe-preview \
  --env-arg port=3102 \
  --env-arg "control_dir=$CONTROL" \
  --auto-approve
```

| Argument | Meaning |
|---|---|
| `sbxenv.yaml` | Read the environment recipe in this directory. |
| `name=wad-recipe-preview` | Name the sandbox we are creating. |
| `port=3102` | Publish its app on host port 3102. |
| `control_dir=$CONTROL` | Supply our host workshop-data path; MCP will use it later. This does not mount it. |
| `--auto-approve` | Apply the recipe without another confirmation prompt. |

You have created the environment. The recipe alone has not delivered our app or
Bean, so there is no board to open yet. Remove this empty practice environment
before the launcher creates the real task environment:

```bash
# HOST — confirm removal of this practice sandbox when prompted
sbx rm wad-recipe-preview
```

`sbx env create` creates a new sandbox; it refuses an existing name. This cleanup
also releases port 3102 for the actual run.

## 4. Add the small amount of host glue

Open `../support/launch`. Its job is to:

1. Read a Bean and export the committed app source.
2. Run the **same `sbx env create` command above**, using the recipe's arguments.
3. Send the app and task into SBX and run the supplied app-startup helper there.

That is the host harness. It does not implement tasks or coordinate the agents.
`chapter.env` picks the task and default port; `launch` is a short entry point to
this shared script. For this chapter, `TASK=wad-101` selects the warm-up.

Now let the launcher create and prepare the task environment:

```bash
# HOST — terminal A, in chapters/my-02
./launch wad-ch-02 "$WARMUP"
```

The first argument names the sandbox. The second selects your committed app from
chapter 1 (or its completed shortcut). The script runs the same SBX creation command with
`name=wad-ch-02`, then supplies the app and task and starts the board.
It exports committed files, so ask the chapter-1 agent to commit first if needed.

Leave this terminal open: its final SBX connection keeps the sandbox awake.
The script is waiting, not supervising the task.

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

We only inspected the app here. Keep using the same `$WARMUP` source for the next
chapters; no copying or new environment-variable settings are needed.

Next: [kits and ACR](../02.5-acr/README.md).
