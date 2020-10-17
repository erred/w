---
description: thoughts on terraform after using it for a few more things
title: terraform second impressions
---

### _terraform_

Control anything with an API.

#### _pain_ points

Terraform is good for spinning up new things,
or changing a few key values.
Less nice for managing something that already exists,
since you have to hunt down the things that already exist
and try to write a matching config.

google/google-beta provider specific complaints:
google-managed service accounts aren't importable or manageable through terraform,
best you can do is record their email as a variable and manage iam bindings.
iam policies feel all over the place.
cloud run requires setting an image,
but this changes every time you release a new version,
making it unsafe to re-apply config especially if you have CI/CD setup

azure devops complaints: anything touching github from the service connection
to setting up a build definition to trigger on github push doesn't work.
Things like secure files aren't even available.
Having terraform only manage half your stuff seems worse than all-or-nothing.

github complaints:
pretty much useless, when thing i want to do most is move/rename repos

hcl complaints:
the block vs map syntax makes the language feel inconsistent.
