# 0. Get ready to build a software factory

We start by installing Docker Sandboxes and preparing access to a coding agent. We
also download the workshop's sample application and demo coding tasks. The app
tracks service incidents and includes a browser UI, API and PostgreSQL database.
Our agents will extend it as we build the factory. By the end of setup, you will
have the sample source code, demo tasks and workshop tools on your laptop. Docker Desktop is not needed;
the database container will run inside a sandbox.

## 1. Install the host prerequisites

These instructions use macOS on Apple silicon.
Use the [standalone SBX installation instructions](https://docs.docker.com/ai/sandboxes/install/).
No host Docker engine or Docker Desktop is required. Docker containers will run
**inside** the sandbox.

Install the release-candidate channel and host utilities:

```bash
# HOST
brew install docker/tap/sbx@rc git jq coreutils gh
sbx version
sbx login
```

`brew install` installs the host tools: SBX for environments, Git for source,
jq for reading JSON, coreutils for checksums, and GitHub CLI for the private download.
`sbx version` tells you which CLI you are running. `sbx login` signs in to Docker;
Claude's model-account login happens separately in chapter 1.

If SBX is already installed, check which executable `command -v sbx` selects before
changing it. This material targets **v0.45.0-rc2**. A floating Homebrew RC channel
may now deliver a later release; use the assets and installation guidance on the
[pinned release](https://github.com/docker/sbx-releases/releases/tag/v0.45.0-rc2)
to install that version. Keep one SBX executable on your PATH.

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

Later, we will add a second coding assistant, Pi, which needs access to its chosen provider. We check that in chapter 03, when we
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
It downloads two things from our GitHub release: the sample application with saved
versions for later exercises, and a tool that will let agents access the workshop's
task tracker from inside a sandbox. We will connect that tool—the Beans MCP
server—in chapter 05.
It verifies their checksums. Nothing starts running yet.

```bash
# HOST — from the workshop repository
./scripts/get-materials.sh
```

The sample application lives in `.local/app`, and the task-access tool in `dist/`.
You do not need to build either tool or install Go.

### Prepare this terminal with one command

```bash
# HOST — run this in each new terminal, from the workshop repository
source ./scripts/workshop-env.sh
```

This sets the workshop paths and prepares your editable starter app in
`.local/warmup-app`. It also creates `host-only.txt` beside that directory for the
file-isolation exercise. It preserves an existing app working copy. We **source** the script
so its variables remain available in your current terminal. Later instructions use three path variables:

| Variable | Location it names |
|---|---|
| `$WORKSHOP` | This workshop repository. |
| `$WARMUP` | Your working copy of the sample app for the first coding exercise. |
| `$CONTROL` | The workshop tools and demo-task backlog on your host. |

The script supplies these paths so you do not need to retype them. In each new host
terminal, return to this repository and source the script again before using them.

## How to use the terminals and chapter directories

A code block says **HOST**, **CLAUDE**, or **SANDBOX**. Do not run application commands on the
host. Keep one terminal connected to each live sandbox; use another for host
commands. In this SBX build, background processes alone do not prevent idle stop.
The later launcher deliberately stays in the foreground for this reason.

At the start of a chapter, pause to name the problem the new component solves.
After an exercise, look at its effect before running the next command. Commands
marked CLAUDE or PI go into that assistant's interface; a sandbox shell is a normal
terminal inside the environment. `exit` leaves that shell and returns to the host.

You will create `chapters/my-*` directories to edit recipes. Keep them directly
under `chapters/`, because the supplied launch wrapper and relative kit paths
expect that layout. The numbered directories are completed reference assemblies.
Their `./launch` commands let you catch up using a supplied app fixture.

**Ready when:** SBX reports the intended version, Docker login works, your agent account is ready, and the app and its fixture tags are present.

Next: [run an agent manually](../01-agent/README.md).
