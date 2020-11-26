---
description: copying containers between registries
title: skopeo copy
---

### _skopeo_

Ever needed to move a container image somewhere?
Or add a new tag?
There's always the trusty:

```sh
docker pull $src
docker tag $src $dst
docker push $dst
```

But this is inefficient and requires you to download the entire image
and maybe upload it again if it's to a different registry.

Instead, you could use _skopeo_:

```sh
skopeo copy docker://$src docker://$dst
```

This is a single command and only needs to operate on the image manifests,
you can pass image layers between registries directly.

How?
[docker image manifests](https://docs.docker.com/registry/spec/manifest-v2-2/#image-manifest)
support a `urls` field for retrieving layers from a remote,
so you can move images between datacenters
without passing it through your own shitty network.
