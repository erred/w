---
description: encrypting and decrypting files with gitattributes
title: gitattributes age encrypt
---

### _gitattributes_

What is that file for? Per path settings.

One of those settings is `filter` which can transform a file
when it's staged `clean` and when its checked out `smudge`.

#### _age_

[age](https://age-encryption.org/)
Because why not have a sane file encrypting tool?

#### _setup_

##### _gitconfig_

Yes, this global config.
But your secrets are global too, surely can't hurt much.

Use ascii armoured because binary files look weird on the web.

```txt
# ~/.config/git/config
[filter "ageencrypt"]
    clean = age -r age14mg08panez45c6lj2cut2l8nqja0k5vm2vxmv5zvc4ufqgptgy2qcjfmuu -a -
    smudge = age -d -i ~/.ssh/age.key -
    required = true
```

#### _gitattributes_

This is per repo config

```gitattributes
secret.yaml filter=ageencrypt
```
