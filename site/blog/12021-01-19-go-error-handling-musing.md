---
description: random thoughts on go error handling
title: go error handling musings
---

### _errors_

_premise_:
Errors are _values_.
Unless some fundamental constraint has been violated,
in which case feel free to _panic_,
errors are just some other state to be handled.

_pain points_:
as far as I can tell, there are 2 main complaints:

- error handling is verbose and repetitive
- errors can be ignored

#### _verbose_

For the first issue, people usually propose some syntax sugar
to shorten `if err != nil { ... }`

How about a slightly more generic solution that checks against the zero value,
then we can coopt the short circuiting behaviour of `&&` to make the common case one line.
A new unary operator `!!` evaluates to true if the operand is the zero value of its type, else false.

```go
_, err = f() // currently valid
err != nil && return _, _, fmt.Errorf("foo: %w", err)

_, err = g() // proposed, pre
!!err && return _, _, fmt.Errorf("bar: %w", err)

_, err = h() // proposed, post
err!! && return _, _, fmt.Errorf("fizz: %w", err)

// equivalent function
func z(v interface{}) bool {
       return reflect.ValueOf(v).IsZero() // but that reflect penalty...
}

_, err = i()
z(err) && return _, _, fmt.Errorf("buzz: %w", err)
```

Also really want `_` to mean the zero value for any type...

I also thought the idea of assigning to an error handler was a good start,
though it does have some issues, like implicit passing of arguments,
handler having to be defined in function scope, and non local returns,
I guess with some extra thought this is where the check/handle came from.

```go

func f() (err error) {
        handler := func(other arg, err error) {
                if err != nil {
                        err = fmt.Errorf("some extra context: %w", err)
                        return // non local return???
                }
        }

        var o arg

        _, handler(o) = g()
}

func handler(err error, arg string) {
        if err != nil {
                return
        }
}

```

### _ignorable_

Usually the people focusing in the second camp want sum types,
forcing you to handle errors.
People usually deal with this now with linters,
and to be honest, I don't really think this is as much of a concern
