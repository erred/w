--- title
reinstall arch v2.md
--- description
let's try again
--- main


what??
it's only been 3 weeks
and it broke

**note:**
pam with u2f is finicky
and there's a good chance you'll get locked out
and just removing the line from `/etc/pam.d/system-auth` doesn't help

# wiki

follow the [Installation Guide](https://wiki.archlinux.org/index.php/Installation_guide)

# changes

## partition / mounting

```
mount /dev/nvme0n1p3 /mnt
mkdir /mnt/{boot,efi}
mount /dev/nvme0n1p1 /mnt/efi
mount --bind /mnt/efi/EFI/arch /mnt/boot
```

1. edit `/etc/fsab` after generating to fix bind mount
2. edit `/efi/loader/entries/arch.conf` to fix UUID

## pacstrap

```
pacstrap /mnt base base-devel linux linux-firmware intel-ucode \
  wpa_supplicant go{,-tools} htop man-{db,pages} neovim git python{,2,-neovim} \
  reflector exa ripgrep aria2 sudo docker openssh zsh zsh-completions yubikey-manager \
  sway xf86-video-intel swaylock mako i3status bemenu grim slurp playerctl \
  brightnessctl alsa-utils kitty xorg-server-xwayland noto-fonts{,-emoji,-cjk} ttf-ibm-plex
```

## user
```
groupadd -r sudo
useradd -m -G vidoe,input,sudo,docker user
passwd user
EDITOR=nvim visudo
```

### as user
```
git clone https://aur.archlinux.org/yay-bin.git
cd yay-bin && makepkg -si
yay -S wl-clipboard-x11 tag-ag
cd ~
git clone https://github.com/seankhliao/pkgbuilds
cd pkgbuilds/sway-service && makepkg -si
cd ~
rm -rf .config 
git clone https://github.com/seankhliao/config .config
sudo ln -s $(pwd)/.config/zsh/zshenv /etc/zsh/zshenv
cd ~
mkdir -p data/{down,xdg/{nvim/{backup,undo},zsh}}
```

## system
after reboot
```
systemctl enable --now wpa_supplicant@wlp58s0
systemctl enable --now pcscd
# edit /etc/resolv.conf
```
