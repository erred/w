---
description: protobuf encoding x3
title: proto encoding
---

### _protobuf_ encoding

So many docs, but they never show us what it actually looks like

#### _proto_

```proto
syntax="proto3";

message Msg {
  string one = 1;
  int64 two = 2;
  bool three = 3;
  float four = 4;
  repeated string five = 5;
  map<string, int64> six = 6;
  E seven = 7;
}

enum E {
  UNKNOWN = 0;
  KNOWN = 1;
}
```

#### _go_

```go
msg := &Msg{
        One:   "een",
        Two:   2,
        Three: true,
        Four:  4.4,
        Five:  []string{"vijf", "fünf"},
        Six:   map[string]int64{"zes": 6},
        Seven: E_KNOWN,
}
```

#### _prototext_

```txt
PROTO TEXT MARSHAL
one:"een"  two:2  three:true  four:4.4  five:"vijf"  five:"fünf"  six:{key:"zes"  value:6}  seven:KNOWN

PROTO TEXT FORMAT
one:  "een"
two:  2
three:  true
four:  4.4
five:  "vijf"
five:  "fünf"
six:  {
  key:  "zes"
  value:  6
}
seven:  KNOWN
```

#### _protojson_

```txt
PROTO JSON MARSHAL
{"one":"een", "two":"2", "three":true, "four":4.4, "five":["vijf", "fünf"], "six":{"zes":"6"}, "seven":"KNOWN"}

PROTO JSON FORMAT
{
  "one":  "een",
  "two":  "2",
  "three":  true,
  "four":  4.4,
  "five":  [
    "vijf",
    "fünf"
  ],
  "six":  {
    "zes":  "6"
  },
  "seven":  "KNOWN"
}
```

#### _protowire_

```txt
PROTO MARSHAL

een%�@*vijf*fünf2
zes8

PROTO MARSHAL HEX
0a 03 65 65 6e 10 02 18 01 25 cd cc 8c 40 2a 04 76 69 6a 66 2a 05 66 c3 bc 6e 66 32 07 0a 03 7a 65 73 10 06 38 01
```
