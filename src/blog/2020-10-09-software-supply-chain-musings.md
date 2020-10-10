---
description: starting point for thinking about software supply chain security
title: software supply chain musings
---

### _software_ supply chain

attacker end goals:

- use compute resources: cryptominers
- extract data: keys, wallets, user/application data
- peristent access

#### _chain_

assuming trusted local dev environment

- code
  - permissions to modify access control
  - written by project developers
    - developers hold trusted commit keys
    - block unauthorized code
      - transport layer: https / ssh
      - commit level: signed
  - dependencies, imported
    - pinned, audited versions
    - content addressed or vendored
- continuous integration
  - permissions to modify the pipeline
  - security scans
    - source code level
    - dependency versions
    - built artifacts scan
  - compiler / packaging
    - trusted not to insert backdoors?
    - reproducible builds
    - also signed?
      - verify the artifacts came through trusted pipeline
  - ci system holds trusted keys for pushing artifacts
- continuous deployment
  - permissions to modify the pipeline
  - only deploy trusted artifacts, signed?
  - push trigger vs pull:
    - push from ci / same system as ci: ci compromise == cd compromise, but faster
    - watch and pull from artifact store: safer, slower
  - cd system holds trusted keys for deploying to production
- execution environment
  - permissions to access the environment
  - permissions to modify access control
  - environment needs to be kept up to date
  - only run artifacts with a clean audit trail, signed?
  - execution environment holds trusted keys for accessing application data

#### _other_

- trusted environment
  - for dev and prod
  - root of trust in hardware
  - secure boot + ...
