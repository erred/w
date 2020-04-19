---
description: using k8s secrets
title: kubernetes secrets
---

### _using_ secrets

#### environment

mount all secret key-values as env

```yaml
# podspec
spec:
  containers:
    - ...
      envFrom:
        - secretRef:
            name: name-of-secret
```

mount as env, per key-value

```yaml
# podspec
spec:
  containers:
    - ...
      env:
        - name: NAME_OF_ENV
          valueFrom:
            secretKeyRef:
              name: name-of-secret
              key: key-in-secret
```

#### files

mount as files in dir

```yaml
# podspec
spec:
  containers:
    - ...
      volumeMounts:
        - name: name-of-volume
          mountPath: /etc/foo
  volumes:
    - name: name-of-volume
      secret:
        secretName: name-of-secret
```

mount as files, per key-value

```yaml
# podspec
spec:
  containers:
    - ...
      volumeMounts:
        - name: name-of-volume
          mountPath: /etc/foo
  volumes:
    - name: name-of-volume
      secret:
        secretName: name-of-secret
        items:
          - key: key-in-secret
            path: path-to-mount-as
```

#### serviceaccounts

`imagePullSecrets`: useful for pulling private images

`secrets`: useless as far as i can tell

```yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: service-account-name
secrets:
  - name-of-secret
imagePullSecrets:
  - name-of-secret
```

#### pod presets

apply same config to all matching pods

see above for env / volume format

```yaml
apiVersion: settings.k8s.io/v1alpha1
kind: PodPreset
metadata:
  name: allow-database
spec:
  selector:
    matchLabels:
      label-key: label-value
  env: ...
  volumeMounts: ...
  volume: ...
```

### secret _types_

what types can we use,
apparently they're (almost) all for use with the k8s api

`data` is _base64 encoded_ values,
replace with `stringData` (plaintext) for convenience

start with the [source code](https://github.com/kubernetes/kubernetes/blob/master/pkg/apis/core/types.go#L4806)

#### Opaque

default type

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: name-of-secret
type: Opaque
data:
  user-defined-key: data
```

#### kubernetes.io/service-account-token

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: name-of-secret
  annotations:
    kubernetes.io/service-account.name: name of ServiceAccount
    kubernetes.io/service-account.uid: uid of ServiceAccount
type: kubernetes.io/service-account-token
data:
  token: token
  kubernetes.kubeconfig: kubeconfig (optional)
  ca.crt: root certificate (optional)
  namespace: default namespace (optional)
```

#### kubernetes.io/dockercfg

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: name-of-secret
type: kubernetes.io/dockercfg
data:
  .dockercfg: ~/.dockercfg
```

#### kubernetes.io/dockerconfigjson

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: name-of-secret
type: kubernetes.io/dockerconfigjson
data:
  .dockerconfigjson: ~/.docker/config.json
```

#### kubernetes.io/basic-auth

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: name-of-secret
type: kubernetes.io/basic-auth
data:
  username: username
  password: password
```

#### kubernetes.io/ssh-auth

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: name-of-secret
type: kubernetes.io/ssh-auth
data:
  ssh-privatekey: private key
```

#### kubernetes.io/tls

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: name-of-secret
type: kubernetes.io/tls
data:
  tls.crt: certificate (public)
  tls.key: key (private)
```
