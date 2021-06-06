<!DOCTYPE html>
<html lang="en">
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width,minimum-scale=1,initial-scale=1">
  <title>{{ .Title }}</title>

  {{- if ne .GTMID "" }}
  <script>
    (function (w, d, s, l, i) {
      w[l] = w[l] || []; w[l].push({ "gtm.start": new Date().getTime(), event: "gtm.js" });
      var f = d.getElementsByTagName(s)[0], j = d.createElement(s), dl = l != "dataLayer" ? "&l=" + l : "";
      j.async = true; j.src = "https://www.googletagmanager.com/gtm.js?id=" + i + dl;
      f.parentNode.insertBefore(j, f);
    })(window, document, "script", "dataLayer", "{{ .GTMID }}");
  </script>
  {{- end }}

  <link rel="canonical" href="{{ .URLCanonical }}">
  <link rel="manifest" href="/manifest.json">

  <meta name="theme-color" content="#000000">
  <meta name="description" content="{{ .Description }}">

  <link rel="icon" href="https://seankhliao.com/favicon.ico">
  <link rel="icon" href="https://seankhliao.com/static/icon.svg" type="image/svg+xml" sizes="any">
  <link rel="apple-touch-icon" href="https://seankhliao.com/static/icon-192.png">

  <style>
* {
  box-sizing: border-box;
}
:root {
  background: #000;
  color: #eceff1;
  font: 18px "Inconsolata", monospace;
}

@font-face {
  font-family: "Inconsolata";
  font-style: normal;
  font-weight: 400;
  font-display: swap;
  src: local("Inconsolata"), local("Inconsolata-Regular"),
    url(https://seankhliao.com/static/inconsolata-var.woff2) format("woff2-variations"),
    url(https://seankhliao.com/static/inconsolata-400.woff2) format("woff2");
}
@font-face {
  font-family: "Inconsolata";
  font-style: normal;
  font-weight: 700;
  font-display: swap;
  src: local("Inconsolata Bold"), local("Inconsolata-Bold"),
    url(https://seankhliao.com/static/inconsolata-var.woff2) format("woff2-variations"),
    url(https://seankhliao.com/static/inconsolata-700.woff2) format("woff2");
}
@font-face {
  font-family: "Lora";
  font-style: normal;
  font-weight: 400;
  font-display: swap;
  src: local("Lora"), local("Lora-Regular"),
    url(https://seankhliao.com/static/lora-var.woff2) format("woff2-variations"),
    url(https://seankhliao.com/static/lora-400.woff2) format("woff2");
}
@font-face {
  font-family: "Lora";
  font-style: normal;
  font-weight: 700;
  font-display: swap;
  src: local("Lora Bold"), local("Lora-Bold"),
    url(https://seankhliao.com/static/lora-var.woff2) format("woff2-variations"),
    url(https://seankhliao.com/static/lora-700.woff2) format("woff2");
}

/* ===== layout general ===== */
body {
  {{ if .Compact }}
  grid: 3.5vh 3.5vh / 1fr repeat(3, minmax(90px, 280px)) 1fr;
  {{ else }}
  grid: 20vh 60vh / 1fr repeat(3, minmax(90px, 280px)) 1fr;
  {{ end }}
  display: grid;
  gap: 0 1em;
  margin: 0;
  padding: 1vmin;

  /* ==override newtab page == */
  background: #000;
  color: #eceff1;
  font: 18px "Inconsolata", monospace;
}

body > * {
  grid-column: 2 / span 3;
}

/* ===== layout header ===== */
h1 {
  {{ if .Compact }}
  font-size: 3vmin;
  grid-area: 1 / 4 / span 1 / span 1;
  {{ else }}
  font-size: 4.5vmin;
  grid-area: 1 / 4 / span 1 / span 2;
  {{ end }}
  margin: 0;
  place-self: end;
}
h2 {
  color: #999;
  {{ if .Compact }}
  font-size: 2.5vmin;
  grid-area: 2 / 4 / span 1 / span 1;
  {{ else }}
  font-size: 3.5vmin;
  grid-area: 2 / 4 / span 1 / span 2;
  {{ end }}
  place-self: start end;
  text-align: right;
}

hgroup {
  {{ if .Compact }}
  font: 700 2.5vmin "Lora", serif;
  grid-area: 1 / 2 / span 2 / span 1;
  {{ else }}
  font: 700 5vmin "Lora", serif;
  grid-area: 1 / 1 / span 2 / span 2;
  {{ end }}
  margin: 0;
  place-self: end start;
}
hgroup a {
  display: grid;
  {{ if .Compact }}
  grid: repeat(2, 3vmin) / repeat(8, 4vmin);
  {{ else }}
  grid: repeat(2, 10vmin) / repeat(8, 10vmin);
  {{ end }}
  place-content: center center;
}
hgroup *:nth-child(n + 5) {
  grid-row: 2 / span 1;
}
/* ===== full bleed ===== */
footer,
iframe,
pre,
table,
picture {
  grid-column: 1 / span 5;
  margin: 0.25em -1vmin 2em;
}
picture img {
  width: 100%;
  margin: auto;
}

/* ===== layout main ===== */
h3,
h4,
picture {
  {{ if .Compact }}
  margin: 1em 0 0.25em 0;
  {{ else }}
  margin: 25vh 0 0.25em 0;
  {{ end }}
}

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
  margin: 0 0 1em 0;
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
  text-decoration: underline 1px #707070;
}
a:hover {
  color: #a06be0;
  transition: color 0.16s;
  text-decoration: underline 1px #a06be0;
}

h1 a,
h1 a:hover,
h1 a:visited,
hgroup a,
hgroup a:hover,
hgroup a:visited {
  color: inherit;
  text-decoration: none;
}

ul {
  list-style: none;
  margin: 0;
}
ul > * {
  margin: 0.5em;
  line-height: 1.5em;
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
  background: #303030;
  font: 1em "Inconsolata", monospace;
  padding: 0.1em;
}
pre {
  background: #303030;
  overflow-x: scroll;
  padding: 1em;
}
pre::-webkit-scrollbar {
  display: none;
}
pre code {
  padding: 0;
}

iframe {
  margin: auto;
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

table {
  border-collapse: collapse;
  border-style: hidden;
}
th,
td {
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

/* ===== gtm ===== */
noscript iframe {
  height: 0;
  width: 0;
  display: none;
  visibility: hidden;
}
    {{ .Style }}
  </style>

  {{- if ne .GTMID "" }}
  <noscript><iframe src="https://www.googletagmanager.com/ns.html?id={{ .GTMID }}" height="0" width="0" style="display: none; visibility: hidden"></iframe></noscript>
  {{- end }}

  <h1>{{ .H1 }}</h1>
  <h2>{{ .H2 }}</h2>

  <hgroup>
    {{ if .Compact }}
    <a href="https://seankhliao.com/?utm_medium=sites&utm_source={{.URLCanonical}}">
    {{ else }}
    <a href="/">
    {{ end }}
      <span>S</span><span>E</span><span>A</span><span>N</span>
      <em>K</em><em>.</em><em>H</em><em>.</em>
      <span>L</span><span>I</span><span>A</span><span>O</span>
    </a>
  </hgroup>

  {{ .Main }}

  <footer>
    <a href="https://seankhliao.com/{{ if .Compact }}?utm_medium=sites&utm_source={{.URLCanonical}}{{ end }}">home</a>
    |
    <a href="https://seankhliao.com/blog/{{ if .Compact }}?utm_medium=sites&utm_source={{.URLCanonical}}{{ end }}">blog</a>
    |
    <a href="https://github.com/seankhliao">github</a>
  </footer>
</html>
