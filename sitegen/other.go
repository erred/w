package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
)

func (o *options) convertImgs() error {
	for i, imgArgs := range defaultImgArgs {
		out, err := exec.Command("convert", imgArgs...).CombinedOutput()
		if err != nil {
			return fmt.Errorf("options.convertImgs: %d: %w\n%s", i, err, out)
		}
	}
	return nil
}

func (o *options) deploy() error {
	cmd := exec.Command("firebase", "-P", o.gcpProject, "deploy")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("options.deploy: %w", err)
	}
	return nil
}

type amper interface {
	setAMP()
}

// writeTemplate takes a target path (without dst) and executes the named template
// twice (once for normal, once for amp)
func (o *options) writeTemplate(fp, tmpl string, data interface{}) error {
	f, err := openWrite(filepath.Join(o.dst, fp))
	if err != nil {
		return fmt.Errorf("options.writeTemplate: %v", err)
	}
	defer f.Close()
	o.templates.ExecuteTemplate(f, tmpl, data)

	if d, ok := data.(amper); ok {
		d.setAMP()
		f, err = openWrite(filepath.Join(o.dst, "amp", fp))
		if err != nil {
			return fmt.Errorf("options.writeTemplate: %v", err)
		}
		defer f.Close()
		o.templates.ExecuteTemplate(f, tmpl, data)
	}

	return nil
}

func (o *options) copyFile(fp string, done *sync.WaitGroup) {
	if done != nil {
		defer done.Done()
	}
	fps := strings.Split(fp, "/")
	f1, err := os.Open(fp)
	if err != nil {
		log.Println("options.parsePages: copy open f1", fp, err)
		return
	}
	defer f1.Close()

	fps[0] = o.dst
	f2, err := openWrite(filepath.Join(fps...))
	if err != nil {
		log.Println("options.parsePages: copy open f2", filepath.Join(fps...), err)
		return
	}
	defer f2.Close()
	io.Copy(f2, f1)
}

func openWrite(fn string) (*os.File, error) {
	err := os.MkdirAll(filepath.Dir(fn), 0755)
	if err != nil {
		return nil, fmt.Errorf("write: %w", err)
	}
	f, err := os.Create(fn)
	if err != nil {
		return nil, fmt.Errorf("write: %w", err)
	}
	return f, nil
}

func normalizeURL(u string) string {
	return strings.TrimSuffix(strings.TrimSuffix(u, ".html"), "index")
}

// func (o *options) getFonts() error {
// 	res, err := http.Get(o.fontURL)
// 	if err != nil {
// 		return fmt.Errorf("options.getFonts: %w", err)
// 	} else if res.StatusCode < 200 || res.StatusCode > 299 {
// 		return fmt.Errorf("options.getFonts: %d %s", res.StatusCode, res.Status)
// 	}
// 	defer res.Body.Close()
// 	buf := bytes.NewBufferString(`{{ define "fontcss" }}`)
// 	buf.ReadFrom(res.Body)
// 	buf.WriteString(`{{ end }}`)
//
// 	o.templates, err = o.templates.New("fontcss").Parse(buf.String())
// 	if err != nil {
// 		return fmt.Errorf("options.getFonts: %w", err)
// 	}
// 	return nil
// }
