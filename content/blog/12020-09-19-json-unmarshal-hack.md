---
description: shadowing types in unmarshaling
title: json unmarshal hack
---

### _go_ json unmarshal

Sometimes you want to perform some extra validatoon logic in unmarshaling json,
or you want to override some types.
Here's one way to do it

#### _validation_

You want to modify some fields, perform some extra calcs
but the json unmarshaling bits remain the same?

You shadow the type, losing its methods (UnmarshalJSON),
use the standard unmarshaling logic on it,
then copy the data back and perform your modifications.

No need for awkward `MyInt`s in your struct

```go
package main

import (
        "encoding/json"
        "fmt"
)

func main() {
        j := []byte(`{"int": -5, "string": "foobar"}`)

        type PlainS S

        var plainS PlainS
        _ = json.Unmarshal(j, &plainS)
        fmt.Println(plainS)            //{-5 0 foobar}

        var s S
        _ = json.Unmarshal(j, &s)
        fmt.Println(s)                 // {1 -50 foobar}

}

type S struct {
        Int           int
        calculatedInt int
        String        string
}

func (s *S) UnmarshalJSON(b []byte) error {
        type S2 S
        var s2 S2
        err := json.Unmarshal(b, &s2)
        if err != nil {
                return err
        }

        // perform modifications
        s2.calculatedInt = s2.Int * 10
        if s2.Int < 0 {
                s2.Int = 1
        }
        *s = S(s2)
        return nil
}
```

#### _type_ overrides

Sometimes you just need to override some unmarshaling logic
for some fields
but you can't or don't have access to add `UnmarshalJSON` to the types.

You shadow the type, you then create a new type embedding the shadowed type
with the overrides you want.

```go
package main

import (
        "encoding/json"
        "fmt"
)

func main() {
        j := []byte(`{"int": 5, "string": "foobar"}`)

        type PlainS S

        var plainS PlainS
        _ = json.Unmarshal(j, &plainS)
        fmt.Println(plainS)            //{ foobar}

        var s S
        _ = json.Unmarshal(j, &s)
        fmt.Println(s)                 // {five foobar}

}

type S struct {
        Int    string
        String string
}

func (s *S) UnmarshalJSON(b []byte) error {
        type S2 S
        type S3 struct {
                S2
                Int int
        }
        var s3 S3
        err := json.Unmarshal(b, &s3)
        if err != nil {
                return err
        }

        // perform type convs
        if s3.Int == 5 {
                s3.S2.Int = "five"
        }
        *s = S(s3.S2)
        return nil
}
```
