# 6. Join the team when a requirement needs a human

What should happen when the task leaves a product decision unanswered? For our next
feature, the team must ask what to do with an old resolution note when reopening an
incident. We will use SSH to enter the existing sandbox, read the question and
record our answer through the same file-message convention. The team can then
continue with a requirement we supplied, rather than guessing one.

## 1. Set up a task with a real choice

On the host, carry your recipe and team forward:

```bash
mkdir -p "$WORKSHOP/chapters/my-06"
cd "$WORKSHOP/chapters/my-06"
cat ../my-05/sbxenv.yaml > sbxenv.yaml
cat ../my-05/team.tsv > team.tsv
for file in chapter.env PROMPT.md launch; do cat "../06-human/$file" > "$file"; done
chmod +x launch
```

The copies keep the infrastructure and provider choices you assembled. Open
`chapter.env`: the task is now `wad-103`. Open `PROMPT.md`: it asks the coordinator
to bring the ambiguous requirement to you **before** implementation.

The question is what to do with the old resolution note when reopening an incident.
Keeping it records the previous conclusion; clearing it avoids presenting an open
incident as resolved. History can still preserve the old note. Both readings are
plausible, so the team must ask.

For this exercise, start from the supplied app checkpoint that already implements
assignment and resolution:

```bash
# HOST — terminal A, in chapters/my-06
./launch wad-ch-06
```

This is a deliberate teaching checkpoint, not your earlier agents' output. If you
saved a completed chapter-05 result and want to continue it instead, use
`./launch wad-ch-06 "$WORKSHOP/.local/feature-app"` **instead of** that command.
Leave the launch terminal open.

## 2. Open the SSH doorway

In terminal B on the host:

```bash
sbx setup ssh
ssh wad-ch-06.sbx
```

The first command prepares the host's SSH integration for SBX. The second opens a
shell in this named sandbox using that integration. You do not need to find a VM
IP address or set up an SSH server in the app. Read more in
[SBX integrations](https://docs.docker.com/ai/sandboxes/integrations/).

Look around from the SSH shell:

```bash
# SANDBOX, over SSH
cd ~/work
export FACTORY_DIR="$HOME/work/factory"
herdr agent list
```

You should recognize the coordinator, developer and QA. This is the same environment
the launcher started, not a fresh copy. The export points our message tools at its
existing shared conversation directory.

## 3. Let the team ask its question

Once the roles have finished introducing themselves:

```bash
./bin/assign
```

As in chapter 05, this delivers the task to the coordinator. Now read its progress:

```bash
herdr agent read coordinator --lines 40
handoff inbox human --json
```

The first command shows what it is doing. The second shows messages addressed to
you without consuming them. Wait for the decision request; if it has not arrived,
read the conversation again after the agent's turn finishes.

Read both options before choosing. For this walkthrough we will **clear the current
resolution note and retain the historical entry**. You can choose the other behavior,
but the recorded decision and your reply must agree.

## 4. Record the choice so every role can refer to it

The task asks for a decision file as well as a reply. That makes the chosen behavior
available to both the implementer and reviewer after the conversation moves on.
First find the current attempt number:

```bash
attempt=$(handoff state | jq -r .attempt)
echo "$attempt"
```

`handoff state` describes this run. `jq` selects its attempt number; the shell keeps
it in `attempt`. You should see a number such as `1`. If you do not, stop and ask
the instructor to help inspect the run state before writing an answer.
We put it in the filename so a later attempt can have its own answer.

Write the decision in readable JSON:

```bash
mkdir -p "$FACTORY_DIR/decisions"
cat > "$FACTORY_DIR/decisions/wad-103-a${attempt}.json" <<EOF
{
  "task": "wad-103",
  "attempt": $attempt,
  "decision": "clear",
  "rationale": "An open incident has no current resolution; history keeps the old note."
}
EOF
```

`mkdir` prepares the decisions directory. The `cat` block writes everything between
it and `EOF` into the named file, substituting the attempt number. The meaningful
parts are the choice and its reason; the surrounding fields identify the task/run.

Now tell the coordinator where to find your answer:

```bash
handoff send --to coordinator --from human --kind decision --body \
  "Clear resolutionNote on reopen; history keeps the old note. The decision is in $FACTORY_DIR/decisions/wad-103-a${attempt}.json. Continue implementation and review."
crew-notify coordinator
```

This is the same leave-a-message/notify pair you used in chapter 04. The new message
kind is `decision`. You are responding to the existing team, not restarting it.

## 5. Watch the team use your answer

Read the coordinator's next turn and then the developer's:

```bash
herdr agent read coordinator --lines 40
herdr agent read developer --lines 40
```

Look for the chosen behavior being relayed and implemented. When QA reviews, it
should use the same choice. Open <http://127.0.0.1:3107> when the team refreshes the
app and try reopening a resolved incident. The important result is that your
answer changed what the team built.

Type `exit` to return to the host. The launch terminal still keeps SBX running.
When the team finishes and commits, save the app with:

```bash
# HOST
sbx cp wad-ch-06:/home/agent/work/app "$WORKSHOP/.local/reopen-app"
```

As in chapter 05, this retrieves the private source copy. End the launcher with
Ctrl-C, then `sbx stop wad-ch-06` when done. A product decision belongs in the team's
conversation; a network-policy change still belongs on the host.

Next: [use the factory for another job](../07-factory/README.md).
