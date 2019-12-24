--- title
experiments with SXG
--- description
experimenting with HTTP Signed Exchanges
--- main

### _HTTP_ Signed Exchanges

SXG signs an individual request/response pair
to authenticate the resource originally came from somewhere.

see: [WICG/webpackage/go/signedexchange](https://github.com/WICG/webpackage/tree/master/go/signedexchange)

#### generate _private_ key:

for a self signed cert:

1. create private key
2. generate certificate signing request
3. self sign certificate

###### prime256v1 ecdsa

```
openssl ecparam -name prime256v1 -genkey -out sxg.prime256v1.key
openssl req -new -sha256 -key sxg.prime256v1.key -out sxg.prime256v1.csr -subj '/CN=seankhliao.com/O=Test/C=US'
openssl x509 -req -days 90 -in sxg.prime256v1.csr -signkey sxg.prime256v1.key -out sxg.prime256v1.pem -extfile <(echo -e "1.3.6.1.4.1.11129.2.1.22 = ASN1:NULL\nsubjectAltName=DNS:seankhliao.com")
```

###### ed25519

```
openssl genpkey -algorithm ed25519 -out sxg.ed25519.key
openssl req -new -sha256 -key sxg.ed25519.key -out sxg.ed25519.csr -subj '/CN=seankhliao.com/O=Test/C=US'
openssl x509 -req -days 90 -in sxg.ed25519.csr -signkey sxg.ed25519.key -out sxg.ed25519.pem -extfile <(echo -e "1.3.6.1.4.1.11129.2.1.22 = ASN1:NULL\nsubjectAltName=DNS:seankhliao.com")
```
