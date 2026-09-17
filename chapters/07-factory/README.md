# 7. Use the factory on your own project

So far, everyone has worked on the same sample application. Now bring a project
that interests you and give the team a small, real change to make. The factory
still needs the same three inputs: **source code, a task, and instructions for
working on that project**. The sandbox, kits, model choices, MCP connection and
handoff workflow can stay as they are.

This is open-lab time. You could try a small change in Testcontainers or another
open-source project you know, or use one of your own repositories. Choose a change
you can explain and check: a regression test and fix, clearer validation, or a
small CLI feature. Knowing what a good result looks like will help you assess the
team's work.

## 1. Bring a different repository

Save any work you want from the previous sandbox and stop it before continuing;
our environment file still publishes host port 3102.

From the workshop repository root, clone your chosen project into a new directory.
Replace the URL below with its clone URL:

```bash
# HOST — from the workshop repository
mkdir -p projects
git clone https://github.com/OWNER/REPOSITORY.git projects/my-project
```

Choose a project whose purpose you understand, so you can judge the result. The
agents will read its documentation and work out its dependencies, setup and tests.
A library may have no application server at all; its useful outcome is a reviewed
code change and passing tests.

The launcher will transfer the committed source at `HEAD` into `~/work/app` in a
new sandbox. That path is simply where our team expects the project. Your host
clone stays separate, and its Git history is not transferred. Commit any local
changes you want included before launching.

## 2. Write a task the team can verify

Create `factory/TASK.md` in your editor. Describe the change for someone seeing the
project for the first time. Use this outline, replacing each bracketed section:

```markdown
# [Short title of your change]

## Problem
[What happens now? Give a concrete example.]

## Wanted behavior
[What should happen after the change?]

## How to check it
[One or two observable examples or focused tests that demonstrate the result.]

## Scope
[What should the team leave alone?]
```

Add this task to the same workshop backlog the MCP server already exposes:

```bash
./.local/chapters/bin/beans \
  --config .local/chapters/beans/.beans.yml \
  --beans-path .local/chapters/beans/.beans \
  create "My project's first factory task" \
  --type task --status todo --body-file factory/TASK.md
```

`--body-file` supplies your requirements. Beans prints the new task's ID; keep it
for the next step. This remains a task in your local workshop backlog, not an
issue or pull request in the upstream project. The agents will read it through
the MCP connection you already configured.

## 3. Tell the factory which project workflow to use

The earlier chapters automatically installed the sample app's Node dependencies
and started its PostgreSQL database and web server. Another project needs its own
setup. Use these launcher settings and starting prompt:

```bash
cp chapters/07-factory/chapter.env factory/chapter.env
cp chapters/07-factory/PROMPT.md factory/PROMPT.md
```

Open `factory/chapter.env` and replace `TASK=replace-with-your-task-id` with the ID
Beans just printed. Notice the two settings that matter here:

- `MODE=mcp` keeps task retrieval and result notes connected to your host backlog.
- `PROJECT_SETUP=project` skips the sample app's dependency installation and server
  startup. The launcher supplies the source and starts the team; the agents follow
  this project's setup instructions.

Read `factory/PROMPT.md`. It asks the team to inspect the repository, determine its
dependencies, get it running, and establish how to test it before making changes.
The agents can install tools and run supporting containers inside SBX—for example,
Node and a PostgreSQL container if that is what the project needs. You do not need
to prescribe those choices yourself.

Add any project context the agents cannot discover from the source: a useful
reproduction, a convention your team follows, or something to leave alone. Your
Bean describes **what to change**; the prompt describes **how the team should
approach the project**. The supplied prompt is enough to start exploring.

For a library or CLI, the team should demonstrate the result with a focused test
or command. For a web app, the prompt asks it to start the server on
`0.0.0.0:8080`, reachable through the existing mapping at
<http://127.0.0.1:3102>. Publishing a port does not start a server.

Keep `factory/sbxenv.yaml` and `factory/team.tsv`. The same coding guidance, kits,
roles and MCP connection will be used. If a needed download is blocked, the team
reports it and you decide whether to allow it. Once you know which tools you want
on every run, you can package their installation in a kit. Repeated project
onboarding instructions can likewise become a skill distributed through ACR.

## 4. Launch the team with your source

```bash
# HOST — terminal A, from the workshop repository
./scripts/launch-factory.sh wad-my-project projects/my-project
```

The first argument names the new sandbox; the second selects your repository.
Underneath, the launcher still calls:

```bash
sbx env create factory/sbxenv.yaml --env-arg name=wad-my-project --auto-approve
```

Then it transfers your committed source, prompt and task ID, installs the coding
guidance and starts the roles. The launcher does that sequence for you; do not run
the second command again. Leave terminal A open. No sample-app database or web
server is started in project mode.

In another host terminal, join through SSH:

```bash
ssh wad-my-project.sbx
```

Inside the sandbox, inspect the setup and submit the task once the roles finish
their introductions:

```bash
# SANDBOX — over SSH
cd ~/work
cat task-id
cat PROMPT.md
herdr agent list
./bin/assign
```

The coordinator asks the developer to fetch your task through MCP. Follow the
conversation as the team learns the project, implements the change and reviews it:

```bash
herdr agent read coordinator --lines 40
herdr agent read developer --lines 40
herdr agent read qa --lines 40
```

If it needs a product decision, use the human handoff from chapter 06. If a network
request is blocked, inspect it and decide on the host whether to allow that
specific destination, as in chapter 05. The environment's access boundaries still
apply when the project changes.

## 5. Inspect and keep the result

Ask what evidence supports the result: which commit did QA review, which checks
ran, and what remains unverified? For a web app, try the changed behavior through
the published port. For a library or CLI, read the focused test results or run the
demonstration inside the sandbox.

When the team finishes, inspect its saved work:

```bash
# SANDBOX
cd ~/work/app
git status --short
git log -1 --oneline
```

The developer should also have appended a result note to your Bean. Exit SSH and
read it on the host, replacing `YOUR-TASK-ID` with your actual task ID:

```bash
# HOST — from the workshop repository
./.local/chapters/bin/beans \
  --config .local/chapters/beans/.beans.yml \
  --beans-path .local/chapters/beans/.beans show YOUR-TASK-ID

sbx cp wad-my-project:/home/agent/work/app ./.local/my-project-result
```

The copy keeps the result in a new host directory, leaving your original clone
unchanged. Review the diff against that clone before carrying changes into your
normal contribution workflow. This exercise does not publish or submit anything
to the upstream project.

End terminal A with Ctrl-C, then run `sbx stop wad-my-project`. You now have the
same factory working with your own source and requirements. For the next job,
change the task and choose the committed project version you want it to start from.

Next: [explore what else SBX can do](../08-presenter/README.md), with the presenter.
