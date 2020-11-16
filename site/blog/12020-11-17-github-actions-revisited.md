---
description: looking at it again
title: github actions revisited
---

### _actions_

CI/CD integrated with github,
can't be that bad, right?

#### _good_ things

More event/trigger types.
Now you can trigger actions on most (all?) github repository events.

Artifacts/caching thing.
Mitigates the fact that you have a clean instance on every run.
Cache keys may still be confusing.

#### _not_ good things

Self hosted runner is unreliable in a container,
eg. in kubernetes.
They seem fine in a vm.

Conditionals controlling evaluation are confusing,
difficult to test
and will silently evaluate non-existing fields to nothing without error.
Advice is to use them as sparingly as possible.

Private actions are not supported,
workaround is to clone the private repo and use it as a local action.

Composite actions cannot compose other actions,
just shell scripts.
So you pretty much have to write your complex actions from scratch
in order to reuse them.

The previous 2 points mean the possibility of workflow reuse is fairly low,
get used to copy-pasting around large chunks of yaml.

I'm starting to miss Tekton
