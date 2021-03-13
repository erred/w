---
description: secrets for docker build
title: docker build secrets
---

### _docker_ build

It's 12020 and we finally have a supported way of using secrets
when building docker containers.

#### _secret_ files

Enable [buildkit](https://docs.docker.com/develop/develop-images/build_enhancements/#new-docker-build-secret-information)
and `RUN --mount=type=secret,id=mysecret cat /run/secrets/mysecret`.
Specify on the command line with `docker build --secret id=mysecret,src=secret.txt`
Useful for, say, go with `netrc`:

```Dockerfile
RUN --mount=type=secret,id=netrc,dst=/root/.netrc go get example.com/my/private/repo
```

#### _ssh_ agent forwarding

someone else can probably [explain it better](https://medium.com/@tonistiigi/build-secrets-and-ssh-forwarding-in-docker-18-09-ae8161d066)
since I don't use ssh-agent

tldr is `docker build --ssh ...`:

```Dockerfile
RUN  --mount=type=ssh go get example.com/my/private/repo.git
```

#### _other_

still waiting for a way to mount build caches
