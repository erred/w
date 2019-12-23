--- title
custom initramfs
--- description
debloat your initramfs?
--- main


### Notes

[these notes](https://hootiegibbon.gitlab.io/blog/2018/10/02/CustomInitramfs.html)
are good,
read them

#### tldr

- boot with kernel option `break=postmount`
- note modules: `lsmod | awk 'NF==3{print $1}{}'`
- create custom `mkinitcpio.conf` where:
  - `MODULES` is the above array
  - `BINARIES` is the fscks for your filesystem
  - `FILES` ignored
  - `HOOKS` only `base` is necessary

#### results

for my XPS 13

```
MODULES=(ext4 ahci i915)
BINARIES=(fsck fsck.ext4 e2fsck)
FILES=()
HOOKS=(base)
```
