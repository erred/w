package main

const HeadTemplate = `
<!doctype html>
<html lang="en">
<meta charset="utf-8">
<meta name="viewport" content="width=device-width,minimum-scale=1,initial-scale=1">
<meta name="theme-color" content="#000000">
<script async src="https://www.googletagmanager.com/gtag/js?id=UA-114337586-4"></script>
<script>window.dataLayer = window.dataLayer || [];function gtag() {dataLayer.push(arguments);};gtag("js", new Date());gtag("config", "UA-114337586-4");</script>
<link rel="stylesheet" href="https://fonts.googleapis.com/css?family=Inconsolata:400,700&display=swap">
<link rel="stylesheet" href="./base.css">

<link rel="icon" type="image/png" sizes="512x512" href="/icon-512.png" />
<link rel="apple-touch-icon" href="/icon-512.png" />
<link rel="shortcut icon" href="/favicon.ico" />

<title>{{ with .Title }}{{ . }} |{{ end }}blog | seankhliao</title>

<link rel="canonical" href="https://blog.seankhliao.com{{ with .URL }}/{{ . }}{{ end }}">
<meta name="description" content="{{ .Description }}">
`
const FootTemplate = `
</html>
`

const PostTemplate = `
{{ template "head" . }}

<hgroup>
  <h1>{{ .Title }}</h1>
  <p>
    <a href="https://seankhliao.com">seankhliao</a> / 
    <a href="https://blog.seankhliao.com">blog</a> / 
    <time datetime="{{ .Date }}">{{ .Date }}</time>
  </p>
</hgroup>

{{ .Content }}

{{ template "foot" . }}
`

const IndexTemplate = `
{{ template "head" . }} 
  
<h1><a href="https://seankhliao.com">seankhliao</a> / blog</h1>
<p>Artisanal, hand-crafted blog posts imbued with delayed regrets</p>
  
<ul>
{{ range .Posts }}
    <li><time datetime="{{ .Date }}">{{ .Date }}</time> | <a href="./{{ .URL }}">{{ .Title }}</a></li>
{{ end }}
</ul>

{{ template "foot" . }}
`
