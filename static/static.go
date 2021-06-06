package static

import (
	"embed"
	"io/fs"
)

var (
	//go:embed root/*
	s    embed.FS
	S, _ = fs.Sub(s, "root")
)
