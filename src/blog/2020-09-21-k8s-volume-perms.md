---
description: file permissions in kubernetes volumes
title: k8s volume perms
---

### _k8s_ volume perms

By default, containers run as root
and volumes are mounted as root too.

Everyone tells you to "do security properly"
and don't run containers as root.
To do this in kubernetes,
you can enforce uid ranges / run as non root in `PodSecurityPolicy`,
and select the exact uid in the `securityContext` section of both `Pod` and `container`.

Volumes are not granted such affordances,
instead you get a `fsGroups` setting in both `PodSecurityPolicy` and `securityContext`,
changing the group but not the owner of the volume.
If you just wanted to read/write files, this is enough.

Unfortunately, sometimes programs try to be "smart" about security,
and complain if some file has permissions beyond a single user.
If we remove the group permissions,
our process, running as non root won't have permissions to see the file.
So our only option is to change the owner, which needs root rights,
which blows a hole in our security policies.

We can sort of limit the scope of the hole by doing it in a separate `initContainer`,
but still....
Also, some volumes are mounted read only,
like `secret` (where your confidential files usually reside)
so now you need to copy them to a writable directoy like `emptyDir`,
and then you lose the transparent updates.

Why do the k8s gods not see it fit to let us set filesystem owners...
