---
description: baby steps, one at a time, now you want some extra encryption
title: arch btrfs with encryption
---

### _encryption_

Why not encrypt everything, you think.
Well why not?

So you want something that plays nice with btrfs.
And you don't feel like touching lvm so soon.

#### _options_

dm-crypt/luks on an entire partition, btrfs inside with subvolumes:
almost everything is encrypted.

btrfs on partition, systemd-homed for user:
your user home can be encrypted, nothing else is,
eg: wifi passwords in `/etc/wpa_supplicant/` or `/var/lib/iwd/`,
wireguard passwords somewhere in `/etc/`

#### _setup_

1. create paritions:
   1. EFI partition
   2. root partition
2. create encrypted device
   1. `cryptsetup luksFormat /dev/nvme0n1p2`
   2. `cryptsetup open /dev/nvme0n1p2 root` last arg is the name to be used in `/dev/mapper`
3. create filesystems
   1. `mkfs.fat -F 32 /dev/nvme0n1p1`
   2. `mkfs.btrfs /dev/mapper/root`
4. setup btrfs
   1. `mount -o compress=zstd /dev/mapper/root /mnt`
   2. `mkdir /mnt/arch`
   3. `btrfs subvolume create /mnt/arch/@`
   4. `umount /mnt`
5. mount filesystems
   1. `mount -o compress=zstd,subvol=arch/@ /dev/mapper/root /mnt`
   2. `mkdir /mnt/boot`
   3. `mount /dev/nvme0n1p1 /mnt/boot`
6. `timedatectl set-ntp true`
7. `iwd station wlan0 connect $network`
8. `reflector --save /etc/pacman.d/mirrorlist -f 5 -p https`
9. `pacstrap /mnt base base-devel linux linux-firmware intel-ucode btrfs-progs zsh zsh-completions iwd git neovim`, minimal set of things to work confortable after reboot
10. `genfstab -U /mnt >> /mnt/etc/fstab`, alternatively edit it to use `/dev/mapper/root`
11. `arch-chroot /mnt`
12. `passwd`
13. edit `/etc/mkinitcpio.conf` set `HOOKS=(base systemd autodetect keyboard sd-vconsole modconf block sd-encrypt filesystems fsck)`, sd-vconsole is optional
14. add `/etc/crypttab.initramfs` with `root UUID=...`, this will get included in the initramfs so no kernel cmdlines are needed
15. regenerate initramfs, `mkinitcpio -p` or the foolproof way `pacman -S linux`
16. `bootctl install`
17. `echo "default arch\ntimeout 0\nconsole-mode max" > /boot/loader/loader.conf`
18. create loader entry, alternatively `root=/dev/mapper/root` is reliable enough since it is mapped from a UUID in our `/etc/crypttab.initramfs`
19. continue after rebooting into new system, some tools require the right systemd/dbus setup not available in chroot.

#### _/boot/loader/entries/arch.conf_

```txt
title   Arch Linux
linux   /vmlinuz-linux
initrd  /intel-ucode.img
initrd  /initramfs-linux.img
options root=/dev/mapper/root rootflags=compress=zstd:3,ssd,space_cache,subvol=arch/@ quiet rw
```

#### _notes_

The archwiki will at times refer to a boot partition and a efi/esp partition.
UEFI only cares about the bit that is FAT32 and contains efi things it can load: systemd-boot, kernel EFISTUB, grub, ...
Arch cares about the "boot partition" because by default hooks will dump the resulting files into `/boot`

The lazy way mounts the efi partition at `/boot`, UEFI things are contained in `/boot/EFI`,
the disadvantages being you need a bigger EFI partition to hold the kernel and stuff and they're not really protected.

The more involved way is the encrypted boot: grub (only grub supports this) is installed onto the unencrypted EFI partition,
where it knows how to decrypt the encrypted boot partition and start the kernel from there.
Note this doesn't protect against someone modifying your bootloader.

The correct solution is Secure Boot with your own keys: UEFI firmware verifies the bootloader,
the bootloader verifies the kernel, the kernel verifies its drivers.
It doesn't matter that the bootloader and the kernel sit out in the open unencrypted,
they shouldn't have secrets in them and any attempt to tamper with them would invalidate the signature and break the verification chain.

##### _decrypting_

you can either specify the devices to decrypt in the boot cmdline through the bootloader or as a file which gets built into the initramfs,
cmdlines will override all file mounts.

`/etc/crypttab` is not necessary since the root device (our only one) needs to already be unlocked when it is read.

##### _systemd-homed_

you can use a luks device to back a systemd-homed home,
so you can have a luks in your luks...

alternatively you can back it with a btrfs subvolume in your already
encrypted system.
