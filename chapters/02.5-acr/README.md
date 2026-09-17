# 2.5. Give a fresh worker your team's guidance

A fresh sandbox now has our code and task, but the agent still needs our team's
coding guidance. We will try a tiny kit to learn how kits add tools to an
environment, then use the ACR kit to install a package manager for agent guidance.
Through ACR, we will give Claude a versioned coding policy and review skill, use
that skill together, and only then automate its installation.

## 1. Make a kit whose effect you can see

A [mixin kit](https://docs.docker.com/ai/sandboxes/customize/kits/) adds setup and
access requirements to a sandbox. On the host, make a directory for yours:

```bash
mkdir -p "$WORKSHOP/chapters/my-hello-kit"
```

In your editor, create `chapters/my-hello-kit/spec.yaml` with this content:

```yaml
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
```

Read it from the outside in. `kind: mixin` means “add this to an environment.”
`setup.install` contains commands SBX will run **inside** that environment as the
`agent` user. The three shell lines create a tool directory, write a tiny program
and make it executable. `$HOME` here is the sandbox user's home, not your host home.

Ask SBX to validate your kit, then create a practice environment with it:

```bash
# HOST
sbx kit validate "$WORKSHOP/chapters/my-hello-kit"
sbx create claude --name wad-kit-first --skills off --cpus 4 --memory 8g \
  --kit "$WORKSHOP/chapters/my-hello-kit"
sbx run --name wad-kit-first
```

Validation checks the kit definition. `--kit` adds it to the built-in Claude
configuration. `create` prepares this environment without an app mount; `run`
opens the agent in it. In Claude, try the installed command directly:

```text
!$HOME/.local/bin/workshop-hello
```

You should see `Hello from a kit`. You did not install the command after entering
the sandbox: its creation recipe did. That is what we want for real tooling too.
Exit Claude, then run `sbx stop wad-kit-first` on the host.

## 2. Choose the real capability to distribute

The [ACR kit](https://github.com/shelajev/acr-sbx-kit) installs ACR.
Our [small policy package](https://github.com/shelajev/coding-policy/tree/workshop)
contains four coding rules and one `review-change` skill. Open the policy and read
the rules before giving them to the agent. Which would matter for the warm-up?

The distinction matters: **the kit makes ACR available; the policy package supplies
the content ACR installs**.

Prepare your next recipe on the host:

```bash
mkdir -p "$WORKSHOP/chapters/my-025"
cd "$WORKSHOP/chapters/my-025"
cat ../my-02/sbxenv.yaml > sbxenv.yaml
for file in chapter.env PROMPT.md launch; do cat "../02.5-acr/$file" > "$file"; done
chmod +x launch
```

This carries your environment forward and brings in this chapter's settings.
In your editor, add the following top-level section to `sbxenv.yaml`:

```yaml
kits:
  - source: "git+https://github.com/shelajev/acr-sbx-kit.git#ref=f446b95bccabd879912db174c190ff09377b7c7b"
```

Instead of a local directory, `source` now names a Git repository at a fixed commit.
SBX fetches that kit when creating the sandbox. Open `chapter.env` and change
`USE_ACR=1` to `USE_ACR=0` for this first run. That keeps our helper from installing
the policy automatically—we want to do it ourselves and see what it produces.

Preview your composition:

```bash
sbx env plan sbxenv.yaml --env-arg name=wad-ch-02-5 \
  --env-arg port=3103 --env-arg "control_dir=$CONTROL"
```

Find the kit in the plan. These are the recipe arguments introduced in chapter 02.
Then start it using the same app-delivery wrapper:

```bash
# HOST — terminal A, in chapters/my-025
./launch wad-ch-02-5 "$WARMUP"
```

Leave this terminal open. In terminal B, join the agent:

```bash
# HOST
sbx run --name wad-ch-02-5
```

## 3. Install the policy together with Claude

Tell Claude:

> Work in /home/agent/work/app. We installed ACR through a kit. Run `acr --version`
> and explain what guidance this project currently has, without changing the app.

Paste the following **as a message to Claude**, not into a Bash shell. It asks the
agent to execute the commands; you do not need the `!` prefix used for direct shell mode.

```text
Work in /home/agent/work/app. Please run these commands to install our coding policy,
then show me the generated guidance files:

env -u GH_TOKEN -u GITHUB_TOKEN acr install \
  github:shelajev/coding-policy@b85031eb0c8963b28b63eaa12efcbd34c850d32d \
  --agent claude-code --agent codex --freshness none --non-interactive
env -u GH_TOKEN -u GITHUB_TOKEN acr realize
```

| Part | Why we use it |
|---|---|
| `env -u GH_TOKEN -u GITHUB_TOKEN` | Removes GitHub token variables for this command. This public download needs no token, and a proxy placeholder could otherwise be mistaken for one. |
| `acr install github:...@...` | Selects the policy package at a fixed commit. |
| `--agent claude-code --agent codex` | Requests guidance files for those agent formats. It does not install the assistants. |
| `--freshness none` | Keeps this exercise on the selected policy revision without freshness checks. |
| `--non-interactive` | Uses the choices supplied in the command without another questionnaire. |
| `acr realize` | Materializes the agent-facing guidance in the project. |

No real token is passed by those `env` options, and no stored credential is deleted.
Ask Claude to show you `agents.yaml`, `AGENTS.md` and the installed review skill.
Notice the difference between the package declaration and the files an agent reads.

## 4. Use the guidance, then automate its installation

Continue the conversation:

> Read AGENTS.md and the installed review-change skill. Use the skill to review
> the warm-up against ~/work/task.json. Explain a concrete finding and which rule
> informed it. Do not modify the application.

Follow up on one finding. Can you connect the advice to the rule you read earlier?
The shared skill gives both assistants the same review procedure.

Exit Claude. On the host, edit `chapters/my-025/chapter.env` back to `USE_ACR=1`. Open
`chapters/support/bin/prepare` and find the `acr install` and `acr realize` lines.
On future launches, that helper performs the installation you just tried manually.
The current sandbox already has its policy; editing this file affects future runs.

End terminal A's hold with Ctrl-C, then run `sbx stop wad-ch-02-5`. We can now give a
fresh worker our guidance. Next we will use it with a different assistant.

Next: [bring in Pi](../03-pi/README.md).
