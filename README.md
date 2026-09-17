# Build a software factory with Docker Sandboxes

Give a team of coding agents a task, watch them implement and review it, and step
in when they need a product decision. In this two-hour workshop, you will build
that workflow yourself, starting with one agent running inside Docker Sandboxes.

## What you will build

Your factory will work on an incident triage board: a web application with an API
and a PostgreSQL database. Agents will add features to the board, run commands and
database containers inside their sandbox, review each other's changes, and write
results back to your task backlog.

![The incident triage board you will extend](chapters/images/incident-triage-board.png)

You will leave with sandbox recipes, reusable kits, coding guidance and a team
configuration you can adapt to your own project.

## What you will learn

Each chapter adds something the factory needs:

| Step | What you will learn |
|---|---|
| Run one agent | Use SBX interactively, explore file isolation and shared files, run containers, and open the app through a published port. |
| Repeat the setup | Describe an environment with sbxenv and connect it to a Beans task with a small host script. |
| Share tools and guidance | Build a simple kit, then use the ACR kit to install a coding policy and review skill. |
| Choose your assistants | Add Pi through a kit and configure the assistant, provider and model for each role. |
| Form a team | Use Herdr to manage sessions and pass assignments between a coordinator, developer and reviewer. |
| Connect to host tools | Give the team scoped access to Beans through the MCP gateway and explore network controls. |
| Bring in a human | Join through SSH, answer a product question and let the team continue. |
| Run another job | Reuse your factory, try the changed app and bring the source back to your laptop. |

The workshop closes with presenter demonstrations of live mounts, cloud sandboxes
and organization governance.

## About the tool choices

Docker Sandboxes provides the environment, kits, network controls and MCP gateway
used in this workshop. The third-party tools—including Beans, ACR, Herdr, Pi,
Claude Code and Codex—and the coding policy and team workflow reflect the author's
personal choices for this example. Their inclusion does not imply endorsement or
recommendation by Docker. You can apply the same Docker capabilities with tools
and conventions that suit your own team.

## Get started

Bring a **Mac with Apple silicon**, a terminal, a browser and model access. The
first exercise uses Claude Code with your Claude account. Later, you will configure
Pi's provider access; [chapter 00](chapters/00-setup/README.md#2-have-an-agent-account-ready)
explains the account requirements. Docker containers run inside SBX, so you do not
need Docker Desktop or a host Docker engine.

This repository and its downloadable materials currently require GitHub access.
Installation and authentication are part of the workshop.

**[Start with chapter 00: setup →](chapters/00-setup/README.md)**

Or browse the [chapter guide](chapters/README.md) to see the whole journey.

## How to use this repository

Follow each chapter from the top. You will build your own configuration in
`chapters/my-*`; the numbered chapter directories contain completed reference
configurations. Supplied app checkpoints let you catch up and keep exploring.

| Directory | What is inside |
|---|---|
| `chapters/00-setup` through `07-factory` | Walkthroughs and completed reference configurations |
| `chapters/examples` and `chapters/kits` | Small examples and reusable kits |
| `chapters/support` | Launcher, sandbox helpers and agent role instructions |
| `scripts` | Material downloads and host setup helpers |
| `backlog/seed` | Tasks for the agents to work on |

## License

Workshop source is licensed under [Apache 2.0](LICENSE). Downloaded third-party
components retain their own licenses and notices.
