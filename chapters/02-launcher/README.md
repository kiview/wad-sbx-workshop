# 2. Describe the sandbox, then launch a task

**Goal:** replace repeated creation commands with an environment file, then deliver
one Bean and app source with a readable host script. The host will not supervise
agents, judge results, merge code or mark tasks complete.

## 1. Try the smallest environment file

An [SBX environment file](https://docs.docker.com/ai/sandboxes/configuration/environment-files/)
is a recipe for creating an environment. Write one before using the full assembly:

```bash
# HOST
mkdir -p "$WORKSHOP/chapters/my-first-env"
cat > "$WORKSHOP/chapters/my-first-env/sbxenv.yaml" <<'YAML'
schemaVersion: "1"
name: wad-env-first
agent: claude
sandboxOptions:
  skills: off
  cpus: 4
  memory: 8g
YAML
sbx env plan "$WORKSHOP/chapters/my-first-env/sbxenv.yaml"
sbx env create "$WORKSHOP/chapters/my-first-env/sbxenv.yaml" --auto-approve
sbx exec wad-env-first bash -lc 'whoami; docker version'
sbx stop wad-env-first
```

Compare this with `sbx create claude --name wad-manual --skills off --cpus 4 --memory 8g`. The file
records choices for the next launch. `plan` lets you inspect the resolved recipe;
it does not prove the sandbox or its applications successfully run.

## 2. Introduce the host backlog

[Beans](https://github.com/hmans/beans) stores tasks as files. We use a separate
workshop backlog so agents can later append notes without touching your real tasks.
Install the pinned CLI and seed the four task contracts:

```bash
# HOST
cd "$WORKSHOP"
export FACTORY_CONTROL_DIR="$CONTROL"
scripts/install-beans.sh
scripts/backlog-init.sh --disposable
"$CONTROL/bin/beans" --config "$CONTROL/beans/.beans.yml" \
  --beans-path "$CONTROL/beans/.beans" show wad-101
```

Read the expected behavior in the task. This is the same warm-up you just did.
The `--disposable` marker explicitly identifies this practice backlog; the MCP
adapter will require it before writing notes in chapter 05.

## 3. Make your working recipe

Start a learner directory directly under `chapters/`. The supplied recipe adds
named arguments, resources and a loopback port mapping to your minimal example.

```bash
# HOST
mkdir -p "$WORKSHOP/chapters/my-02"
cd "$WORKSHOP/chapters/my-02"
for f in sbxenv.yaml chapter.env PROMPT.md launch; do
  cat "../02-launcher/$f" > "$f"
done
chmod +x launch
cat sbxenv.yaml
cat chapter.env
```

`name`, `port` and `control_dir` are recipe inputs. `sandboxOptions` sets the VM
resources; `ports` maps the app's 8080 to a host loopback port. There is deliberately
no `workspace` mount. `control_dir` becomes useful for MCP later.
`chapter.env` is input to our Bash launcher, **not SBX configuration**: it chooses
the task, app fixture, app port and which teaching helpers to start.

```bash
# HOST
sbx env plan sbxenv.yaml --env-arg name=wad-ch-02 \
  --env-arg port=3102 --env-arg "control_dir=$CONTROL"
cat launch
cat ../support/launch
```

Read the launcher as four operations: export committed app source; obtain a Bean
snapshot; create the recipe; transfer and prepare files inside SBX. The final
`sleep infinity` holds the client connection open. It is not a task scheduler.
Inside-SBX app setup is supplied in `support/bin/prepare` and `board` so we do not
spend the workshop rewriting npm/database startup commands.

## 4. Run it with your warm-up result

```bash
# HOST — terminal A, in chapters/my-02
WORKSHOP_APP_REPO="$WARMUP" WORKSHOP_APP_REF=HEAD ./launch wad-ch-02
```

Leave this terminal open. In terminal B:

```bash
# HOST
curl -fsS http://127.0.0.1:3102/healthz
sbx exec wad-ch-02 bash -lc 'cat ~/work/task.json; git -C ~/work/app log -1 --oneline'
sbx exec -it -w /home/agent/work/app wad-ch-02 \
  env -u GH_TOKEN -u GITHUB_TOKEN claude --dangerously-skip-permissions
```

Ask the agent to read the task and confirm the warm-up is present; do not implement
it a second time. The launcher creates a new Git repository from the exported
source, so its initial commit ID differs from your original result. Source content
is carried; Git history is not carried by `git archive`.

**Result:** one short host command creates the sandbox, starts the app and supplies
the task. Open <http://127.0.0.1:3102>. A fresh worker still lacks your conventions.
Next we teach it through a kit and a policy package.

Exit Claude, end terminal A's hold with Ctrl-C, then `sbx stop wad-ch-02` before
continuing. We reuse `$WARMUP` for the next review exercises, which do not edit it.
If you made another change, commit and copy it to a new host path first, then set
`WARMUP` to that path in both terminals.

Next: [kits and ACR](../02.5-acr/README.md).
