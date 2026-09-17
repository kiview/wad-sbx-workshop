# 1. Let an agent act inside a sandbox

**Goal:** run an autonomous coding agent, see that the host filesystem is separate,
start a database container inside SBX, and open the app through a published port.
We start with direct commands so the later recipe has something concrete to replace.
See [Docker Sandboxes](https://docs.docker.com/ai/sandboxes/).

## 1. Create a sandbox without mounting your checkout

```bash
# HOST — terminal A
sbx create claude --name wad-manual --skills off --cpus 4 --memory 8g
sbx exec -it wad-manual bash
```

Leave this shell open. `create` without a workspace path gives us the isolated
starting point we want. Do not replace it with `sbx run .`: that would deliberately
share the current directory. Inside the new shell:

```bash
# SANDBOX — terminal A
whoami
pwd
docker version
```

The Docker daemon here belongs to the sandbox. Your host does not need one.

## 2. Check the filesystem boundary

In a second host terminal, with the variables from chapter 00:

```bash
# HOST — terminal B
mkdir -p "$CONTROL"
printf 'host-only workshop marker\n' > "$CONTROL/host-marker.txt"
sbx exec wad-manual test ! -e "$CONTROL/host-marker.txt"
echo $?
```

Expected: `0`, meaning the same absolute path is absent inside SBX. We have not
mounted that host directory. This demonstrates the boundary for this file; it is
not a claim that no access can ever be granted through a mount or a tool.

## 3. Transfer the starter and its first task

```bash
# HOST — terminal B
git -C "$APP_REPO" archive app-00-starter | sbx exec -i wad-manual bash -lc \
  'mkdir -p "$HOME/work/app"; tar xf - -C "$HOME/work/app"'
sbx exec -i wad-manual bash -lc 'cat > "$HOME/work/task.md"' \
  < "$WORKSHOP/backlog/seed/wad-101--warm-up-active-filter-count.md"
```

`-i` forwards stdin. Here we send an archive and a task, not a live host mount.
Inside the shell you already opened:

```bash
# SANDBOX — terminal A
cd ~/work/app
git init
git config user.name 'Workshop learner'
git config user.email 'workshop@example.invalid'
git add .
git commit -m 'Workshop starter'
npm ci --no-audit --fund=false
./scripts/db-up.sh
export DATABASE_URL=postgres://board:board@127.0.0.1:55432/board
npm run db:migrate
npm run db:seed
docker ps
PORT=8080 HOST=0.0.0.0 npm run dev
```

Leave the development server running. The database is a real PostgreSQL container;
the agent can also start containers for its own checks inside this environment.

```bash
# HOST — terminal B
sbx ports wad-manual --publish 3101:8080
curl -fsS http://127.0.0.1:3101/healthz
sbx exec -it -w /home/agent/work/app wad-manual \
  env -u GH_TOKEN -u GITHUB_TOKEN claude --dangerously-skip-permissions
```

Open <http://127.0.0.1:3101>. Complete the first-use trust/login prompts if shown.
This permission mode lets the agent execute without approving each action; the
sandbox's filesystem and access policy still apply.

Give Claude this prompt:

> Read ~/work/task.md and implement only wad-101, the active-filter count warm-up.
> Inspect the app first. Run the relevant checks inside this sandbox. Commit the
> change and report what changed and the commit. Do not install a browser or run
> the browser-test suite for this exercise.

Watch the agent use files and tools. In the browser, try a filter and inspect the
count after its change. We are learning the execution environment, not benchmarking
models or running every possible test.

## 4. Save the warm-up for the following chapters

Exit the agent after it finishes, leaving terminal A's server running. Confirm that
its changes are committed. If there are uncommitted changes, ask it to commit its
intended result before copying.

```bash
# HOST — terminal B
sbx exec wad-manual git -C /home/agent/work/app status --short
sbx exec wad-manual git -C /home/agent/work/app log -1 --oneline
sbx cp wad-manual:/home/agent/work/app "$WARMUP"
git -C "$WARMUP" log -1 --oneline
```

`$WARMUP` must be a new destination. If it exists from a previous rehearsal, choose
another path rather than overwriting it. The launcher transfers committed source
from it; copied dependencies are not used. Keep this repository as an artifact,
without running its application code on the host.

After retrieving the work, stop the development server with Ctrl-C in terminal A,
exit its shell, and run `sbx stop wad-manual` on the host.

**Result:** you have a changed app and its committed source, created by an agent
with a working container engine inside SBX. The commands are repetitive. Next we
put the environment into a recipe and the task delivery into a tiny launcher.

**Catch up:** use the supplied `app-01-warmup-solution` fixture in later chapters if
the warm-up takes too long. Omit the `WORKSHOP_APP_REPO`/`WORKSHOP_APP_REF` overrides
there; chapter 02 needs `WORKSHOP_APP_REF=app-01-warmup-solution` for that shortcut.

Next: [recipes and the host launcher](../02-launcher/README.md).
