---
description: here at it again
title: reinstall arch v4
---

### _install_

1. setup disks, network, and base install from livecd,
   refer to eg [yesterday](/blog/12020-11-08-arch-dm-crypt-btrfs).
   - sanity check things that were installed with pacstrap:
     `base base-devel linux linux-firmware intel-ucode iwd neovim btrfs-progs zsh zsh-completions sudo git`
2. Prepare network stuff
   1. `echo "[Match]\nName=wlan\*\n\n[Network]\nDHCP=yes" > /etc/systemd/network/90-wlan.conf`
   2. `echo 'eevee' > /etc/hostname`
3. `systemctl enable systemd-{networkd,resolved,homed} iwd`
4. reboot into new system
5. Network stuff
   1. `iwctl station wlan0 connect XXX`
6. locale stuff
   1. `sed -i 's/#en_US.UTF-8/en_US.UTF-8/' /etc/locale.gen`
   2. `locale-gen`
   3. `localectl set-locale LANG=en_US.UTF-8`
7. time stuff
   1. `timedatectl set-ntp true`
   2. `timedatectl set-timezone Europe/Amsterdam`
   3. `timedatectl set-local-rtc false`
8. user stuff
   1. `useradd sudo`
   2. `echo '%sudo ALL=(ALL) ALL' > /etc/sudoers.d/sudo`
   3. `homectl create arccy --shell=/bin/zsh --storage=subvolume -G adm,sudo`
9. auth stuff
   1. `/etc/pam.d/system-auth`: add `nodelay` to `pam_unix.so` to skip the delays on wrong passwords
   2. `/etc/pam.d/system-auth`: add `auth sufficient pam_u2f.so cue origin=pam://eevee appid=pam://eevee` for 1fa yubikey authentication
   - `pamu2fcfg -i pam://eevee -o pam://eevee > ~/.config/Yubico/u2f_keys` to enroll first key, add `-n` to add more with append
   - replace `eevee` with your hostname
   - if using `systemd-homed` your home directory is not mounted (an maybe is encrypted) when you're not logged in so the pam-u2f won't be able to find the authorized keys...
10. setup pacman
    1. `pacman-key --init`
    2. `pacman-key --populate archlinux`
    3. uncomment options in `/etc/pacman.conf`: `UseSyslog` `Color` `TotalDownload` `VerbosePkgLists` and testing repos
11. get prebuild `yay`
    1. `curl -L https://github.com/Jguer/yay/releases/download/v10.1.0/yay_10.1.0_x86_64.tar.gz | tar xzv`
    2. `./yay -Syu yay`
12. install stuff
    - DE: `brightnessctl grim i3status mako slurp sway swaylock wf-recorder wofi wl-clipboard-x11 xorg-server-xwayland xf86-video-intel`
    - apps: `alacritty google-chrome yubioath-desktop`
    - audio: `pulseaudio pulseaudio-bluetooth pulsemixer`
    - editor: `bash-language-server prettier neovim-git python-pynvim-git dockerfile-language-server-bin`
    - container: `docker kubectl kustomize`
    - tools: `age iw htop jq man-db openssh pam-u2f reflector wireguard-tools`
    - qol: `aria2 exa git-delta-bin ripgrep rsync skim xsv-bin yay-bin`
    - fonts: `noto-fonts-cjk noto-fonts-emoji terminus-font ttf-google-fonts-git`
13. enable more stuff
    - `systemctl enable --now bluetooth`
    - `systemctl enable --now docker.socket pcscd.socket`
    - `systemctl enable --now fstrim.timer`
    - `systemctl --user enable --now mako`
14. other
    - `echo 'FONT=ter-128n' > /etc/vconsole.conf` set large console font with terminus
    - setup bluetooth
    - setup editor plugins
    - clone .config
    - get ssh config
    - setup wireguard
    - shorten timeouts in `/etc/systemd/{system,user}.conf`
    - browser local extensions (custom newtab, bypass-paywalls)
    - go toolchain and tools
