--- title
signing and encrypting
--- description
modern signing and encrypting
--- main

### What

_kill_ gpg

##### _With_...?

[ssh-keygen](http://man7.org/linux/man-pages/man1/ssh-keygen.1.html)
and
[age](https://age-encryption.org)

### tldr

```
# encrypt for me
# use: encrypt file1 file2...
function encrypt() {
  # local pubkey="ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIKnAmz4u5/51kPPsebDCiYTXvuftUORh/TJ4pvN3NvQa"
  local pubkey=age14mg08panez45c6lj2cut2l8nqja0k5vm2vxmv5zvc4ufqgptgy2qcjfmuu
  for f in "$@"; do
    age -r ${pubkey} -o ${f}.age ${f}
  done
}

# decrypt for me
# use: decrypt file1 file2...
function decrypt() {
  # local privkey=$HOME/.ssh/id_ed25519
  local privkey=$HOME/keys/age.key
  for f in "$@"; do
    age -d -i ${privkey} -o ${f} ${f%%.age}
  done
}

# sign by me
# use: sign file1 file2...
function sign() {
  local privkey=$HOME/.ssh/id_ed25519
  for f in "$@"; do
    ssh-keygen -Y sign -f ${privkey} -n signed@seankhliao.com ${f}
  done
}

# verify by me
# use: verify file1 file2...
function verify() {
  local accepted=$HOME/keys/ssh-sign-accepted
  for f in "$@"; do
    ssh-keygen -Y verify -n signed@seankhliao.com -f ${accepted} -I arccy@eevee -s ${f}.sig < ${f}
  done
}
```

#### Encrypt

with _age_

##### Install

no, `go get` doesn't work

```
git clone https://github.com/FiloSottile/age
cd age && go install ./cmd/...
```

##### Keygen

or use ssh keys

```
age-keygen -o age.key
```

##### Encrypt

```
age -r "public key of recipient" -o output.file.age input.file

"age public key": "age14mg08panez45c6lj2cut2l8nqja0k5vm2vxmv5zvc4ufqgptgy2qcjfmuu"
"ssh public key": "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIKnAmz4u5/51kPPsebDCiYTXvuftUORh/TJ4pvN3NvQa"
```

##### Decrypt

```
age -d -i path/to/private.key -o output.file input.file.age
```

#### Sign

with _ssh-keygen_.
why keygen? I don't know

##### Keygen

```
ssh-keygen -t ed25519
```

##### Sign

```
ssh-keygen -Y sign -f path/to/private.key -n file@seankhliao.com input.file

```

- -n: takes an arbitrary string as namespace, recommended: namespace@domain.tld

##### Verify

```
ssh-keygen -Y verify -n file@seankhliao -f accepted.file -I identity -s input.file.sig < input.file
```

- -n: takes an arbitrary string as namespace, recommended: namespace@domain.tld
- -I: takes an identity, needs to match the principals in accepted_signers, not related to identity at the end of pubkey used to sign file
- -f: takes a file of accepted_signers of the following format

```
# comments
user@domain key-type KEYGOESHERE

# certs signed by this CA are accepted
*@domain cert-authority key-type KEYGOESHERE

user@domain namespaces="whitelist,of,namespaces" key-type KEYGOESHERE
```
