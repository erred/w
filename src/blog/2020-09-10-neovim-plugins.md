---
description: plugin ecosystem for neovim Sep. 2020
title: neovim plugins
---

### _neovim_ plugins

#### _plugin_ manager

you have plugins, but how do you maintain all of them?

##### _runtimepath_

These plugin managers only update the `runtimepath`
so the plugin files are loaded correctly.
This extends the native support for single files
to entire directory trees/bundles/packages.
It is the user's responsibility to retrieve and keep up to date the plugin files.

examples:

- vim8 packages
- neovim support for vim8 packages
- [pathogen][pathogen]

##### _full_ manager

These are fully featured:
you give them a line in your config file `init.nvim`,
and they can download, update, clean up packages with a single command.

examples:

- [vim-plug][plug]
- [vundle][vundle]
- [dein][dein]
- [minpac][minpac] makes use of native vim8 packages for both vim/nvim

[pathogen]: https://github.com/tpope/vim-pathogen
[vundle]: https://github.com/VundleVim/Vundle.vim
[plug]: https://github.com/junegunn/vim-plug
[dein]: https://github.com/Shougo/dein.vim
[minpac]: https://github.com/k-takata/minpac
