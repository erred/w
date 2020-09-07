---
description: first impressions on terraform
title: terraform first impressions
---

### _first_ impressions on terraform

[terraform][terraform] is probably the industry standard
in cross platform/vendor infrastructure as code (IaC).
On its [terraform vs other][vs] pages,
it claims not to be a configuration management tool,
instead focusing on provisioning,
though this is only partly true.

#### _what_ is terraform

terraform is a way for you to write config files in hcl,
and have providers compare the config with current live state,
and apply any changes.
In this way, terraform can control anything with an API
to view the current state and update the state.
Like all the configs of your SaaS...
Sounds like configuration management, right?

So why does it say it's not for config management?
because you can't really manage arbitrary file-based configs
through terraform, other than an initial run when it is first provisioned.

#### _day1_ pitfalls

Don't try to rename stuff unless you don't mind things getting deleted and recreated.
Your config is compared with saved and live state to generate the delta.
If you made a mistake, it's safer to clear out the state and start over (resync/reimport).

The providers that terraform uses are of varying quality.
Sometimes they don't offer all the capabilities you think they should have,
and sometimes they hardcode weird assumptions
( _cough_ github _cough_ ).

You also don't appear to be able to create multiple instances of a provider,
ex, I see no way of managing multiple github organizations together
as the organization name is part of the provider config.

You probably don't want everything in a single workspace.
You might think it's a good idea,
you can link things against each other,
but now, every time you want to plan or apply,
you have to wait for it to make a bazillion requests
to check if something changed and another bazillion requests
to actually update the config.
Terraform is not that fast.

[terraform]: https://www.terraform.io/
[vs]: https://www.terraform.io/intro/vs/index.html
