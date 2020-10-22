---
description: do maths on search and replace
title: search and replace maths
---

### _problem_

Say you're using terraform with some poorly designed cloud thing,
and you have a linked list of rules which also each needs their own index
as an attribute
[cloudflare/terraform-provider-cloudflare#187](https://github.com/cloudflare/terraform-provider-cloudflare/issues/187).
And now you have to insert an item somewhere close to the beginning of your 100 item chain.
Do you go edit every number by hand?

#### _(n)vim_

- `#` as delimiter (`|` doesn't work)
- `*` as count, (or `\?` or `\+` or `\{m,n}`)
- `\=` to specify using a [`:h sub-replace-expression`](https://neovim.io/doc/user/change.html#sub-replace-expression)

```
:%s#priority = \(\d*\)#\=printf("priority = %d", submatch(1)+1)#
```

#### _other_ search and replace stuff

one day I will read the entire manual, one day...

- [`:h sub-replace-special`](https://neovim.io/doc/user/change.html#sub-replace-special)
  - `\U...\E` make everything between uppercase
  - `\L...\E` make everything between lowercase
  - `CTRL-V <ENTER>` newline
- [`:h character-classes`](https://neovim.io/doc/user/pattern.html#/character-classes)
