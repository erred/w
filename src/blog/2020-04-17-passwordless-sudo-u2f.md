---
description: sudo without a password
title: passwordless sudo u2f
---

who cares about the _2nd_ factor?
just the hardware key needed!

### steps

#### install pam-u2f

```sh
pacman -S pam-u2f
```

#### add keys

replace `~/.config` with `$XDG_CONFIG_HOME`

origin and id is set to enforce a stable name

```sh
# first key
pamu2fcfd -i pam://hostname -o pam://hostname > ~/.config/Yubico/u2f_keys
# other keys
pamu2fcfd -n -i pam://hostname -o pam://hostname > ~/.config/Yubico/u2f_keys
```

~/.config/Yubico/u2f_keys

```txt
username:xxxxxx..key1..xxxxxx:xxxxxx..key2..xxxxxx
```

#### add authentication method to pam

/etc/pam.d/sudo

- authfile: set alternative location for config file
- cue: prompt to touch device (so you don't wonder why it's stuck doing nothing)

```txt
auth    sufficient    pam_u2f.so cue origin=pam://hostname appid=pam://hostname
...
```
