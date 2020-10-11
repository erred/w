---
description: letting the user linger
title: loginctl linger
---

### _lingering_

systemd, the process lifecycle manager,
there's a global one, and one per user.
The user version starts when you first login,
and exits when you last logout.

If you want to keep things running,
or run things on boot

```sh
sudo loginctl enable-linger $USER
```

and reboot
