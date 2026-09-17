# 2.5. Give every worker the same tools and guidance

Our worker has the application and its task. How does it learn our team's coding
conventions? We could paste them into every conversation. Instead, we will install
a versioned policy and review skill, then make that preparation repeatable.

There are two building blocks: an SBX **kit** prepares the environment; **ACR**
installs the guidance the agent reads in a project.

## 1. Try a kit small enough to understand completely

A [mixin kit](https://docs.docker.com/ai/sandboxes/customize/kits/) adds capabilities
to an existing sandbox environment. In your editor, create
`factory/hello-kit/spec.yaml` with:

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

The `setup.install` block is code SBX runs **inside** the sandbox, as the `agent`
user. It creates a tiny command and makes it executable. `$HOME` is the sandbox
home. This shell code belongs to the kit; you are not running it on your laptop.

Validate the definition in HOST:

```bash
sbx kit validate factory/hello-kit
```

A valid result means SBX accepts its structure. Now try its actual effect in your
SANDBOX tab, starting at the workshop root:

```bash
sbx run claude "$(pwd)/sample-app" --name wad-kit-first --skills off --cpus 4 --memory 8g --kit ./factory/hello-kit
```

`--kit` adds our installation step to the built-in Claude environment. In Claude:

```text
!workshop-hello
```

Expect `Hello from a kit`. The tool was installed during creation, not by your
conversation. Exit with `/exit`. In HOST:

```bash
sbx rm wad-kit-first
```

## 2. Use a kit to supply a real tool

The [ACR kit](https://github.com/shelajev/acr-sbx-kit) installs a package manager
for agent guidance. Our [workshop policy package](https://github.com/shelajev/coding-policy/tree/workshop)
contains four coding rules and a `review-change` skill. Open that package and read
a rule: how would it help someone review the change you just made?

The kit installs the **tool**. The package supplies the **rules and skill**. This
separation lets you update shared guidance without inventing a new sandbox image.

Add this section to `factory/sbxenv.yaml`:

```yaml
kits:
  - source: "git+https://github.com/shelajev/acr-sbx-kit.git#ref=f446b95bccabd879912db174c190ff09377b7c7b"
```

Unlike our local hello kit, this source is a Git repository pinned to a commit.
Every attendee receives the same kit revision. Keep `USE_ACR=0` for now: we want
to install the guidance together before automating it.

In HOST:

```bash
sbx env plan factory/sbxenv.yaml --env-arg name=wad-ch-02-5
```

Find the ACR kit and the unchanged application mount. In your SANDBOX tab:

```bash
./scripts/launch-factory.sh wad-ch-02-5
```

SBX may also ask you to approve the kit's credential bindings. In the linked kit's
`spec.yaml`, find `credentials`: it declares an **optional** service named
`acr-github`, for private repositories, higher API limits and publishing. Its
`inject` entries name the GitHub API, archive and upload hosts that may receive
that credential.

Read those service and destination names in the approval prompt. A binding
describes where a configured credential may be used; it does not mean every
download needs one. Our policy package is public, so you do not need to create a
GitHub token or run a secret-setup command for this exercise.

## 3. Install and inspect shared guidance

Ask Claude:

> Run acr --version. Then run the workshop's install-guidance helper. Show me
> agents.yaml, AGENTS.md and the installed review-change skill. Explain which file
> declares the package and which files tell an assistant how to work.

The helper performs these ACR operations, shown here for explanation:

```text
acr install github:shelajev/coding-policy@b85031eb0c8963b28b63eaa12efcbd34c850d32d --agent claude-code --agent codex --freshness none --non-interactive
acr realize
```

`install` selects the pinned policy and the assistant formats to generate.
`realize` writes the agent-facing guidance. Claude and Codex have generated formats;
Pi can read the shared `AGENTS.md` and skill too. The fixed revision and
`--freshness none` keep everyone on the same exercise; `--non-interactive` uses
these supplied choices. The helper also removes GitHub token placeholders for this
public download; that compatibility detail stays in the helper.

Open `sample-app/AGENTS.md` in your host editor. These are real files in the same
mounted project. Ask Claude:

> Read AGENTS.md and use the installed review-change skill to review the warm-up
> against ~/work/task.json. Explain one rule you checked and the code that supports
> your conclusion. Report a problem only if you find one. Do not edit application code.

Discuss its answer. Shared guidance is useful when you can connect the advice to
an actual rule and change, not just because a file was installed.

## 4. Make guidance part of future preparation

In HOST, set `USE_ACR=1` in `factory/chapter.env`. Future launches install the policy
for a project without existing guidance; this project's installed guidance remains
in its working copy. The helper preserves a project's existing `AGENTS.md` instead
of silently replacing it.

Exit Claude. In HOST:

```bash
sbx env rm factory/sbxenv.yaml --env-arg name=wad-ch-02-5
```

We now have a repeatable environment and shared coding guidance. Next we will add
another assistant through another kit.

Next: [bring in Pi](../03-pi/README.md).
