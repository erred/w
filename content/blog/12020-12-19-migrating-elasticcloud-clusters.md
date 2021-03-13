---
description: moving an elastic cluster
title: migrating elasticcloud clusters
---

### _elasticsearch_ cloud

Say you need to move elastic deployments,
maybe because you don't like being charged by credit card
and would prefer to pay through invoiced billing accounts
in google cloud.
Who cares that it's a little more expensive.

Anyway, you don't care that much about the data
(it'll all be gone after 2 weeks retention has passed),
but you do want to move your index(-patterns),
dashboards, and users over.
What do you do?

#### _elasticdump_

[elasticdump](https://github.com/elasticsearch-dump/elasticsearch-dump)
is a tool for moving indices,
and what do you know,
index-patterns and users are stored as indices.

```sh
# index patterns
docker run --rm -it elasticdump/elasticsearch-dump \
  --input=https://user:password@old-url/.kibana \
  --output=https://user:password@new-url/.kibana \
  --type=data

# roles/users
docker run --rm -it elasticdump/elasticsearch-dump \
  --input=https://user:password@old-url/.security \
  --output=https://user:password@new-url/.security \
  --type=data
```

#### _dashboards_

Somewhere in settings, "Saved Objects", export as json and import it.
