---
description: accounts with 2fa that actually work
style: |-
  main > table {
    grid-column: 1 / span 3;
  }
  table {
    border-collapse: collapse;
    border-style: hidden;
  }
  th, td {
    padding: 0.4em;
    text-align: left;
  }
  th {
    font-weight: 700;
    border-bottom: 0.2em solid #999;
  }
  tr:nth-child(5n) td {
    border-bottom: 0.1em solid #999;
  }
  tbody tr:hover {
    background: #404040;
  }
title: 2fa
---

### _2fa_ by service

#### Goals

- username + password + security key
- backup code as backup
- totp > sms but both a phishable

#### Legend

- yes: available
- no: not available
- opt: optional
- req: required either for 2fa or for security key

#### _Services_

| Service       | SMS | TOTP | Key | Backup Code |
| ------------- | --- | ---- | --- | ----------- |
| Adobe         | opt | yes  | no  | yes         |
| Amazon        | req | yes  | no  | no          |
| Booking       | yes | no   | no  | no          |
| _Cloudflare_  | no  | opt  | yes | yes         |
| Docker        | no  | yes  | no  | yes         |
| _Dropbox_     | opt | req  | yes | yes         |
| _Facebook_    | opt | req  | yes | yes         |
| _GitHub_      | no  | req  | yes | yes         |
| _Google norm_ | opt | opt  | yes | yes         |
| _Google adv_  | no  | no   | yes | no          |
| Instagram     | opt | yes  | no  | yes         |
| Keybase       | no  | no   | no  | yes         |
| LinkedIn      | opt | yes  | no  | yes         |
| Mastadon      | no  | yes  | no  | yes         |
| Microsoft     | opt | yes  | no  | yes         |
| _Namecheap_   | opt | opt  | yes | yes         |
| Paypal        | req | yes  | no  | no          |
| Slack         | no  | yes  | no  | yes         |
| Twitch        | no  | yes  | no  | yes         |
| Twitter       | no  | req  | 1   | yes         |

##### Notes

- Keybase requires another active device
