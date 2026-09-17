# 2. Turn a one-off sandbox into a repeatable work environment

You ran an agent, opened its app in a browser and kept its changes in `sample-app/`.
Now we want the next worker to receive the same environment and a written task.
We will describe SBX in a file, try a small task tracker, then connect the two with
a short host launcher. The application stays in the same mounted directory.

## 1. Describe the environment

An [SBX environment file](https://docs.docker.com/ai/sandboxes/configuration/environment-files/)
records the settings you previously typed into `sbx run`. In HOST, create the
configuration directory:

```bash
mkdir -p factory
```

Create `factory/sbxenv.yaml` in your editor:

```yaml
schemaVersion: "1"
name: wad-env-first
agent: claude
workspace: ../sample-app
sandboxOptions:
  skills: off
  cpus: 4
  memory: 8g
ports:
  - sandbox: 8080
    host: 3102
    protocol: tcp4
    hostIP: 127.0.0.1
```

`agent` selects the built-in Claude environment. `workspace` is relative to this
file: `../sample-app` shares our application, including `.git`. The resource and
skills settings are the choices from chapter 01. `ports` publishes the app's port
8080 at localhost:3102 every time this environment is created. It does not start
the app; the agent does that when needed.

Preview it in HOST:

```bash
sbx env plan factory/sbxenv.yaml
```

Find your sample-app path and the 3102 → 8080 mapping. A plan shows the proposed
setup without creating it. In your SANDBOX tab, still at the repository root:

```bash
sbx env run factory/sbxenv.yaml
```

Read and approve the plan. This creates the environment and opens Claude.
Ask the assistant to inspect the project
and confirm your filter-count change is present. You have reused the actual files
and history, not a fresh copy of the application.

Type `/exit` to leave Claude. In HOST:

```bash
sbx env rm factory/sbxenv.yaml
```

Approve removal. The source stays in `sample-app/`. Removing this practice
sandbox lets us apply new configuration on the next creation.

## 2. Put work in a task tracker

[Beans](https://github.com/hmans/beans) is a command-line issue tracker. Its tasks
are Markdown files: a person can read them in an editor, and an agent can read
requirements, update notes and use the CLI without a browser integration.

In an ordinary project, `beans init` creates the task store, `beans create` adds
work, and `beans show` reads it. Our backlog will stay **on the host**, outside the
mounted app. Later, MCP will expose selected operations to the agents.

In HOST, install the pinned CLI and initialize the supplied workshop tasks:

```bash
./scripts/install-beans.sh
./scripts/backlog-init.sh --disposable
```

The first helper downloads Beans. The second creates our four demo tasks under
`.local/chapters/beans/`. `--disposable` marks this as practice data that our later
MCP server may append notes to.

Try the CLI:

```bash
./scripts/beans list
./scripts/beans show wad-101
./scripts/beans create "Try a small improvement" --type task --status todo
```

`list` shows the backlog; `show` prints the warm-up requirements; `create` adds a
new task and prints its ID. Use `./scripts/beans show YOUR-ID` with that ID to
read it. Open `.local/chapters/beans/.beans/` in your editor to see the Markdown
behind the commands.

`scripts/beans` is a small wrapper around the real Beans CLI. It selects this
workshop's executable and backlog so each command need not repeat their paths.
It does not run an agent. We will keep using the supplied tasks for the walkthrough.

## 3. Connect a task to a sandbox

We want a reusable environment with a different name for each chapter. At the top
of `factory/sbxenv.yaml`, replace the `name` and `workspace` values and add `args`.
Keep `sandboxOptions` and `ports` as they are:

```yaml
schemaVersion: "1"
name: "${{ env.args.name }}"
agent: claude
workspace: "${{ env.args.app }}"
args:
  name:
    required: true
  app:
    default: ../sample-app
```

The name is an input now. The app path defaults to the same working directory;
chapter 07 will show how to select another project.

Create `factory/chapter.env`:

```text
TASK=wad-101
MODE=manual
USE_ACR=0
SESSION=claude
```

These are our helper's settings, not SBX syntax. They select the task, leave team
startup manual, leave guidance installation off, and open Claude for us.

Create `factory/PROMPT.md`:

```markdown
Read ~/work/task.json and inspect this project. Confirm whether the warm-up
filter-count change is already present; do not implement it twice. Explain
how the current code meets the task and what you would check.
```

This file is the agent's instruction. The environment file provides the workspace;
the task says what the change is; the prompt says how to approach today's exercise.

## 4. Run the small host launcher

Open `chapters/support/launch` if you want to inspect the glue. It reads the selected
Bean, calls SBX to mount the project, supplies the task and workshop helpers, then
opens your session. It does not implement the task, start application services,
or verify the result.

The SBX creation command inside it is equivalent to this **reference example**:

```text
sbx env create factory/sbxenv.yaml --env-arg name=wad-ch-02 --auto-approve
```

The launcher also passes the absolute application path. `--auto-approve` applies
the configuration you just inspected. Preview that configuration in HOST:

```bash
sbx env plan factory/sbxenv.yaml --env-arg name=wad-ch-02
```

Then in the SANDBOX tab, at the workshop root:

```bash
./scripts/launch-factory.sh wad-ch-02
```

The name selects this sandbox. The launcher opens Claude in the mounted project.
Tell it:

> Read ~/work/PROMPT.md and follow those instructions. Explain what you found.

You should recognize your own warm-up change. In your host editor, the same files
and Git history are still in `sample-app/`.

Type `/exit` to end the session. In HOST:

```bash
sbx env rm factory/sbxenv.yaml --env-arg name=wad-ch-02
```

Next: [give each worker shared guidance through a kit](../02.5-acr/README.md).
