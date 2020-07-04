---
description: but you can't just dump text over http and expect people to read it
title: hyped blogging
---

### _hyped_ blogging

Apparently it's not enough to just convert
thoughts -> text/markdown -> html and dump it over https to people's browsers.
Nooo, the interwebs bemoan _centralization_ and the lack of _discoverability_,
and in true internet form,
dream up a thousand new protocols to keep blogging on life support.

[xkcd: standards](https://xkcd.com/927/)

Note: not using any of the below

#### _old tech_ - not very connected

It works, why change it? they say

- _rss_: designed when people thought XML was a good idea,
  people have **strong** opinions on what you SHOULD put in here
- _atom_: a protocol upgrade to _rss_
- _webring_: sites with links pointing to each other,
  traditionally with forward/back in a ring,
  occasionally with a central directory

#### _IndieWeb_ - loosely connected

Self host, but still sort of connect to each other?
Abstract ideas of being "people focused",
no real tech standards.

- _webmention_: centralized server watches sites,
  if watched sites link to you, you get notification based on link hidden in your site
- _microformat_: overload your site with html classes and hrefs
  so it can be turned into super verbose json
- _micropub_: http/microformat based protocol for content management systens
- _syndication_: sites directly publish content from other sites, not just a link

#### _Fediverse_ - connected clusters

Expand beyond blogging!
Social Media!
Does this really make sense for decentralization?
Mostly clones of popular services

- _activitypub_: the OO people got their hands on HTTP/JSON,
  publish / subscribe in a decentralized way(?),
  main protocol for fediverse
- _XMPP_: chat protocol, in XML!
- _mastodon_: twitter clone
- _pixelfed_: instagram clone
- _peertube_: youtube clone

#### _other_ tech

Sometime derided by the same people as above,
but I think these are truly decentralized,
decoupling content from the serving protocol.

- _WebSub_: PubSubHubbub, extended rss/atom, push content with subscriber webhooks
- _AMP_: first step in content first, even if people hate it
- _WebPackage_: trusted web content bundled together, no longer tied to http!
  uses SXG / WebBundle
- _Signed HTTP Exchange_: SXG, resource with a signature to trust origin,
  you know the content was from the domain at some point, even if you didn't retrieve it directly
- _WebBundle_: resources bundled together
