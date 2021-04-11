#### _sources_

- [wiki][exrep]: experience reports
- [reddit][reddit]: first 100 results from google
- [go-nuts][nuts] search "modules" up to 2020
- [issues][issues]: label modules, also [google search][issues-search]
- random blogs i found

#### _reports_

summary / comments

- [abusing go as generic depenency manager][nuts23] 2020, ...
- [checksum mismatch][reddit30] 2018, lots of reasons
- [CI and mod cache][reddit13] 2018-2020, it's complicated
- [Cockroach migration][cockroach] 2021, dep->modules, go.sum lines, tooling versioning, c/proto files
- [codegen][nuts21] 2020, no, pregenerate
- [contribution / fork flow][nuts16] 2021, more examples?
- [Cothority][cothority] 2019, people are no longer forced to place repos in predictable locations
- [dev dependencies][reddit28] 2018, not really possible
- [Dgraph / Badger][badger1] 2019, want v2 but don't want to force users to rewrite import path
- [Dgraph / Badger Round2][badger2] 2019, consider more than API compat in semver
- [docker share module cache][nuts4] 2021, need docker support
- [exclude directories][i37724] 2020, add a go.mod
- [friction log][nuts6] 2021, overloaded go get
- [gitlab subgroups][i34094] 2019, fixed
- [go create][nuts15] 2021, unclear problem statement
- [Go directive][godirect] 2021, better docs needed?
- [godoc with modules][reddit6] 2019, it works now
- [gopath / code layout][nuts12] 2021, `GOPATH/src/*/` interdependencies?
- [GopherJS Migration][gopherjs] 2018, should've done everything in a single commit
- [Hero][hero] 2021, references mostly 2018-2019 issues, issues supporting modules / GOPATH simultaneously
- [How to fork][reddit1] 2020, people seem to get this wrong pretty often, also [1][reddit23]
- [ignore dependency's test dependencies][reddit10] 2018, lazy module loading?
- [internal, relative?][nuts9] 2021, internal is enforced
- [k8s and module cycles][nuts11] 2021, k8s and prometheus, also [1][twitter1]
- [license check][reddit25] 2020, go list -> go license
- [Linux distro Go version][distro] 2021, stable distros are slow
- [local debug][nuts2] 2021, local replace, also [1][nuts13] [2][nuts22] [3][nuts24]
- [Local only code][reddit2] 2019, module names != code lives on internet, also [1][reddit16] [2][reddit20]
- [mod graph confusion][i40513] 2020, replaces source code not version
- [module and ci][dockerci] 2019, download only
- [modules and gopls][nuts8] 2021, need to exmplain module path, gopls woes, also [1][nuts19]
- [Multiple module development][reddit3] 2018-2020, this definitely needs better tooling, also [1][reddit5] [2][reddit8] [3][reddit11] [4][reddit12] [5][reddit22]
- [need example][nuts3] 2021, example workflow, local only
- [new major versions][i40323] proposal
- [no commands to work across module boundaries][nuts7] 2021, unfortunate, also [1][nuts17]
- [non root tags][nuts20] 2020, docs
- [painful][nuts1] 2021, multiple interdependent flows, more docs for independent / small teams
- [patching][reddit27] 2019, complicated by caddy doing nonstandard
- [private bitbucket][reddit19] 2020, GOPRIVATE
- [private gitlab][reddit24] 2018, .netrc
- [private, ssh][nuts14] 2021, more visible docs? also [1][nuts18] [2][nuts25]
- [project structuring][reddit15] 2020, not exactly modules problem
- [protobuf][reddit26] 2018, precompile protos, also [1][reddit29]
- [proxy 410 gone][reddit18] 2019, invalid file name, private code
- [proxy cache issues][i38065] ongoing, set GOPROXY
- [proxy operations problems][reddit14] 2020, proxy dependent?
- [Rally / Docker rename][rally] 2017, docker/docker -> moby/moby, solved by modules requiring canonical import path
- [SamWhited migrating multiple projects][samwhited] 2019, tooling could definitely be improved
- [Show chain of dependencies][reddit4] 2021, `go mod why` is anemic
- [simple test programs][nuts5] 2021, go run \*.go
- [Tooling][both] 2019, tooling support for modules is only halfway there
- [too many go.mod][reddit21] 2021, as title
- [Too optimized for remote repos][possible] 2021, tooling changes should stabilize soon, `@ref`, losing clone into `GOPATH/src`
- [unversioned not latest major][nuts10] 2021, SIV
- [use branched][reddit17] 2018, @ref allows branches
- [vendor and modules][reddit7] 2019, it automatically uses vendor now
- [visualize dependency graph][reddit9] 2019, needs better tooling

[badger1]: https://discuss.dgraph.io/t/go-modules-on-badger-and-dgraph/4662
[badger2]: https://dgraph.io/blog/post/serialization-versioning/
[both]: https://brandon.dimcheff.com/2019/04/go-modules-the-best-of-both-worlds/
[cockroach]: https://www.cockroachlabs.com/blog/dep-go-modules/
[cothority]: https://gist.github.com/ineiti/4a4a1798876225f7a553a13120d705ae
[distro]: https://utcc.utoronto.ca/~cks/space/blog/programming/GoModuleSupportNeed
[dockerci]: https://evilmartians.com/chronicles/speeding-up-go-modules-for-docker-and-ci
[exrep]: https://github.com/golang/go/wiki/ExperienceReports#modules
[godirect]: https://utcc.utoronto.ca/~cks/space/blog/programming/GoModulesGoVersionWhy
[gopherjs]: https://gist.github.com/myitcv/79c3f12372e13b0cbbdf0411c8c46fd5
[hero]: https://github.com/KateGo520/Hero/issues/1
[i34094]: https://github.com/golang/go/issues/34094
[i37724]: https://github.com/golang/go/issues/37724
[i38065]: https://github.com/golang/go/issues/38065
[i40323]: https://github.com/golang/go/issues/40323
[i40513]: https://github.com/golang/go/issues/40513
[issues]: https://github.com/golang/go/issues?page=10&q=label%3Amodules
[issues-search]: https://www.google.com/search?q=site:github.com/golang/go/issues+modules+problems
[nuts10]: https://groups.google.com/g/golang-nuts/c/aOvjBRUWJPA/m/c3CqVI9iFwAJ
[nuts11]: https://groups.google.com/g/golang-nuts/c/FAO6x5AhfPg/m/dPvO_5r1FgAJ
[nuts12]: https://groups.google.com/g/golang-nuts/c/KL0VwEN--k0/m/BYewzynlFQAJ
[nuts13]: https://groups.google.com/g/golang-nuts/c/9MfGXLmRu8w/m/D2gm_viYBAAJ
[nuts14]: https://groups.google.com/g/golang-nuts/c/dp96wZbHtvs/m/qybVYz1WBAAJ
[nuts15]: https://groups.google.com/g/golang-nuts/c/0VgPQbQEKdU/m/qFKLoVQpAgAJ
[nuts16]: https://groups.google.com/g/golang-nuts/c/gjM1zVnd7Ek/m/w8lOn9v-AQAJ
[nuts17]: https://groups.google.com/g/golang-nuts/c/B-gvL92b2Vo/m/EqzEzewoDgAJ
[nuts18]: https://groups.google.com/g/golang-nuts/c/0uMSmt_TnKM/m/oB_TqrFcAgAJ
[nuts19]: https://groups.google.com/g/golang-nuts/c/KXmc4v2ay4k/m/jjiWOV20BgAJ
[nuts1]: https://groups.google.com/g/golang-nuts/c/_BqV6Rk15UA/m/ns4y8jbxBgAJ
[nuts20]: https://groups.google.com/g/golang-nuts/c/7Z6U5aKxaJI/m/-Trvp6sxBgAJ
[nuts21]: https://groups.google.com/g/golang-nuts/c/PPmCyg4T1hY/m/uqt1f9sCBgAJ
[nuts22]: https://groups.google.com/g/golang-nuts/c/9-5aDopSGvo/m/dLGsOtnQBQAJ
[nuts23]: https://groups.google.com/g/golang-nuts/c/21xRZmknkQQ/m/kX_JlOSHCwAJ
[nuts24]: https://groups.google.com/g/golang-nuts/c/ga1XPbquXL4/m/DBhDCNG3AAAJ
[nuts25]: https://groups.google.com/g/golang-nuts/c/lIIxpRmAuYY/m/pF402MkmCAAJ
[nuts2]: https://groups.google.com/g/golang-nuts/c/WbbVeO321ak/m/3QmWi5vdBgAJ
[nuts3]: https://groups.google.com/g/golang-nuts/c/HOfo5INo3nM/m/C0e9fPduAQAJ
[nuts4]: https://groups.google.com/g/golang-nuts/c/l7oXXpfmqUo/m/BGrzFqpWBgAJ
[nuts5]: https://groups.google.com/g/golang-nuts/c/bxbe9vI6Duc/m/LNtAMC3EBQAJ
[nuts6]: https://groups.google.com/g/golang-nuts/c/RZ1qGp8REYg/m/QKe8QMofCwAJ
[nuts7]: https://groups.google.com/g/golang-nuts/c/JAmfHLMN2XE/m/EK5lIRoICgAJ
[nuts8]: https://groups.google.com/g/golang-nuts/c/2Xcfb4f7ans/m/M6Eg50DhBwAJ
[nuts9]: https://groups.google.com/g/golang-nuts/c/d3ZMjah6VGE/m/HULNtwu9BgAJ
[nuts]: https://groups.google.com/g/golang-nuts/search?q=modules
[possible]: https://utcc.utoronto.ca/~cks/space/blog/programming/GoModuleBuildsWhatPossible
[rally]: https://www.rallyhealth.com/coding/docker-moby-go-dependencies
[reddit10]: https://www.reddit.com/r/golang/comments/98hbvk/go_modules_how_to_deal_with_test_dependencies/
[reddit11]: https://www.reddit.com/r/golang/comments/kku3ec/local_development_between_2_go_modules/
[reddit12]: https://www.reddit.com/r/golang/comments/ejsgl0/using_local_development_modules_without_pushing/
[reddit13]: https://www.reddit.com/r/golang/comments/9p2xti/go_modules_cache_location/
[reddit14]: https://www.reddit.com/r/golang/comments/f243lc/go_modules_proxy/
[reddit15]: https://www.reddit.com/r/golang/comments/igyry9/go_modules_and_project_structuring/
[reddit16]: https://www.reddit.com/r/golang/comments/fix89j/using_modules_locally_without_publishing_to_vcs/
[reddit17]: https://www.reddit.com/r/golang/comments/9ahlvo/use_branches_with_go_modules/
[reddit18]: https://www.reddit.com/r/golang/comments/d2n5s0/error_410_gone_when_switching_to_modules_in_go_113/
[reddit19]: https://www.reddit.com/r/golang/comments/fmrmbq/help_unable_to_setup_go_modules_with_private/
[reddit1]: https://www.reddit.com/r/golang/comments/j8pqms/go_modules_making_me_rage_how_do_i_fork_a_module/
[reddit20]: https://www.reddit.com/r/golang/comments/kcc2td/importing_a_module_from_a_local_path/
[reddit21]: https://www.reddit.com/r/golang/comments/kqf3ui/modules_hellp_private_modules_monorepo_module/
[reddit22]: https://www.reddit.com/r/golang/comments/lvqgln/how_do_you_do_rapid_development_with_modules/
[reddit23]: https://www.reddit.com/r/golang/comments/guf22q/golang_modules_and_developing_in_a_fork_to/
[reddit24]: https://www.reddit.com/r/golang/comments/9els7j/go_module_with_private_gitlab_repos/
[reddit25]: https://www.reddit.com/r/golang/comments/g9mrtn/managing_licenses_with_go_modules/
[reddit26]: https://www.reddit.com/r/golang/comments/9s3512/go_modules_and_importing_protobuf_files/
[reddit27]: https://www.reddit.com/r/golang/comments/bxnqcb/applying_patches_to_go_module_dependencies/
[reddit28]: https://www.reddit.com/r/golang/comments/9weqd8/does_go_modules_support_something_like_dev_only/
[reddit29]: https://www.reddit.com/r/golang/comments/gdv1ah/go_modules_and_proto/
[reddit2]: https://www.reddit.com/r/golang/comments/ah0w1q/modules_and_local_imports/
[reddit30]: https://www.reddit.com/r/golang/comments/9u4zsj/go_modules_checksum_mismatch/
[reddit3]: https://www.reddit.com/r/golang/comments/9gwqg4/how_to_handle_working_on_multiple_modules_at_once/
[reddit4]: https://www.reddit.com/r/golang/comments/lkenmn/modules_how_do_i_show_the_chain_of_dependencies/
[reddit5]: https://www.reddit.com/r/golang/comments/jdwuyy/how_to_work_with_multiple_modules/
[reddit6]: https://www.reddit.com/r/golang/comments/aspm32/godoc_with_modules_outside_gopath/
[reddit7]: https://www.reddit.com/r/golang/comments/b9osrj/modules_and_vendoring/
[reddit8]: https://www.reddit.com/r/golang/comments/juizr6/what_is_the_current_115_best_practice_to_work_on/
[reddit9]: https://www.reddit.com/r/golang/comments/bdtrti/best_way_to_visualize_library_dependencies_with/
[reddit]: https://www.google.com/search?q=site%3Areddit.com%2Fr%2Fgolang+modules
[samwhited]: https://blog.samwhited.com/2019/01/supporting-go-modules/
[twitter1]: https://twitter.com/rakyll/status/1348723364894961666
