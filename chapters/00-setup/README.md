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

## 2. Check model access before starting the team

The supplied main path uses Claude Code plus Pi on Anthropic. Have a working
Anthropic credential for Pi as well as access to Claude Code. A Claude subscription
login does not by itself prove that Pi can make Anthropic API calls. Chapter 03
checks an actual Pi response before you build a team.

Configure credentials using SBX's agent login flow or its interactive secret store;
never put a key into a recipe, role brief or committed file. For an API credential:

```bash
# HOST — enter the value at the prompt, not as a command argument
sbx secret set anthropic
```

If your only provider is different, use the instructor's tested profile for that
provider. The supplied baseline is not a promise that replacing a model name makes
every harness support every login method. Everyone will configure independent roles;
multiple paid providers are not required. Google/Gemini and Codex are demonstrated
in the [mixed-model variation](../03-pi/MIXED-MODELS.md). A quota error is an account
limit, not evidence of a broken handoff protocol.

## 3. Clone this private workshop and get its materials

GitHub CLI access to the private repository is needed during rehearsal. Authenticate
with `gh auth login` if needed, then:

```bash
# HOST
gh repo clone shelajev/wad-sbx-workshop
cd wad-sbx-workshop
./scripts/get-materials.sh
```

The download contains a Git bundle of the separate sample application (including
fixture tags) and a prebuilt MCP adapter. Checksums are verified before installation.
The app is cloned into `.local/app`; release binaries go into `dist/`. Both are
ignored working data. You do not need Go, host Docker or another visible source repo.

In **each host terminal**, enter the workshop repository and set these variables:

```bash
# HOST — replace this one path with your checkout
cd /path/to/wad-sbx-workshop
export WORKSHOP="$PWD"
export APP_REPO="$WORKSHOP/.local/app"
export CONTROL="$WORKSHOP/.local/chapters"
export FACTORY_CONTROL_DIR="$CONTROL"
export WARMUP="$WORKSHOP/.local/warmup-app"
mkdir -p "$CONTROL"
git -C "$APP_REPO" tag --list 'app-*'
```

You should see `app-00-starter`, `app-01-warmup-solution` and
`app-02-feature-solution` among the fixture tags. These are teaching checkpoints,
not claimed output from your own agents.

## How to use the terminals and chapter directories

A code block says **HOST** or **SANDBOX**. Do not run application commands on the
host. Keep one terminal connected to each live sandbox; use another for host
commands. In this SBX build, background processes alone do not prevent idle stop.
The later launcher deliberately stays in the foreground for this reason.

You will create `chapters/my-*` directories to edit recipes. Keep them directly
under `chapters/`, because the supplied launch wrapper and relative kit paths
expect that layout. The numbered directories are completed reference assemblies.
Their `./launch` commands let you catch up using a supplied app fixture.

**Ready when:** SBX reports the intended version, Docker login works, a provider
credential is available, and the app and its fixture tags are present.

Next: [run an agent manually](../01-agent/README.md).
