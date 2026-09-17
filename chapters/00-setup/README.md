# 0. Get ready to build a software factory

**Goal:** have standalone Docker Sandboxes and a usable coding-agent account, then
download the sample application. We install the other tools when we introduce
them. This workshop assembles Docker technologies; it does not install a single
product called “Docker's Agentic Platform.”

The application is an incident triage board: a browser UI, an API and PostgreSQL.
Our agents build it. They are not incident-response agents operating production.
The factory ends with a team implementing a task, reviewing it and writing a note
back to the host backlog.

## 1. Install the host prerequisites

The rehearsed host is macOS on Apple silicon. Other platforms need a separate
rehearsal; a downloadable SBX binary alone does not verify these workshop scripts.
Use the [standalone SBX installation instructions](https://docs.docker.com/ai/sandboxes/install/).
No host Docker engine or Docker Desktop is required. Docker containers will run
**inside** the sandbox.

For the macOS rehearsal, install the release-candidate channel and host utilities:

```bash
# HOST
brew install docker/tap/sbx@rc git jq coreutils gh
sbx version
sbx login
```

If SBX is already installed, check which executable `command -v sbx` selects before
changing it. This material targets **v0.45.0-rc2**. A floating Homebrew RC channel
may now deliver a later release; use the assets and installation guidance on the
[pinned release](https://github.com/docker/sbx-releases/releases/tag/v0.45.0-rc2)
when reproducing this rehearsal. Do not install two conflicting CLI versions.

You also need a browser, SSH, a Bash-compatible terminal, internet access and enough
available memory for a 4-CPU/8-GB sandbox. Run one main chapter sandbox at a time
when resources are limited. Git, jq and sha256sum must be on the host PATH:

```bash
# HOST
command -v sbx git jq sha256sum ssh
```

## 2. Have an agent account ready

For chapter 1, use your Claude subscription: we will sign in with `/login` inside
Claude Code. You do not need to create an API key for that chapter. If SBX already
has an Anthropic API credential configured, it may use that instead.

Later, Pi needs access to its chosen provider. We check that in chapter 03, when we
introduce Pi. Having a Claude subscription does not automatically give another
assistant Anthropic API access. The team configuration supports separate providers,
with a one-provider route when that is what you have available.

## 3. Clone this private workshop and get its materials

While the workshop is private, use GitHub CLI with an account that can access it.
Run `gh auth login` if needed, then clone:

```bash
# HOST
gh repo clone shelajev/wad-sbx-workshop
cd wad-sbx-workshop
```

### Download the application and the prebuilt tool

`get-materials.sh` is **a helper supplied by this workshop**, not an SBX command.
It downloads two things from our GitHub release: the incident-board application
with its chapter checkpoints, and the prebuilt Beans MCP server used in chapter 05.
It verifies their checksums. Nothing starts running yet.

```bash
# HOST — from the workshop repository
./scripts/get-materials.sh
```

The application lives in `.local/app`, and the MCP download in `dist/`.
You do not need to build either tool or install Go.

### Prepare this terminal with one command

```bash
# HOST — run this in each new terminal, from the workshop repository
source ./scripts/workshop-env.sh
```

This sets the workshop paths and prepares your editable starter app in
`.local/warmup-app`. It preserves an existing working copy. We **source** the script
so its variables remain available in your current terminal. Later instructions use
`$WORKSHOP` for this repository and `$WARMUP` for your working app; you do not need
to type their full paths or copy a list of exports.

## How to use the terminals and chapter directories

A code block says **HOST**, **CLAUDE**, or **SANDBOX**. Do not run application commands on the
host. Keep one terminal connected to each live sandbox; use another for host
commands. In this SBX build, background processes alone do not prevent idle stop.
The later launcher deliberately stays in the foreground for this reason.

You will create `chapters/my-*` directories to edit recipes. Keep them directly
under `chapters/`, because the supplied launch wrapper and relative kit paths
expect that layout. The numbered directories are completed reference assemblies.
Their `./launch` commands let you catch up using a supplied app fixture.

**Ready when:** SBX reports the intended version, Docker login works, your agent account is ready, and the app and its fixture tags are present.

Next: [run an agent manually](../01-agent/README.md).
