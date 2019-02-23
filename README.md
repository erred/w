# com-seankhliao

handcrafted personal static site

[![License](https://img.shields.io/github/license/seankhliao/com-seankhliao.svg?style=for-the-badge&maxAge=31536000)](LICENSE)
[![GoDoc](https://img.shields.io/badge/godoc-reference-5272B4.svg?style=for-the-badge&maxAge=31536000)](https://godoc.org/github.com/seankhliao/com-seankhliao)
[![Build](https://badger.seankhliao.com/i/github_seankhliao_com-seankhliao)](https://badger.seankhliao.com/l/github_seankhliao_com-seankhliao)

## About

who am i?

## Todo

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
