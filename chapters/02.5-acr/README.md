# 2.5. Give every fresh worker the same guidance

**Goal:** understand a mixin kit by making a tiny one, then add ACR and a small
versioned coding policy. We want fresh agents to know our conventions without
pasting the same instructions into every conversation.

## 1. Try a kit small enough to understand completely

A [kit](https://docs.docker.com/ai/sandboxes/customize/kits/) describes setup and
access requirements. A **mixin** adds something to a sandbox's base configuration.
Create this kit on the host:

```bash
# HOST
mkdir -p "$WORKSHOP/chapters/my-hello-kit"
cat > "$WORKSHOP/chapters/my-hello-kit/spec.yaml" <<'YAML'
schemaVersion: "2"
kind: mixin
name: workshop-hello
setup:
  install:
    - user: agent
      command: |
        mkdir -p "$HOME/.local/bin"
        printf '#!/bin/sh\necho "Hello from a kit"\n' > "$HOME/.local/bin/workshop-hello"
        chmod +x "$HOME/.local/bin/workshop-hello"
YAML
sbx kit validate "$WORKSHOP/chapters/my-hello-kit"
sbx create claude --name wad-kit-first --skills off --cpus 4 --memory 8g \
  --kit "$WORKSHOP/chapters/my-hello-kit"
sbx exec wad-kit-first /home/agent/.local/bin/workshop-hello
sbx stop wad-kit-first
```

Expected: `Hello from a kit`. The install command runs **inside** the sandbox as
`agent`, not on your host. This kit needs no extra network access. Its completed
reference is in `chapters/examples/hello-kit`.

## 2. Replace the toy capability with ACR

Our real requirement is a policy plus a reusable review skill. The
[ACR kit](https://github.com/shelajev/acr-sbx-kit) installs the tool that distributes
these. The [small workshop policy](https://github.com/shelajev/coding-policy/tree/workshop)
contains four rules and one review skill; it does not bring an orchestration system.

Build on your previous recipe, adding a `kits` section:

```bash
# HOST
mkdir -p "$WORKSHOP/chapters/my-025"
cd "$WORKSHOP/chapters/my-025"
for f in sbxenv.yaml chapter.env PROMPT.md launch; do cat "../my-02/$f" > "$f"; done
chmod +x launch
cat >> sbxenv.yaml <<'YAML'
kits:
  - source: "git+https://github.com/shelajev/acr-sbx-kit.git#ref=f446b95bccabd879912db174c190ff09377b7c7b"
YAML
# Keep USE_ACR=0 for this first launch: install the policy yourself below.
cat ../02.5-acr/PROMPT.md > PROMPT.md
sbx env plan sbxenv.yaml --env-arg name=wad-ch-02-5 \
  --env-arg port=3103 --env-arg "control_dir=$CONTROL"
```

A Git source pins the kit so new environments get the same tool setup. The kit
installs ACR; it is not itself the coding-policy package. Launch the recipe:

```bash
# HOST — terminal A, in chapters/my-025
WORKSHOP_PORT=3103 WORKSHOP_APP_REPO="$WARMUP" WORKSHOP_APP_REF=HEAD \
  ./launch wad-ch-02-5
```

In terminal B:

```bash
# HOST
sbx exec -it -w /home/agent/work/app wad-ch-02-5 bash
```

Now install and materialize the policy:

```bash
# SANDBOX
acr --version
env -u GH_TOKEN -u GITHUB_TOKEN acr install \
  github:shelajev/coding-policy@b85031eb0c8963b28b63eaa12efcbd34c850d32d \
  --agent claude-code --agent codex --freshness none --non-interactive
env -u GH_TOKEN -u GITHUB_TOKEN acr realize
cat agents.yaml
cat AGENTS.md
cat .claude/skills/acr__shelajev__coding-policy__review-change/SKILL.md
```

`install` records the package and selected agent targets; `realize` produces the
agent-facing files. The public package needs no GitHub login. We unset token
variables for these commands because a proxy placeholder is not a public GitHub
credential. We neither print nor copy real secrets.

## 3. Use the skill, then make installation repeatable

```bash
# SANDBOX — still in ~/work/app
env -u GH_TOKEN -u GITHUB_TOKEN claude --dangerously-skip-permissions
```

Ask:

> Read AGENTS.md and the installed review-change skill. Use that skill to review
> the warm-up against ~/work/task.json. Do not modify the app. Report the rules
> applied and concrete findings. Do not install a browser for this exercise.

Look for the policy and skill being read; files merely existing is not evidence
of their use. Exit the agent and sandbox shell when done.

For subsequent fresh launches, enable the supplied in-sandbox preparation step:

```bash
# HOST — terminal B
cat "$WORKSHOP/chapters/02.5-acr/chapter.env" \
  > "$WORKSHOP/chapters/my-025/chapter.env"
cat "$WORKSHOP/chapters/support/bin/prepare"
```

`USE_ACR=1` makes preparation run the commands you just tried. It takes effect on
the **next** launch. We have changed a local configuration file, not a running
script. The current sandbox already has its manually installed policy.

**Result:** the recipe installs a kit and preparation installs the pinned skills
package. End the terminal A hold and `sbx stop wad-ch-02-5`. Next we reuse this
guidance in a second assistant.

Next: [Pi and model configuration](../03-pi/README.md).
