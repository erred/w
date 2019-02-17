# com-seankhliao

[![Build](https://img.shields.io/badge/endpoint.svg?url=https://badger.seankhliao.com/r/github_seankhliao_com-seankhliao)](https://console.cloud.google.com/cloud-build/builds?project=com-seankhliao&query=source.repo_source.repo_name%20%3D%20%22github_seankhliao_com-seankhliao%22)
[![License](https://img.shields.io/github/license/seankhliao/com-seankhliao.svg?style=for-the-badge)](LICENSE)

handcrafted personal static site

## TODO

### CI / CD

- [ ] Use a static site generator
- [ ] cache site build
- [x] convert images
- [x] deploy to firebase
- [x] purge cloudflare cache

### Compression

Setup compression / minification

- js / css / html minification
- HTTP gzip / brotli compression

### Content Security Policy

Enforce content verification

- Try with Content-Security-Policy-Report-Only
- integrate into build process (meta) / server (header)
- Setup report collection infra

### Feature Policy

Experimental, enforce feature (non)-usage

### Head tags

- `<meta charset="">`
- `<meta name="viewport" content="">`
- `<link rel="preload" href="" as="">`
  - preconnect: setup tcp connection
  - preload: mandatory high priority content
  - prefetch: resource used lated
  - prerender: next nav target
- `<title></title>`
- `<meta name="description" content="">`
- `<meta name="theme-color" content="">`
- `<link rel="canonical" href="">`
- `<link rel="manifest" href="">`
- `<link rel="icon" href="">`

### Search Engine Optimization

- Decide on domain / subdomain
- No ES6, googlebot
- Sitemaps / autosubmit

#### AMP

decide

#### Open Graph

for Facebook

#### JSON-LD

for Google

- Article
- Breadcrumb
- Logo
- Q&A page
- Social Profile
