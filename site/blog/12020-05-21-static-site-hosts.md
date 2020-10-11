---
description: quick comparison of static site hosts
title: static site hosts
---

### _static_

You have _html_, _css_, _js_, and some images you want the world to see.
You want to serve it under your own domain and https is a must.

What are your choices?

#### _platforms_

The first 6 are all acceptable choices.

All can be paired with ex cloudflare for ipv6 + custom domain + https (CDN-client)

<!-- prettier-ignore -->
| platform      | upload          | https         | custom domain | header | free limits                  | pricing / month               |
| ------------- | --------------- | ------------- | ------------- | ------ | ---------------------------- | ----------------------------- |
| _vercel_      | js cli, gh app  | yes           | yes           | yes    | unlimited                    | opt. $20: more builds         |
| _netlify_     | js cli          | yes           | yes           | yes    | unlimited store, 100GB net   | $0.20/GB net                  |
| _firebase_    | js cli          | yes           | yes           | yes    | 1GB store, 10GB net          | $0.25/GB store, $0.15/GB net  |
| surge.sh      | js cli          | yes           | yes           | paid   | unlimited?                   | opt. $30                      |
| _github_      | git push        | yes           | yes           | no     | unlimited?                   | opt. github pro               |
| gitlab        | git push        | yes           | yes           | no     | unlimited?                   | opt. gitlab                   |
| gcs           | gcloud cli, web | load balancer | yes           | no     | 5GB store, 1GB net NA-others | $0.026/GB store, $0.12/GB net |
| s3            | aws cli, web    | cloudfront    | yes           | no     | none                         | $0.023/GB store, $0.09/GB net |
| neocities.org | ruby cli, web   | paid          | paid          | no     | 1GB store, 200GB net         | opt. $5: 50GB store, 3TB net  |

##### _notes_

- vercel: also includes app hosting
- firebase: supports functions
- s3: buckets and pricing is regional
- gcs: buckets and pricing is multi-regional, regional storage pricing \$0.02/GB

#### _others_

- glitch.com: primarily app hosting, poor documentation, feature incomplete (ex contact support to ...)
- heroku: primarily app hosting, meh documentation, stagnated?
