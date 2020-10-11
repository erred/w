---
description: oneliners to generate self signed certs with openssl
title: openssl oneline curve certs
---

### _self_ signed certs

you need a cert,
it doesn't need the strict requirements of WebPKI,
and you're tired of all the online blogs that use RSA
and outdated OpenSSL that doesn't support `-addext`

#### _ECDSA_

This should be simple and have wide support

```sh
openssl req -x509 -nodes -days 7300 \
  -newkey ec:<(openssl ecparam -name prime256v1) \
  -subj "/O=weechat/CN=localhost" \
  -addext "subjectAltName=DNS:localhost,IP:0.0.0.0" \
  -keyout relay.pem -out relay.pem
```

This will generate a merged key/cert,
some old stuff (like weechat) wants this,
else separate the `-keyout` and `-out`.

#### _Ed25519_

More exotic, but at least Go and Curl work

```sh
openssl req -x509 -nodes -days 7300 \
  -key <(openssl genpkey -algorithm ED25519 | tee relay.pem) \
  -subj "/O=weechat/CN=localhost" \
  -addext "subjectAltName=DNS:localhost,IP:0.0.0.0" \
  -keyout relay.pem -out relay.pem
```

This gives `req` a fully formed key instead of just the params,
which gives us some issues with outputing the key.

This will also generate a merged key/cert,
don't ask me why the private key needs to be written twice
(`tee` and `-keyout`), otherwise it will only output a cert.
For separate key/cert, remove the `-keyout` and rename the `tee`

#### _inspect_ the cert

inspect the generated cert with

```sh
openssl x509 -in relay.pem -text -noout
```

or from a running server,
where `-servername` is used for `SNI`
and `connect` is for the destination

```sh
openssl s_client -showcerts \
  -servername localhost \
  -connect 127.0.0.1:8080 </dev/null \
  | openssl x509 -text -noout
```
