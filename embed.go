package sitedata

import (
	"embed"
	"text/template"
)

//go:embed public/*
var S embed.FS

//go:embed templates/*
var t embed.FS

var T = template.Must(template.Must(template.ParseFS(t, "templates/*")).ParseFS(S, "public/base.css"))
