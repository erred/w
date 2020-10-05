package render

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"text/template"

	"go.seankhliao.com/com-seankhliao/v12/render/style"
)

type Options struct {
	In        string
	Out       string
	Template  *template.Template
	URLBase   string
	URLLogger string

	Analytics  bool
	EmbedStyle bool
}

func (o *Options) InitFlags(fs *flag.FlagSet) {
	fs.StringVar(&o.In, "in", "src", "input directory")
	fs.StringVar(&o.Out, "out", "public", "output directory")
	fs.StringVar(&o.URLBase, "base", "https://seankhliao.com", "base url")
	fs.StringVar(&o.URLLogger, "logger", "https://statslogger.seankhliao.com/beacon", "statslogger url")
	fs.BoolVar(&o.Analytics, "analytics", true, "include analytics")
	fs.BoolVar(&o.EmbedStyle, "embedstyle", false, "use inlined css")
}

func Process(o Options) error {
	pages, err := processInput(o)
	if err != nil {
		return fmt.Errorf("process: %w", err)
	}
	if len(pages) > 1 {
		pages, err = processFill(pages, o.Out)
		if err != nil {
			return fmt.Errorf("process: %w", err)
		}
	}
	err = processOutput(o, pages)
	if err != nil {
		return fmt.Errorf("process: %w", err)
	}
	return nil
}

func processInput(o Options) ([]*Page, error) {
	info, err := os.Stat(o.In)
	if err != nil {
		return nil, fmt.Errorf("stat file %s: %w", o.In, err)
	}
	var pages []*Page
	if !info.IsDir() {
		page, err := NewPageFromFile(o.In)
		if err != nil {
			return nil, fmt.Errorf("process input file %s: %w", o.In, err)
		}
		pages = []*Page{page}
	} else {
		err = filepath.Walk(o.In, walker(o.In, &pages))
		if err != nil {
			return nil, fmt.Errorf("process input dir %s: %w", o.In, err)
		}
	}
	for i := range pages {
		pages[i].URLBase = o.URLBase
		pages[i].URLLogger = o.URLLogger
		pages[i].URLAbsolute = canonical(strings.TrimPrefix(pages[i].name, o.In))
		pages[i].URLCanonical = o.URLBase + pages[i].URLAbsolute
		pages[i].Analytics = o.Analytics
		if pages[i].name != o.In {
			r, err := filepath.Rel(o.In, pages[i].name)
			if err == nil {
				pages[i].name = filepath.Join(o.Out, r)
			}
		}
	}
	return pages, nil
}

func processFill(pages []*Page, out string) ([]*Page, error) {
	sort.Slice(pages, func(i, j int) bool { return pages[i].name > pages[j].name })

	blogindex, buf := 0, strings.Builder{}
	buf.WriteString("<ul>\n")
	for i, p := range pages {
		if strings.Contains(p.name, "/blog/") {
			if filepath.Base(p.name) != "index.html" {
				pages[i].Date = filepath.Base(p.name)[:10]
				pages[i].Header = blogHeader(pages[i].Date)
				buf.WriteString(blogLink(p.Date, p.URLAbsolute, p.Title))
			} else {
				blogindex = i
			}
		}
		if filepath.Ext(p.name) == ".html" {
			pages[i].Main = imgHack(p.Main)
		}
	}

	buf.WriteString("</ul>\n")
	pages[blogindex].Main += "\n" + buf.String()
	pages[blogindex].Header = blogIndexHeader()

	// create sitemap
	all := make([][]byte, len(pages))
	for i := range all {
		all[i] = []byte(pages[i].URLCanonical)
		// + "?utm_source=sitemap&utm_medium=txt&utm_campaign=sitemap.txt"
	}
	p, err := NewPage(filepath.Join(out, "sitemap.txt"), bytes.Join(all, []byte("\n")), true)
	if err != nil {
		return nil, fmt.Errorf("fill sitemap.txt: %w", err)
	}
	pages = append(pages, p)

	return pages, nil
}

func processOutput(o Options, pages []*Page) error {
	for _, p := range pages {
		os.MkdirAll(filepath.Dir(p.name), 0o755)
		f, err := os.Create(p.name)
		if err != nil {
			return fmt.Errorf("create file %s: %w", p.name, err)
		}
		if p.pass {
			_, err = f.Write(p.data)
		} else {
			if o.EmbedStyle {
				p.Style = style.StyleGohtml + "\n" + p.Style
			}

			err = o.Template.ExecuteTemplate(f, "LayoutGohtml", p)
		}
		f.Close()
		if err != nil {
			return fmt.Errorf("write file %s: %w", p.name, err)
		}
	}
	return nil
}

func walker(base string, pt *[]*Page) func(p string, i os.FileInfo, err error) error {
	pages := *pt
	return func(p string, i os.FileInfo, err error) error {
		if err != nil {
			return err
		} else if i.IsDir() {
			return nil
		}
		page, err := NewPageFromFile(p)
		if err != nil {
			return fmt.Errorf("walk %s: %w", p, err)
		}
		pages = append(pages, page)
		*pt = pages
		return nil
	}
}

func canonical(p string) string {
	if p[0] != '/' {
		p = "/" + p
	}
	if strings.HasSuffix(p, ".html") {
		p = strings.TrimSuffix(strings.TrimSuffix(p, ".html"), "index")
		if p == "" {
			p = "/"
		}
		if p[len(p)-1] != '/' {
			p = p + "/"
		}
	}
	return p
}

func blogIndexHeader() string {
	return `<h2><a href="/blog/">b<em>log</em></a></h2>
<p>Artisanal, <em>hand-crafted</em> blog posts imbued with delayed <em>regrets</em></p>`
}

func blogHeader(date string) string {
	return fmt.Sprintf(`<h2><a href="/blog/">b<em>log</em></a></h2>
<p><time datetime="%s">%s</time></p>`, date, date)
}

func blogLink(date, urlabsolute, title string) string {
	return fmt.Sprintf(`<li><time datetime="%s">%s</time> | <a href="%s">%s</a></li>`+"\n",
		date, date, urlabsolute, title)
}

func imgHack(html string) string {
	r := regexp.MustCompile(`<h4><img src="(.*?).webp" alt="(.*?)"></h4>`)
	return r.ReplaceAllString(html, `
<picture>
        <source type="image/webp" srcset="$1.webp">
        <source type="image/jpeg" srcset="$1.jpg">
        <img src="$1.png" alt="$2">
</picture>
`)
}
