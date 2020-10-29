---
description: down with GPG!
title: git signing signify
---

### _git_ signing

Ah, _git_, well known for its not very friendly ux,
integrates with another piece of software
that is also often derided for having bad ux,
namely _gpg_.

So you can sign tags, sign commits, and verify them.
But who wants to use GPG?
Why not use something else like minisign or signify
or SSH keys?

Well you can, sort of.
Git right now hardcodes some GPG stuff so you need to work around it.
Specifically it expects an out of band communication channel `--status-fd=X`
for getting notified of
`[GNUPG:] SIG_CREATED`, `[GNUPG:] GOODSIG`, `[GNUPG:] BADSIG`
and while you can return any text as the signature for the signing operations,
the verifying operations will look for a `-----BEGIN/END PGP SIGNATURE-----`
block that it can extract before the verifying program is called.

Note: if you do this, it is unlikely that anyone else will be able to
parse/verify your signatures and forget about it working in Github/etc.

So what you can do is smuggle your signature in the block
and write to the status fd as appropriate.
This is an ugly hack because you will need those fake PGP lines,
but you can put whatever inside.
This is demonstrated in a
[script by Leah Neukirchen](https://leahneukirchen.org/dotfiles/bin/git-signify) replicated below.

```sh
#!/bin/sh
# git-signify [GIT COMMAND] - use git with signify(1)
#
# First, you need to set the signing key for the repo, e.g.
#   git config --local user.signingKey ~/.signify/cwm
# This will use cwm.sec and cwm.pub.
#
# Then you can use
#   gpg signify commit -S
#   gpg signify verify-commit
#
#   gpg signify tag -s
#   gpg signify verify-tag
#
# You also can set this script as "gpg.program" to use signify
# automatically.
#
# To the extent possible under law, Leah Neukirchen has waived
# all copyright and related or neighboring rights to this work.
# http://creativecommons.org/publicdomain/zero/1.0/

getkey() {
    key=$(git config user.signingKey)
    if [ -z "$key" ]; then
        echo "git-signify: no user.signingKey defined!" 1>&2
        exit 7
    fi
}

while :; do
case "$1" in
-bsau)
    getkey
    echo "-----BEGIN PGP SIGNATURE----- (really git-signify)"
    {
        signify -S -s "$key.sec" -m - -x -
        if [ $? -eq 0 ] && [ -n "$statusfd" ]; then
            printf '\n[GNUPG:] SIG_CREATED ' >/dev/fd/$statusfd
        fi
    } | sed "s/: .*/: verify with git-signify and ${key##*/}.pub/"
    echo "-----END PGP SIGNATURE-----"
    exit 0
    ;;
--verify)
    getkey
    sed -i '/-----.* PGP SIGNATURE-----/d' "$2"
    if signify -V -p "$key.pub" -m - -x "$2" 1>&2; then
        echo "[GNUPG:] GOODSIG "
        exit 0
    else
        r=$?
        echo "[GNUPG:] BADSIG "
        exit $r
    fi
    ;;
--status-fd=*)
    statusfd=${1#--status-fd=}
    shift
    ;;
--*)
    # ignore all other arguments
    shift
    ;;
*)
    exec git -c "gpg.program=$0" "$@"
    ;;
esac
done
```
