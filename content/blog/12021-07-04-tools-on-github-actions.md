---
title: tools on github actions
descriptions: that pesky clean environment every time
---

### _tools_ on github actions

[Github Actions](https://github.com/features/actions),
your favorite arbitrary code execution platform with a long startup time.

Anyway, within you CI process, usually you will want some non default tools.
If you were [self hosting](https://docs.github.com/en/actions/hosting-your-own-runners/about-self-hosted-runners)
the runners (executors), that isn't a problem,
just install them on the runner.
If you're using the cloud ones though...

The "recommended" way is to have an action install/configure the tool every time.
This results in something like:

```yaml
jobs:
  build:
    steps:
      - use: toolx/setup
      - use: tooly/setup
      - use: toolz/setup
      - run: toolx ...
      - run: tooly ...
      - run: toolz ...
```

Which is very much meh.
You install the tool every time, slowing down the entire process for what?
And if you have a lot of repos keeping everything in sync is not fun.

#### _container_

What you could instead is bundle all your tools into a container image...

```yaml
jobs:
  build:
    container:
      image: ghcr.io/my/builtools
    steps:
      - run: toolx ...
      - run: tooly ...
      - run: toolz ...
```

This does mean everything shares the same tools / versions,
but that's also an advantage since updating needs only be done in one place.
but you still use a different container for a step if it requires it.
