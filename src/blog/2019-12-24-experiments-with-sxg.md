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

apparently it _must_ be prime256v1?

**TODO**: try other keys

```
openssl ecparam -name prime256v1 -genkey -out sxg.prime256v1.key
openssl req -new -sha256 -key sxg.prime256v1.key -out sxg.prime256v1.csr -subj '/CN=seankhliao.com/O=Test/C=US'
openssl x509 -req -days 90 -in sxg.prime256v1.csr -signkey sxg.prime256v1.key -out sxg.prime256v1.pem -extfile <(echo -e "1.3.6.1.4.1.11129.2.1.22 = ASN1:NULL\nsubjectAltName=DNS:seankhliao.com")
```

#### generate cbor

ignore warnings

```
gen-certurl -pem sxg.prime256v1.pem -ocsp <(echo ocsp) > cert.cbor
```

#### write _go_ to sign stuff

```
import (
        "github.com/WICG/webpackage/go/signedexchange"
        "github.com/WICG/webpackage/go/signedexchange/version"
)
```

##### create a _Signer_

```
signer := &signedexchange.Signer{
        Date: time.Now(),

        // Max, see section 3.5
        // https://wicg.github.io/webpackage/draft-yasskin-http-origin-signed-responses.html#name-signature-validity
        Expires:     time.Now().Add(7 * 24 * time.Hour),

        // sxg.prime256v1.pem
        // certs, err := signedexchange.ParseCertificates(certsb)
        Certs:       certs,

        // url to locate cert.cbor
        // certURL, err := url.Parse("")
        CertUrl:     certURL,

        // still don't know what this is for
        ValidityUrl: validityURL,

        // sxg.prime256v1.key
        // privKey, err := signedexchange.ParsePrivateKey("")
        PrivKey:     privKey,
}
```

##### create an _Exchange_ per resource

```
// uri: location of resource to be signed (fallback)
// resHeader: response headers (content type mandatory?)
// payload: content
e := signedexchange.NewExchange(version.Version1b3, uri, http.MethodGet, http.Header{}, http.StatusOK, resHeader, payload)

// encode and sign
e.MiEncodePayload(4096)
e.AddSignatureHeader(signer)
```
