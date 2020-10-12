package render

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"go.seankhliao.com/com-seankhliao/v13/internal/style"
	"sigs.k8s.io/yaml"
)

var (
	mdParser = goldmark.New(goldmark.WithExtensions(extension.Table))
)

type PageData struct {
	// mandatory
	Title        string
	Description  string
	URLCanonical string
	Main         string

	// optional
	H1               string
	H2               string
	Style            string
	Date             string
	DisableAnalytics bool
	EmbedStyle       bool
	BlogPost         bool
}

type PageInfo struct {
	URLCanonical string

	BlogPost bool
	Title    string
	Date     string
}

func ProcessDir(src, dst, baseURL string, disableAnalytics, embedStyle bool) error {
	var pis []PageInfo
	err := filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}

		dstFile, u := pageName(src, dst, baseURL, path)
		pi, err := ProcessFile(path, dstFile, u, disableAnalytics, embedStyle)
		if err != nil {
			return err
		}
		pis = append(pis, pi)
		return nil
	})
	if err != nil {
		return fmt.Errorf("ProcessDir dir=%s walk: %w", src, err)
	}

	sort.Slice(pis, func(i, j int) bool {
		return pis[i].URLCanonical > pis[j].URLCanonical
	})

	pi, err := blogIndex(dst, baseURL, pis)
	if err != nil {
		return fmt.Errorf("ProcessDir blog: %w", err)
	}
	pis = append(pis, pi)

	err = sitemap(dst, pis)
	if err != nil {
		return fmt.Errorf("ProcessDir sitemap: %w", err)
	}

	return nil
}

func ProcessFile(src, dst, u string, disableAnalytics, embedStyle bool) (PageInfo, error) {
	if path.Ext(src) != ".md" {
		return PageInfo{}, fmt.Errorf("ProcessFile ext=%s unknown ext", path.Ext(src))
	}
	b, err := ioutil.ReadFile(src)
	if err != nil {
		return PageInfo{}, fmt.Errorf("ProcessFile src=%s read: %w", src, err)
	}

	pd := PageData{
		URLCanonical:     u,
		DisableAnalytics: disableAnalytics,
		EmbedStyle:       embedStyle,
	}

	// extract header indo
	if bytes.HasPrefix(b, []byte(`---`)) {
		parts := bytes.SplitN(b, []byte(`---`), 3)
		err := yaml.Unmarshal(parts[1], &pd)
		if err != nil {
			return PageInfo{}, fmt.Errorf("ProcessFile src=%s parse header: %w", src, err)
		}
		b = parts[2]
	}

	// extract extra info
	pi := PageInfo{
		URLCanonical: u,
	}
	for _, ps := range strings.Split(src, "/") {
		if ps == "blog" {
			pi.BlogPost = true
			pi.Title = pd.Title
			pi.Date = path.Base(u)[:11]
			pd.Date = path.Base(u)[:11]
			pd.H1 = `<a href="/blog/">b<em>log</em></a>`
			pd.H2 = fmt.Sprintf(`<time datetime="%s">%s</time>`, pi.Date, pi.Date)
			break
		}
	}

	// render markdown
	var buf bytes.Buffer
	err = mdParser.Convert(b, &buf)
	if err != nil {
		return PageInfo{}, fmt.Errorf("ProcessFile src=%s render md: %w", src, err)
	}
	pd.Main = buf.String()

	// render template
	os.MkdirAll(path.Dir(dst), 0o755)
	f, err := os.Create(dst)
	if err != nil {
		return PageInfo{}, fmt.Errorf("ProcessFile dst=%s create: %w", dst, err)
	}
	defer f.Close()
	err = style.Template.ExecuteTemplate(f, "layout", pd)
	if err != nil {
		return PageInfo{}, fmt.Errorf("ProcessFile dst=%s render html: %w", dst, err)
	}

	return pi, nil
}

func blogIndex(dst, urlBase string, pis []PageInfo) (PageInfo, error) {
	dstPath := path.Join(urlBase, "blog")
	dstFile := path.Join(dst, "blog", "index.html")

	pd := PageData{
		Title:        "blog | seankhliao",
		Description:  "list of things i wrote",
		URLCanonical: dstPath,
		H1:           `<a href="/blog/">b<em>log</em></a>`,
		H2: `Artisanal, <em>hand-crafted</em> blog posts
imbued with delayed <em>regrets</em>`,
		Style: `
ul li {
        white-space: nowrap;
        overflow: hidden;
        text-overflow: ellipsis;
}`,
	}

	var buf strings.Builder
	buf.WriteString(`
<h3><em>B</em>log</h3>
<p>web log of things that never made sense,
maybe someone will find this useful</p>

<ul>`)
	for _, pi := range pis {
		if pi.BlogPost {
			relPath, _ := filepath.Rel(dstPath, pi.URLCanonical)
			buf.WriteString(fmt.Sprintf(`<li><time datetime="%s">%s</time> | <a href="%s">%s</a></li>`+"\n", pi.Date[1:], pi.Date, relPath+"/", pi.Title))
		}
	}
	buf.WriteString("</ul>")

	pd.Main = buf.String()

	f, err := os.Create(dstFile)
	if err != nil {
		return PageInfo{}, fmt.Errorf("blogIndex dst=%s create: %w", dstFile, err)
	}
	err = style.Template.ExecuteTemplate(f, "layout", pd)
	if err != nil {
		return PageInfo{}, fmt.Errorf("blogIndex dst=%s render html: %w", dstFile, err)
	}

	return PageInfo{
		URLCanonical: dstPath,
	}, nil
}

func sitemap(dst string, pis []PageInfo) error {
	urls := make([]string, len(pis))
	for i := range pis {
		urls[i] = pis[i].URLCanonical
	}
	sitemap := strings.Join(urls, "\n")
	dstFile := path.Join(dst, "sitemap.txt")

	err := ioutil.WriteFile(dstFile, []byte(sitemap), 0o644)
	if err != nil {
		return fmt.Errorf("generateSitemap dst=%s write: %w", dstFile, err)
	}
	return nil
}

func pageName(srcDir, dstDir, urlBase, file string) (dstFile, fullURL string) {
	relFile, _ := path.Rel(srcDir, file)
	noext := strings.TrimSuffix(relFile, ".md")
	withext := noext + ".html"
	return path.Join(dstDir, withext), path.Join(urlBase, noext) + "/"
}
