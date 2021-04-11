---
title: sharing git hooks
description: getting those git hooks to everyone
---

### _git_ hooks

git has [hooks](https://git-scm.com/docs/githooks),
useful for executing scripts on some events,
like making sure everyone ran that formatter...
Anyway, the problem is that it lives outside of version control
and probably for good reason:
if it's a symlink and references the main working space,
checking out a different ref might change what hooks are run.

But you don't care,
anyway, there are (or were) several projects to install them:

- [icefox/git-hooks](https://github.com/icefox/git-hooks) bash, has been deleted along with all repos by the user...
- [git-hooks/git-hooks](https://github.com/git-hooks/git-hooks) Go
- [rycus86/githooks](https://github.com/rycus86/githooks) bash, comes with [blog post](https://blog.viktoradam.net/2018/07/26/githooks-auto-install-hooks/)
- [gabyx/githooks](https://github.com/gabyx/githooks) Go, port of rycus86/githooks

but the easiest way if you don't care about the symlink problem
(or change to `cp` and pester collaborators to rerun every time you update it if you do)
might just be something like this in a Makefile:

```Makefile
install-hooks:
        rm -rf .git/hooks
        ln -s ../hooks .git/hooks
        chmod +x hooks/*
```
