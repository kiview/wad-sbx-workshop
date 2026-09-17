# 4. Let the assistants hand work to each other

We have two assistants, but someone still has to give them roles and pass work
between them. We will add [Herdr](https://github.com/herdrdev/herdr) through a kit
to manage their sessions. Our file-message helpers preserve assignments and replies;
a separate notification wakes the recipient. You will start the team manually,
send a small request and follow its handoff before letting the launcher start
future teams for you.

### Another way to build a team

[Docker Agent](https://docs.docker.com/ai/docker-agent/) is Docker's open-source
framework for defining agent teams in YAML. Each agent can have its own model,
instructions and tools, and delegate work to other agents. It is available as a
standalone binary, so Docker Desktop is not required.

You could build an SBX-based team around Docker Agent by packaging it in a kit and
configuring its tools and provider access. Here we use Herdr to coordinate Claude
Code and Pi sessions, following the author's own setup. The focus is the same:
give your chosen team an environment, shared guidance and controlled access to tools.

## 1. Give the environment a session manager

Open `chapters/kits/herdr/spec.yaml`. Like the Pi kit, it records a pinned executable,
its installation and download access. This mixin supplies Herdr; it does not itself
decide what the agents should do.

Prepare the next working directory on the host:

```bash
mkdir -p "./chapters/my-04"
cat chapters/my-03/sbxenv.yaml > chapters/my-04/sbxenv.yaml
for file in chapter.env PROMPT.md launch team.tsv; do cat "chapters/04-team/$file" > "chapters/my-04/$file"; done
chmod +x chapters/my-04/launch
```

You are carrying forward your recipe and adding this chapter's settings and team
configuration. In your editor, append to the existing `kits` list:

```yaml
  - source: ../kits/herdr
```

Your recipe now combines ACR for guidance, Pi for another assistant and Herdr for
sessions. Before starting it, edit `chapters/my-04/chapter.env`: change `MODE=team` to
`MODE=manual`. That asks our preparation helper to install everything and start the
app, while leaving **you** to start the team in this first exercise.

## 2. Decide who should do what

Open `chapters/my-04/team.tsv`. It has one tab-separated row for each role:

| Role | Responsibility | Initial harness |
|---|---|---|
| coordinator | Route work and questions; wait for replies | Pi |
| developer | Change the app and report what was tested | Claude Code |
| qa | Review the developer's change and request corrections | Claude Code |

Put the working provider/model choice from chapter 03 into the Pi row. The default
uses Anthropic for all roles. The [mixed-provider variation](../03-pi/MIXED-MODELS.md)
lets you use separate providers when available.

Now read the short files in `chapters/support/roles/`. Notice that the coordinator
is told not to edit app code. QA reviews a particular commit, so it knows which
version the developer is asking it to review. These responsibilities come from
instructions, not from the model's name.

Keep Claude as the developer for the main path: in the next chapter that role will
use the configured MCP gateway client.

## 3. Start the environment, then start the team

```bash
# HOST — terminal A, from the workshop repository
./chapters/my-04/launch wad-ch-04 "$(pwd)/sample-app"
```

The launcher uses your updated SBX recipe and starts the app as before. In terminal B:

```bash
# HOST
sbx exec -it wad-ch-04 bash
```

You now have a shell inside SBX. Set the message directory and look at the supplied
team-startup helper:

```bash
# SANDBOX shell
cd ~/work
export FACTORY_DIR="$HOME/work/factory"
cat bin/start-team
```

The export tells the message tools where this team's shared files live. Preparation
also made `handoff` and `crew-notify` available on the sandbox's command path, so we
can use them by name. In the
script, find `herdr server`, `herdr agent start` and `herdr agent prompt`. They start
the session service, launch each configured harness, and deliver its role brief.
The remaining lines connect those operations to the rows of `team.tsv`.

Run the helper from the directory you are in:

```bash
./bin/start-team
```

Then look at the sessions:

```bash
herdr agent list
herdr agent read coordinator --lines 30
```

`list` shows who is present. `read` shows the end of an agent's actual terminal, so
you can see its role acknowledgment. Look at `developer` and `qa` the same way.
Wait for their introductions before giving them work.

## 4. Leave a message, then wake its recipient

We will ask the team a small question before entrusting it with a feature. In your
sandbox shell, leave this message for the coordinator:

```bash
handoff send --to coordinator --from human --kind assignment --body \
  'Ask developer to list the test scripts in package.json. Ask QA to check that list. Send me the combined answer. Do not change the app.'
```

The options say who receives it, who sent it, what kind of message it is, and its
text. This command writes a message file. See it waiting:

```bash
handoff inbox coordinator --json
```

The inbox is a view of the stored messages. The agent has not been interrupted just
because a file appeared. Now notify it through Herdr:

```bash
crew-notify coordinator
```

That is the second half of delivery: give the existing session a turn so it can
read its inbox. The coordinator will use these same helpers to contact developer
and QA. Herdr handles the sessions; files carry the conversation.

## 5. Follow the conversation

Take a moment to read what the coordinator is doing:

```bash
herdr agent read coordinator --lines 40
```

Try the same command with `developer` and `qa`. You should be able to follow the
request, the developer's answer, and QA checking it. When they have replied:

```bash
handoff inbox human --json
```

Read the result. Did the coordinator combine the answers? Did QA actually inspect
the app? If a session is still working, give it time. A busy status alone tells you
less than its conversation; read the messages to follow the work.

If a wakeup says the agent is busy, let it finish before retrying. If delivery is
uncertain, read the terminal before sending anything again.

## 6. Make this repeatable

Exit the sandbox shell. On the host, edit your `chapters/my-04/chapter.env` back to
`MODE=team`. On a future launch, preparation will call `start-team` for you. You
have just performed the step that is being automated.

We have a team that can discuss work. It still receives a copied task and cannot
write a result back to the host backlog. Next we will connect that boundary with MCP.
End the launcher with Ctrl-C, then run `sbx stop wad-ch-04` on the host.

Next: [connect the host backlog](../05-mcp/README.md).
