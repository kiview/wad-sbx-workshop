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
brew install docker/tap/sbx@rc git jq coreutils
sbx version
sbx login
```

`brew install` installs the host tools: SBX for environments, Git for source,
jq for reading JSON, and coreutils for checksums. Downloads use `curl`.
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

## 3. Clone the workshop and get its materials

Download the workshop repository, then enter its directory:

```bash
# HOST
git clone https://github.com/shelajev/wad-sbx-workshop.git
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

The application checkpoints are stored in `.local/app`. The prebuilt task-access
tool is downloaded to `dist/`; we will install it for the factory in chapter 05.

### Find your workshop files

The download script also prepares `sample-app/`, your working copy for the coding
exercises. It keeps an existing working copy if you run the download again.

**Run host commands from the workshop repository root.** The sample application's
code is in `sample-app/`. In each new host terminal, open the workshop repository.

## Where to enter commands

Each command block tells you where it belongs:

- **HOST**: your laptop's terminal, in the workshop repository.
- **CLAUDE** or **PI**: the coding assistant's input.
- **SANDBOX**: a shell inside the sandbox. Type `exit` to return to the host.

The instructions will tell you when to open another terminal and which files to
create or edit. Application commands run inside the sandbox.

You are ready to start when SBX is installed, you have signed in to Docker, and
`sample-app/` contains the downloaded application.

Next: [run an agent manually](../01-agent/README.md).
