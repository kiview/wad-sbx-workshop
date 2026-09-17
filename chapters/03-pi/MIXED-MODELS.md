# Role, harness, provider, model

Everyone sets up separate roles. The default uses one Anthropic model so it works
with one provider. The presenter variation below runs Pi on Google, Claude as the
developer and Codex as QA. The developer stays gateway-capable in the MCP chapters.

The supplied sandbox kit `../../kits/multi-provider` declares provider credentials
and installs Codex. OAuth declarations belong in a sandbox kit, not a mixin; Pi and
Herdr remain separate mixins. No API key belongs in this repository. Configure the
`anthropic`, `openai` and `gemini` services in SBX. Existing `google` secrets can be
selected with the kit's `google_secret_service` argument when invoking SBX directly.
The default example expects `gemini`.

For the basic non-MCP team chapter:

```bash
# HOST · from chapters/04-team
./launch-mixed my-mixed-team
```

Inspect the three sessions and repeat the short handoff exercise from chapter 04.
Changing a model in `team.tsv` affects a newly launched team, not an already running
session. Use an available model ID for your account. One-provider learners edit the
same columns and keep supported combinations; no second subscription is invented.

For MCP chapters, start with the normal Claude base. Porting the full mixed-provider
recipe there also needs gateway client registration in the custom base kit; the
example above intentionally covers chapter 04, which has no gateway dependency.

The source workshop recorded successful basic chapter runs on SBX v0.45.0-rc2.
Prior full-factory provider smoke results are not proof of a new recipe. The available
Google credentials previously returned a monthly spending-limit error; a valid key
alone does not demonstrate a successful Gemini response.

On the tested SBX build, `sbx env plan` accepts a custom sandbox-kit string but
`sbx env create` rejects it as an unknown agent. The mixed launcher therefore uses
the supported `sbx create KIT --kit MIXIN` command. This is six additional host
lines; the normal chapters still teach sbxenv. It is a recorded CLI limitation,
not a claim that the plan proves custom-agent creation works.
