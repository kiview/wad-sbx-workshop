# 7. Take the factory with you

**Goal:** identify the working system you built, keep its output, and know how to
give it another task. There is no new infrastructure component in this chapter.

## 1. Inspect the result of your main run

Use chapter 05's completed assignment/resolution run, or chapter 06's later result.
You should be able to point to each of these:

| Evidence | Where to look |
|---|---|
| Working app behavior | Published browser port of the corresponding live sandbox |
| Committed source | `git log -1` and `git status` inside `~/work/app` |
| Real implementation/review conversation | Herdr terminal output and handoff files |
| Host task write-back | `beans show wad-102` or `wad-103` on the host |
| Reproducible infrastructure | Your latest `sbxenv.yaml`, `chapter.env`, `team.tsv` |
| Reusable guidance | Pinned ACR kit and coding-policy package |

Retrieve the source before stopping its sandbox, as in chapter 05. Keep the copied
Git repository and your `chapters/my-*` recipes. The factory does not need a host
acceptance gate, merge bot or task-closing system. You can add independent
verification later if your real workflow calls for it.

## 2. Try another job during open lab

This step is for the remaining workshop time, not a prerequisite for claiming that
the earlier MCP feature run worked. It reuses exactly the same infrastructure.

```bash
# HOST
mkdir -p "$WORKSHOP/chapters/my-07"
cd "$WORKSHOP/chapters/my-07"
cat ../my-06/sbxenv.yaml > sbxenv.yaml
cat ../my-06/team.tsv > team.tsv
for f in chapter.env PROMPT.md launch; do cat "../07-factory/$f" > "$f"; done
chmod +x launch
cat chapter.env
```

`TASK=wad-104` selects a service filter. The default app fixture is the known
assignment/resolution solution. To continue your own chapter-05 implementation:

```bash
# HOST — terminal A
WORKSHOP_APP_REPO="$WORKSHOP/.local/feature-app" WORKSHOP_APP_REF=HEAD \
  ./launch wad-ch-07
```

Use `.local/reopen-app` instead if you saved your chapter-06 result there. If you
need the supplied fixture, run `./launch wad-ch-07` without overrides. Choose one
path, not all three. The starting source must contain the features the task assumes.

In terminal B, wait for role initialization, then submit once:

```bash
# HOST
sbx exec wad-ch-07 /home/agent/work/bin/assign
sbx exec wad-ch-07 herdr agent list
"$CONTROL/bin/beans" --config "$CONTROL/beans/.beans.yml" \
  --beans-path "$CONTROL/beans/.beans" show wad-104
```

When the team reports completion, open <http://127.0.0.1:3108>, filter by
`checkout-api`, and inspect the result count. Compare the Bean note's commit with
the commit in SBX and inspect QA's findings. Save the result:

```bash
# HOST — after the team finishes and commits
sbx exec wad-ch-07 git -C /home/agent/work/app status --short
sbx cp wad-ch-07:/home/agent/work/app "$WORKSHOP/.local/service-filter-app"
```

## 3. Adapt it to your own application

Replace the app source, startup helper and task contracts. Keep responsibilities
clear: the host chooses a task and starts an environment; the sandbox team
implements/reviews; MCP supplies deliberately scoped host tools; a human answers
questions. Edit roles/models in `team.tsv` and skills in the policy package without
rewriting the launcher. Add tools through small kits when they are needed.

## 4. Stop without losing the work

Copy committed source before cleanup. End each launch terminal with Ctrl-C, then
stop only your named sandbox, for example:

```bash
# HOST
sbx stop wad-ch-07
```

A stop preserves sandbox files but terminates live processes. It is not a saved
agent-session checkpoint. Keep the host backlog and copied repositories. Do not
remove unrelated sandboxes or reset global policies during workshop cleanup.

**You now have:** repeatable sandbox creation, container execution, shared skills,
configurable agent roles, team coordination, a scoped host-task bridge and a human
intervention path.

Next: [presenter extensions](../08-presenter/README.md), time permitting.
