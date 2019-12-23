--- title
standalone coredns
--- description
running coredns in standalone mode
--- main


so you need DNS...

# _Corefile_

by default at `/etc/coredns/Corefile`

```
. {
  # still not sure if this works
  forward . 8.8.8.8

  # nice to have
  loop
  whoami

  # logging
  errors
  log
}
```

# configure from file

you have a zone
(in that bind format)

```
$TTL 30
@     IN SOA  ns1 admin.example.com. 0000000001 30 30 360000 30

@     IN A    1.1.1.1
sub   IN TXT  "hello world"
```

your options are

## plugin/auto

discover zones from a directory and serve,

```
. {
  auto example.com {
    # zone files named db.(your.zone)
    directory /path/to/directory
  }
}

. {
  auto {
    # use non default regex
    directory /path/to/directory regex[extract](zonename.*).from.file {1}
  }
}
```

## plugin/file

serve a single file,
zone _must_ be in the server block or in the file plugin
else it probably won't work as expected

```
. {
  file /path/to/zone/file example.com
}

example.com {
  file /path/tp/zone/file
}
```

# DNSSEC

in for a world of pain

## plugin/dnssec

- supports RSA/ECDSA/ED25519 (probably,
  [issue #3379](https://github.com/coredns/coredns/issues/3379),
  [pr #3380](https://github.com/coredns/coredns/pull/3380)
  ) keys
- supports KSK/ZSK split
- NSEC only, no NSEC3 (key types still valid)
- signed on demand (zone transfer will not show a signed zone)

Corefile:

```
. {
  dnssec {
    key file /path/to/key1 /path/to/key/2
  }
}
```

## plugin/sign

- supports RSA/ECDSA/ED25519 (probably,
  [issue #3379](https://github.com/coredns/coredns/issues/3379),
  [pr #3380](https://github.com/coredns/coredns/pull/3380)
  ) keys
- No KSK/ZSK split, must be a KSK key
- NSEC only, no NSEC3 (key types still valid)
- entire zone is presigned

Corefile:

```
. {
  file /path/to/signed/db.(your.zone).signed
  sign /path/to/unsigned/your.unsigned.zone {
    key file /path/to/key
    directory /path/to/signed
  }
}
```
