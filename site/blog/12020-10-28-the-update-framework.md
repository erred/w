---
description: "framework" is so arbitrary
title: the update framework
---

### _the_ update framework

The Update Framework [TUF][tuf] is
a specification (and implementation)
of a process for generating and validating metadata
to ensure content is fresh and available.

#### _longer_

TUF aims to be flexible and support old stuff
so it is as underspecified as it is specified,
preferring to leave a lot of details to implementors.

From the client perspective (thing asking for updates),
TUF provides 2 main functions:
_polling_ for new updates
and _fetching_ those updates.
The [spec][spec] describes a set of roles (with matching keys),
metadata files for the roles,
and a process to be followed to either get updates
or detect a denial of service or other attack,
starting with just a (possibly outdated) copy of the root metadata.
This is implemented as a python library [tuf][python]
and a go library [tuf-go][go],
as well as a client-server [notary][notary]
for storing and serving the metadata.

[tuf]: https://theupdateframework.io/
[spec]: https://github.com/theupdateframework/specification/blob/master/tuf-spec.md
[python]: https://github.com/theupdateframework/tuf
[go]: https://github.com/theupdateframework/go-tuf
[notary]: https://github.com/theupdateframework/notary
