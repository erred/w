---
description: who wants to write all those terraform files by hand
title: importing terraform
---

### _terraform_

It's okay when you start with terraform from a clean slate,
but what if you already have things that you now want to manage
with [terraform](https://www.terraform.io/)?

Tried importing my GCP project

#### _terraformer_

[GoogleCloudPlatform/terraformer](https://github.com/GoogleCloudPlatform/terraformer)
made by google people works with some other providers,
but not an official product.
At the time of writing 0.8.10, (only) supports terraform 0.12,
even though 0.13 has been out for a long time and 0.14 has just been released.

Outputs hcl that isn't proper syntax: `google-provider = = {`.
Has a bug where it endlessly loops on cloud monitoring api.
Doesn't actually list all your resources.

But at least it works, good for auditing what you have.

#### _terracognita_

[cycloidio/terracognita](https://github.com/cycloidio/terracognita)
Chokes on resources with apis not enabled,
cba to go in and exclude every resource.
Also needs a region for some reason
