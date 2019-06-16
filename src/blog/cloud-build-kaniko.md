title = cloud build kaniko
date = 2019-06-01
desc = using Kaniko Executor with Google Cloud Build

---

# kani-what?

Maybe for some reason you don't like using `docker build`,
you want to be cool and use `gcr.io/kaniko-project/executor`
if you're not building things in subdirectories,
you're all set.
Else, read on.

## _options_

cloudbuild.yaml `dir`: this sets the working directory,
defaults to the project root (`/workspace`)
kaniko will pick the `Dockerfile` from here
(if not overridden by the `-f` flag)

kaniko-project/executor `-f`: Dockerfile,
defaults to `Dockerfile`
relative paths are from the cloudbuild working directory

kaniko-project/executor: `-c`: kaniko build context,
defaults to the project root (`/workspace`)
also must be an absolute path (?)

## building in a subdir

to get the equivalent of `dir: subdir` and `docker build .`

you need `dir: subdir` and `kaniko-project/executor -c /workspace/subdir`
