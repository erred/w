<!DOCTYPE html>
<html lang="en">
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width,minimum-scale=1,initial-scale=1">
  <title>{{ .Title }}</title>

  {{- if not .DisableAnalytics }}
  <script>
    (function (w, d, s, l, i) {
      w[l] = w[l] || []; w[l].push({ "gtm.start": new Date().getTime(), event: "gtm.js" });
      var f = d.getElementsByTagName(s)[0], j = d.createElement(s), dl = l != "dataLayer" ? "&l=" + l : "";
      j.async = true; j.src = "https://www.googletagmanager.com/gtm.js?id=" + i + dl;
      f.parentNode.insertBefore(j, f);
    })(window, document, "script", "dataLayer", "GTM-TLVN7D6");
  </script>
  {{- end }}

  {{- if not .EmbedStyle }}
  <link rel="stylesheet" href="/base.css" crossorigin>
  {{- end }}

  <link rel="canonical" href="{{ .URLCanonical }}">
  <link rel="manifest" href="/manifest.json">

  <meta name="theme-color" content="#000000">
  <meta name="description" content="{{ .Description }}">

  <link rel="icon" href="https://seankhliao.com/favicon.ico">
  <link rel="icon" href="https://seankhliao.com/static/icon.svg" type="image/svg+xml" sizes="any">
  <link rel="apple-touch-icon" href="https://seankhliao.com/static/icon-192.png">

  <style>
    {{- if .EmbedStyle }}
    {{ template "base.css" }}
    {{- end }}
    {{ .Style }}
  </style>

  {{- if not .DisableAnalytics }}
  <noscript><iframe src="https://www.googletagmanager.com/ns.html?id=GTM-TLVN7D6" height="0" width="0" style="display: none; visibility: hidden"></iframe></noscript>
  {{- end }}

  <h1>{{ .H1 }}</h1>
  <h2>{{ .H2 }}</h2>

  <hgroup>
    <a href="/">
      <span>S</span><span>E</span><span>A</span><span>N</span>
      <em>K</em><em>.</em><em>H</em><em>.</em>
      <span>L</span><span>I</span><span>A</span><span>O</span>
    </a>
  </hgroup>

  {{ .Main }}

  <footer>
    <a href="https://seankhliao.com/">home</a>
    |
    <a href="https://seankhliao.com/blog/">blog</a>
    |
    <a href="https://github.com/seankhliao">github</a>
  </footer>
</html>
