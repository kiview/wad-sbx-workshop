Ask the developer to retrieve the task named in ~/work/task-id through the Beans
MCP get_task tool and save it to ~/work/task.json. Do not invent its requirements.
The developer should inspect the mounted project's documentation, install its
dependencies and start any required services inside SBX, then implement the task.
Ask QA to review the exact commit against the task and coding policy.
After review, refresh the web app on 0.0.0.0:8080 and leave it running.
Ask the developer to append a Beans result note with the change, commit and actual
check outcomes. Leave the task open. Send the human a summary and how to try it.
