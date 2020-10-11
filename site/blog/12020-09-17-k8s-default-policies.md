---
description: k8s config files are long, use common default policies to shorten them
title: k8s default policies
---

### _common_ config

Kubernetes config files are long/verbose enough as it is
and a lot of the time it's just the same boilerplate repeated over and over.

#### _todo_

[PodPreset](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#podpresetspec-v1alpha1-settings-k8s-io) currently in alpha, for volumes and env

[NetworkPolicy](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#networkpolicy-v1-networking-k8s-io) needs support from CNI

Wishlist: readiness/liveliness probe and securitycontext

#### _resources_

[LimitRange](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#limitrange-v1-core)
can set per namespace defaults and limits on resource requests and limits.
This is nice in that it fails closed,
ie the default is applied unless explicitly overrriden.

Note: policies do not get retroactively applied.

```yaml
apiVersion: v1
kind: LimitRange
metadata:
  name: default-limitrange
spec:
  limits:
    - type: Container
      default:
        cpu: 100m
        memory: 100Mi
      defaultRequest:
        cpu: 50m
        memory: 50Mi
      # min:
      # max:
```

#### _security_

[PodSecurutyPolicy](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#podsecuritypolicy-v1beta1-policy)
is tied to service accounts, one can be bound to the default service account.
This is less nice as specifying a different service account
will bypass this policy.
Instead policies have to be enforced by an
[Admission Controller](https://kubernetes.io/docs/reference/access-authn-authz/admission-controllers/).

Note: [default and additional capabilities](https://docs.docker.com/engine/reference/run/#runtime-privilege-and-linux-capabilities)

```yaml
apiVersion: policy/v1beta1
kind: PodSecurityPolicy
metadata:
  annotations:
    seccomp.security.alpha.kubernetes.io/allowedProfileNames: "runtime/default"
    seccomp.security.alpha.kubernetes.io/defaultProfileName: "runtime/default"
    apparmor.security.beta.kubernetes.io/allowedProfileNames: "runtime/default"
    apparmor.security.beta.kubernetes.io/defaultProfileName: "runtime/default"
  name: default-security
spec:
  # # user
  runAsGroup:
    rule: "MustRunAs"
    ranges:
      - min: 1
        max: 65535
  runAsUser:
    rule: "MustRunAsNonRoot"
  fsGroup:
    rule: "MustRunAs"
    ranges:
      - min: 1
        max: 65535
  supplementalGroups:
    rule: "MustRunAs"
    ranges:
      - min: 1
        max: 65535
  # # privilege
  privileged: false
  allowPrivilegeEscalation: false
  defaultAllowPrivilegeEscalation: false
  # # host
  hostIPC: false
  hostNetwork: false
  hostPID: false
  hostPorts: []
  # allowedHostPaths: []
  # allowedProcMountTypes: []
  # allowedUnsafeSysctls: []
  # forbiddenSysctls: []
  seLinux:
    rule: "RunAsAny"
  # # capabilities
  # allowedCapabilities:
  # defaultAddCapabilities:
  requiredDropCapabilities:
    - ALL
    # # Default set from Docker, without DAC_OVERRIDE or CHOWN
    # - FOWNER
    # - FSETID
    # - KILL
    # - SETGID
    # - SETUID
    # - SETPCAP
    # - NET_BIND_SERVICE
    # - NET_RAW
    # - SYS_CHROOT
    # - MKNOD
    # - AUDIT_WRITE
    # - SETFCAP
  # # filesystem
  readOnlyRootFilesystem: false
  volumes: [] # may want to relax this for configMap / persistentVolumeClaim / ...
---
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: default-security
rules:
  - apiGroups:
      - policy
    resources:
      - podsecuritypolicies
    verbs:
      - use
    resourceNames:
      - default-security
---
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: RoleBinding
metadata:
  name: default-security
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: default-security
subjects:
  - kind: ServiceAccount
    name: default
```
