---
title: istio gateway mtls
description: mtls notes
---

### _gateway_

[Istio](https://istio.io/latest/)
[Gateway](https://istio.io/latest/docs/reference/config/networking/gateway/)
is the way things get into your service mesh (cluster).
Now, maybe you want to push authentication into the service mesh.

Original plan: use [pomerium](https://pomerium.io/)
as a forward auth service like i did for nginx/traefik.
Problem: istio does do forward auth,
but not in a way that pomerium supports
and I don't want to proxy everything through pomerium.

Plan B: use the entire [Ory](https://www.ory.sh/) stack.
Problem: Ory is API first, need to implement my own UI.
Probably save this for later.

Plan C: use mTLS everywhere.

#### _mTLS_

mTLS works on hosts/domains. so you can't just protect a few paths...
So if you were wondering between a single `*.example.com` gateway
vs a gateway for each service, here's your reason.

##### _cert_ manager

Get those certs. I'm using
[cert-manager](https://cert-manager.io/)
[ca](https://cert-manager.io/docs/configuration/ca/).

One way is to:

- create a self-signing issuer
- create a ca certificate using said issuer
- create a ca using said certificate
- create mtls certificates using said ca

###### _issuer_

```yaml
# self signed CA root provisioner
apiVersion: cert-manager.io/v1
kind: Issuer
metadata:
  name: private-ca-root
spec:
  selfSigned: {}
---
# self signed CA
apiVersion: cert-manager.io/v1
kind: Issuer
metadata:
  name: private-ca
spec:
  ca:
    secretName: private-ca-root
```

###### _certs_

```yaml
# CA root cert
apiVersion: cert-manager.io/v1
kind: Certificate
metadata:
  name: private-ca-root
spec:
  secretName: private-ca-root
  duration: 87600h # 10y
  renewBefore: 8760h # 1y
  privateKey:
    algorithm: ECDSA
  isCA: true
  commonName: cluster27-private-ca
  issuerRef:
    name: private-ca-root
---
# cert for gateways
apiVersion: cert-manager.io/v1
kind: Certificate
metadata:
  name: k-mtls
spec:
  secretName: k-mtls
  duration: 2160h # 3m
  renewBefore: 720h # 1m
  privateKey:
    algorithm: ECDSA
  usages:
    - server auth
  dnsNames:
    - "k.seankhliao.com"
    - "*.k.seankhliao.com"
  issuerRef:
    name: private-ca
---
# cert for clients
apiVersion: cert-manager.io/v1
kind: Certificate
metadata:
  name: client-eevee
spec:
  secretName: client-eevee
  keystores:
    pkcs12:
      create: true
      passwordSecretRef:
        name: p12-password
        key: password
  duration: 8760h # 1y
  renewBefore: 2160h # 3m
  privateKey:
    algorithm: ECDSA
  usages:
    - client auth
  emailAddresses:
    - eevee@seankhliao.com
  issuerRef:
    name: private-ca
```

##### _gateway_ setup

Now to use it. On our gateway:

```yaml
apiVersion: networking.istio.io/v1alpha3
kind: Gateway
metadata:
  name: httpbin
spec:
  selector:
    istio: ingressgateway
  servers:
    - port:
        name: https-httpbin
        number: 443
        protocol: HTTPS
      hosts:
        - "httpbin2.k.seankhliao.com"
      tls:
        mode: MUTUAL
        credentialName: k-mtls
        minProtocolVersion: TLSV1_3
    - port:
        name: http-httpbin
        number: 80
        protocol: HTTP
      hosts:
        - "httpbin2.k.seankhliao.com"
      tls:
        httpsRedirect: true
```

##### _client_ setup

So to get our devices to trust and use private ca...

##### _arch_ linux

get the `ca.crt` from any one of them,

```sh
sudo trust anchor --store ca.crt
```

Get the `keystore.p12` from `client-x`,
import in browser [chrome://settings/certificates](chrome://settings/certificates)

##### _android_

get `ca.crt` from any one of them,
Settings > Security > Encryption & Credentials > Install a certificate > CA certificate

Get the `keystore.p12` from `client-y`,
Settings > Security > Encryption & Credentials > Install a certificate > VPN & app user certificate
