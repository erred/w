---
title: nomad questions
descriptions: open questions about using nomad
---

### _nomad_

[nomad](https://www.nomadproject.io/) is a scheduler by
[hashicorp](https://www.hashicorp.com/):
you give it [job](https://www.nomadproject.io/docs/job-specification/job) definitions,
it schedules them onto your nodes.

The ux for managing jobs feels fairly barebones,
you submit it once and hope for the best?
The other thought I had was managing it through terraform,
but even though it's both hcl, it's just a giant string
in their [provider...](https://registry.terraform.io/providers/hashicorp/nomad/latest/docs/resources/job)

About the job config,
I don't think they do a good enough job of documenting
which features are available under which conditions.
The config options are hierarchical but it doesn't seem enough?

About workflow,
nomad and their new [waypoint](https://www.waypointproject.io/)
both seem to think your application will only ever need a few env vars as config.
Maybe in some idea world that would be true,
but in this reality I need to pass config files somehow
and nomad doesn't have a good story for that.
