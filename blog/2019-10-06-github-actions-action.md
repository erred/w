wrting actions for github actions

---

footguns,
footguns everywhere

## workflows

they finally added a note saying
that the default `GITHUB_TOKEN`
doesn't work with `docker.pkg.github.com`

everything is an env,
no nested declarations,
copy and paste things everywhere

## actions

`action.yml` _MUST_ end in `.yml` and not `.yaml`,
otherwise it won't be used.
So confusing since workflows work perfectly fine with `.yaml`

built in a clean environment on an as needed basis

env from the following are injected at runtime:

- github default env: `GITHUB_SHA`, `GITHUB_REPOSITORY`, etc...
- workflow global env
- workflow job env
- workflow step env
- with/input as `INPUT_WITH-ID`

you are honestly better off just copy and pasting the damn thing everytime you need it,
it's just as long and just as fragile,
especially if have special input/output.
