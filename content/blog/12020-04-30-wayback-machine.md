---
description: archiving with the internet archive
title: wayback machine
---

### _Wayback_ Machine

we go a long way back

#### _TLS_ and wayback

So if you're a security conscious,
you've configured your website to use the mozilla's
_Modern_ Configuarion through
[SSL Configuration Generator](https://ssl-config.mozilla.org/),
which means minimum TLS 1.3

Unfortunately, the
[Internet Archive](https://archive.org/)
[Wayback Machine](https://web.archive.org/),
only looks back at history, not forward into modernity
and doesn't support such newfangled security mechanisms.

_Solution?_ temporarily degrade your security, TLS 1.2 should work
**as of 2020-04-30**
while you make the wayback machine archive your site

#### _Submitting_ pages

if you have a sitemap,
this oneliner should submit everything in it,
84 pages took ~230 secs.

```sh
xargs -I '{}' -n 1 curl -Lv https://web.archive.org/save/'{}' < sitemap.txt
```
