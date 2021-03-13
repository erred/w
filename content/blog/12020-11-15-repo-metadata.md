---
description: all the things that aren't code in a repo
title: repo metadata
---

### _repo_ metadata

(Long lived) code is usually hosted in a repo somewhere.
Sometimes that can be public,
or maybe just shared with a few people
(include your future self as one of those people).
Now you need to add some metadata to your repo
so people can understand / work better together
(ie share common standards),
and litter the root of your repo with a bunch of files.

#### _README_

Probably the first thing people see,
note this document may also be copied into other places,
such as documentation sites pkg.go.dev or container registries docker.io .

There are lots of ideas as to what should go in here,
I think a _short intro_ as to what it does,
an _example_ of what it's like to use the thing,
and how to _get started_ are the things I look for.
Other sections such as contributing/license/architecture/...
usually already have their own files.

##### _badges_

Because we need little blocks of color
making a dozen requests every time you view the page
and they have to fight with caching.
At least they're good for _at a glance_ and a consistent place to put third party links.

There are a lot of badges from a lot of places, ex [shields.io](https://shields.io/)

The common ones I see are:

- License
- Version
- Documentation link
- Community link

Other ones include

- Code Quality
- CI status
- Package manager / platforms
- Downloads / counters
- Publications
- Support / donations link

#### _LICENSE_

The only file mandatory for making things public.

Should be something that tools such as SPDX can recognize.

##### _CONTRIBUTING_

Related to license, how to get setup to build, patch, and contribute upstream.

##### _CODE_OF_CONDUCT_

Related to contributing, how to interact with the "community"

##### _templates_ for issues and prs

Related to contributing,
Github/other code host specific,
templates to follow so maintainers don't have to chase you to get more info.

#### _CHANGELOG_

No standard format exists.

Single file or a directory of files.
Handwritten ones are good especically if they highlight the important bits,
but are tiresome to write.
Worth it for big projects with infrequent releases.
Generated ones are noisy but maybe good enough for continuous releases
and if you enforce good commit messages.

#### _SECURITY_

Who do you call if you found a critical bug?

Need somewhere to put retractions/CVEs...

##### _semver_

[semver](https://semver.org/)

Related to releasing, use semver to communicate stability
of public interfaces / project maturity.

#### _tooling_

Because each repo is a special snowflake,
we need more files to tell our tools how it should work in this repo.

##### _.editorconfig_

Some ini-like file that reduces the editor specific files a little bit.
Right now supports spacing/tabs.
Probably needs to support a lot more for it not to be yet another file.

##### _.gitattributes_

Tell git how to treat different files, add optional filters.

##### _.\*ignore_

Tell our tools to ignore other files that just have to be there for other purposes.

##### _build_ configuration

This is usually a good thing, having a documented, repeatable way of building code.
Also includes the CI/CD configuration
that is unfortunately not really portable across platforms.
