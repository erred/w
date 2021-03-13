---
description: tally of ssh login failures
title: ssh login failures
---

### _login_ fail

So I have/had a server with a public address,
and of course people try to login.

With SSH key only login + no root login,
I didn't run `fail2ban`, because why would you need it?

```txt
Jun 03 16:11:20 nevers sshd[1418181]: pam_tally2(sshd:auth): Tally overflowed for user root
```

Ah yes, overflows only took 5 months.
And they clutter up the journal anyways

#### _time_ period

6 months, from clean install

`head /var/log/pacman.log`:

```txt
[2020-01-07 09:18] [PACMAN] Running 'pacman -r /mnt -Sy --cachedir=/mnt/var/cache/pacman/pkg --noconfirm base base-devel linux linux-firmware intel-ucode zsh git docker sudo go go-tools htop man-db man-pages neovim python python-neovim prettier reflector exa ripgrep aria2 opemssh zsh-completions kitty-terminfo'
[2020-01-07 09:18] [PACMAN] synchronizing package lists
[2020-01-07 09:18] [PACMAN] Running 'pacman -r /mnt -Sy --cachedir=/mnt/var/cache/pacman/pkg --noconfirm base base-devel linux linux-firmware intel-ucode zsh git docker sudo go go-tools htop man-db man-pages neovim python python-neovim prettier reflector exa ripgrep aria2 openssh zsh-completions kitty-terminfo'
[2020-01-07 09:18] [PACMAN] synchronizing package lists
[2020-01-07 09:21] [ALPM] transaction started
...
```

#### _pam_ tally

which users do people/bots try?

`pam_tally2 --reset`:

```txt
Login           Failures Latest failure     From
root            65534    06/21/20 13:28:30  54.37.68.66
bin              1003    06/21/20 12:13:08  117.50.77.220
daemon          14205    06/21/20 04:41:31  14.18.61.73
mail            14322    06/21/20 10:57:08  139.213.220.70
ftp              4321    06/21/20 12:41:24  27.34.251.60
http              169    06/14/20 07:59:06  198.46.242.223
uuidd              31    05/15/20 20:59:46  182.61.108.39
dbus              109    06/11/20 02:08:55  54.38.158.17
ntp                42    05/27/20 06:48:57  195.231.1.153
polkitd            58    05/30/20 22:35:46  178.128.13.87
grafana           107    06/20/20 17:26:02  106.13.147.89
prometheus        145    06/21/20 12:39:08  83.17.166.241
znc                60    06/21/20 01:56:39  103.1.179.94
dhcp               16    06/09/20 00:54:43  61.154.14.234
mysql            3858    06/21/20 13:02:41  61.111.32.137
cacti             350    06/21/20 13:01:34  46.164.143.82
colord             50    06/20/20 10:56:08  59.63.212.100
avahi              72    06/13/20 15:44:31  49.233.88.126
git              7746    06/21/20 13:28:26  182.74.25.246
systemd-network    12    05/16/20 22:53:50  78.118.109.44
gerrit             56    06/20/20 19:56:08  1.71.140.71
nobody           2851    06/21/20 11:17:43  148.70.35.211
```
