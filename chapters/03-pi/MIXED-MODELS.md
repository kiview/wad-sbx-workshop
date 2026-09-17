# Choose a different model for each responsibility

A developer and a reviewer can be separate conversations with the same model. With
more provider access, they can also use different models. The configuration should
let us make that choice without rewriting how work moves through the team.

This variation belongs alongside chapter 04. If you have one provider, still read
and edit the role table, then keep the working choices from the main exercise.
The presenter can demonstrate the broader combination.

## Read the proposed team

Open `chapters/04-team/team-mixed.tsv`:

| Role | Assistant | Provider |
|---|---|---|
| Coordinator | Pi | Google |
| Developer | Claude Code | Anthropic |
| QA | Codex | OpenAI |

The last column of each file row supplies the model ID. Choose models your accounts
can actually use. Changing a row affects a newly started team, not a session already
running. Consider the distinction: replacing QA's model should not change how the
coordinator sends it a review request.

## See what the environment must supply

Open `chapters/kits/multi-provider/spec.yaml`. This is a **sandbox kit**, not a mixin:
it declares the base environment and provider credential bindings, and installs
Codex. Pi and Herdr are still separate mixins layered on top.

Use the instructor's configured SBX provider services for this variation. The kit
expects `anthropic`, `openai` and `gemini`; an existing `google` service can be selected
with its `google_secret_service` kit argument. Real keys and OAuth tokens do not
belong in the TSV file. These bindings are why selecting a model involves more than
just changing its name in an assistant command.

## Start the variation and explore it

From the host:

```bash
cd "$WORKSHOP/chapters/04-team"
./launch-mixed my-mixed-team
```

Open the short `launch-mixed` script first if you want to see its choices. It selects
the sandbox kit, this role table and browser port 3205, then uses the shared launcher.
For this RC, custom sandbox kits use the `sbx create KIT --kit MIXIN` path: environment
planning accepts them, but `sbx env create` did not create that custom agent in our
rehearsal. This is why the variation uses a different creation command.

Leave the launcher open. In another host terminal:

```bash
sbx exec -it my-mixed-team bash
```

Inside the sandbox:

```bash
export FACTORY_DIR="$HOME/work/factory"
herdr agent list
herdr agent read coordinator --lines 30
herdr agent read developer --lines 30
herdr agent read qa --lines 30
```

Read each role's introduction. You are looking at three assistant sessions using
the role choices from the table. Repeat the message-and-wakeup exercise in chapter
04: the communication mechanism stays the same even though the models differ.

If Gemini reports a quota limit, use the working one-provider team for the exercise;
the current presenter Google account has previously hit that limit. Changing the
handoff protocol cannot provide more account quota.

Exit the shell, end the launch terminal and stop `my-mixed-team` when done. Return
to the normal recipe for chapter 05: the custom-base variation still needs additional
gateway-client configuration before it can replace that chapter's tested MCP route.
