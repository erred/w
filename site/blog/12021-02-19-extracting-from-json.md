---
description: extracting partial elements from json
title: extracting from json
---

### _json_

You have a JSON shaped problem,
and now you want to write a short query for getting something out of it.

```json
{
  "a": [
    {
      "x": 100
    },
    {
      "x": 200
    }
  ]
}
```

#### json _pointer_

[RFC 6901](https://tools.ietf.org/html/rfc6901) introduces JSON Pointer:
_tldr:_ `/` separated key names, arrays are 0 indexed.
No special syntax for all array elements etc...

```jsonpointer
# 200
/a/1/x
```

Not very flexible, but it is used in JSON Patch: [RFC 6902](https://tools.ietf.org/html/rfc6902)

#### json _path_

[JSON Path](https://goessner.net/articles/JsonPath/) is basically XPath for JSON.

_tldr:_ `.` separated, `..` recursive descent, `[::]` array access, `?()` script filter (eg on child attributes)

```jsonpath
# 100, 200
$..x
$.a[:].x

# [{ "x": 100 }]
$.a[?(@.x == 100)]
```

#### _jq_

[jq](https://stedolan.github.io/jq/)
is a more general purpose utility that can do advanced queries
and reshape data as necessary.
_tldr_: `.` separated, `..` recursive descent, `[]` array access, `|` pipe and use functions for more.

```jq
# 100, 200
..x
.a[].x

# { "x": 100 }
.a[] | select(.x == 100)
```

#### _cel_

[cel](https://github.com/google/cel-go)
([spec](https://github.com/google/cel-spec))
is a minimal language for evaluating single expressions.
With [macros](https://github.com/google/cel-spec/blob/master/doc/langdef.md#macros)
it becomes viable to query nested data structures.

```cel
# 100, 200
o.a.map(n, n.x)

# { "x": 100 }
o.a.filter(n, n.x == 100)
```
