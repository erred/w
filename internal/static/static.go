package static

import (
	"embed"
	"io/fs"
	"text/template"
)

//go:embed root/*
var s embed.FS

var S, _ = fs.Sub(s, "root")

//go:embed layout.tpl
var l string

// T contains the parsed templates
var T = template.Must(template.New("layout").Parse(l))
