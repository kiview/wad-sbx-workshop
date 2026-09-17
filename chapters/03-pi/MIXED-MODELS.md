# Demonstration: different models for different responsibilities

The role table separates responsibility from assistant, provider and model.
**Use `gemini-3.8-flash` for every Gemini role in this workshop.** The
presenter can show Pi using Gemini for coordination, Claude for development and
Codex for review. Attendees still configure the roles even if all their working
sessions use one provider.

## What changes, and what does not?

Compare `chapters/04-team/team.tsv` with `team-mixed.tsv`:

| Role | Assistant | Provider |
|---|---|---|
| Coordinator | Pi | Google |
| Developer | Claude Code | Anthropic |
| QA | Codex | OpenAI |

The fourth column selects each model. The role briefs and file-message protocol
remain the same. A different model does not require a new way to send a review request.

Open `chapters/kits/multi-provider/spec.yaml`. This is a **sandbox kit**, which
selects a starting image and declares provider credential bindings, rather than
just a mixin. Read it in these pieces:

1. `sandbox` selects the base image and entrypoint.
2. `credentials` describes each provider's request hosts and credential injection.
   API-key and OAuth routes differ: a subscription login must not be forced into
   API-key mode by setting a placeholder variable indiscriminately.
3. `permissions.network` allows the endpoints these assistants need.
4. `setup.install` installs Codex and prepares the assistant configuration appropriate
   to the selected authentication mode. Pi and Herdr remain separate mixins.

The real credentials live in SBX's host credential store, not the role table. The
presenter's setup needs usable `anthropic`, `openai` and `gemini` services; the kit's
`google_secret_service` argument can select an existing `google` service instead.

## Presenter walkthrough

Finish and remove the normal chapter sandbox first. This demonstration uses the
same mounted `sample-app/` and port 3102, so it runs on its own. In the SANDBOX tab:

```bash
./chapters/04-team/launch-mixed wad-mixed
```

The helper uses `sbx create` with the custom sandbox kit and the ACR, Pi and Herdr
mixins, then opens the same workshop shell. Try:

```bash
crew ask "Ask developer to list the project's test commands, ask QA to check the answer, then report to me. Do not change application code."
crew watch
```

Show `crew logs coordinator` and `crew logs qa` to identify the actual assistants.
The observable result is a cross-provider handoff, not merely three model names
in a file. If a provider lacks quota, explain the unavailable route and use the
working team configuration; do not present a failed call as a successful model swap.

Exit the shell when finished. In HOST:

```bash
sbx rm wad-mixed
```

## Carry the mixed team into the MCP exercise

The custom kit also configures Claude's MCP client. To demonstrate the chapter-05
connection with this same team, use the supplied MCP reference configuration:

```bash
# SANDBOX tab — before connecting, at the workshop root
./chapters/04-team/launch-mixed wad-mixed-mcp 05-mcp
```

The second argument chooses the chapter's completed settings and prompt. This
creates one new sandbox on the same `sample-app/`. Inside its shell, try a small
read-only request before a full feature:

```bash
crew ask "Ask developer to read wad-101 through MCP and list package.json test scripts. Ask QA to check the list independently. Report the combined answer; do not change project files or write a task note."
crew watch
```

The result should include a real host task read by Claude, an independent Codex
check, and Pi's combined answer. The gateway-capable role remains Claude; the other
roles communicate with it through the same file-message helpers.

This RC's `sbx env` accepts built-in agent names, not a custom sandbox kit as its
`agent` value. The presenter helper therefore uses `sbx create KIT --kit MIXIN` and
attaches the host server with `--static-mcp`. These are different ways to compose
the same SBX building blocks. The reference helper uses the ACR, Pi and Herdr kits;
it does not read extra kits from your evolving `factory/` configuration.

Exit the shell, then in HOST:

```bash
sbx rm wad-mixed-mcp
sbx mcp rm wad-mixed-mcp-beans
```

The main attendee path continues with `factory/` and the built-in Claude environment.
