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

These commands assume [Homebrew](https://brew.sh) is installed. Install the
release-candidate channel and host utilities:

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
to install that version. Use `sbx version` again afterward; the version for this workshop is `v0.45.0-rc2`.

You also need a browser, SSH, a Bash-compatible terminal, internet access and enough
available memory for a 4-CPU/8-GB sandbox. Run one main chapter sandbox at a time
when resources are limited. Git, jq and sha256sum must be on the host PATH:

```bash
# HOST
command -v sbx git jq sha256sum ssh
```

You should see a path for each of the five tools above. A missing path means that
tool is not available in this terminal; fix its installation before continuing.

## 2. Have an agent account ready

For chapter 1, use your Claude subscription: we will sign in with `/login` inside
Claude Code. You do not need to create an API key for that chapter. If SBX already
has an Anthropic API credential configured, it may use that instead.

In chapter 03 you will also try Pi, an open-source assistant. The main Pi example
uses an **Anthropic API key**, configured on the host through SBX. A Claude
subscription alone does not supply that API access. If you only have a Claude
subscription, you can build all three team roles with Claude Code; chapter 04
provides that configuration. You can still install Pi and explore its interface,
or pair with someone who has provider access for the Pi conversation. Chapter 03
marks the model exercises to follow with a partner or presenter, then brings
everyone back together for team setup. You will still build the kits and configure
each role yourself.

With additional accounts, you can choose different models for development and
review. Configure the roles even if you initially give them the same model.
Never put actual keys or OAuth tokens in the workshop's files.

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

## Set up your two terminals

Open two tabs or windows, each at the workshop repository root. Name them if your
terminal supports it:

- **HOST** stays on your laptop. Use it for SBX controls, task tracking and editing
  the workshop configuration.
- **SANDBOX** starts as another host shell at the same repository root. You run
  `sbx run` or our launcher there; it becomes your connection to an agent or a
  shell inside SBX. Keep it open while work is running. Exiting returns to the host.

Code blocks say **HOST**, **SANDBOX tab — before connecting**, **CLAUDE**, **PI**,
or **SANDBOX shell**. A prompt shown as a quotation is something to say to the
assistant, not a shell command. Edit configuration files with your normal editor.

We use one sandbox at a time and one application directory, `sample-app/`, through
chapters 01–06. The directory is mounted read/write: edits, commits and deletions
inside it are visible on your host. The rest of this workshop repository—including
the host task backlog—is outside that mount. Run application code and containers
inside SBX, not on the host.

A chapter's sandbox is temporary; your source and its Git history remain. We will
remove each sandbox before creating the next one with its new capabilities.

You are ready when `sbx version` reports the intended version, Docker login is
complete, and you can open `sample-app/README.md` in your editor.

Next: [run one agent](../01-agent/README.md).
