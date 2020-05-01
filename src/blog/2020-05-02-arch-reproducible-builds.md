---
description: reproducible builds for arch linux
title: arch reproducible builds
---

### _Arch Linux_ Reproducible Builds

Reason I don't use Gentoo / Linux from Scratch:
I don't want to build everything from source
(on my puny xps 13).
But have you ever wondered if the prebuilr binary packages
you get through the repos are actually what they claim to be?

[Reproducible Builds](https://reproducible-builds.org/)
is an effort to change that,
making build artifacts byte-for-byte identical and accountable.
[Arch Linux](https://tests.reproducible-builds.org/archlinux/archlinux.html)
is in there, and [here](https://reproducible.archlinux.org/),
with tooling like [repro](https://github.com/archlinux/archlinux-repro)
and [rebuilderd](https://github.com/kpcyrd/rebuilderd)
to make it easier for end users to automate verifications.

Verifying means taking the distributed package,
getting the sources,
and building it yourself to compare the results.

It's all still alpha quality software,
features are missing and may be buggy.
I [run](https://rebuilder.seankhliao.com/) /
[ran](https://web.archive.org/web/20200501111141/https://rebuilder.seankhliao.com/)
rebuilderd for [core] at some point.

#### _notes_

as of the time of writing:

- systemd-nspawn dependency makes running in a container more difficult (?)
- not split packages aware, rebuilds gcc 6 time
- config file options in constant flux
