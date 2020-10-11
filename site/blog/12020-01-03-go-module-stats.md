---
description: stats from scraping proxy.golang.org
title: go module stats
---
### _stats_ on go modules

by default, go now uses a [proxy](https://proxy.golang.org)
which also has an [index](https://index.golang.org),
which is on by default.

So now we have a nice way
of knowing what modules people are using.

#### so*...*

as of _2020-01-03_ there are:

- 99189 unique modules
- 687248 unique module versions
- of which 606951 are proxied
- avg: 6.93, max: 5113

#### host names

we can look at which host name segments are popular (top 20):

```
          com:  94288 95.06%
       github:  91885 92.64%
          org:   1478  1.49%
       gitlab:   1477  1.49%
        gopkg:   1250  1.26%
           in:   1250  1.26%
    bitbucket:    900  0.91%
           io:    700  0.71%
          git:    606  0.61%
        gitee:    452  0.46%
           go:    420  0.42%
          k8s:    203  0.20%
          dev:    200  0.20%
         code:    183  0.18%
          net:    165  0.17%
        gitea:    112  0.11%
 cloudfoundry:    109  0.11%
          xyz:    101  0.10%
           co:    100  0.10%
           ht:     95  0.10%
```

#### _go_ in module names

people like to put [-]go[-] somewhere in the names:

```
prefix: go-:  11849 11.95%
 prefix: go:   8993  9.07%

suffix: -go:   3366  3.39%
 suffix: go:   2054  2.07%
```

#### suffixed _vcs_

so we know everyone is on github,
but what about those that use suffixes?

_TODO_: implement the `go-get=1` query

```
 suffix: .svn:      0  0.00%
 suffix: .bzr:      0  0.00%
 suffix: .git:    142  0.14%
 suffix:  .hg:      2  0.00%
```

#### pitfalls

this is only public code

_410 Gone_ is kinda annoying,
people delete modules (understandable),
and people leak private module names
(which choke up badly written code).