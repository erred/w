package sitedata

import (
	"embed"
	"text/template"
)

// S contains the static files
//
//go:embed public/*
var S embed.FS

//go:embed templates/*
var t embed.FS

// T contains the parsed templates
var T = template.Must(template.Must(template.ParseFS(t, "templates/*")).ParseFS(S, "public/base.css"))
