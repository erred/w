---
description: secrets and gitops and k8s
title: k8s gitops secrets
---

### _gitops_

The idea that you can track your desired state
(of entire clusters) in git and some tooling will realize it for you.
Problem is secrets.
Dumping plain secrets into git is a bad idea,
especially if the repo is going to be public.

#### _options_

##### [_bitnami-labs/sealed-secrets_](https://github.com/bitnami-labs/sealed-secrets)

Secrets are encrypted with a public key and stored as a CRD.
A controller in cluster will watch CRDs and decrypt with a private key,
creating appropriate Secrets.

##### [_external-secrets/kubernetes-external-secrets_](https://github.com/external-secrets/kubernetes-external-secrets)

Secrets are stored with cloud platform specific secret managers (or Hashipcorp Vault).
A controller watches for CRDs and retrieves the secrets.

Supports: aws, hashicorp vault, azure, alibaba, gcp

###### _alternatives_

- [_ContainerSolutions/externalsecret-operator_](https://github.com/ContainerSolutions/externalsecret-operator): aws, gcp, gitlab
- [mumoshu/aws-secret-operator](https://github.com/mumoshu/aws-secret-operator): aws

##### [_mozilla/sops_](https://github.com/mozilla/sops) based solutions

SOPS is the generic secret encryption tooling,
additional tooling is needed to integrate with k8s tooling.
Secrets are encrypted either with GPG key or cloud KMS keys.
Tooling wrappers/plugins are used to decrypt the contents
before being rendered/used as plain vars.

###### _examples_

- [_zendesk/helm-secrets_](https://github.com/zendesk/helm-secrets): for helm, `helm secrets install/upgrade/template/...`
- [_viaduct-ai/kustomize-sops_](https://github.com/viaduct-ai/kustomize-sops#argo-cd-integration): for kustomize, uses the slightly questionable kustomize/go [plugins](https://github.com/kubernetes-sigs/kustomize/blob/master/examples/secretGeneratorPlugin.md)

For helm specifically,
individual helm vars are encrypted with GPG or cloud KMS.
`helm secrets ...` wraps normal helm commands to decode the secrets values
which can be used as normal helm vars in templating.

##### git

use git filters/commands to encrypt/decrypt files before git operations.
