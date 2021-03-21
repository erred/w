---
title: k8s cluster metrics
description: where to get your stats from
---

### _k8s_

The DIY PaaS, and today I'm thinking about how to measure CPU and memory
in this giant jigsaw puzzle.

#### _collecting_ data

##### _cluster_ level

The kubernetes control plane exposes about itself through its apiservers.
Additionally, [kube-state-metrics](https://github.com/kubernetes/kube-state-metrics)
generates metrics about the things inside the cluster.

_TODO:_
decide if we need to replace the `instance` label to a stable name in the case of multiple instances

```yaml
scrape_configs:
  - job_name: kubernetes-apiservers
    kubernetes_sd_configs:
      - role: endpoints
    scheme: https
    tls_config:
      ca_file: /var/run/secrets/kubernetes.io/serviceaccount/ca.crt
    bearer_token_file: /var/run/secrets/kubernetes.io/serviceaccount/token
    relabel_configs:
      - source_labels:
          - __meta_kubernetes_namespace
          - __meta_kubernetes_service_name
          - __meta_kubernetes_endpoint_port_name
        action: keep
        regex: default;kubernetes;https

  # note kube-state-metrics also has an alternate port with metrics about itself
  - job_name: kube-state-metrics
    kubernetes_sd_configs:
      - role: endpoints
    relabel_configs:
      - source_labels:
          - __meta_kubernetes_namespace
          - __meta_kubernetes_service_name
          - __meta_kubernetes_endpoint_port_name
        action: keep
        regex: kube-state-metrics;kube-state-metrics;kube-state-metrics
```

##### _node_ level

Kubernetes kubelets expose both their own metrics
and the metrics on pods runnin on their node

```yaml
scrape_configs:
  - job_name: kubernetes-nodes
    kubernetes_sd_configs:
      - role: node
    scheme: https
    tls_config:
      ca_file: /var/run/secrets/kubernetes.io/serviceaccount/ca.crt
      insecure_skip_verify: true
    bearer_token_file: /var/run/secrets/kubernetes.io/serviceaccount/token
    relabel_configs:
      - action: labelmap
        regex: __meta_kubernetes_node_label_(.+)

  - job_name: kubernetes-cadvisor
    kubernetes_sd_configs:
      - role: node
    scheme: https
    metrics_path: /metrics/cadvisor
    tls_config:
      ca_file: /var/run/secrets/kubernetes.io/serviceaccount/ca.crt
      insecure_skip_verify: true
    bearer_token_file: /var/run/secrets/kubernetes.io/serviceaccount/token
    relabel_configs:
      - action: labelmap
        regex: __meta_kubernetes_node_label_(.+)

  - job_name: node-exporter
    kubernetes_sd_configs:
      # remember to create a service to capture expose this
      # using role: node is also possible if you expose a hostport
      - role: endpoints
    relabel_configs:
      - source_labels:
          - __meta_kubernetes_namespace
          - __meta_kubernetes_service_name
          - __meta_kubernetes_endpoint_port_name
        action: keep
        regex: node-exporter;node-exporter;node-exporter
      - action: replace # rename the instance from the discovered pod ip (we're using endpoints) to the node name
        target_label: instance
        source_labels:
          - __meta_kubernetes_pod_node_name
```

##### _dumping_ metrics

prometheus has a useful `/federate` endpoint you can use to dump out everything
after relabelling, example query:

```txt
http://localhost:8080/federate?match[]={job=~".*"}
```

##### _example_

repo: [testrepo-cluster-metrics](https://github.com/seankhliao/testrepo-cluster-metrics)

#### _future_

In the future we might be able to run
`N+1` (`N` = number of nodes) instances of
[opentelemetry-collector](https://github.com/open-telemetry/opentelemetry-collector),
`N` as agents (DaemonSet) and `1` as gateway,
replacing the need for `N` node-exporter, `1` kube-state-metrics, as well as `N+1` tracing
collectors and `N+1` logging collectors.
As it currently stands, it still needs some more work to export the metrics in a stable manner,
and maybe some extra exporters to write directly to storage.
