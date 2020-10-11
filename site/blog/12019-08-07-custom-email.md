---
description: setup email in 2019
title: custom email
---
Email is annoyingly hard,
if you want it to be _reliable_

I've setup mine to work with [Mailgun](https://www.mailgun.com/)
and [Firebase](https://firebase.google.com/) (for transactional emails)

### _Receiving_

This is easy

```
MX    seankhliao.com    mxa.mailgun.org
MX    seankhliao.com    mxb.mailgun.org
```

### Sending

This is complicated to do safely and securely

#### SPF

_Sender Policy Framework_

Identifies which mail servers are allowed to send email from the domain

```
TXT     seankhliao.com          v=spf1 include:_spf.firebasemail.com include:mailgun.org ~all
```

#### DKIM

_Domain Keys Identied Mail_

Email hashes are signed and verified with public-private keys,
public keys identified by the `DKIM-Signature` header,
`...; s=selector; d=example.com` corresonds to the key at `selector._domainkey.example.com`

```
CNAME   firebase1._domainkey    mail-seankhliao-com.dkim1._domainkey.firebasemail.com
CNAME   firebase2._domainkey    mail-seankhliao-com.dkim2._domainkey.firebasemail.com
TXT     smtp._domainkey         v=DKIM1; k=rsa; p=...
```

#### DMARC

_Domain-based Message Authentication, Reporting & Conformance_

Specifies the policies (monitor, quarantine **spam**, reject) that should be applied when verifying against `SPF` and `DKIM`

Set another `TXT some.doman._report._dmarc v=DMARC1` if you want to receive reports for `somae.domain` at antoher domain

```
TXT     _dmarc                  v=DMARC1; p=reject; adkim=r; aspf=r; ruf=mailto:...; rua=mailto:...; fo=1; ri=86400; pct=100; rf=afrf;
```