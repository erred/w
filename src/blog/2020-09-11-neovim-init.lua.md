---
description: early look at lua in neovim
title: neovim init.lua
---

### _neovim_

[neovim][nvim] has partial support for [lua][lua],
specifically [LuaJIT][luajit].
This looks like a possibly great way to replace vimscript,
because why not?

tested with `NVIM v0.5.0-678-ga621c45ba`

example: [init.lua][conf], [permalink][confp]

#### _init.lua_

Proposed but not yet implemented.

#### _options_

All your `set XXX`, `set XYZ=abc`
can be represented as `vim.o.XXX`, `vim.wo.YYY`, `vim.bo.ZZZ`.
`vim.o.*` won't error with unknown options
but also (currently) won't work with all options.

#### _plugins_

Why not use the new built in support for loading vim packages?
this uses [minpac][minpac].

The differences in calling method appear to be because they
are builtin vs functions.

```lua
vim.cmd('packadd minpac')
vim.fn['minpac#add']('k-takata/minpac', {type = 'opt'})
```

#### _augroups_ and functions

augroups live in `nvim_exec` and functions in `nvim_command`.
Not quite sure why they need to be different.

#### _keymaps_

`nvim_set_keymap` works, but not sure what `cnoreabbrev` turns into.

[nvim]: https://neovim.io/
[lua]: http://www.lua.org/
[luajit]: https://luajit.org/
[conf]: https://github.com/seankhliao/config/blob/master/nvim/init.lua
[confp]: https://github.com/seankhliao/config/blob/6193c3e30610ec872626f0f4b40e3e623a025004/nvim/init.lua
[minpac]: https://github.com/k-takata/minpac
