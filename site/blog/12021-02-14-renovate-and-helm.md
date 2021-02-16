---
description: ugly yaml hacks
title: renovate and helm
---

### _renovate_

At $work we run [renovate](https://github.com/renovatebot/renovate)
to keep our docker images up to date,
and for us that also means updating the helm charts to use newer images.
We also use "upstream" Helm charts
(a decision we're considering walking back on).

Renovate supports helm values files with the structure below,
note the parent must be called `image`.
(lifted from [docs](https://docs.renovatebot.com/modules/manager/helm-values/))

```yaml
image:
  repository: "some-docker/dependency"
  tag: v1.0.0
  registry: registry.example.com # optional key, will default to "docker.io"
```

#### _helm_

"upstream" helm charts are all over the place in terms of quality and standards.
Thankfully, most people seem to follow the basic template created by `helm create`.

Unfortunately, Helm has a best practices doc that says different things,
in particular [flat or nested values](https://helm.sh/docs/chart_best_practices/values/#flat-or-nested-values).
Example of flat:

```yaml
imageRepository: "some-docker/dependency"
imageTag: v1.0.0
imageRegistry: registry.example.com
```

So what if your helm chart does this?

#### _options_

renovate has a general purpose [regex](https://docs.renovatebot.com/modules/manager/regex/)
manager module that you could use to match the keys.
But that's more config to keep somewhere.

YAML has anchors, and references, and Helm thankfully ignores unused values,
so you could instead do:

```yaml
app:
  image:
    repository: &repo "some-docker/dependency"
    tag: &t v1.0.0
    registry: &reg registry.example.com # optional key, will default to "docker.io"

imageRepository: *repo
imageTag: *t
imageRegistry: *reg
```

renovate would recognize the nested structure, and helm would dereference the value
when it parses the yaml so it all comes out correct.

I was only midly surprised my coworker agreed to this hack...
