---
description: renaming files with regex
title: renaming files
---

### _renaming_

Maybe if I write this down I will actually remember next time

#### _shell_ expansion

uses [parameter expansion](http://zsh.sourceforge.net/Doc/Release/Expansion.html#Parameter-Expansion)

quick:

- `${var#pattern}` shortest prefix pattern trim, `##` for longest
- `${var%pattern}` shortest suffix pattern trim, `%%` for longest
- `${var/pattern/repl}` replace 1st pattern with repl
- `${var//pattern/repl}` replace all pattern with repl

```sh
for f in *.txt; do
  mv $f "${f%txt}md"
done

# or with find for more precise targetting
find . -name '*.txt' -exec sh -c 'f={}; mv $f ${f%txt}md'
```

##### _zmv_

a zsh module in [other functions](http://zsh.sourceforge.net/Doc/Release/User-Contributions.html#Other-Functions)

- uses globbing (instead of regex)
- `()` to capture for reference with `$n`
  (a variable where parameter expansion can be applied)
- `$f` for filename

```sh
zmv '(*).txt' '$1.md'
```

#### _rename_

the one that comes with `util-linux` on Arch,
takes a pattern and what to replace it with and some files,
super undocumented

```sh
rename md txt *
```

#### _perl-rename_

basically sed for file names

```sh
perl-rename 's/(.*).md/\1.txt/' *
```
