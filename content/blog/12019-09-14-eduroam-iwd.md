---
description: eduroam, iwd version
title: eduroam iwd
---

My [previous attempts](/blog/12019-09-04-eduroam/)
to get eduroam wifi with [wpa_supplicant](https://wiki.archlinux.org/index.php/WPA_supplicant)
more or less worked.
But some recent updates seem to have made everything unstable again.
So why not try [iwd](https://wiki.archlinux.org/index.php/Iwd)

#### Correction 2019-09-16

It doesn't work at another location
back to `wpa_supplicant`

#### _UvA_ (University of Amsterdam)

The only account I have access to right now,
they supposedly use `TTLS` with `MSCHAPV2` for phase2,
**which works as described for `wpa_supplicant`**,
but `iwd` is weird and the error messages are beyond useless even with debugging turned on.
Following the advice of some archlinux forum post to "play around with the eap method",
`PEAP` works :facepalm:

```
[Security]
EAP-Method=PEAP
EAP-Identity=anonymous@uva.nl
EAP-PEAP-Phase2-Method=MSCHAPV2
EAP-PEAP-Phase2-Identity=...@uva.nl
EAP-PEAP-Phase2-Password=...

[Settings]
Autoconnect=true
```

