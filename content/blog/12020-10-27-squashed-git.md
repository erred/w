---
description: getting the right merge locally
title: squashed git
---

### _squashed_ git

extending [git notes](/blog/12019-06-06-notes-git/)

Don't know about you,
but my git history usually contains a bunch of errors
that isn't going to be useful to anyone.
So I think having everything related to a change in a single commit
is a very good idea (also linear history!).

#### _git_

git, being git, has quite a few ways to achieve the same effect.

#### _rebase_

`git rebase -i main` from the working branch
and change all the commits _after the first_
to `f` or `s`
(`f` will likely let you skip handling a bunch of conflicts).

#### _reset_

`git reset --soft main` and `git commit -m "..."`.
Works by forgetting all history and
creating a new commit with the final results.

#### _other_ notes

`git checkout -B branch` sets branch to point to whereever you are now
