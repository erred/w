---
description: private go module proxies
title: go module proxies
---

### _go_ module proxies

what are the options besides the default public [proxy.golang.org](https://proxy.golang.org/) if you want to run something private?

#### _self_ hostable

Athens is the most "feature complete" but goproxy/goproxy is much simpler with almost the same functionality,
and it's easier to embed/customize.

| project                        | hosted                     | direct | upstream proxy | exclude | sumdb | cache control | access control |
| ------------------------------ | -------------------------- | ------ | -------------- | ------- | ----- | ------------- | -------------- |
| _[gomods/athens][athens]_      | [azure][athensazure]       | yes    | yes            | yes     | proxy | no            | no             |
| _[goproxy/goproxy][goproxycn]_ | [goproxy.cn][goproxycnweb] | yes    | yes            | yes     | proxy | no            | no             |
| [goproxyio/goproxy][goproxyio] | [goproxy.io][goproxyioweb] | yes    | yes            | ?       | proxy | no            | no             |
| [thumbai/thumbai][thumbai]     | [thumbai.app][thumbaiweb]  | yes    | yes            | ?       | ?     | no            | no             |

[athens]: https://github.com/gomods/athens
[athensazure]: https://athens.azurefd.net/
[goproxycn]: https://github.com/goproxy/goproxy
[goproxycnweb]: https://goproxy.cn/
[goproxyio]: https://github.com/goproxyio/goproxy
[goproxyioweb]: https://goproxy.io/
[thumbai]: https://github.com/thumbai/thumbai
[thumbaiweb]: https://thumbai.app/

##### _notes_

- athens is probably the most "feature complete"
- goproxy/goproxy is the simplest and can be customized / embedded as a handler the easiest
- athens and goproxy/goproxy both have support for a wide range of storage backends
- goproxyio looks overly complex yet not configurable
- thumbai does not look maintained

#### _others_

not ready / commercial / part of another project

- [Gitlab](https://docs.gitlab.com/ee/user/packages/go_proxy/) more features, beta
- [JFrog Artifactory](https://www.jfrog.com/confluence/display/JFROG/Go+Registry) and the hosted public version [gocenter.io](https://search.gocenter.io/)
