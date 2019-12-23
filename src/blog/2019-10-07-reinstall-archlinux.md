--- title
reinstall archlinux
--- description
the long and arduous journey of reinstalling arch
--- main


*why?*
you ask,
it has gotten cluttered and full of dubious configs,
and arch released the new `base` metapackage today,
so what better time for a fresh start

### clean device

#### make backups

  - and check them, i forgot, so i have none
  - the most painful bit is *losing ssh keys*

#### create install usbs

  - check that they work properly

#### wipe device

  - now is the time to make changes to bios

### install

the wiki is the holy text you must read and understand

#### partitioning

```
/dev/nvme0n1p1  2G        /efi  EFI system part
/dev/nvme0n1p2  10G             alpine, todo
/dev/nvme0n1p3  rest      /     arch
```

#### pacstrap

on bare metal dell xps 13 *9350* with intel *8265* wifi

```
pacstrap /mnt base linux linux-firmware intel-ucode wpa_supplicant git sudo neovim python-neovim
```

#### bootloader

systemd-boot

1. `bootctl --path=/efi install`
2. `mkdir /efi/EFI/arch`
3. `mv /boot/* -t /efi/EFI/arch`
4. `mount --bind /efu/EFI/arch /boot`
5. add bind mount to `/etc/fstab`
6. add entry to `/efi/loader/entries/`, paths will be like `/EFI/arch/vmlinuz-linux`

### reboot

pray you didn't forget anything important, 
*like me*,
or it's back to the install usb, mount and arch-chroot to fix stuff

#### user

`groupadd sudo`, `useradd -a -G user,video,input,sudo user`, `passwd user`, edit `/etc/sudoers`

#### networking

1. `ip link set wlan0 up`
2. `wpa_passphrase ssid passwd > /etc/wpa_supplicant/wpa_supplicant-wlan0.conf`
3. `systemctl enable --now wpa_supplicant@wlan0`
4. add basic entry to `/etc/systemd/network/30-any.network`
5. `systemctl enable --now systemd-networkd`
6. edit `/etc/systemd/resolved.conf` for preferred dns
7. `systemctl enable --now systemd-resolved` 

#### major changes

~install `kmscon`~: doesn't play nice with sway

#### get *yay* back

`git clone https://aur.archlinux.org/yay-bin.git` 
and `cd yay-bin` 
then `makepkg -si`

**note:** this probably requires the `base-devel` stuff

### custom stuff

#### getting custom stuff back

`git clone https://github.com/seankhliao/config .config`

and then *disable* starting the wm in `zsh/zprofile` 
and *disable* enforcing ssh keys in `git/config`

#### reinstalling more stuff

or just use the `xps-system`, `gui-env`, `dev-tools` packages [here](https://github.com/seankhliao/pkgbuilds)

##### visual environemnt
  - sway
  - xf86-video-intel
  - i3status
  - bemenu
  - mako
    - `systemctl --user enable mako`
  - grim, slurp
  - playerctl, alsa-utils
  - brightnessctl
  - kitty
  - google-chrome-dev
  - xorg-server-xwayland
  - wl-clipboard(-x11)

##### tools
  - reflector
  - exa
  - ripgrep
  - tag-ag
  - htop
  - aria2
  - neovim-plug-git
  - man-db, man-pages
  - openssh

#### editing configs

##### `/etc/makepkg.conf`
```
-march=native
-mtune=native
OPTIONS=(!strip !zipman)
INTEGRITY_CHECK=(sha256)
PKGEXT='.pkg.tar'
SRCEXT='.src.tar'
```

##### `/etc/pacman.conf`
```
### Misc options
UseSyslog
Color
TotalDownload
CheckSpace
VerbosePkgLists
```

##### `/etc/pam.d/system-auth`
```
auth      required  pam_unix.so     try_first_pass nullok nodelay
```

##### `/etc/systemd/logind.conf`
```
HoldoffTimeoutSec=10s
```

##### `/etc/systemd/system.conf`
```
RebootWatchdogSec=10s
ShutdownWatchdogSec=10s
DefaultTimeoutStartSec=10s
DefaultTimeoutStopSec=10s
```
##### `/etc/systemd/user.conf`
```
DefaultTimeoutStartSec=10s
DefaultTimeoutStopSec=10s
```

#### ssh keys

```
ssh-keygen -t ed25519
```
