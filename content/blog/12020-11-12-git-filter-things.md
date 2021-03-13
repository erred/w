---
description: the overloaded word filter in git
title: git filter things
---

#### _git_ filter-branch

[git-filter-branch](https://git-scm.com/docs/git-filter-branch)
oldest (and not recommended) tool here for rewriting history.
Maybe still useful in some cases?

#### _git_ filter-repo

[git-filter-repo](https://github.com/newren/git-filter-repo)
Python script that's faster than filter-branch and does mostly the same stuff.
Use to: repo->subdirectory, subdirectory->repo, remove file from history, ...

#### _git_ clone --filter

docs in `man git-rev-list` Object traversal.
`--filter=blob:none` skips copying the actual files (blobs),
retrieving them on-demand with other commands, eg `checkout <ref> -- <path>`.
`--filter=tree:0` skips blobs and trees (things that hold the blobs together eg directories).

Right now the main gains are in (really) big repos since the on-demand checkout
is sequentially per file.
