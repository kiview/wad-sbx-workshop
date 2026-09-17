# 4. Turn assistants into a team

**Goal:** make the agents coordinate inside SBX. Herdr owns their sessions; role
briefs define responsibilities; file messages carry durable assignments and reviews.
The host remains a launcher.

## 1. Add the orchestration kit

Read the [Herdr kit](../kits/herdr/spec.yaml), then extend your existing recipe:

```bash
# HOST
cat "$WORKSHOP/chapters/kits/herdr/spec.yaml"
sbx kit validate "$WORKSHOP/chapters/kits/herdr"
mkdir -p "$WORKSHOP/chapters/my-04"
cd "$WORKSHOP/chapters/my-04"
cat ../my-03/sbxenv.yaml > sbxenv.yaml
cat >> sbxenv.yaml <<'YAML'
  - source: ../kits/herdr
YAML
for f in chapter.env PROMPT.md launch team.tsv; do cat "../04-team/$f" > "$f"; done
chmod +x launch
cat team.tsv
cat ../support/roles/coordinator.md ../support/roles/developer.md ../support/roles/qa.md
cat ../support/bin/start-team
```

Edit `team.tsv` now if chapter 03 required a different available model. Preserve
its tab-separated columns. The baseline is Pi coordinator, Claude developer and
Claude QA. Do not swap the developer to a harness without gateway configuration:
that role will access MCP in the next chapter. The mixed-provider demonstration is
linked from chapter 03 and can repeat this chapter's short relay exercise.

`start-team` starts a headless Herdr server, creates sessions, launches the configured
harnesses and gives each its role. This supplied setup avoids making a terminal
supervisor the focus of the workshop. Read the short script before running it.
The three kits now contribute distinct capabilities: ACR, Pi and Herdr.

## 2. Start the three real sessions

```bash
# HOST — terminal A, in chapters/my-04
WORKSHOP_APP_REPO="$WARMUP" WORKSHOP_APP_REF=HEAD ./launch wad-ch-04
```

In terminal B, check <http://127.0.0.1:3105>. Wait for all three agents to acknowledge
their roles, then do the relay below. This is deliberately a small task: proving
communication should not require implementing a whole feature.

## Inspect and exercise coordination

```bash
# HOST — terminal B
cd "$WORKSHOP/chapters/my-04"
cat team.tsv
cat ../kits/herdr/spec.yaml
cat ../support/roles/*.md
sbx exec wad-ch-04 herdr agent list
sbx exec -it wad-ch-04 bash
```

Inside the sandbox, export the shared location and send a small relay exercise:

```bash
# SANDBOX
export FACTORY_DIR="$HOME/work/factory"
handoff send --to coordinator --from human --kind assignment --body \
  'Handoff exercise only: ask developer to read the app package.json and send its test script names back; ask QA to confirm those names. Report both replies to human. Do not implement the feature yet.'
crew-notify coordinator
handoff inbox human --json
herdr agent read coordinator --lines 30
```

Allow model turns to finish before inspecting the result. The expected evidence is
three real sessions plus messages from developer and QA, not merely a green Herdr
status. `handoff` writes durable files; `crew-notify` wakes the recipient through
Herdr. Codex status detection can be unreliable, so inspect files and terminal output.

To start the real feature using the pre-MCP task snapshot, run `assign` in the sandbox.
You can instead leave this as a quick coordination exercise and start the feature in
chapter 05. Do not wait for the whole feature twice.

If a notification reports that the recipient is busy, inspect its terminal, let it
finish, then retry `crew-notify ROLE`. If delivery was unconfirmed, inspect before
sending another prompt. The helper retains a claim to avoid duplicate prompts.

**Result:** a Pi coordinator routes work between developer and QA. The policy supplies
coding guidance; Herdr supplies sessions; files supply durable handoffs.
**Next limitation:** tasks are still copied into the sandbox, with no host write-back.


## Check the end state and move on

In the human inbox, find the developer's script names and QA's confirmation. If
needed, inspect each agent's terminal with `herdr agent read ROLE --lines 40`.
`handoff inbox` observes messages; `handoff read` consumes pending deliveries.
A file message and a wakeup are both needed: writing a file does not wake a model.
These are shell tools; a harness's unrelated built-in messaging feature will not
reach this team.

Do not start `assign` as well as the relay unless you want an extra feature run.
Chapter 05 is where we start the main assignment/resolution feature. Once the relay
is complete, leave the app unchanged and stop this sandbox after ending its hold:
`sbx stop wad-ch-04`.

Next: [MCP and access controls](../05-mcp/README.md).
