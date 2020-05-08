---
description: http metadata
title: http headers
---

### _headers_

so much metadata to read/set

from the point of view of a server

request , **response**

#### generic

- **Server**: info about the server `nginx`
- User-Agent: info about the client `Go-http-client/1.1`

#### content

- **Content-Disposition**: download instead of display `attachment; filename="file.jpg"`
-

##### _negotiation_

- Accept: content types to accept `<MIME_type>/<MIME_subtype>`
- **Content-Type**: `<MIME_type>/</MIME_subtype>`

- Accept-Encoding: compression algorithms to accept `gzip, br`
- **Content-Encoding**: `br`

#### caching

- **Clear-Site-Data**: `"cache", "cookies", "storage", "executionContexts"` or `"*"`
- **Cache-Control**: disable caching: `no-store`, static assets: `public, max-age=604800, immutable`
- **Vary**: headers to match for caching
- **Etag**: `"opaque unique string"`
- If-None-Match: `"list", "of", "etags"` return _304 Not Modified_ or _200_

#### cookies

- Cookie: list of cookies from client
- **Set-Cookie**: cookies to set in client

##### _cross origin_

- Origin: origin the request came from `https://example.com:1234`
- **Access-Control-Allow-Origin**: `*` or echo back the _Origin_ header
- **Access-Control-Allow-Credentials**: `true`
- **Access-Control-Allow-Methods**: `OPTIONS, GET, POST`
- **Access-Control-Allow-Headers**: additional headers to allow in request
- **Access-Control-Expose-Headers**: allow clients to access additional headers
- **Access-Control-Max-Age**: seconds the preflight response is valid for `600`

##### _security_

- **Strict-Transport-Security**: https only comms, 2 year preload `max-age=63072000; includeSubDomains; preload`
- **Content-Security-Policy**: limit origins for resources `default-src 'self'; ...`
- **Content-Security-Policy-Report-Only**: report, do not enforce version of CSP
- **Feature-Policy**: limit apis for page `microphone 'none'; ...`
- Sec-Fetch-Dest: where the response will end up `font`, `image`, ...
- Sec-Fetch-Mode: type of request `navigate`, `same-origin`, `cors`, ...
- Sec-Fetch-Site: relation between source and current origin `cross-site`, `same-origin`. ...
- Sec-Fetch-User: user triggered `?0`, `?1`

##### _privacy_

- Referer: previous page that linked to current request `https://example.com/hello.html`
- **Referrer-Policy**: how much to include in _Referer_ header `no-referrer`, `strict-origin-when-cross-origin`

##### _client hints_

- **Accept-CH**: client hints to accept `Width`
- **Accept-CH-Lifetime**: how long to accept for `86400`

##### _other_

- **Alt-Svc**: availability through other protocols `h3-25=":443"; ma=3600, h2=":443"; ma=3600`
- **Link**: alternate to `<link>` html tag `<https://one.example.com>; rel="preconnect", <https://two.example.com>; rel="preconnect"`
- **Server-Timing**: metrics `db;dur=53, app;dur=47.2`
- **Timing-Allow-Origin**: cross origin timings `https://one.example.com, https://two.example.com`
