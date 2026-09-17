# 6. Enter the factory when the team needs a decision

**Goal:** resolve an ambiguous requirement through SSH and let the same live team
continue. More autonomous execution cannot supply a missing product decision.

## 1. Carry the infrastructure forward

```bash
# HOST
mkdir -p "$WORKSHOP/chapters/my-06"
cd "$WORKSHOP/chapters/my-06"
cat ../my-05/sbxenv.yaml > sbxenv.yaml
cat ../my-05/team.tsv > team.tsv
for f in chapter.env PROMPT.md launch; do cat "../06-human/$f" > "$f"; done
chmod +x launch
cat chapter.env
cat PROMPT.md
```

There is no new kit or gateway here. The new capability is your SSH connection and
a decision protocol. The task is `wad-103`: reopen a resolved incident.

For predictable classroom timing, the default below uses the supplied
`app-02-feature-solution` fixture. This is an **explicit app checkpoint**, not a
claim that your chapter-05 implementation was copied. It already has assignment
and resolution, so everyone reaches the same ambiguous reopening requirement.

```bash
# HOST — terminal A, in chapters/my-06; predictable fixture path
./launch wad-ch-06
```

If your chapter-05 result is complete and copied, you may use it instead of the
fixture. Choose this command **instead of** the preceding launch:

```bash
# HOST — alternative, same terminal A
WORKSHOP_APP_REPO="$WORKSHOP/.local/feature-app" WORKSHOP_APP_REF=HEAD \
  ./launch wad-ch-06
```

Leave the launcher running. In terminal B, open <http://127.0.0.1:3107>, wait for
role initialization, and continue below. SSH connects to this same live sandbox;
it does not create another environment. See
[SBX integrations](https://docs.docker.com/ai/sandboxes/integrations/).

## Ask for the missing decision

This fixture already has assignment/resolution. `wad-103` deliberately leaves one
question open: should reopening clear `resolutionNote` or keep it? Read `PROMPT.md`.
The coordinator must ask you rather than choose. Start the task:

```bash
# HOST
sbx exec wad-ch-06 /home/agent/work/bin/assign
sbx setup ssh
ssh wad-ch-06.sbx
```

Once connected, inspect the real sessions and request:

```bash
# SANDBOX over SSH
export FACTORY_DIR="$HOME/work/factory"
herdr agent list
handoff read human --json
herdr agent read coordinator --lines 40
```

Wait until the team has actually asked the question; inspect its terminal and check
again if no request is present yet. Then record the decision in the path required
by the task and send it to the coordinator:

```bash
# SANDBOX over SSH
attempt=$(handoff state | jq -r .attempt)
mkdir -p "$FACTORY_DIR/decisions"
jq -n --argjson attempt "$attempt" \
  '{task:"wad-103",attempt:$attempt,decision:"clear",rationale:"An open incident has no current resolution; history keeps the old note."}' \
  > "$FACTORY_DIR/decisions/wad-103-a${attempt}.json"
handoff send --to coordinator --from human --kind decision --body \
  "Clear resolutionNote on reopen. History preserves the old note. Decision: $FACTORY_DIR/decisions/wad-103-a${attempt}.json. Resume implementation and review."
crew-notify coordinator
```

You are entering the same live sandbox, not restarting the factory. Inspect messages,
let the team continue, and later inspect the board and Bean note. If a spontaneous
question takes too long, the prompt explicitly requests this decision first; label
that as a rehearsed exercise, not an unscripted discovery.

**Result:** a human requirement is recorded and the same team resumes. Access-policy
changes still belong on the host; SSH supplies a product decision, not extra authority.


## Finish the intervention

Wait until an actual decision request is present before sending the decision.
If `handoff read human` is empty because the model is still working, inspect its
terminal and check again later. Do not keep resending the assignment.

After the decision, inspect `herdr agent read coordinator --lines 40` and the
subsequent developer/QA messages. The proof is that the team acknowledges your
chosen behavior and resumes, not merely that SSH connected. Keep the SSH session
open while observing, then `exit` to the host. The launcher still holds the sandbox.

When the task finishes, inspect the `wad-103` Bean using the same host Beans command
as chapter 05, substituting the task ID. Commit/copy any result you want to keep
before ending the hold and stopping the sandbox. For example, copy to the new path
`$WORKSHOP/.local/reopen-app`.

**Expected end state:** a recorded human decision and a resumed team; the app's
reopen behavior follows that decision when implementation completes.

Next: [inspect and reuse the factory](../07-factory/README.md).
