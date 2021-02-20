---
description: an opionionated and commented deployment manifest for generic apps
title: opionionated k0s deployment
---

### _kubernetes_ manifests

YAML engineer reporting in.

_note:_ yaml is long and repetitive,
I'm still not sure if I'm happy I introduced [yaml anchors](https://yaml.org/spec/1.2/spec.html#anchor//)
to my team. tldr, the 2 docs below are equivalent, anchors do not carry accross documents (`---`):

```yaml
name: &name foo
somewhere:
  else:
    x: *name
---
name: foo
somewhere:
  else:
    x: foo
```

#### metadata

Every object has them: names, labels, annotations.
They even have a [recommended set of labels](https://kubernetes.io/docs/concepts/overview/working-with-objects/common-labels/)

```yaml
metadata:
  name: foo
  annotations:
    # stick values that you don't want to filter by here,
    # such as info for other apps that read service definitions
    # or as a place to store data to make your controller stateless
  labels:
    # sort of duplicates metadata.name
    app.kubernetes.io/name: foo

    # separate multiple instances, not really necessary if you do app-per-namespace
    app.kubernetes.io/instance: default

    # you might not want to add this on everything (eg namespaces, security stuff)
    # since with least privilege you can't change them
    # and they don't really change that often(?)
    app.kubernetes.io/version: "1.2.3"

    # the hardest part is probably getting it to not say "helm" when you don't actually use helm
    app.kubernetes.io/managed-by: helm

    # these two aren't really necessary for single deployment apps
    #
    # the general purpose of "name", eg name=envoy component=proxy
    app.kubernetes.io/component: server
    # what the entire this is
    app.kubernetes.io/part-of: website
```

#### _namespace_

The hardest part about namespaces is your namespace allocation policy,
do you:

- dump everything in a single namespace (`default`?)
- give each team their own namespace
- give each app their own namespace
- give each app revision their own namespace

[Hierarchical Namespaces](https://github.com/kubernetes-sigs/multi-tenancy/tree/master/incubator/hnc)
might help a bit, making the latter ones more tenable but still, things to think about.

Currently I'm in the "each app their own namespace" camp,
and live with the double names in service addresses

```yaml
apiVersion: v1
kind: Namespace
metadata:
  name: foo
```

#### _ingress_

The least common denominator of L4/L7 routing...

```kubernetes
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: foo
spec:
  # for if you run multiple ingress controllers
  ingressClassName: default

  rules:
      # DNS style wildcards only
    - host: "*.example.com"
      http:
        paths:
          - path: /
            pathType: Prefix # or Exact, prefix uses path segment matching
            backend:
              service:
                name: foo
                port:
                  name: http
                  # number: 80

  tls:
    secretName: foo-tls
    hosts:
      - "*.example.com"
```

#### _service_

```yaml
apiVersion: v1
kind: Service
metadata:
  name: foo
spec:
  # change as needed
  type: ClusterIP

  # only for type LoadBalancer
  externalTrafficPolicy: Local

  # for statefulsets that need peer discovery,
  # eg. etcd or cockroachdb
  publishNotReadyAddresses: true

  ports:
    - appProtocol: opentelemetry
      name: otlp
      port: 4317
      protocol: TCP
      targetPort: otlp # name or number, defaults to port

  selector:
    matchLabels:
      # these 2 should be enough to uniquely identify apps,
      # note this value cannot change once created
      app.kubernetes.io/name: foo
      app.kubernetes.io/instance: default
```

#### _serviceaccount_

note: while it does have a `spec.secrets` field, it currently doesn't really do anything useful.

```yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: foo
  annotations:
    # workload identity for attaching to GCP service accounts in GKE
    iam.gke.io/gcp-service-account: GSA_NAME@PROJECT_ID.iam.gserviceaccount.com
```

#### _app_

##### _deployment_

Use only if you app is truly stateless:
no PersistentVolumeClaims unless it's `ReadOnlyMany`,
even then PVCs still restrict the nodes you can run on.

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: foo
spec:
  # don't set if you plan on autoscaling
  replicas: 1

  # stop cluttering kubectl get all with old replicasets,
  # your gitops tooling should let you roll back
  revisionHistoryLimit: 3

  selector:
    matchLabels:
      # these 2 should be enough to uniquely identify apps,
      # note this value cannot change once created
      app.kubernetes.io/name: foo
      app.kubernetes.io/instance: default

  # annoyingly named differently from StatefulSet or DaemonSet
  strategy:
    # prefer maxSurge to keep availability during upgrades / migrations
    rollingUpdate:
      maxSurge: 25% # rounds up
      maxUnavailable: 0

    # Recreate if you want blue-green style
    # or if you're stuck with a PVC
    type: RollingUpdate

  template: # see pod below
```

##### _statefulset_

If your app has any use for persistent data, use this,
even if you only have a single instance.
Also gives you nice DNS names per pod.

```yaml
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: foo
spec:
  # or Parallel for all at once
  podManagementPolicy: OrderedReady
  replicas: 3

  # stop cluttering kubectl get all with old replicasets,
  # your gitops tooling should let you roll back
  revisionHistoryLimit: 3

  selector:
    matchLabels:
      # these 2 should be enough to uniquely identify apps,
      # note this value cannot change once created
      app.kubernetes.io/name: foo
      app.kubernetes.io/instance: default

  # even though they say it must exist, it doesn't have to
  # (but you lose per pod DNS)
  serviceName: foo

  template: # see pod below

  updateStrategy:
    rollingUpdate: # this should only be used by tooling
    type: RollingUpdate

  volumeClaimTemplates: # see pvc below
```

##### _daemonset_

```yaml
apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: foo
spec:
  # stop cluttering kubectl get all with old replicasets,
  # your gitops tooling should let you roll back
  revisionHistoryLimit: 3

  selector:
    matchLabels:
      # these 2 should be enough to uniquely identify apps,
      # note this value cannot change once created
      app.kubernetes.io/name: foo
      app.kubernetes.io/instance: default

  template: # see pod below

  updateStrategy:
    rollingUpdate:
      # make it faster for large clusters
      maxUnavailable: 30%
    type: RollingUpdate
```

##### _pod_

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: foo
spec:
  containers:
    - name: foo
      args:
        - -flag1=v1
        - -flag2=v2
      envFrom
        - configMapRef:
            name: foo-env
            optional: true
          prefix: APP_
      image: docker.example.com/app:v1
      imagePullPolicy: IfNotPresent

      ports:
        - containerPort: 4317
          name: otlp
          protocol: TCP

      # do extra stuff
      lifecycle:
        postStart:
        preStop:

      startupProbe: # allow a longer startup
      livenessProbe: # stay alive to not get killed
      readinessProbe: # stay alive to route traffic

      securityContext:
        allowPrivilegeEscalation: false
        capabilities:
          add:
            - CAP_NET_ADMIN
        privileged: false
        readOnlyRootFilesystem: true

      resources:
        # ideally set after running some time and profiling actual usage,
        # prefer to start high and rachet down
        requests:
          cpu: 500m
          memory: 128Mi
        limits:
          cpu: 1500m
          memory: 512Mi

      volumeMounts: # as needed

  # don't inject env with service addresses/ports
  # not many things use them, they clutter up the env
  # and may be a performance hit with large number of services
  enableServiceLinks: false

  # do create PriorityClasses and every pod one,
  # helps with deciding which pods to kill first
  priorityClassName: critical

  securityContext:
    fsGroup: 65535
    runAsGroup: 65535
    runAsNonRoot: true
    runAsUser: 65535 # may conflict with container setting and need for $HOME

  serviceAccountName: foo

  terminationGracePeriodSeconds: 30

  volumes: # set as needed
```

###### _scheduling_

theres is some overlap in managing pod scheduling, especially around where they run:

- `affinity`: these only let you select to either run 0 or unlimited pods per selector
- `affinity.nodeAffinity`: general purpose choose a node
- `affinity.podAffinity`: general purpose choose to schedule next to things
- `affinity.podAntiAffinity`: general purpose choose not to schedule next to things
- `nodeSelector`: shorthand for choosing nodes with labels
- `tolerations`: allow scheduling on nodes with taints
- `topologySpreadConstraints`: choose how many to schedule in a single topology domain

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: foo
spec:
  affinity:
    nodeAffinity:
      requiredDuringSchedulingIgnoredDuringExecution:
        nodeSelectorTerms: # OR
          # has to be pool-0
          - matchExpressions: # AND
              - key: cloud.google.com/gke-nodepool
                operator: In
                values:
                  - pool-0
      preferredDuringSchedulingIgnoredDuringExecution
        # prefer zone us-central1-a
        - weight: 25
          preference:
            - matchExpressions: # AND
              - key: topology.kubernetes.io/zone
                operator: In
                values:
                  - us-central1-a

    podAffinity:
      preferredDuringSchedulingIgnoredDuringExecution:
        # prefer to be on the same node as a bar
        - weight: 25
          podAffinityTerm:
            labelSelector:
              matchLabels:
                app.kubernetes.io/name: bar
                app.kubernetes.io/instance: default
            topologyKey: kubernetes.io/hostname

    podAntiAffinity:
      requiredDuringSchedulingIgnoredDuringExecution: # AND
        # never schedule in the same region as buzz
        - labelSelector:
            matchLabels:
              app.kubernetes.io/name: buzz
              app.kubernetes.io/instance: default
          topologyKey: topology.kubernetes.io/region


  topologySpreadConstraints: # AND
    # limit to 1 instance per node
    - maxSkew: 1
      labelSelector:
        matchLabels:
          app.kubernetes.io/name: foo
          app.kubernetes.io/instance: default
      topologyKey: kubernetes.io/hostname
      whenUnsatisfiable: DoNotSchedule # or ScheduleAnyway
```

###### _command_

- none: docker default
- `args`: Docker entrypoint + container args
- `command`: container command
- `command` and `args`: container command + container args

#### _persistentvolumeclaim_

```yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: foo
spec:
  accessModes: ReadWriteOnce # ReadOnlyMany or ReadWriteMany (rare)

  dataSource: # prepopulate with data from a VolumeSnapshot or PersistentVolumeClaim

  resources:
    requests:
      storage: 10Gi

  # bind to existing PV
  selector: matchLabels

  storageClassName: ssd

  volumeMode: Filesystem # or Block
```

#### _horizontalpodautoscaler_

```yaml
apiVersion: autoscaling/v2beta2
kind: HorizontalPodAutoscaler
metadata:
  name: foo
spec:
  behavior: # fine tune when to scale up / down

  maxReplicas: 5
  minReplicas: 1

  metrics:
    -  # TODO

  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: foo
```

#### _poddisruptionbudget_

```yaml
apiVersion: policy/v1beta1
kind: PodDisruptionBudget
metadata:
  name: foo
spec:
  # when you have a low number of replicas
  # ensure you can disrupt them
  maxUnavailable: 1

  # allows for more disruptions
  minAvailable: 75%

  selector:
    matchLabels:
      app.kubernetes.io/name: foo
      app.kubernetes.io/instance: default
```
