---
description: why is finding a good diff tool so hard?
title: diff tools
---

### _diff_ tools

I use git in a terminal, and I want:

- side by side diff
- colorful (red/green) diff
- no bloated dependencies

basically what Github has

#### _note_

git diff calls external tools with 7 args:

```txt
path old-file old-hex old-mode new-file new-hex new-mode
```

meaning in most cases you need a basic wrapper script
to pass just `$2` and `$5`

```sh
#!/bin/sh
your-diff-tool "$2" "$5"
```

#### _options_

_tldr_: side by side only with python or vim/nvim

- [ccdiff][ccdiff]
  - ~ perl, git uses it so meh
  - \- no side by side view
- [cdiff][cdiff]
  - fork of ydiff, unknown differences
- [colordiff][colordiff]
  - only adds color to a normal diff
- [cwdiff][cwdiff]
  - wrap wdiff with color
- [diff-so-fancy][diff-so-fancy]
  - \+ word diff
  - ~ perl, git uses it so meh
  - \- no side by side view
- [dwdiff][dwdiff]
  - only does word diff
- [icdiff][icdiff]
  - \+ side by side
  - \+ partial diff
  - \- python
- [nvim -d][nvim -d]
  - \+ side by side
  - \+ editor tooling
  - \- difficult to close diffs of multiple files
- [vim -d][vim -d]
  - also known as vimdiff
  - \+ side by side
  - \+ editor tooling
  - \- difficult to close diffs of multiple files
- [wdiff][wdiff]
  - only does word diff
- [ydiff][ydiff]
  - \+ side by side
  - \+ word diff
  - \- python

[ydiff]: https://github.com/ymattw/ydiff
[cdiff]: https://github.com/amigrave/cdiff
[colordiff]: https://www.colordiff.org/
[icdiff]: https://github.com/jeffkaufman/icdiff
[diff-so-fancy]: https://github.com/so-fancy/diff-so-fancy
[vim -d]: http://vimdoc.sourceforge.net/htmldoc/diff.html
[nvim -d]: https://neovim.io/doc/user/diff.html
[wdiff]: https://www.gnu.org/software/wdiff/
[ccdiff]: https://metacpan.org/pod/ccdiff
[dwdiff]: https://os.ghalkes.nl/dwdiff.html
[cwdiff]: https://github.com/junghans/cwdiff
