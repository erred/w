---
description: do not recommend systemd-homed
title: thoughts on systemd-homed
---

### _systemd-homed_

I have systemd-homed setup to use a subvolume (btrfs).
Unfortunately, even though it uses no encryption
it still means that your home directory data is in `/home/user.homedir`
which gets mounted on `/home/user` on login.
This breaks some things,
such as ssh login which is expected.

What isn't expected is for it to break my single factor yubikey login flow.
I already moved the authorized keys list to somewhere in `/etc`.
The first login on boot will always fail, requiring me to type my password,
subsequent logins work fine.
Attached are the logs for one such login:

```txt
Nov 18 20:14:10 eevee systemd-homed[473]: arccy: changing state inactive → activating-for-acquire
Nov 18 20:14:10 eevee systemd-homework[582]: None of the supplied plaintext passwords unlocks the user record's hashed passwords.
Nov 18 20:14:10 eevee systemd-homed[473]: Activation failed: Required key not available
Nov 18 20:14:10 eevee systemd-homed[473]: arccy: changing state activating-for-acquire → inactive
Nov 18 20:14:10 eevee systemd-homed[473]: Got notification that all sessions of user arccy ended, deactivating automatically.
Nov 18 20:14:10 eevee systemd-homed[473]: Home arccy already deactivated, no automatic deactivation needed.
Nov 18 20:14:13 eevee systemd-homed[473]: arccy: changing state inactive → activating-for-acquire
Nov 18 20:14:13 eevee systemd-homework[583]: Provided password unlocks user record.
Nov 18 20:14:13 eevee systemd-homework[583]: Read embedded .identity file.
Nov 18 20:14:13 eevee systemd-homework[583]: Provided password unlocks user record.
Nov 18 20:14:13 eevee systemd-homework[583]: Reconciling embedded user identity completed (host and embedded version were identical).
Nov 18 20:14:13 eevee systemd-homework[583]: Recursive changing of ownership not necessary, skipped.
Nov 18 20:14:13 eevee systemd-homework[583]: Synchronized disk.
Nov 18 20:14:13 eevee systemd-homework[583]: Everything completed.
Nov 18 20:14:13 eevee systemd-homed[473]: Home arccy is signed exclusively by our key, accepting.
Nov 18 20:14:13 eevee systemd-homed[473]: arccy: changing state activating-for-acquire → active
```

The other thing that broke was groups.
The usual `usermod -a -G docker arccy` seemed to half work,
so does `homectl update arccy --member-of docker,adm,sudo,arccy`
since it doesn't have a append flag.
`userdbctl` will happily report that I am now a member of the `docker` group,
but `groups` doesn't and neither do other processes.
Relogin doesn't work,
maybe a restart of systemd-homed might have fixed it,
but it was easier to just reboot.

All in all, not a pleasant experience.
