---
description: HTTPS troubles with GCP and Cloudflare
title: GCP HTTPS Cloudflare
---

### _https_

Secure all the things!
HTTPS everywhere.
But how do you chain 2 not very compatible things together
in an unbroken chain of HTTPS?

#### _Google_ Cloud Storage

There are 3 ways to access things on GCS:

- HTTPS: `storage.googleapis.com/b-example-com/object` with a bucket named `b-example-com`
- HTTPS: `b-example-com.storage.googleapis.com/object` with a bucket named `b-example-com`
- HTTP: `b.example.com/object` with a bucket named `b.example.com` and a CNAME to `c.storage.googleapis.com`

You'll notice there's no support for HTTPS with your custom domain,
the supported way is with a load balancer.

With Cloudflare and page rules,
you can have HTTPS on a custom domain without paying for a load balancer.

- `CNAME b.example.com b-example-com.storage.googleapis.com`
- Page rule `Host Override` to rewrite the host to `b-example-com.storage.googleapis.com`

Note: you can't use a Resolve Override because it won't pass cloudflare validation

#### _Google_ Cloud Run

Cloud run gives you a stable domain per app at `your-app.a.run.app`.
You can map custom domains to it,
and Google will tell you to map a CNAME to ghs.googlehosted.com.
But if you have Cloudflare proxying turned on,
it will forever be stuck provisioning the certificate,
likely because Cloudflare [hijacks `/.well-known/`](https://stackoverflow.com/a/62047503).
If you turn off proxying temporarily,
it will succeed now, but fail silently in 90 days when the letsencrypt certificate expire,
and if you have Cloudflare strict origin SSL on, you'll get HTTP 525.

Your best bet is to do the same (also recommended by a Google engineer for a [no downtime switch](https://stackoverflow.com/a/59260376)):

- `CNAME your-app.example.com your-app.a.run.app`
- Page rule `Host Override` to rewrite the host to `your-app.a.run.app`

Note: you can't use a Resolve Override because it won't pass cloudflare validation
