Read your role brief. Coordinate through handoff files and crew-notify.
Ask the developer to retrieve the task named in ~/work/task-id through Beans MCP
and save it to ~/work/task.json before implementing anything.

This is a different project, available in ~/work/app. First ask the developer to read its README,
contributor guide, agent instructions, dependency manifests and build configuration.
Determine the required runtimes, dependencies and supporting services. Install and
start what is needed inside SBX, using containers for services where appropriate.
Discover and run a relevant baseline check; report any existing failure separately
from failures caused by the change. Then implement only the requested change. QA reviews the
exact commit and runs the relevant checks. Do not replace existing project guidance.

Project setup: work out and follow the repository's development setup inside SBX;
the human does not need to supply installation commands.
Use the sandbox's Docker engine if tests need containers. Do not use the workshop
sample's board helper or assume that this project uses Node or PostgreSQL.
If a dependency, toolchain or network request is unavailable, report the exact need
to the human. Do not claim a skipped check passed.

Preview: for a library or CLI, demonstrate the change with a focused test or CLI
example; no web server is required. For a web app, use its documented start command
and listen on 0.0.0.0:8080 to use the existing host port 3102 mapping.

After review, ask the developer to append a Beans result note with the commit,
summary, actual check outcomes and QA findings. Leave the task open. Ask the human
about ambiguous requirements rather than choosing a product answer.
