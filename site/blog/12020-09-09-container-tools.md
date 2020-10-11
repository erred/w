---
description: break the docker cli monoculture!
title: container tools
---

#### _tools_

You're tired of the docker monoculture
or you really don't like docker's daemon model.
There are other tools.

#### _image_ management

| command            | build            | pull                | push                | tag                     | rm                 |
| ------------------ | ---------------- | ------------------- | ------------------- | ----------------------- | ------------------ |
| [docker][docker]   | `docker build .` | `docker pull $img`  | `docker push $img`  | `docker tag $img $tag`  | `docker rmi $img`  |
| [podman][podman]   | `podman build .` | `podman pull $img`  | `podman push $img`  | `podman tag $img $tag`  | `podman rmi $img`  |
| [buildah][buildah] | `buildah bud .`  | `buildah pull $img` | `buildah push $img` | `buildah tag $img $tag` | `buildah rmi $img` |
| [img][img]         | `img build .`    | `img pull $img`     | `img push $img`     | `img tag $img $tag`     | `img rmi $img`     |
| [ctr][ctr]         | -                | `ctr i pull $img`   | `ctr i push $img`   | `ctr i tag $img $tag`   | `ctr i rm $img`    |
| [crictl][crictl]   | -                | `crictl pull $img`  | -                   | -                       | `crictl rmi $img`  |

##### _others_

- docker has 2 build backends, the default and buildx
- [skopeo][skopeo] has a copy command that can probably do the same as push/pull
- [kaniko][kaniko] can build (and pull/push) all in one step, mostly meant for CI

#### _container_ management

| command          | run                        | stop               | rm               | exec                    | cp                    | status                     |
| ---------------- | -------------------------- | ------------------ | ---------------- | ----------------------- | --------------------- | -------------------------- |
| [docker][docker] | `docker run $img`          | `docker stop $con` | `docker rm $con` | `docker exec $con $cmd` | `docker cp $src $dst` | `docker ps` / `docker top` |
| [podman][podman] | `podman run $img`          | `podman stop $con` | `podman rm $con` | `podman exec $con $cmd` | `podman cp $src $dst` | `podman ps` / `podman top` |
| [ctr][ctr]       | `ctr run $img $id`         | -                  | `ctr c rm`       | -                       | -                     | -                          |
| [crictl][crictl] | `crictl run $cconf $pconf` | `crictl stop $con` | `crictl rm $con` | `crictl exec $con $cmd` | -                     | `crictl ps`                |

##### _example_ full command

###### _docker_ and podman

they share the same flags

```sh
docker run --rm -it \
  --name some-container \
  --runtime=runsc \
  -e MY_ENV=1 \
  -v $(pwd)/data:/data \
  -p 80:8080 \
  $img
```

###### _ctr_

less capable

```sh
ctr run --rm -t \
  --runc-binary=runsc \
  --env MY_ENV=1 \
  --mount src=$(pwd)/data,dst=/data \
  $img
```

##### _multi_

- [docker-compose][docker-compose]: `docker-compose up/down`
- [podman-compose][podman-compose]: `podman-compose up/down`
- kubernetes single node: k3s, kind
- crictl can run mutliple containers together in a pod

#### inspect

- [dive][dive]: interactive terminal container layer browser / file diff
- [container-diff][diff]: high level container diff (eg packages)
- docker/podman diff: show file has changed between layers

[buildah]: https://github.com/containers/buildah
[crictl]: https://github.com/kubernetes-sigs/cri-tools
[ctr]: https://github.com/containerd/containerd
[diff]: https://github.com/GoogleContainerTools/container-diff
[dive]: https://github.com/wagoodman/dive
[docker]: https://github.com/docker/docker
[docker-compose]: https://github.com/docker/compose
[img]: https://github.com/genuinetools/img
[kaniko]: https://github.com/GoogleContainerTools/kaniko
[podman]: https://github.com/containers/podman
[podman-compose]: https://github.com/containers/podman-compose
[skopeo]: https://github.com/containers/skopeo
