package process

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"go.seankhliao.com/w/v16/render"
)

type Options struct {
	GTMID     string
	Canonical string
	Raw       bool
	Compact   bool
}

func Dir(o Options, dst, src string) error {
	var pis []pageInfo
	err := filepath.WalkDir(src, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, p)
		if err != nil {
			return fmt.Errorf("not relative: %w", err)
		}

		// create directories
		if d.IsDir() {
			dstf := filepath.Join(dst, rel)
			err = os.MkdirAll(dstf, 0o755)
			if err != nil {
				return fmt.Errorf("create dir %s: %w", dstf, err)
			}
			return nil
		}

		// process file
		urlcanonical, dstf := pageName(o.Canonical, rel)
		o2 := o
		o2.Canonical = urlcanonical
		// o2.Compact = strings.Contains(rel, "/") // big header only for root entries
		dstf = filepath.Join(dst, dstf)
		pi, err := process(o2, dstf, p)
		if err != nil {
			return fmt.Errorf("process %s: %w", p, err)
		}
		pis = append(pis, pi)
		return nil
	})
	if err != nil {
		return fmt.Errorf("walk %s: %w", src, err)
	}

	pi, err := blogIndex(o, filepath.Join(dst, "/blog/index.html"), pis)
	if err != nil {
		return fmt.Errorf("blog index: %w", err)
	}
	pis = append(pis, pi)

	sm := sitemap(pis)
	err = os.WriteFile(filepath.Join(dst, "sitemap.txt"), sm, 0o644)
	if err != nil {
		return fmt.Errorf("write sitemap: %w", err)
	}
	return nil
}

type pageInfo struct {
	URLCanonical string
	Title        string
	Date         string
}

var dateRe = regexp.MustCompile(`^\d{5}-\d{2}-\d{2}`)

func process(o Options, dst, src string) (pageInfo, error) {
	fin, err := os.Open(src)
	if err != nil {
		return pageInfo{}, fmt.Errorf("open %s: %w", src, err)
	}
	defer fin.Close()
	fout, err := os.Create(dst)
	if err != nil {
		return pageInfo{}, fmt.Errorf("create %s: %w", dst, err)
	}
	defer fout.Close()

	date := dateRe.FindString(filepath.Base(src))
	var h1 string
	if date != "" {
		h1 = `<a href="/blog/">b<em>log</em></a>`
	}

	ro := &render.Options{
		MarkdownSkip: o.Raw,
		Data: render.PageData{
			URLCanonical: o.Canonical,
			GTMID:        o.GTMID,
			Compact:      o.Compact,
			Date:         date,
			H1:           h1,
			H2:           date,
		},
	}

	err = render.Render(ro, fout, fin)
	if err != nil {
		return pageInfo{}, fmt.Errorf("render %s: %w", src, err)
	}
	return pageInfo{
		URLCanonical: o.Canonical,
		Title:        ro.Data.Title,
		Date:         ro.Data.Date,
	}, nil
}

func File(o Options, dst, src string) error {
	_, err := process(o, dst, src)
	return err
}

func blogIndex(o Options, dst string, pis []pageInfo) (pageInfo, error) {
	sort.Slice(pis, func(i, j int) bool {
		return pis[i].URLCanonical > pis[j].URLCanonical
	})

	buf := &bytes.Buffer{}
	buf.WriteString(`
<h3><em>B</em>log</h3>
<p>we<em>b log</em> of things that never made sense,
maybe someone will find this useful</p>
<ul>`)
	for _, pi := range pis {
		if !strings.Contains(pi.URLCanonical, "/blog/") {
			continue
		}
		fmt.Fprintf(buf, `<li><time datetime="%s">%s</time> | <a href="%s">%s</a></li>`+"\n",
			pi.Date[1:], pi.Date, strings.TrimPrefix(pi.URLCanonical, o.Canonical), pi.Title,
		)
	}
	buf.WriteRune('\n')

	ro := &render.Options{
		MarkdownSkip: true,
		Data: render.PageData{
			URLCanonical: o.Canonical + "/blog/",
			GTMID:        o.GTMID,
			Title:        "blog | seankhliao",
			Description:  "list of things i wrote",
			H1:           `<a href="/blog/">b<em>log</em></a>`,
			H2: `Artisanal, <em>hand-crafted</em> blog posts
imbued with delayed <em>regrets</em>`,
			Style: `
ul li {
        white-space: nowrap;
        overflow: hidden;
        text-overflow: ellipsis;
}`,
		},
	}

	fout, err := os.Create(dst)
	if err != nil {
		return pageInfo{}, fmt.Errorf("create %s: %w", dst, err)
	}

	err = render.Render(ro, fout, buf)
	if err != nil {
		return pageInfo{}, fmt.Errorf("render %s: %w", dst, err)
	}
	return pageInfo{}, nil
}

func sitemap(pis []pageInfo) []byte {
	buf := &bytes.Buffer{}
	for _, pi := range pis {
		buf.WriteString(pi.URLCanonical)
		buf.WriteRune('\n')
	}
	return buf.Bytes()
}

func pageName(urlBase, rel string) (canonicalURL, filename string) {
	urlBase = strings.TrimSuffix(urlBase, "/")
	rel = strings.TrimPrefix(rel, "/")
	rel = strings.TrimSuffix(rel, ".md")
	if strings.HasSuffix(rel, "index") {
		canonicalURL = strings.TrimSuffix(rel, "index")
	} else {
		canonicalURL = rel + "/"
	}
	filename = rel + ".html"
	return urlBase + "/" + canonicalURL, filename
}
