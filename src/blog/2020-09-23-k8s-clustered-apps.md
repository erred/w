---
description: starting clustered applications in k8s
title: k8s clustered apps
---

### _clustered_

So you want to run a clustered thing in k8s?
Likely a database, using [raft](https://raft.github.io/)
or similar.

Use a statefulset: you get your own persistentvolume per pod,
and you get your own stable, addressable hostname.
This can be retrieved either as the `HOSTNAME` env var (possibly unstable?),
or set by custom env var with fieldref.

Use cert-manager: who wants to futz around with csrs

use `publishNotReadyAddresses: true` on a headless service to get name resolution
before pods are ready, pods need to see each other before they are ready.

#### _etcd_

Notes: data dir should be changed,
Assumes a CA is available and called _internal-ca_.

Just works (I think)

```yaml
apiVersion: cert-manager.io/v1
kind: Certificate
metadata:
  name: etcd-certs
spec:
  secretName: etcd-certs
  duration: 2160h
  renewBefore: 360h
  dnsNames:
    - "localhost"
    - "etcd"
    - "etcd.default"
    - "etcd.default.svc"
    - "etcd.default.svc.cluster.local"
    - "*.etcd-headless"
    - "*.etcd-headless.default"
    - "*.etcd-headless.default.svc"
    - "*.etcd-headless.default.svc.cluster.local"
  ipAddresses:
    - "127.0.0.1"
    - "::1"
  issuerRef:
    name: internal-ca
    kind: ClusterIssuer
---
apiVersion: v1
kind: Secret
metadata:
  name: etcd
  labels:
    app.kubernetes.io/name: etcd
type: Opaque
data:
  etcd-root-password: "eDgzelB1aVlsUQ=="
---
apiVersion: v1
kind: Service
metadata:
  name: etcd-headless
  labels:
    app.kubernetes.io/name: etcd
spec:
  type: ClusterIP
  clusterIP: None
  publishNotReadyAddresses: true
  ports:
    - name: client
      port: 2379
      targetPort: client
    - name: peer
      port: 2380
      targetPort: peer
  selector:
    app.kubernetes.io/name: etcd
---
apiVersion: v1
kind: Service
metadata:
  name: etcd
  labels:
    app.kubernetes.io/name: etcd
spec:
  type: ClusterIP
  ports:
    - name: client
      port: 2379
      targetPort: client
    - name: peer
      port: 2380
      targetPort: peer
  selector:
    app.kubernetes.io/name: etcd
---
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: etcd
  labels:
    app.kubernetes.io/name: etcd
spec:
  selector:
    matchLabels:
      app.kubernetes.io/name: etcd
  serviceName: etcd-headless
  podManagementPolicy: Parallel
  replicas: 3
  updateStrategy:
    type: RollingUpdate
  template:
    metadata:
      labels:
        app.kubernetes.io/name: etcd
    spec:
      securityContext:
        fsGroup: 1001
        runAsUser: 1001
      containers:
        - name: etcd
          image: docker.io/bitnami/etcd:3.4.13-debian-10-r22
          imagePullPolicy: "IfNotPresent"
          command:
            - etcd
          env:
            - name: POD_NAME
              valueFrom:
                fieldRef:
                  fieldPath: metadata.name
            - name: ETCDCTL_API
              value: "3"
            - name: ETCD_NAME
              value: "$(POD_NAME)"
            - name: ETCD_DATA_DIR
              value: /bitnami/etcd/data
            - name: ETCD_ADVERTISE_CLIENT_URLS
              value: "https://$(POD_NAME).etcd-headless.default.svc.cluster.local:2379"
            - name: ETCD_LISTEN_CLIENT_URLS
              value: "https://0.0.0.0:2379"
            - name: ETCD_INITIAL_ADVERTISE_PEER_URLS
              value: "https://$(POD_NAME).etcd-headless.default.svc.cluster.local:2380"
            - name: ETCD_LISTEN_PEER_URLS
              value: "https://0.0.0.0:2380"
            - name: ALLOW_NONE_AUTHENTICATION
              value: "yes"
            - name: ETCD_ROOT_PASSWORD
              valueFrom:
                secretKeyRef:
                  name: etcd
                  key: etcd-root-password
            - name: ETCD_INITIAL_CLUSTER
              value: "etcd-0=https://etcd-0.etcd-headless.default.svc.cluster.local:2380,etcd-1=https://etcd-1.etcd-headless.default.svc.cluster.local:2380,etcd-2=https://etcd-2.etcd-headless.default.svc.cluster.local:2380"
            - name: ETCD_INITIAL_CLUSTER_STATE
              value: new
            - name: ETCD_CLIENT_CERT_AUTH
              value: "true"
            - name: ETCD_TRUSTED_CA_FILE
              value: /var/secret/tls/ca.crt
            - name: ETCD_CERT_FILE
              value: /var/secret/tls/tls.crt
            - name: ETCD_KEY_FILE
              value: /var/secret/tls/tls.key
            - name: ETCD_PEER_CLIENT_CERT_AUTH
              value: "true"
            - name: ETCD_PEER_TRUSTED_CA_FILE
              value: /var/secret/tls/ca.crt
            - name: ETCD_PEER_CERT_FILE
              value: /var/secret/tls/tls.crt
            - name: ETCD_PEER_KEY_FILE
              value: /var/secret/tls/tls.key
          ports:
            - name: client
              containerPort: 2379
            - name: peer
              containerPort: 2380
            - name: metrics
              containerPort: 2381
          livenessProbe:
            httpGet:
              path: /health
              port: 2381
          readinessProbe:
            httpGet:
              path: /health
              port: 2381
          volumeMounts:
            - name: certs
              mountPath: /var/secret/tls
            - name: data
              mountPath: /bitnami/etcd
      volumes:
        - name: certs
          secret:
            secretName: etcd-certs
        - name: data
          emptyDir: {}
```

#### _cockroachdb_

Notes: will complain if certs have wider perms than `rwx------`,
which will cause issues if running as non-root in k8s (uses fsGroups to keep volume owner as root).
Slightly modified from official manifests
(change service names, certs, data dir for kind).
Assumes a CA is available and called _internal-ca_.

Nodes need a manual action to join,
could be a Job but need to time it right (after cert signing, nodes started):

```sh
kubectl exec -it cockroachdb-0 -- /cockroach/cockroach init --certs-dir=/cockroach/cockroach-certs
```

manifest:

```yaml
apiVersion: cert-manager.io/v1
kind: Certificate
metadata:
  name: cockroachdb-node
spec:
  secretName: cockroachdb-node
  duration: 2160h
  renewBefore: 360h
  dnsNames:
    - "node"
    - "localhost"
    - "cockroachdb"
    - "cockroachdb.default"
    - "cockroachdb.default.svc"
    - "cockroachdb.default.svc.cluster.local"
    - "*.cockroachdb-headless"
    - "*.cockroachdb-headless.default"
    - "*.cockroachdb-headless.default.svc"
    - "*.cockroachdb-headless.default.svc.cluster.local"
  ipAddresses:
    - "127.0.0.1"
    - "::1"
  issuerRef:
    name: internal-ca
    kind: ClusterIssuer
---
apiVersion: cert-manager.io/v1
kind: Certificate
metadata:
  name: cockroachdb-client-root
spec:
  secretName: cockroachdb-client-root
  duration: 2160h
  renewBefore: 360h
  commonName: root
  usages:
    - client auth
  issuerRef:
    name: internal-ca
    kind: ClusterIssuer
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: cockroachdb
  labels:
    app: cockroachdb
---
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: Role
metadata:
  name: cockroachdb
  labels:
    app: cockroachdb
rules:
  - apiGroups:
      - ""
    resources:
      - secrets
    verbs:
      - get
---
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: RoleBinding
metadata:
  name: cockroachdb
  labels:
    app: cockroachdb
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: cockroachdb
subjects:
  - kind: ServiceAccount
    name: cockroachdb
    namespace: default
---
apiVersion: v1
kind: Service
metadata:
  name: cockroachdb
  labels:
    app: cockroachdb
spec:
  ports:
    - port: 26257
      targetPort: 26257
      name: grpc
    - port: 8080
      targetPort: 8080
      name: http
  selector:
    app: cockroachdb
---
apiVersion: v1
kind: Service
metadata:
  name: cockroachdb-headless
  labels:
    app: cockroachdb
  annotations:
    prometheus.io/scrape: "true"
    prometheus.io/path: "_status/vars"
    prometheus.io/port: "8080"
spec:
  ports:
    - port: 26257
      targetPort: 26257
      name: grpc
    - port: 8080
      targetPort: 8080
      name: http
  publishNotReadyAddresses: true
  clusterIP: None
  selector:
    app: cockroachdb
---
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: cockroachdb
spec:
  serviceName: "cockroachdb-headless"
  replicas: 3
  selector:
    matchLabels:
      app: cockroachdb
  template:
    metadata:
      labels:
        app: cockroachdb
    spec:
      serviceAccountName: cockroachdb
      affinity:
        podAntiAffinity:
          preferredDuringSchedulingIgnoredDuringExecution:
            - weight: 100
              podAffinityTerm:
                labelSelector:
                  matchExpressions:
                    - key: app
                      operator: In
                      values:
                        - cockroachdb
                topologyKey: kubernetes.io/hostname
      containers:
        - name: cockroachdb
          image: cockroachdb/cockroach:v20.1.5
          imagePullPolicy: IfNotPresent
          ports:
            - containerPort: 26257
              name: grpc
            - containerPort: 8080
              name: http
          livenessProbe:
            httpGet:
              path: "/health"
              port: http
              scheme: HTTPS
            initialDelaySeconds: 30
            periodSeconds: 5
          readinessProbe:
            httpGet:
              path: "/health?ready=1"
              port: http
              scheme: HTTPS
            initialDelaySeconds: 10
            periodSeconds: 5
            failureThreshold: 2
          volumeMounts:
            - name: datadir
              mountPath: /cockroach/cockroach-data
            - name: certs
              mountPath: /cockroach/cockroach-certs/ca.crt
              subPath: ca.crt
            - name: certs
              mountPath: /cockroach/cockroach-certs/node.crt
              subPath: tls.crt
            - name: certs
              mountPath: /cockroach/cockroach-certs/node.key
              subPath: tls.key
            - name: client
              mountPath: /cockroach/cockroach-certs/client.root.crt
              subPath: tls.crt
            - name: client
              mountPath: /cockroach/cockroach-certs/client.root.key
              subPath: tls.key
          env:
            - name: COCKROACH_CHANNEL
              value: kubernetes-secure
            - name: POD_NAME
              valueFrom:
                fieldRef:
                  fieldPath: metadata.name
            - name: GOMAXPROCS
              valueFrom:
                resourceFieldRef:
                  resource: limits.cpu
                  divisor: "1"
            - name: MEMORY_LIMIT_MIB
              valueFrom:
                resourceFieldRef:
                  resource: limits.memory
                  divisor: "1Mi"
          command:
            - /bin/sh
            - -exc
            - >
              /cockroach/cockroach
              start
              --logtostderr=WARNING
              --certs-dir=/cockroach/cockroach-certs
              --advertise-host=$(POD_NAME).cockroachdb-headless.default
              --http-addr=0.0.0.0
              --join=cockroachdb-0.cockroachdb-headless.default,cockroachdb-1.cockroachdb-headless.default,cockroachdb-2.cockroachdb-headless.default
      terminationGracePeriodSeconds: 60
      volumes:
        - name: datadir
          emptyDir: {}
        - name: certs
          secret:
            secretName: cockroachdb-node
            defaultMode: 256
        - name: client
          secret:
            secretName: cockroachdb-client-root
            defaultMode: 256
  podManagementPolicy: Parallel
  updateStrategy:
    type: RollingUpdate
```
