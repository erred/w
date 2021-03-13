---
description: what am i doing?
title: junior sre at non tech
---

### _junior_

So me:
fresh out of graduate school in Amsterdam
in the middle of a global pandemic.
What am I going to do?
After 78+ applications
and interviews at a dozen companies,
I got exactly 1 offer.
SRE, but not at a tech company,
wonder how it will go.

#### _getting_ started

2 months in, what am I actually doing?
Including me, it's a team of 5,
supporting 80(?) devs and 2 kubernetes clusters.

Onboarding was a bit chaotic,
day 1 got me full access to our production cluster,
but not cloudflare or even our full codebase
(I still can't see everything).
I leaned heavily on [`ripgrep`](https://github.com/BurntSushi/ripgrep)
and [`t`](https://github.com/seankhliao/t)
to find about things.

#### _doing_ what

So what do I do day to day?
Look at Jira, pick out a task,
get [annoyed at MS Teams](/blog/12020-12-01-ms-teams/),
and wonder why I have so many meetings.

Fiddle with DNS and Page Rules in [Cloudflare](https://www.cloudflare.com/)
through [Terraform](https://www.terraform.io/),
migrate stuff from [Ansible](https://www.ansible.com/)
to [Helm](https://helm.sh/) (if you have a choice [avoid Helm](/blog/12020-12-02-avoid-helm/)),
play around with [Grafana](https://grafana.com/) dashboards,
and if I'm lucky, write some minor service in [Go](https://golang.org/)
to glue together other services (the rest of our devs use Java),
Oh, and investigate more processes and additions to tooling,
and consider migrating from [Jenkins](https://www.jenkins.io/) to something else.
