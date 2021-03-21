---
title: development environment
description: what do I need to write code
---

### _dev_

So apparently I don't write as much code as I try to get things running,
wonder how this affects my environment.

#### _host_

It's [Arch Linux](https://archlinux.org/),
because I don't want to have to worry about packages
being too far out of date from upstream
or that a package might not even be available.
Also, it's pretty no-bullshit in terms of what the distrobution forces on you.

While not technically part of the host,
[docker](https://www.docker.com/) is invaluable for having self-contained,
disposable environments that could also easily be shared.
Yes, I know about [NixOS](https://nixos.org/) but who has time for that?

The 2 types of windows I have open are:
terminal ([alacritty](https://github.com/alacritty/alacritty))
and browser ([chrome](https://www.google.com/chrome/),
now with wayland and pipewire enabled, sort of janky).

[zsh](https://www.zsh.org/) is a decent balance between usable and exotic,
preferably with
[autosuggestions](https://github.com/zsh-users/zsh-autosuggestions),
[syntax-highlighting](https://github.com/zsh-users/zsh-syntax-highlighting),
[history-substring-search](https://github.com/zsh-users/zsh-history-substring-search),
and [completions](https://github.com/zsh-users/zsh-completions).

#### _text_ editor

[Neovim](https://neovim.io/) it is,
running `-git` for
[tree-sitter](https://github.com/tree-sitter/tree-sitter) based highlighting
and with [coc.nvim](https://github.com/neoclide/coc.nvim) for completion.
The other plugins are:
[cpg/vim-fahrenheit](https://github.com/fcpg/vim-fahrenheit) for a warm color scheme,
[mhinz/vim-signify](https://github.com/mhinz/vim-signify) for git gutter,
[tyru/caw.vim](https://github.com/tyru/caw.vim) for commenting shortcut,
and [sheerun/vim-polyglot](https://github.com/sheerun/vim-polyglot) for extra languages(?),
to be honest, not sure if it has any appreciable effect.

On top of this, a whole host of [Language servers](https://langserver.org/)
to power completion, as well as [prettier](https://prettier.io/)
for some formatting.

#### _other_ tools

[git](https://git-scm.com/) is still the VCS of choice for now,
with a few custom aliases
and occasional use of GitHub's [cli](https://github.com/cli/cli).
Also [delta](https://github.com/dandavison/delta) for highlighted / side-by-side diffs.

[exa](https://github.com/ogham/exa) is a replacement `ls` for color and git status,
[bat](https://github.com/sharkdp/bat) for colorizing command output,
[htop](https://htop.dev/) for managing processes
(no, none of the newer ones are significantly better).
[ripgrep](https://github.com/BurntSushi/ripgrep) for replacing grep,
with my own wrapper [t](https://github.com/seankhliao/t)
adding results as shortcuts.

[curl](https://curl.se/) continues to be the tool of choice for most network things,
and I prefer [drill](https://linux.die.net/man/1/drill) for DNS things.

For wrangling data:
[jq](https://stedolan.github.io/jq/) for json
and [xsv](https://github.com/BurntSushi/xsv) for CSV / TSV,
along with the usual suspects of:
`sed` / `tr` / `sort` | `uniq` | `head` | `tail` and very rarely `awk`.

Finally, [kubernetes](https://kubernetes.io/):
[kubectl](https://kubectl.docs.kubernetes.io/) is mandatory,
and [krew](https://github.com/kubernetes-sigs/krew) will help install some plugins as subcommands under different names, though then they don't have completion....
[kubectx](https://github.com/ahmetb/kubectx) for changing default context (cluster) and namespace,
[kustomize](https://kubernetes-sigs.github.io/kustomize/) for generating manifests
(`kubectl ... -k` was stuck at `2.0.3` for the longest time but `1.21` will get `4.0.5`),
[skaffold](https://skaffold.dev/) for rebuild-deploy loops
and finally [k9s](https://github.com/derailed/k9s) for a navigatable dashboard.
