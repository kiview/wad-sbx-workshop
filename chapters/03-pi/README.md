# 3. Bring a second assistant into the workspace

Claude can work on the app and read our policy. Now we want a choice of assistant
and model for the roles we will create. We will add [Pi](https://pi.dev), an
open-source coding assistant, through another kit. You will enter the sandbox,
start Pi yourself and have a conversation about the same repository. This gives us
a chance to understand its model configuration before asking it to join a team.

## 1. Decide what belongs in the kit

Open `chapters/kits/pi/spec.yaml` in your editor. Find these three parts:

| Part | Why this assistant needs it |
|---|---|
| `pi_version` | A fresh environment should install the same tool version. |
| `permissions.network.allow` | Installation and model calls need access to specific hosts. |
| `setup.install` | Download, check and install the executable as the sandbox user. |

The checksum and CPU-architecture handling are supplied implementation details. The
important choice is to record installation and access in a kit, so each new worker
gets the capability without somebody installing Pi by hand again.

A [template](https://docs.docker.com/ai/sandboxes/customize/templates/) is the starting
filesystem. A [kit](https://docs.docker.com/ai/sandboxes/customize/kits/) adds setup
and access requirements. We are extending the existing Claude environment, rather
than building a new image during the workshop.

## 2. Add Pi to your existing recipe

Keep your environment recipe and update the exercise instructions:

```bash
for file in chapter.env PROMPT.md; do
  cat "chapters/03-pi/$file" > "factory/$file"
done
```

This updates the task settings and prompt. Your environment recipe keeps the ACR
kit you added in the previous chapter.

In your editor, add this entry to the **existing** `kits` list in `factory/sbxenv.yaml`:

```yaml
  - source: ../chapters/kits/pi
```

Keep the ACR entry. You are composing two capabilities, not replacing one with the
other. The relative path is resolved from this recipe's directory.

Preview the recipe:

```bash
sbx env plan factory/sbxenv.yaml --env-arg name=wad-ch-03
```

Only the sandbox name changes; the recipe supplies the port and installed kits. Find the Pi kit in the plan. Then create the app
environment using our existing launcher:

```bash
# HOST — terminal A, from the workshop repository
./scripts/launch-factory.sh wad-ch-03
```

This still calls `sbx env create`, transfers your committed app and starts its
services. The additional kit makes Pi available. Leave this terminal open.

## 3. Enter the sandbox and start Pi yourself

In terminal B on the host:

```bash
sbx exec -it wad-ch-03 bash
```

`exec` starts a process in the existing sandbox; `-it` gives you an interactive
terminal; `bash` is the shell we want to use. From now until you type `exit`, commands
in this terminal run **inside SBX**.

```bash
# SANDBOX shell
cd ~/work/app
pi --version
```

`cd` puts you in the same app the developer will use. The version output tells you
which executable the kit installed. Next we will open its actual interface.

### Give Pi access to its model

For the workshop's Anthropic route, SBX must already have a usable Anthropic API
credential. If you need to set one, do it in another **host** terminal with
`sbx secret set anthropic`; enter the real key at the prompt. Claude subscription
login and another assistant's API access are separate things.

Back in your sandbox shell:

```bash
export ANTHROPIC_API_KEY=proxy-managed
pi --provider anthropic --model claude-sonnet-5
```

The export gives Pi the placeholder expected by our SBX credential-proxy setup.
It is not your real key: the proxy supplies that when the request leaves SBX.
The variable belongs to this sandbox shell; do not put it in your host shell profile.

`pi` opens the interactive assistant. `--provider` selects the service and `--model`
selects the initial model. Accept its project-trust prompt if shown. Inside Pi,
type `/model` to explore the available choices. Select a model your account can use
and note its ID—we will put that choice into the team configuration next.
Pi's [interactive commands](https://github.com/earendil-works/pi/blob/v0.85.1/packages/coding-agent/README.md#commands)
include `/quit` to return to your shell.

Pi also supports its own provider login flows. This workshop's automated team uses
the SBX proxy configuration above; a separate Pi login is not automatically carried
to a newly created sandbox. If you have only subscription access, work with the
instructor on the provider route before continuing to the automated team.

## 4. Have a conversation about the same application

Ask Pi:

> Look around this project and explain what this sample application does. Do not change
> anything yet. Where is the active-filter result count calculated?

Let it inspect the files and answer. Follow up in your own words if its explanation
is unclear. Then ask:

> Read AGENTS.md and
> .claude/skills/acr__shelajev__coding-policy__review-change/SKILL.md.
> Use that review skill on the warm-up described in ~/work/task.json. Explain one
> concrete thing you checked against our policy. Do not modify the app.

You are using a different assistant in the same environment with the same guidance.
Compare its explanation with Claude's. Does it point to actual code? Can you follow
why it considers the warm-up correct? These are useful questions for a reviewer,
not just a check that the API returned text.

If you get an authentication or quota error, pause here and resolve the provider
choice with the instructor. Installing another orchestrator will not fix account
access. The [mixed-provider variation](MIXED-MODELS.md) shows the intended
Pi/Google, Claude/Anthropic and Codex/OpenAI combination when accounts permit it.

## 5. Turn your choices into roles

Type `/quit` to leave Pi, then `exit` to leave the sandbox shell. Note down the
working model ID you selected. Read `chapters/04-team/team.tsv` as an example—leave
this reference file unchanged. In the next chapter you will edit your own copy at
`factory/team.tsv`. Each row separates four choices:

| Choice | Question it answers |
|---|---|
| Role | What responsibility does this agent have? |
| Harness | Which assistant program runs it? |
| Provider | Which service supplies its model? |
| Model | Which specific model should that role use? |

A coordinator can run in Pi while a developer runs in Claude Code. A QA role can use
another provider, or the same one if that is all you have. Everyone configures the
roles; extra subscriptions determine how diverse the live team can be.

We now have assistants, but you are still carrying information between them by hand.
Next we will give them sessions, responsibilities and a way to talk to each other.
End terminal A's hold with Ctrl-C, then run `sbx stop wad-ch-03` on the host.

Next: [turn the assistants into a team](../04-team/README.md).
