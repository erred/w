package render

import (
	"bytes"
	_ "embed"
	"fmt"
	"io"
	"regexp"
	"sync"
	"text/template"

	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
	mhtml "github.com/tdewolff/minify/v2/html"
	"github.com/tdewolff/minify/v2/js"
	"github.com/tdewolff/minify/v2/svg"
	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"go.seankhliao.com/w/v16/picture"
	"sigs.k8s.io/yaml"
)

var (

	//go:embed layout.tpl
	layoutTpl  string
	layoutName = "layout"
	T          = template.Must(template.New(layoutName).Parse(layoutTpl))
)

var defaultMarkdown = goldmark.New(
	goldmark.WithExtensions(extension.Table, meta.Meta, picture.Picture),
	goldmark.WithParserOptions(parser.WithAutoHeadingID()),
	goldmark.WithRendererOptions(html.WithUnsafe()),
)

var (
	defaultMinifyOnce sync.Once
	defaultMinifyM    *minify.M
	defaultMinify     = func() *minify.M {
		defaultMinifyOnce.Do(func() {
			m := minify.New()
			m.AddFunc("text/html", mhtml.Minify)
			m.AddFunc("text/css", css.Minify)
			m.AddFunc("image/svg+xml", svg.Minify)
			m.AddFuncRegexp(regexp.MustCompile("^(application|text)/(x-)?(java|ecma)script$"), js.Minify)
			// m.AddFuncRegexp(regexp.MustCompile("[/+]json$"), json.Minify)
			// m.AddFuncRegexp(regexp.MustCompile("[/+]xml$"), xml.Minify)
			defaultMinifyM = m
		})
		return defaultMinifyM
	}
)

type Options struct {
	Markdown     goldmark.Markdown
	MarkdownSkip bool // skip markdown processing

	Minify       *minify.M
	Template     *template.Template
	TemplateName string

	Data PageData
}

func (o *Options) init() {
	if o.Markdown == nil {
		o.Markdown = defaultMarkdown
	}
	if o.Minify == nil {
		o.Minify = defaultMinify()
	}
	if o.Template == nil {
		o.Template = T
	}
	if o.TemplateName == "" {
		o.TemplateName = layoutName
	}
}

type PageData struct {
	// mandatory
	URLCanonical string
	Date         string // for blog posts
	GTMID        string // for analytics
	Compact      bool

	// Extracted
	Title       string
	Description string
	H1          string // default to Title
	H2          string // default to Description
	Style       string

	// Filled
	Main string
}

func (d *PageData) FromMap(md map[string]interface{}) {
	d.Title, _ = md["title"].(string)
	d.Description, _ = md["description"].(string)
	d.H1, _ = md["h1"].(string)
	d.H2, _ = md["h2"].(string)
	d.Style, _ = md["style"].(string)

	if d.H1 == "" {
		d.H1 = d.Title
	}
	if d.H2 == "" {
		d.H2 = d.Description
	}
}

func Render(o *Options, w io.Writer, r io.Reader) error {
	o.init()

	b, err := io.ReadAll(r)
	if err != nil {
		return fmt.Errorf("markdown: read: %w", err)
	}

	buf1 := &bytes.Buffer{}
	if o.MarkdownSkip {
		n, md, err := extractHeader(b)
		if err != nil {
			return fmt.Errorf("extract header: %w", err)
		}
		buf1.Write(b[n:])
		o.Data.FromMap(md)
	} else {
		mdCtx := parser.NewContext()
		err = o.Markdown.Convert(b, buf1, parser.WithContext(mdCtx))
		if err != nil {
			return fmt.Errorf("markdown: convert: %w", err)
		}
		o.Data.FromMap(meta.Get(mdCtx))
	}
	buf1.WriteRune('\n')
	buf1.WriteString(o.Data.Main)
	o.Data.Main = buf1.String()

	buf2 := &bytes.Buffer{}
	err = o.Template.ExecuteTemplate(buf2, o.TemplateName, o.Data)
	if err != nil {
		return fmt.Errorf("markdown: template: %w", err)
	}

	err = o.Minify.Minify("text/html", w, buf2)
	if err != nil {
		return fmt.Errorf("markdown: minify: %w", err)
	}
	return nil
}

func extractHeader(b []byte) (int, map[string]interface{}, error) {
	if !bytes.HasPrefix(b, []byte("---\n")) {
		return 0, nil, nil
	}
	i := bytes.Index(b[4:], []byte("---\n"))
	if i == -1 {
		return 0, nil, fmt.Errorf("no metadata ending")
	}
	var m map[string]interface{}
	err := yaml.Unmarshal(b[4:i], &m)
	if err != nil {
		return 0, nil, fmt.Errorf("unmarshal metadata: %w", err)
	}
	return i + 4, m, nil
}
