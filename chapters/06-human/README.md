# 6. Join the team when it needs a human decision

Our team can implement a clear task. What should it do when two interpretations
are reasonable? We will ask it to add reopening a resolved incident, but leave one
product choice for us: should reopening clear the current resolution note?

The coordinator will ask. We will connect through SSH, answer in ordinary language,
and watch the same team continue. SSH provides access to the running environment;
our role instructions tell the team when to involve a human.

## 1. Give it a task with a real choice

Keep using `sample-app/`, including your completed chapter-05 feature. If you skipped
that feature, the [catch-up instructions](../README.md#catch-up) can install a completed
checkpoint before you continue.

In HOST, change only `TASK=wad-103` in `factory/chapter.env`. Keep `MODE=mcp`,
`USE_ACR=1` and `SESSION=shell`. Keep the existing prompt and append this paragraph
to `factory/PROMPT.md`:

```markdown
Before implementing reopening, ask the human whether to clear the current
resolutionNote or retain it. Explain both interpretations and wait for an answer.
The coordinator records the human's choice for the developer and QA, then resumes
implementation and review. Preserve the historical resolution entry either way.
```

The existing prompt reads the task ID supplied by `TASK`, so it now asks for
`wad-103`. The added paragraph tells the team which decision belongs to you.
The environment, tools and team are the same.
In the SANDBOX tab, at the workshop root:

```bash
./scripts/launch-factory.sh wad-ch-06
```

Then in its sandbox shell:

```bash
crew submit
crew watch
```

Expect a question describing both options. If you added the browser-access kit in
chapter 05, it is also installed in this newly created environment.

## 2. Reach the existing team through SSH

Leave the SANDBOX session open. In HOST:

```bash
sbx setup ssh
ssh wad-ch-06.sbx
```

`setup ssh` configures the host's SSH integration. The second command opens another
shell in this exact sandbox, not a new agent team. You do not need to find an IP
address or install an SSH server in the app. See [SBX integrations](https://docs.docker.com/ai/sandboxes/integrations/).

Your HOST tab is temporarily an SSH session. Inside it:

```bash
crew status
```

You should see the same task and product question. The message tools know where
the current team's files live; reconnecting does not require exports or new setup.

## 3. Give a product answer

For this walkthrough, clear the current note while retaining history:

```bash
crew reply "Clear resolutionNote when reopening. An open incident has no current resolution, but retain the old note in history. Continue implementation and review using that decision."
crew watch
```

You supply the decision and its reason. The helper delivers it to the coordinator;
the coordinator records the task-specific decision and passes it to the developer
and QA. You do not need to assemble a JSON message or know a run's attempt number.

Look for the choice being used in implementation and review. When the team reports
that the app is ready, open <http://127.0.0.1:3102> and reopen a resolved incident.
Does the current note clear while history retains the old resolution?

Ctrl-C leaves the message view. Type `exit` to close SSH and restore your HOST tab.
Your original SANDBOX tab remains connected; SSH was another doorway into the
same running team. The code changes are already in `sample-app/`.

## 4. Finish this environment

Once the team is finished, leave its message view and type `exit` in the original
SANDBOX tab. In HOST:

```bash
sbx env rm factory/sbxenv.yaml --env-arg name=wad-ch-06
sbx mcp rm wad-ch-06-beans
```

The second command removes this sandbox's host MCP registration.

A product decision can travel through the team's conversation. An access-policy
change still requires the operator's host controls. You have now used both.

Next: [use the factory on another project](../07-factory/README.md).
