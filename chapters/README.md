# From one agent to a software factory

Build one environment a capability at a time. Keep the same application in
`sample-app/` and the evolving configuration in `factory/`. Each chapter explains
what is missing, introduces a building block, and lets you use it.

| Chapter | What you add | What you can observe |
|---|---|---|
| [00: setup](00-setup/README.md) | SBX, accounts and workshop materials | Two terminals and the sample source on your laptop |
| [01: one agent](01-agent/README.md) | Isolated execution and a mounted project | An agent runs containers, changes code and opens the app through a port |
| [02: repeatable environment](02-launcher/README.md) | sbxenv and a Beans task | The same project and its task in a newly created sandbox |
| [02.5: shared guidance](02.5-acr/README.md) | A simple kit, then ACR | A reusable installation and a policy-backed review |
| [03: another assistant](03-pi/README.md) | Pi and provider configuration | Try assistants and models against the same code and guidance |
| [04: team](04-team/README.md) | Herdr, roles and file messages | A request passes from coordinator to developer to QA |
| [05: host tools](05-mcp/README.md) | MCP gateway and scoped access | Agents read a host task and append the reviewed result |
| [06: human intervention](06-human/README.md) | SSH and a product answer | The existing team resumes using your decision |
| [07: reuse](07-factory/README.md) | Another project and task | The same factory works in a different mounted repository |
| [Presenter extensions](08-presenter/README.md) | Runtime mounts, cloud and governance | Additional capabilities demonstrated by the presenter |

The third-party tools and workflow are the author's choices, not Docker
endorsements. [About the tool choices](../README.md#about-the-tool-choices).

## The working arrangement

**HOST** stays at the workshop repository root. **SANDBOX** starts there too, then
connects to the agent or shell. Keep that connection open while agents and services
run: a sandbox may stop after its last client disconnects. A shell, Claude session
or SSH connection provides that client; a separate keepalive terminal is unnecessary.

Use one workshop sandbox at a time. At the end of a chapter, leave its session and
remove its environment with the command shown. This stops its processes and removes
its private runtime files. Your mounted app and Git history remain on the host.
The next worker discovers dependencies and starts services inside its new sandbox.
Database-container data is temporary unless you explicitly arrange persistence;
our sample app can recreate its demo data from its setup instructions.

Edit host configuration with your normal editor. Run application code inside SBX.
If you need another diagnostic shell while an assistant is occupied, open a third
tab and use `sbx exec -it SANDBOX-NAME bash`, substituting the current name. That
shell is for a concrete investigation; it is not a permanent part of the workflow.

## Catch-up

The numbered directories contain completed reference configurations. A catch-up
command puts the chosen configuration in **the same `factory/` directory**, so the
next incremental chapter starts from the right state. Exit and remove the current
sandbox before changing configuration or replacing its mounted project.

For example, to start chapter 03 with the ACR/Pi environment already assembled:

```bash
# HOST
./scripts/use-chapter.sh 03-pi
```

This saves your previous configuration under `.local/` and updates `factory/`.
It leaves `sample-app/` unchanged. In your SANDBOX tab:

```bash
./scripts/launch-factory.sh wad-ch-03
```

Continue at chapter 03's provider setup. The finished configuration already enables
guidance installation and opens the shell.

If you also need a completed application checkpoint, select it explicitly. To join
the SSH chapter without doing the earlier feature:

```bash
# HOST
./scripts/use-chapter.sh 06-human app-02-feature-solution
```

This also replaces `sample-app/` with the completed feature, saving its previous
contents under `.local/` and printing the backup location. Then launch `wad-ch-06`
and submit its task as chapter 06 describes. The path you work in remains `sample-app/`.

## Repeating a chapter

To recreate a chapter after changing kits or environment settings, exit its session
and use that chapter's `sbx env rm` command in HOST. Run its launcher again with the
same name. Your app's current edits remain; use a checkpoint only if you want to
replace them. Removing a sandbox also removes its in-sandbox conversation, so read
any result or question you need before removal.

If creation failed, inspect `sbx ls` first. Remove only the named workshop environment
if it exists; then retry its launch. Do not remove unrelated sandboxes.
