# com-seankhliao

[![License](https://img.shields.io/github/license/seankhliao/com-seankhliao.svg?style=flat-square)](LICENSE)
![Version](https://img.shields.io/github/v/tag/seankhliao/com-seankhliao?sort=semver&style=flat-square)

handcrafted personal static site

## About

Who am i?

## style

```txt
--black: #000000;
--primary: #a06be0;
--gray: #999;
--white: #eceff1;
--mono: "Inconsolata", monospace;
--serif: "Lora", serif;
```

## csp

```txt
Report-To: {"group": "csp-endpoint", "max_age": 10886400, "endpoints": [{"url":"https://statslogger.seankhliao.com/json"}]}
```

```txt
default-src 'self';
upgrade-insecure-requests;
connect-src https://statslogger.seankhliao.com https://www.google-analytics.com;
font-src https://static.seankhliao.com;
img-src *;
object-src 'none';
script-src-elem 'nonce-deadbeef2' 'nonce-deadbeef3' 'nonce-deadbeef4' https://static.seankhliao.com https://www.google-analytics.com https://ssl.google-analytics.com https://www.googletagmanager.com;
sandbox allow-scripts;
style-src-elem 'nonce-deadbeef1' https://static.seankhliao.com;
report-to csp-endpoint;
report-uri https://statslogger.seankhliao.com/json";
```

### notes

#### gtm csp

```txt
script-src: 'unsafe-inline' https://www.googletagmanager.com
img-src: www.googletagmanager.com
```

#### ga csp

```txt
script-src: https://www.google-analytics.com https://ssl.google-analytics.com
img-src: https://www.google-analytics.com
connect-src: https://www.google-analytics.com
```

#### gtm preview csp

```txt
script-src: https://tagmanager.google.com
style-src: https://tagmanager.google.com https://fonts.googleapis.com
img-src: https://ssl.gstatic.com https://www.gstatic.com
font-src: https://fonts.gstatic.com data:
```
