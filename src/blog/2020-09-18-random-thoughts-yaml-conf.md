---
description: musing on yaml as a configuration file format
title: random thoughts yaml conf
---

### _yaml_

The internet loves to hate it.
Untyped, whitespace sensitive, need quotes for random strings, ...

I think it's passable, just wish it was a bit less "powerful" in special strings.
My editor yells at me when it autocompletes and inserts tabs :facepalm:
(`&anchor node` and `*anchor` to reuse things may be powerful but sure looks confusing)
(type tags look horrible)

#### _alternatives_

##### _ini_

Too unspecified / no standard

##### _json_

Ugly keys, no comments

##### _toml_

Okay for key value pairs (ini like)
but has weird array of objects syntax, needs repeating the array name,
and doesn't have a good looking way of indenting nested heirarchies.

##### _xml_

You hate humanity.

##### cue, nix, jsonnet, hcl, dhall

These are more like programming languages,
powerful, but harder to see the result without evaluating it.
For some the output might change depending on the environment.
