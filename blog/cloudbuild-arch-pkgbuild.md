title = cloudbuild arch pkgbuild
date = 2019-06-04
desc = building arch PKGBUILDs in cloudbuild

---

# What?

I wanted a private private arch repository
My thinking was:

- git repo of [PKGBUILD](https://wiki.archlinux.org/index.php/PKGBUILD)s
- Build in the cloud on push
- Store in cloud blob storage

## Infra

- git: [Github](https://github.com)
- build / CI/CD: GCP [Cloud Build](https://cloud.google.com/cloud-build/)
- storage / hosting: GCP [Cloud Storage](https://cloud.google.com/storage/)

## git

Easy.

## Triger build on push

Mostly reliable [build triggers](https://cloud.google.com/cloud-build/docs/running-builds/automate-builds)

## Building

Arch uses [makepkg](https://wiki.archlinux.org/index.php/Makepkg)

makepkg doesn't run as root,
so `sudo -u nobody makepkg --needed --noconfirm`

but makepkg needs pacman to resolve dependencies,
so first `echo 'nobody ALL=(ALL) NOPASSWD: /bin/pacman' >> /etc/sudoers`

also as nobody you don't have write permissions for the mounted `/workspace`,
copy to/from `/tmp`

## Deploying

Remember to pull in your `repo.db.tar`,
update with [repo-add](https://wiki.archlinux.org/index.php/Pacman/Tips_and_tricks#Custom_local_repository)
and push both the new database files and packages to hosting

and firgure out a way to clean up old packages

## Notes

Arch with `base` and `base-devel` installed is 1.6GB,
even with GCP's network and caching,
it is still a pain to download it onto build workers
