# 3. Give the factory a choice of assistant

A coding team needs roles, but those roles need not all use the same assistant or
model. SBX is a useful place to try those choices: a new assistant and its tools can
live inside the sandbox instead of becoming another installation on your laptop.
You can give it a real project, constrain its access and see how it behaves.

We will install [Pi](https://pi.dev), an open-source coding assistant, and
have a conversation with it about the same application and coding policy.

We will keep the project, task and guidance the same, so the comparison is about
how the assistant works rather than about a different starting point.

This introduces three separate choices: the assistant program, its model provider,
and the particular model. Later we will assign those choices to team roles.

## 1. Understand what the Pi kit supplies

Open `chapters/kits/pi/spec.yaml`. Read these parts before using the whole file:

| Part | Purpose |
|---|---|
| `args.pi_version` | Select a repeatable version of the assistant. |
| `permissions.network.allow` | Permit the download hosts and model API endpoints it needs. |
| `setup.install` | Download the right build, verify its checksum and install it as the sandbox user. |

The longer architecture and checksum code makes that small installation recipe
portable. You do not need to rewrite it. The important choices are **what to
install**, **where it may connect**, and **which user runs setup**.

A [template](https://docs.docker.com/ai/sandboxes/customize/templates/) supplies the
starting filesystem. A [kit](https://docs.docker.com/ai/sandboxes/customize/kits/)
adds capabilities and access requirements. We can combine kits on the existing
Claude template without building our own image.

## 2. Compose it with the guidance kit

Add this entry under the existing `kits` list in `factory/sbxenv.yaml`:

```yaml
  - source: ../chapters/kits/pi
```

Keep the ACR entry. Relative kit paths are resolved from the environment file,
so this points from `factory/` to the supplied Pi kit.

Set `SESSION=shell` in `factory/chapter.env`. We want a shell where we can start Pi
ourselves; `MODE=manual` and `USE_ACR=1` stay as they are.

In HOST:

```bash
sbx env plan factory/sbxenv.yaml --env-arg name=wad-ch-03
```

Check that both kits and the same application mount appear. In the SANDBOX tab:

```bash
./scripts/launch-factory.sh wad-ch-03
```

The launcher opens an interactive shell inside SBX, already in the mounted project.
It keeps this sandbox session open while you work. Try:

```bash
# SANDBOX shell
pi --version
```

You should see the version selected by the kit. Now we will use the actual assistant.

## 3. Connect an assistant to a provider

Claude subscription login and Pi's API access are separate. If you have an
Anthropic API key, register it from HOST:

```bash
sbx secret set anthropic
```

Enter the real key at the secure prompt, not in a command or file. Skip this if the
service is already configured. SBX keeps the real credential on the host and injects
it into permitted model requests. Allowing a network host and providing a credential
solve different problems: the request needs both permission and authentication.

In the SANDBOX shell, start Pi with the provider and model you want to try:

```bash
pi --provider anthropic --model claude-sonnet-5
```

Our kit's Pi wrapper supplies a placeholder API value when this route needs one;
the SBX proxy supplies the real credential. Open the kit's wrapper section if you
want to see that mechanism. No real key is written into the project.

Accept the workspace trust prompt if shown. In Pi, type `/model` to see available
models and choose one your account can use. Keep that model ID for the next chapter.

**Only have a Claude subscription?** You can explore Pi's interface, then use the
three-Claude team configuration in chapter 04. For a real Pi model conversation,
pair with someone who has API access or use the presenter's demonstration. Do not
enter a subscription token as an API key.

## 4. Use the same guidance with a different assistant

Ask Pi:

> Look around this project and explain what our sample application does. Where is
> the filter-result count calculated? Do not change anything yet.

Then:

> Read AGENTS.md and use the installed review-change skill. Review the warm-up
> described in ~/work/task.json. Explain one rule you checked against actual code.
> Do not modify the application.

Follow up in your own words. Does its answer help you understand the change? You
have changed the assistant while keeping the application, task and guidance.

## 5. Use this environment to try another model

Pi separates the assistant interface from the model behind it. Type `/model` to
inspect the choices your provider offers. If you have access to another model,
select it and ask the same review question. For Gemini in this workshop, always
use **`gemini-3.8-flash`**; the mixed-provider exercise uses that model.

Compare something concrete: does the answer identify the right code, apply the
shared rule, and support its conclusion? A different style of answer is not by
itself evidence of a better review. For a clean comparison, start a fresh conversation
with `/new` before repeating the request; that keeps the first answer out of context.
Keep the working provider/model ID you want for your team.

This is a reusable way to evaluate tools: put an assistant's installation and
network requirements in a kit, create its sandbox, give it a bounded task, and
inspect what it actually does. When you remove the sandbox, its installed tools
and running processes go away. The mounted project and its edits remain, so use
a separate checkout when an experiment should not affect work you want to keep.

Type `/quit` to return to the sandbox shell. The [Pi command reference](https://github.com/earendil-works/pi/blob/v0.85.1/packages/coding-agent/README.md#commands)
explains its other interactive controls. Type `exit` to return to the host shell
in this tab. In HOST:

```bash
sbx env rm factory/sbxenv.yaml --env-arg name=wad-ch-03
```

Next: [give the assistants roles and communication](../04-team/README.md).
