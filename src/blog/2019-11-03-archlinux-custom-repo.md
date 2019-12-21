custom self updating arch linux repo

---

# what is it?

ugly hack to prebuild [AUR](https://aur.archlinux.org/)
packages and host them on Github
(this will likely have to change in the future if repo size is gets too big)

# How?

1. prebuild a docker image with `base-devel` and `yay-bin`
2. use `yay` to install packages
3. create a repository with `repo-add`
4. push back to repo

## Build

this runs on [Github Actions](https://github.com/seankhliao/arch-repo/blob/master/.github/workflows/workflow.yaml)
and also uses Github Package Repository for the docker image

# use

The project is at [seankhliao/arch-repo](https://github.com/seankhliao/arch-repo)

Use with something like the following in `pacman.conf`

```
[seankhliao]
SigLevel = Optional TrustAll
Server = https://raw.githubusercontent.com/seankhliao/arch-repo/master/pkgs/
```

# is this legal?

maybe not?
