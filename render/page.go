package render

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"sigs.k8s.io/yaml"
)

type Page struct {
	// passthrough, don't process
	pass bool
	name string
	data []byte

	// Optional yaml front matter,
	// also supplied to ExecuteTemplate
	Date        string
	Description string
	Header      string
	Style       string
	Title       string

	// filled, supplied to ExecuteTemplate
	Main         string // html content
	URLAbsolute  string // start from /
	URLBase      string // https://... without trailing /
	URLCanonical string // URLBase + URLAbsolute
	URLLogger    string

	Analytics bool
}

func NewPageFromFile(fpath string) (*Page, error) {
	b, err := ioutil.ReadFile(fpath)
	if err != nil {
		return nil, fmt.Errorf("read file %s: %w", fpath, err)
	}
	var pass bool
	if filepath.Ext(fpath) != ".md" {
		pass = true
	}
	return NewPage(fpath, b, pass)
}

// NewPage reates a new page from filename and file contents
// if data starts with `---` it is assumed to be followed by
// a yaml front matter then markdown
func NewPage(name string, data []byte, pass bool) (*Page, error) {
	p := Page{
		pass: pass,
		data: data,
		name: name,
	}
	if !pass && bytes.HasPrefix(data, []byte(`---`)) {
		parts := bytes.SplitN(p.data, []byte(`---`), 3)
		err := yaml.Unmarshal(parts[1], &p)
		if err != nil {
			return nil, err
		}
		p.data = parts[2]
	}
	if filepath.Ext(p.name) == ".md" {
		p.name = strings.TrimSuffix(p.name, "md") + "html"

		var buf bytes.Buffer
		err := goldmark.New(goldmark.WithExtensions(extension.Table)).Convert(p.data, &buf)
		if err != nil {
			return nil, err
		}
		p.Main = buf.String()
	}
	return &p, nil
}
