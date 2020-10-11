---
description: basic notes on zsh completion system
title: zsh completion
---

### _completion_

`TAB` all the way.

#### _zsh_

[zsh manual](http://zsh.sourceforge.net/Doc/Release/Completion-System.html)

##### _compinit_

initialize the completion system

```sh
compinit -d ${XDG_CACHE_HOME:-$HOME}/.zcompdump
```

##### _compdef_

using a shell function generate completions

```sh
compdef function_to_use commands to complete

# when used as the first line of an autoloaded function found during compinit
#compdef ...

# reuse completions created for command2 for command1
compdef command1=command2
# ex:
alias tf=terraform
compdef tf=terraform
```

#### _bash_

use bash completion in zsh

[bash manual](https://www.gnu.org/software/bash/manual/html_node/Programmable-Completion-Builtins.html)

##### _bashcompinit_

call after `compinit`

```sh
bashcompinit
```

##### _complete_

how to complete commands

```sh
complete -o option -[CF] command

# commonly used to not add extra space
-o nospace

complete -C command_output_used_as_completion command
complete -F function_output_used_as_completion command
```
