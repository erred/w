---
description: it still doesn't work
title: wpa eap tls
---
### Get TLS chain from Wifi

create a monitoring interface

```bash
sudo iw dev wlp58s0 interface add wmon type monitor
sudo ip link set wmon up
```

wireshark filter for `eap` and (re)auth against the AP,
Check for `Server Hello`

_Note_: if you get `Unrecognized Protocol` in place of `Server Hello`,
just try again...

Export the certificates and start thinking about what to do with them :shrug: