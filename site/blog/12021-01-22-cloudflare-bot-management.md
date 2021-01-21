---
description: managing bots with cloudflare
title: cloudflare bot management
---

### _cloudflare_

$job has a bot problem,
(or maybe a shitty backend problem),
[cloudflare bot management](https://www.cloudflare.com/products/bot-management/)
to the rescue?

#### _background_

There are typically 3 types of resources we serve:
page, static resource, api endpoint.
A page is the initial request for a user,
this could be a static html page or some dynamically generated one.
Static resources are usually served from a CDN
and api endpoints serve both the web frontend and mobile apps.
What's important is that for pages, we can serve challenges
(cloudflare making you wait a bit or making you solve a captcha).
For static resources, we probably don't care that much,
and for api endpoints, because the requests are made programmatically,
there's not much we can do besides allow/block.

##### _products_

Cloudflare has quite a few blocking products that sit in front of their CDN,
workers, and your servers. The order they happen in can be seen:

![cloudflare rule priority](https://developers.cloudflare.com/firewall/static/5bbd49ff428e4d4f2f3d3b8688a8a8f3/29007/firewall-rules-order-and-priority-1.png)

##### _not_ enough

So why are the existing products not enough?
The existing tools (IP, user agent, WAF, rate limiters...) are like scalpels,
precise for their functions, but not very flexible to changes.
Ex: IP and rate limit blocks are easy to get around by just spinning up new VMs in clouds,
or if you really want (and don't mind IPs from [dubious sources](https://luminati.io/proxy-types/rotating-residential-ips))
with a large pool of proxies.

Sure, with enough investigation, you can tailor blocks using the existing tools
to more or less only block specific attacks, but that's a lot of work for a fragile solution,
and you inevitably end up with a bunch of rules that are no longer relevent as script kiddies move on.

#### _bot_ management

If the existing tools were fishing rods and bait,
bot management is like a dragnet.
Requests are classed by heuristics, machine learning
(they refuse to say what sort of machine learning is fast enough to run on every request),
and behavioural analysis,
and you get a 1-100 score you can use as part of your firewall rules.

- Advantages: it catches a lot of things.
  The landscape shifts pretty often and due to the whitelisting style of rules you write,
  it's slightly easier to stay ahead.
- Disadvantages: it catches too much stuff.
  It's a really broad brush and if you're not careful you can easily break something in prod,
  as you accidentally block something another team uses
  (and good luck getting marketing to list out all the tools they use).
  Also, hope your internal services and mobile apps set proper user agents for you to filter...

Where previously your rules are more building up ORs of potentially bad sources,
ex (`cf.client.bot` is available to everyone as verified good bots):

```txt
(
  (
    ip.geoip.asn in { xxxx yyyy zzzz }
    or http.user_agent contains "fuzz"
    or ip.src in { "a.b.c.d/32" "q.r.s.t/31" }
  )
  and not cf.client.bot
  and http.host eq "example.com"
)
```

With bot management it is more digging out exceptions for known good sources:

```txt
(
  cf.bot_management.score eq 1
  and not cf.bot_management.static_resource
  and not cf.bot_management.verified_bot
  and not ip.src in $office
  and not ip.src in $cloud
  and not ip.src in $trusted_third_party
  and not http.user_agent contains "company"
  and not http.user_agent contains "seo.tool.we.use"
  and http.host eq "example.com"
)
```

With terraform it's slightly better since we can use
variables and string interpolation to manage a bunch of rules,
but it doesn't help with investigating all those potentially legitimate requests you might block

```terraform
resource "cloudflare_filter" "ex" {
  # ...
  expression = <<-EOF
  (
    cf.bot_management.score eq 1
    and not cf.bot_management.static_resource
    and not cf.bot_management.verified_bot
    ${join("\n  ", formatlist("and not ip.src in %s", local.known_ip_lists))}
    ${join("\n  ", formatlist("and not http.user_agent contains %s", local.known_user_agents))}
    and http.host eq "example.com"
  )
  EOF
}
```
