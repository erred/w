package main

import (
	"fmt"
	"text/template"
)

func (o *options) parseTemplates() error {
	var err error
	o.templates = template.New("")
	for name, tmpl := range rawTemplates {
		o.templates, err = o.templates.New(name).Parse(tmpl)
		if err != nil {
			return fmt.Errorf("options.parseTemplates: %w", err)
		}
	}
	return nil
}

var (
	rawTemplates = map[string]string{
		"layout-main": ` {{ define "layout-main" }}
<!doctype html>
{{ if .AMP }}<html ⚡>{{ else }}<html lang="en">{{ end }}
<head>
        {{ template "head" . }}
</head>
<body>
        {{ template "body-amp" . }}
        <header>
                {{ template "header-name" }}
                {{ .Header }}
        </header>
        <main>
                {{ .Main }}
        </main>
        {{ template "footer" . }}
</body>
</html>
{{ end }}
`,

		"layout-blogpost": `{{ define "layout-blogpost" }}
<!doctype html>
{{ if .AMP }}<html ⚡>{{ else }}<html lang="en">{{ end }}
<head>
        {{ template "head" . }}
</head>
<body>
        {{ template "body-amp" . }}
        <header>
                {{ template "header-name" }}
                <h2><a href="/blog">b<em>log</em></a></h2>
                <p><time datetime="{{ .Date }}">{{ .Date }}</time><p>
        </header>
        <main>
                <h3>{{ .Title }}</h3>
                {{ .Main }}
        </main>
        {{ template "footer" . }}
</body>
</html>
{{ end }}
`,
		"layout-blogindex": `{{ define "layout-blogindex" }}
<!doctype html>
{{ if .AMP }}<html ⚡>{{ else }}<html lang="en">{{ end }}
<head>
        {{ template "head" . }}
</head>
<body>
        {{ template "body-amp" . }}
        <header>
                {{ template "header-name" }}
                <h2><a href="/blog">b<em>log</em></a></h2>
                <p>Artisanal, <em>hand-crafted</em> blog posts imbued with delayed <em>regrets</em></p>
        </header>
        <main>
                <ul>
                {{ range .Posts }}
                        <li><time datetime="{{ .Date }}">{{ .Date }}</time> | <a href="/{{ if .AMP }}amp/{{ end }}blog/{{ .URL }}">{{ .Title }}</a></li>
                {{ end }}
            </ul>
        </main>
        {{ template "footer" . }}
</body>
</html>
{{ end }}
`,
		"body-amp": `{{ define "body-amp" }}
{{ if .AMP }}
<amp-analytics type="gtag" data-credentials="include">
<script type="application/json">
{
  "vars": {
    "gtag_id": "{{ .GAID }}",
    "config": {
      "{{ .GAID }}": { "groups": "default" }
    }
  }
}
</script>
</amp-analytics>
{{ end }}
{{ end }}
`,

		// <link rel="preconnect" href="https://securetoken.googleapis.com" crossorigin>
		// <link rel="preconnect" href="https://api.seankhliao.com" crossorigin>
		// <script defer src="https://www.google.com/recaptcha/api.js?render=6LcYb6gUAAAAAAugirVvpmFxt8b-j3RfwcNDrchn"></script>

		// /* navigator.serviceWorker && window.addEventListener("load", () => {navigator.serviceWorker.register("./sw.js");}); */
		// window.addEventListener('load', ()=>{
		//     grecaptcha.ready(()=> {
		//         grecaptcha.execute('6LcYb6gUAAAAAAugirVvpmFxt8b-j3RfwcNDrchn', {
		//             action: 'pageview'
		//         }).then(token=> {
		//             fetch('https://api.seankhliao.com/recaptcha.v3', {
		//                 method: 'POST',
		//                 mode: 'cors',
		//                 body: token
		//             }).then(res =>{console.log(res)});
		//         });
		//     });
		// });

		"head": `{{ define "head" }}
<meta charset="utf-8">
<meta name="viewport" content="width=device-width,minimum-scale=1,initial-scale=1">

<title>{{ .Title }}</title>

{{ if not .AMP }}
<link rel="preconnect" href="https://www.gstatic.com" crossorigin>
<script defer src="https://www.googletagmanager.com/gtag/js?id=UA-114337586-1"></script>
{{ end }}

<link rel="canonical" href="{{ .URLCanonical }}" />
<link rel="amphtml" ref="{{ .URLAMP }}" />
<link rel="manifest" href="/manifest.json" />
<link rel="alternate" type="application/atom+xml" title="seankhliao.com - Atom Feed" href="https://seankhliao.com/feed.atom" />

<meta name="theme-color" content="#000000">
<meta name="description" content="{{ .Description }}">

<link rel="icon" type="image/png" sizes="512x512" href="/icon-512.png" />
<link rel="apple-touch-icon" href="/icon-512.png" />
<link rel="shortcut icon" href="/favicon.ico" />

<style {{ if .AMP }}amp-custom{{ end }}>
{{ template "css-fonts" }}
{{ template "css" }}
{{ .Style }}
</style>
{{ if .AMP }}
<script async src="https://cdn.ampproject.org/v0.js"></script>
<script async custom-element="amp-analytics" src="https://cdn.ampproject.org/v0/amp-analytics-0.1.js"></script>

<style amp-boilerplate>body{-webkit-animation:-amp-start 8s steps(1,end) 0s 1 normal both;-moz-animation:-amp-start 8s steps(1,end) 0s 1 normal both;-ms-animation:-amp-start 8s steps(1,end) 0s 1 normal both;animation:-amp-start 8s steps(1,end) 0s 1 normal both}@-webkit-keyframes -amp-start{from{visibility:hidden}to{visibility:visible}}@-moz-keyframes -amp-start{from{visibility:hidden}to{visibility:visible}}@-ms-keyframes -amp-start{from{visibility:hidden}to{visibility:visible}}@-o-keyframes -amp-start{from{visibility:hidden}to{visibility:visible}}@keyframes -amp-start{from{visibility:hidden}to{visibility:visible}}</style>
<noscript><style amp-boilerplate>body{-webkit-animation:none;-moz-animation:none;-ms-animation:none;animation:none}</style></noscript>
{{ else }}
<script>
"use strict";
window.dataLayer = window.dataLayer || [];
function gtag(){dataLayer.push(arguments);}
gtag("js", new Date());
gtag("config", "{{ .GAID }}");
</script>
{{ end }}
{{ end }}
`,
		"header-name": `{{ define "header-name" }}
<h1>
    <a href="/amp">
        <span>S</span>
        <span>E</span>
        <span>A</span>
        <span>N</span>
        <em>K</em>
        <em>.</em>
        <em>H</em>
        <em>.</em>
        <span>L</span>
        <span>I</span>
        <span>A</span>
        <span>O</span>
    </a>
</h1>
{{ end }}
`,
		"footer": `{{ define "footer" }}
<footer>
    <a href="{{ if .AMP }}/amp{{ end }}/">home</a>
    |
    <a href="{{ if .AMP }}/amp{{ end }}/privacy">privacy</a>
    |
    <a href="{{ if .AMP }}/amp{{ end }}/terms">terms</a>
    |
    <a href="https://github.com/seankhliao/com-seankhliao">github</a>
</footer>
{{ end }}
`,

		"css": `{{ define "css" }}
* {
  box-sizing: border-box;
}
:root {
  /*
  --black: #000000;
  --primary: #a06be0;
  --gray: #999;
  --white: #eceff1;
  --mono: "Inconsolata", monospace;
  --serif: "Lora", serif;
*/
  background: #000;
  color: #eceff1;
  font: 18px "Inconsolata", monospace;
}

/* ===== layout general ===== */
body {
  display: flex;
  flex-flow: column nowrap;
  margin: 0;
}

/* ===== layout header ===== */
header {
  display: grid;
  grid:
    [r1-s] ".    .    .    title title" 1fr [r1-e]
    [r1-s] ".    .    .    tag   tag" 1fr [r2-e]
    [r3-s] ".    .    .    .     ." 1fr [r3-e]
    [r4-s] "logo logo logo logo  logo" 1fr [r4-e]
    / 1fr 1fr 1fr 1fr 1fr;
  height: 80vh;
  padding: 2vmin;
}
header h1 {
  font: 700 5vmin "Lora", serif;
  grid-area: logo;
  margin: 0;
  place-self: stretch start;
}
header h1 a {
  display: grid;
  grid: repeat(2, 10vmin) / repeat(8, 10vmin);
  place-content: center center;
}
header h1 *:nth-child(n + 5) {
  grid-row: 2 / span 1;
}
header h2 {
  font-size: 4.5vmin;
  grid-area: title;
  margin: 0;
  place-self: end;
}
header p {
  color: #999;
  font-size: 3.5vmin;
  grid-area: tag;
  place-self: start end;
  text-align: right;
}

/* ===== layout main ===== */
main {
  display: grid;
  grid: auto-flow / 1fr minmax(280px, 840px) 1fr;
  grid-gap: 0 1em;
  margin: 20vh 0;
}
main > * {
  grid-column: 2 / span 1;
}
main > picture,
main > pre {
  grid-column: 1 / span 3;
}
picture img {
  width: 100%;
}

h3,
h4,
h5,
h6 {
  margin: 1.5em 0 0.25em 0;
}
h3 {
  font-size: 2.441em;
}
h4 {
  font-size: 1.953em;
}
h5 {
  font-size: 1.563em;
}
h6 {
  font-size: 1.25em;
}
p {
  line-height: 1.5;
  margin: 0 0 0.5em 0;
}

/* ===== layout footer ===== */
footer {
  margin: 10vh auto 3vh;
}

/* ===== style ===== */
a,
a:visited {
  color: inherit;
  font-weight: 700;
  text-decoration: underline;
}
a:hover {
  color: #a06be0;
  transition: color 0.16s;
}

header a,
header a:hover,
header a:visited {
  color: inherit;
  text-decoration: none;
}

ul {
  list-style: none;
  margin: 0;
}
ul > * {
  margin: 0.5em;
}
ul > li:before {
  content: "»";
  margin: 0 1ch 0 -3ch;
  position: absolute;
}

ol > * {
  line-height: 1.75em;
}

blockquote {
  margin: 1em;
  padding: 0.25em 1em;
  border-left: 1ch solid #999;
}

code {
  background: #404040;
  font: 1em "Inconsolata", monospace;
  padding: 0.25em 0.5em;
}
pre {
  background: #404040;
  overflow-x: scroll;
  padding: 1em;
}
pre::-webkit-scrollbar {
  display: none;
}
pre code {
  padding: 0;
}

em {
  color: #a06be0;
  background-color: unset;
  font-style: normal;
  font-weight: 700;
}
time {
  color: #999;
}

/* ===== hide recaptcha ===== */
.grecaptcha-badge {
  visibility: hidden;
}
{{ end }}
`,
		"css-fonts": `{{ define "css-fonts" }}
@font-face {
  font-family: 'Inconsolata';
  font-style: normal;
  font-weight: 400;
  font-display: swap;
  src: local('Inconsolata Regular'), local('Inconsolata-Regular'), url(https://fonts.gstatic.com/s/inconsolata/v18/QldKNThLqRwH-OJ1UHjlKGlZ5q0.ttf) format('truetype');
}
@font-face {
  font-family: 'Inconsolata';
  font-style: normal;
  font-weight: 700;
  font-display: swap;
  src: local('Inconsolata Bold'), local('Inconsolata-Bold'), url(https://fonts.gstatic.com/s/inconsolata/v18/QldXNThLqRwH-OJ1UHjlKGHiw71p5_k.ttf) format('truetype');
}
@font-face {
  font-family: 'Lora';
  font-style: normal;
  font-weight: 400;
  font-display: swap;
  src: local('Lora Regular'), local('Lora-Regular'), url(https://fonts.gstatic.com/s/lora/v14/0QIvMX1D_JOuMwr7Jg.ttf) format('truetype');
}
@font-face {
  font-family: 'Lora';
  font-style: normal;
  font-weight: 700;
  font-display: swap;
  src: local('Lora Bold'), local('Lora-Bold'), url(https://fonts.gstatic.com/s/lora/v14/0QIgMX1D_JOuO7HeNtxunw.ttf) format('truetype');
}
{{ end }}
`,
		"css-loader": `{{ define "css-loader" }}
.loader {
        display: grid;
        grid-gap: 2px;
        grid-template: repeat(3, 1fr) / repeat(3, 1fr);
        height: 48px;
        margin: 5em auto;
        place-items: stretch;
        width: 48px;
}
.loader div {
        background: #333;
        animation: loader 1.8s infinite ease-in-out;
}
.loader div:nth-child(5) {
        background: #a06be0;
}
.loader div:nth-child(7) {
        animation-delay: 0s;
}
.loader div:nth-child(4),
.loader div:nth-child(8) {
        animation-delay: 0.3s;
}
.loader div:nth-child(1),
.loader div:nth-child(5),
.loader div:nth-child(9) {
        animation-delay: 0.6s;
}
.loader div:nth-child(2),
.loader div:nth-child(6) {
        animation-delay: 0.9s;
}
.loader div:nth-child(3) {
        animation-delay: 1.2s;
}
@keyframes loader {
        0%,
        70%,
        100% {
                transform: scale3D(1, 1, 1);
        }
        35% {
                transform: scale3D(0, 0, 1);
        }
}
{{ end }}
`,
		"loader": `{{ define "loader" }}
<div class="loader">
  <div></div>
  <div></div>
  <div></div>
  <div></div>
  <div></div>
  <div></div>
  <div></div>
  <div></div>
  <div></div>
</div>
{{ end }}
`,

		"layout-gomod": `{{ define "layout-gomod" }}
<!doctype html>
<html lang="en">
<head>
        <meta name="go-import" content="github.com/{{ .Old }} git {{ .New }}">
        <meta name="go-source" content="{{ .New }} https://github.com/{{ .Old }} https://github.com/{{ .Old }}/tree/master{/dir} https://github.com/{{ .Old }}/tree/master{/dir}/{file}#L{line}">
        <meta http-equiv="Refresh" content="5; url='godoc.org/{{ .New }}'">

        <link rel="preconnect" href="https://www.gstatic.com" crossorigin>
        <script defer src="https://www.googletagmanager.com/gtag/js?id=UA-114337586-1"></script>
        <script>
"use strict";
window.dataLayer = window.dataLayer || [];
function gtag(){dataLayer.push(arguments);}
gtag("js", new Date());
gtag("config", "UA-114337586-1");
        </script>
    </head>
</html>
{{ end }}
`,
	}
)
