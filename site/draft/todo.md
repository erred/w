# Ideas

- bazel
- smtp

  - postfix
  - exim
  - https://github.com/albertito/chasquid
  - sendmail
  - opensmtpd
  - qmail
  - https://github.com/flashmob/go-guerrilla
  - https://github.com/mhale/smtpd
  - https://github.com/emersion/go-smtp
  - https://mailinabox.email/
  - https://www.iredmail.org/
  - https://modoboa.org/en/
  - https://mailu.io/1.7/
  - https://github.com/tomav/docker-mailserver
  - https://wildduck.email/
  - https://archiveopteryx.org/
  - http://www.dbmail.org/
  - http://www.elasticinbox.com/
  - https://mailcow.github.io/mailcow-dockerized-docs/
  - https://james.apache.org/

- tsdb
  - influxdb
  - opentsdb
  - timescaledb
  - graphite

## k8s plan

table stakes?

- istio
- cert-manager
- otel-collector
- promtail
- loki
- prometheus
- jaeger
- grafana

### knative stack

per request scaling, no sync to git state

- knative
- tekton

### argo stack

sync to git state, old scaling

argo rollouts
argo events
argo cd
argo workflows

### learn

https://cloud.google.com/blog/topics/anthos/introducing-the-anthos-developer-sandbox
https://go.qwiklabs.com/qwiklabs-free?utm_source=google&utm_medium=lp&utm_campaign=GKE

### secure boot

https://keylime.dev/
https://safeboot.dev/
https://media.defense.gov/2020/Sep/15/2002497594/-1/-1/0/CTR-UEFI-SECURE-BOOT-CUSTOMIZATION-20200915.PDF/CTR-UEFI-SECURE-BOOT-CUSTOMIZATION-20200915.PDF
https://mjg59.dreamwidth.org/35742.html
https://github.com/google/go-tpm-tools
https://github.com/google/go-attestation
