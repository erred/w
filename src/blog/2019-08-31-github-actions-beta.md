--- title
github actions beta
--- description
first impressions on github actions beta
--- main


A full 10 months after clicking "Join Beta",
I'm finally in.

Docs are scarce and already outdated,
it's now _yaml_ instead of hcl.

# links

- `workflow.yaml` syntax: [docs](https://help.github.com/en/articles/workflow-syntax-for-github-actions)
- `action.yml` syntax: [docs](https://help.github.com/en/articles/metadata-syntax-for-github-actions)

# Build stuff

Right now,
I like docker containers,
so build on push:

```
on: push
jobs:
  job1:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@master
      - run: docker build -t some_image .
      - run: docker login registry -u username -p $TOKEN
        env:
          TOKEN: ${{ secrets.SOME_SECRET }}
      - run: docker push some_image

```

# model

you get a vm with a bunch of [outdated stuff installed](https://help.github.com/en/articles/software-in-virtual-environments-for-github-actions)

run commands with `jobs.jobid.steps[x].run:`

run actions with `jobs.jobid.steps[x].uses:`

## actions

js or a container with extra metadata for input/output validation,
refer to by:

- user / repo (/ path) @ ref
- . (current repo root) / path / to / action / dir
- docker://some-image:tag

# notes

secrets are filtered from log output,
**yay**?
but now how do we debug?

apparently there's no way to set global substitutions?

also doesn't integrate well with the package registry beta

there are 3 ways to run a container as a step:

1. `run: docker run some-image`
2. `uses: docker://some-image`
3. `uses: user/repo` (assuming a Dockerfile and action.yml exists)
