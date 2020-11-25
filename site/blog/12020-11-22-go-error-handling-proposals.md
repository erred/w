---
description: overview of error handling proposals in go
title: go error handling proposals
---

### _error_ handling

```go
if err != nil {
        // handle error
}
```

Error handling seems to be a recurring theme in go,
but most proposals get nowhere.

[Go2ErrorHandlingFeedback](https://github.com/golang/go/wiki/Go2ErrorHandlingFeedback)

#### _proposals_

baseline code

```go
x, err := foo()
if err != nil {
        return nil, err
}
```

- [result type](https://github.com/golang/go/issues/19991):
  box of `value|err`

##### with error handling

note all the ones that claim to use "plain functions" as error handlers have an implicit double return

###### _predeclared_ handlers

- [`handle err { return _, wrap(err) }` ; `x := check foo()`](https://go.googlesource.com/proposal/+/master/design/go2draft-error-handling.md): block scoped
- [`handle err { return _, wrap(err) }` ; `x, # := foo()`](https://gist.github.com/oktalz/f04f36a3c2f61af22c7a6e06095d18eb)
- [`handle func(err error) (T, error) { return _, wrap(err) }` ; `x, ? := foo()`](https://github.com/rockmenjack/go-2-proposals/blob/master/error_handling.md)
- [`watch err { if err != nil { return _, wrap(err) } }` ; `x, err != foo()`](https://github.com/golang/go/issues/40821)
- [`expect err != nil { return _, wrap(err) }` ; `x, err := foo()`](https://github.com/golang/go/issues/32804)
- [`with { return _, wrap(err) }` ; `handle err { x, err := foo() }`](https://github.com/golang/go/issues/32795)

###### _call_ specific handler

- [`x, #err := foo()` ; `catch err { return _, wrap(err) }`](https://github.com/golang/go/issues/27519): tagged error handlers
- [`x, @err := foo()` ; `err: return _, wrap(err)`](https://gist.github.com/dpremus/3b141157e7e47418ca6ccb1fc0210fc7): labels and `goto`
- [`grab err() { if err != nil { return _, wrap(err) } }` ; `x, err() := foo()`](https://didenko.github.io/grab/grab_worth_it_0.1.1.html#12): assign to handler
- [`err := inline(err error){ if err != nil { return _, wrap(err) } }` ; `x, err := foo()`](https://github.com/gooid/gonotes/blob/master/inline_style_error_handle.md): assign to handler
- [`err := func(err error) (T, error) { return _, wrap(err) }` ; `x, #err := foo()`](https://gist.github.com/the-gigi/3c1acfc521d7991309eec140f40ccc2b): block scoped
- [`_check = func(err error) (T, error) { return _, wrap(err)}` ; `x, ?err := foo()`](https://gist.github.com/8lall0/cb43e1fa4aae42bc709b138bda02284e): not fully formed idea on return
- [`handler := func(err error) (T, error) { return^2 _, wrap(err) }` ; `x, handler(err) := foo()`](https://github.com/golang/go/issues/32473): note multilevel return

##### _wrapping_

some rely on `wrap` being smart and passing through `nil` (so not `fmt.Errorf`),

- [`x, err := foo()` ; `reflow _, wrap(err)`](https://github.com/golang/go/issues/21146): implicit `err != nil` and return
- [`x, err := foo()` ; `refuse _, wrap(err)`](https://gist.github.com/alexhornbake/6a4c1c6a0f2a063da6dda1bf6ec0f5f3)
- [`x, err := foo()` ; `pass _, wrap(err)`](https://github.com/golang/go/issues/37141)
- [`x, err := foo()` ; `ereturn _, wrap(err)`](https://github.com/golang/go/issues/38349)
- [`x, err := foo()` ; `err ?: return _, wrap(err)`](https://github.com/golang/go/issues/32946)
- [`x, err := foo()` ; `on err, return _, wrap(err)`](https://github.com/golang/go/issues/32611)
- [`x, err := foo()` ; `err ? { return _, wrap(err) }`](https://github.com/golang/go/issues/33067)
- [`x, err := foo()` ; `onErr { return _, wrap(err) }`](https://github.com/golang/go/issues/32946)
- [`x, err := foo()` ; `if err != nil { return _, wrap(err) }`](https://github.com/golang/go/issues/33113), [also](https://github.com/golang/go/issues/27135): if ... on 1 line
- [`x, err := foo()` ; `if !err { return _, wrap(err) }`](https://gist.github.com/fedir/50158bc351b43378b829948290102470)
- [`x, err := foo()` ; `if err { return _, wrap(err) }`](https://github.com/golang/go/issues/26712)
- [`x, err := foo()` ; `if err? { return _, wrap(err) }`](https://github.com/golang/go/issues/32845)
- [`x, err := foo(); if err != nil { return _, wrap(err) }`](https://gist.github.com/jozef-slezak/93a7d9d3d18d3fce3f8c3990c031f8d0), [also](https://github.com/golang/go/issues/27450): everything on 1 line
- [`x, err := foo() /* err */ { return _, wrap(err) }`](https://github.com/gooid/gonotes/blob/master/inline_style_error_handle.md), [also](https://github.com/golang/go/issues/41908)
- [`x, err := foo() ?? { return _, wrap(err) }`](https://github.com/golang/go/issues/37243)
- [`x, err := foo(); err.return wrap(err)`](https://github.com/golang/go/issues/39372)
- [`x := foo() or err: return _, wrap(err)`](https://github.com/golang/go/issues/33029)
- [`x := foo() ?err return _, wrap(err)`](https://github.com/golang/go/issues/33074)
- [`x := check wrap() foo()`](https://gist.github.com/jozef-slezak/93a7d9d3d18d3fce3f8c3990c031f8d0), [also](https://gist.github.com/morikuni/bbe4b2b0384507b42e6a79d4eca5fc61)
- [`x := foo() ? wrap()`](https://gist.github.com/gregwebs/02479eeef8082cd199d9e6461cd1dab3)
- [`x := foo() or wrap()`](https://github.com/golang/go/issues/36338)
- [`x := foo() || wrap(err)`](https://github.com/golang/go/issues/21161)
- [`x := foo() on_error err fail wrap(err)`](https://medium.com/@peter.gtz/thinking-about-new-ways-of-error-handling-in-go-2-e56d116952f1)
- [`x := foo() onerr return _, wrap(err)`](https://github.com/golang/go/issues/32848)
- [`x := try(foo(), wrap)`](https://github.com/golang/go/issues/32853)
- [`x := collect(&err, foo(), wrap)`](https://github.com/golang/go/issues/32880)
- [`try x, err := foo() { return _, wrap(err) }`](https://github.com/golang/go/issues/39890)

##### _return_

can use `defer` for wrapping

- [`x := try(foo())`](https://go.googlesource.com/proposal/+/master/design/32437-try-builtin.md)
- [`x := must(foo())`](https://github.com/golang/go/issues/32219): panic instead of return
- [`x := foo!()`](https://github.com/golang/go/issues/21155)
- [`x := foo()?`](https://gist.github.com/yaxinlx/1e013fec0e3c2469f97074dbf5d2e2c0), [also](https://github.com/golang/go/issues/39451)
- [`x := #foo()`](https://github.com/golang/go/issues/18721)
- [`x := guard foo()`](https://github.com/golang/go/issues/31442)
- [`x := must foo()`](https://gist.github.com/VictoriaRaymond/d70663a6ec6cdc59816b8806dccf7826)
- [`x, # := foo()`](https://github.com/golang/go/issues/22122): panic instead of return
- [`x, ? := foo()`](https://github.com/golang/go/issues/42214), [also](https://github.com/golang/go/issues/32601)
- [`x, ! := foo()`](https://gist.github.com/lldld/bf93ca94c24f172e95baf8c123427ace), [also](https://github.com/golang/go/issues/33150), [panic](https://github.com/golang/go/issues/35644)
- [`x, !! := foo()`](https://github.com/golang/go/issues/32884)
- [`x, !err := foo()`](https://github.com/golang/go/issues/14066)
- [`x, ^err := foo()`](https://github.com/golang/go/issues/42318)
- [`x, err? := foo()`](https://github.com/golang/go/issues/36390)
- [`x, err := foo() throws err`](https://github.com/golang/go/issues/32852)
- [`tryfunc func(...){ x := foo() }`](https://github.com/golang/go/issues/32964)
- [`x, err := foo()` ; `check(err)`](https://github.com/golang/go/issues/33233): builtin `if err != nil { return ..., err }` macro
- [`x, err := foo()` ; `catch(err)`](https://github.com/golang/go/issues/32811): builtin `if err != nil { return ..., err }` macro

##### _try..catch_

- [`try { x := foo() } catch(e Exception) { ??? }`](https://www.netroby.com/view/3910): literally try catch
- [`try err != nil { x, err := foo() } except { return _, wrap(err) }`](https://github.com/golang/go/issues/33387)
- [`until err != nil { check x, err := foo() } else { return _, wrap(err) }`](https://gist.github.com/coquebg/afe44e410f883a313dc849da3e1ff34c): insert after every `check`
- [`break err != nil { step: x, err := foo() }`](https://github.com/golang/go/issues/27075): insert after every repeatable label
- [`break err != nil { try x, err := foo() }`](https://github.com/golang/go/issues/27075): insert after every `try`
- [`if x, err := foo(); err != nil { return _, wrap(err) } else { ... } undo { ??? } done { ??? } defer { ??? }`](https://gist.github.com/jansemmelink/235228a0fb56d0eeba8085ab5f8178f3)
- [`handler := func(err error) error { return wrap(err) }` ; `check(handler){ x := foo() }`](https://devmethodologies.blogspot.com/2018/10/go-error-handling-using-closures.html)
- [`check { x := check foo() } handle err { return _, wrap(err) }`](https://gist.github.com/mathieudevos/2bdae70596aca711e50d1f2ff6d7b7cb)
- [`check { x, err1 := foo() } catch err { return _, wrap(err) }`](https://gist.github.com/eau-de-la-seine/9e2e74d6369aef4a76aa50976e34de6d)
- [`check { x, err := foo()` ; `catch: return _, wrap(err) }`](https://github.com/golang/go/issues/32968)
- [`x, err := foo() !!!` ; `fail: return _, wrap(err)`](https://github.com/golang/go/issues/34140)
- [`handle err { x, err := foo() ; case err != nil: return _, wrap(err) }`](https://github.com/golang/go/issues/35086)
- [`try { x := foo() ; if err != nil { return _, wrap(err) } }`](https://github.com/golang/go/issues/35179)

##### _others_

- [`x, (err) := foo()`](https://github.com/golang/go/issues/21732): only assign to LHS if `()` content is not currently `nil`
- [`collect err { x, _! := foo() }` ; `if err != nil { return _, wrap(err) }`](https://github.com/golang/go/issues/25626): err is an value that accumulates errors? does it continue?
- [`errorhandling(err){ x, err := foo() }`](https://github.com/Konstantin8105/Go2ErrorTree): err is an accumulator? messes with types
- [`x, err := foo()` ; `handle err1 || err2 || err3 { return _, wrap(err) }`](https://gist.github.com/Kiura/4826db047e22b7720d378ac9ac642027): shorter if chain?
- [multilevel return](https://github.com/golang/go/issues/35093)
- [multilevel return with return type](https://github.com/golang/go/issues/42811)
- [`returnfrom label, err`](https://gist.github.com/spakin/86ea86ca48aefc78b672636914f4fc23): multilevel return
