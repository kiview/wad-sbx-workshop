# Presenter extensions: where the factory can go next

These demonstrations belong in the final open-lab slot. Attendees already have the
factory; they do not need cloud access, an enrolled organization or a second-sandbox
communication protocol to finish. Rehearse these separately before presenting them.

## Inspect the environment you built

```bash
# HOST — choose one of your existing chapter sandboxes
sbx inspect wad-ch-05
```

Point to the assembled kits, published ports, mounts, access configuration and
sessions. This is useful for explaining what exists and troubleshooting it. It is
not a transcript of every agent action or proof that a task succeeded. Review the
output before projecting it; it can identify local paths and account configuration.

## Runtime mount and unmount

The opening chapter deliberately omitted host mounts. Now demonstrate granting and
revoking access to one new, non-sensitive directory. Keep the target sandbox's
launcher/session open. On the host:

```bash
# HOST
mkdir -p "$WORKSHOP/.local/mount-demo"
printf 'A message from the host\n' > "$WORKSHOP/.local/mount-demo/message.txt"
sbx mount wad-ch-05 "$WORKSHOP/.local/mount-demo:/home/agent/mailbox:ro"
sbx exec wad-ch-05 cat /home/agent/mailbox/message.txt
sbx exec wad-ch-05 sh -c 'echo changed > /home/agent/mailbox/message.txt'
```

Expected: the read succeeds; the write fails because this mount is read-only.
Then revoke it:

```bash
# HOST — CLI spelling is umount
sbx umount wad-ch-05 "$WORKSHOP/.local/mount-demo:/home/agent/mailbox"
sbx exec wad-ch-05 test ! -e /home/agent/mailbox/message.txt
echo $?
```

Expected: `0`, confirming that file is no longer visible there. The host file is
still present. Use the same host path and target when unmounting.

**Extension idea, not a new workshop component:** mount one scratch directory into
two sandboxes, and one can write a message the other reads. A shared directory is
only transport. A real protocol still needs ownership, atomic publication, message
IDs and wakeups. Our factory already has those concerns handled inside one sandbox
by files plus Herdr; do not add a second protocol during this session. A read-write
mount deliberately gives the sandbox write access to that host directory.

## Cloud, on the presenter's account

Before the session, confirm account access, supported commands and costs. Show the
current supported surface with:

```bash
# HOST
sbx --cloud --help
sbx --cloud ls
```

Use an already prepared cloud demo sandbox, substituting its name below:

```bash
# HOST — replace CLOUD_DEMO_NAME with your prepared cloud sandbox name
sbx --cloud exec CLOUD_DEMO_NAME bash -lc 'whoami; pwd'
sbx --cloud ports CLOUD_DEMO_NAME
```

If it has a prepared service listening on 8080, publish it with
`sbx --cloud ports CLOUD_DEMO_NAME --publish 8080` and open the URL returned by the
cloud control plane. Cloud port publishing takes the sandbox port, not the local
`3106:8080` binding syntax. Do not assume that local host paths or the local Beans
stdio process become reachable from a cloud sandbox. Moving the full factory's
host integration to cloud is a separate design exercise.

This is a prepared presenter demo, not a verified chapter-05 cloud migration.
Remove the demo's exposed port when done with
`sbx --cloud ports CLOUD_DEMO_NAME --unpublish 8080`, and stop the cloud sandbox
according to the account's supported lifecycle commands.

## Organization governance

On an enrolled presenter account, show a preconfigured named-tool rule, make the
corresponding harmless allowed/denied call, and inspect the actual evidence.
Identify what is enforced centrally versus what the local adapter itself restricts.
Do not ask attendees to enroll an organization during the workshop. If access is
unavailable, explain the extension without fabricating a live result.
