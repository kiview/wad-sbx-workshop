# 3. Add another assistant with another kit

**Goal:** run Pi beside Claude Code, using the same guidance. Then separate the
concepts of role, harness, provider and model before creating the team.

## 1. Understand what this kit needs

[Pi](https://pi.dev) is an open-source coding assistant that can use different
providers. Its kit needs an executable, download access, and access to the selected
model endpoint. Read the supplied implementation:

```bash
# HOST
cat "$WORKSHOP/chapters/kits/pi/spec.yaml"
sbx kit validate "$WORKSHOP/chapters/kits/pi"
```

Identify the pinned version, network permissions, `agent` install user and checksum
check. Unlike our hello kit, this installs a downloaded tool. The architecture and
checksum boilerplate is supplied; the educational change is adding a second mixin.

A [template](https://docs.docker.com/ai/sandboxes/customize/templates/) supplies the
base filesystem. A kit declares additions/setup/access. Here the built-in Claude
base already supplies the container-enabled environment; Pi is layered onto it.
You do not need Docker on the host or an image build during this workshop. A
prebuilt template is a later way to reduce repeated download time.

## 2. Extend your recipe

```bash
# HOST
mkdir -p "$WORKSHOP/chapters/my-03"
cd "$WORKSHOP/chapters/my-03"
cat ../my-025/sbxenv.yaml > sbxenv.yaml
cat >> sbxenv.yaml <<'YAML'
  - source: ../kits/pi
YAML
for f in chapter.env PROMPT.md launch; do cat "../03-pi/$f" > "$f"; done
chmod +x launch
sbx env plan sbxenv.yaml --env-arg name=wad-ch-03 \
  --env-arg port=3104 --env-arg "control_dir=$CONTROL"
```

The new line is another entry in the existing `kits` list, not a second `kits`
section. Paths are relative to the recipe. Kits and the app are assembled into a
fresh sandbox; we are not assuming a running sandbox hot-reloads edited YAML.

```bash
# HOST — terminal A, in chapters/my-03
WORKSHOP_APP_REPO="$WARMUP" WORKSHOP_APP_REF=HEAD ./launch wad-ch-03
```

## 3. Prove an actual Pi model call

```bash
# HOST — terminal B
sbx exec wad-ch-03 pi --version
sbx exec -it -w /home/agent/work/app wad-ch-03 \
  env -u GH_TOKEN -u GITHUB_TOKEN ANTHROPIC_API_KEY=proxy-managed \
  pi --provider anthropic --model claude-sonnet-5 --approve
```

`proxy-managed` is a placeholder used with SBX's configured credential proxy, not
an API key. Keep that environment override local to this Pi command: Claude Code
uses its own credential setup. Real credentials stay in SBX's credential mechanism.

Ask Pi:

> Read AGENTS.md and
> .claude/skills/acr__shelajev__coding-policy__review-change/SKILL.md.
> Summarize the guidance, then review the warm-up against ~/work/task.json without
> editing files. Do not run browser tests or install a browser.

Expected: an actual model response using the same policy as Claude. An installed
binary or an authentication error is not a successful model call. If the selected
model is unavailable, choose an accessible model and record it in the team config.
If Pi lacks provider access, resolve that before starting three headless sessions;
a working Claude subscription alone is insufficient evidence.

## 4. Configure the team on paper before launching it

```bash
# HOST — after exiting Pi
cat "$WORKSHOP/chapters/04-team/team.tsv"
cat "$WORKSHOP/chapters/04-team/team-mixed.tsv"
```

| Column | Meaning | Example |
|---|---|---|
| role | Responsibility | coordinator, developer, QA |
| harness | Assistant process Herdr launches | pi, claude, codex |
| provider | Model service | anthropic, google, openai |
| model | Specific model used by the role | an available ID for that provider |

Everyone will configure all roles in chapter 04. The baseline uses one provider
and model; you can still learn the separation. Independent providers for development
and review are the intended richer setup when accounts permit it. Follow the
[mixed-model instructions](MIXED-MODELS.md) for the presenter/credential-equipped
variation. Its custom-base launcher is currently a chapter-04 variation; the main
MCP path keeps the tested gateway-capable Claude base.

**Result:** both assistants work inside SBX and can use the policy, but you are
still manually routing work. Exit Pi, end the launch hold and `sbx stop wad-ch-03`.

Next: [Herdr and handoffs](../04-team/README.md).
