# Build a software factory with Docker Sandboxes

A two-hour workshop: start one agent in an isolated environment, then add
repeatable recipes, shared skills, another assistant, Herdr coordination, host task
access through MCP, and human intervention through SSH.

**Start with [chapter 00](chapters/00-setup/README.md).**
Read chapters from the top; each introduces its capability before assembling the
finished version. [All chapters](chapters/README.md) ·
[Presenter extras](chapters/08-presenter/README.md)

This is the private workshop edition. The original development repository retains
tests, historical reports and the larger reference implementation. Here:

| Directory | Purpose |
|---|---|
| `chapters/00-setup` through `07-factory` | Instructions and completed reference states |
| `chapters/examples`, `chapters/kits` | Small teaching examples and mixins |
| `chapters/support` | Shared launcher, sandbox helpers and role briefs |
| `scripts` | Download materials, install host tools and seed Beans |
| `backlog/seed` | Four tasks used by the workshop |
| `kits/multi-provider` | Mixed-provider team variation |

The separate application and prebuilt MCP server download from the private
[materials release](https://github.com/shelajev/wad-sbx-workshop/releases/tag/materials-v0.1.0)
into ignored directories. No host Docker engine, Docker Desktop or Go build is
needed. GitHub CLI access to this private repository is needed to download them.

Work in `chapters/my-*` directories. Generated state and model output stay under
`.local/`; neither belongs in Git. The factory runs app code inside SBX. The host
launcher starts the environment; it does not supervise, independently verify,
merge, or automatically close tasks.

## Rehearsal status

The source chapter assemblies were exercised on macOS/Apple silicon with
SBX v0.45.0-rc2. The new instructions and this trimmed distribution still need a
full learner rehearsal. Other host platforms are not advertised as rehearsed.
Gemini access currently has a known account quota limitation. Cloud and organization
governance remain presenter-only extensions.

This edition derives from workshop commit `4c5e5f3`; its app bundle preserves app
commit `e15230f` and fixture tags. The materials manifest identifies the prebuilt
adapter's build provenance. New packaging does not imply a new adapter build.

Local checks: `chapters/check` validates recipes against the installed SBX CLI and
checks fixture refs after materials download. It does not run a team or model calls.
