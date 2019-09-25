package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

// ImgOptions holds configs for resizing / transforming images
// has external dependencies:
type ImgOptions struct {
	Src string
	Dst string
}

func NewImgOptions(args []string) *ImgOptions {
	var o ImgOptions
	f := flag.NewFlagSet("img", flag.ExitOnError)
	f.StringVar(&o.Src, "src", "img", "source directory")
	f.StringVar(&o.Dst, "dst", "dst", "output directory")
	f.Parse(args)
	return &o
}
func (o *ImgOptions) Exec(opt *Options) error {
	fis, err := ioutil.ReadDir(o.Src)
	if err != nil {
		return fmt.Errorf("ImgOptions.Exec read dir %v: %w", o.Src, err)
	}

	var wg sync.WaitGroup

	for _, fi := range fis {

		if strings.HasPrefix(fi.Name(), "icon.") {
			wg.Add(1)
			go func(fi os.FileInfo) {
				defer wg.Done()

				err := ImgConvertIcon{
					filepath.Join(o.Src,
						fi.Name()),
					o.Dst,
					60,
					[]int{512, 192, 128, 64, 48, 32, 16},
				}.Exec()
				if err != nil {
					log.Printf("ImgOptions.Exec %v: %v\n", fi.Name(), err)
				}
			}(fi)
		} else if filepath.Ext(fi.Name()) == ".svg" {
			wg.Add(1)
			go func(fi os.FileInfo) {
				defer wg.Done()

				err := ImgConvertSvg{
					filepath.Join(o.Src, fi.Name()),
					o.Dst,
					1200,
					"1920x1235",
				}.Exec()

				if err != nil {
					log.Printf("ImgOptions.Exec %v: %v\n", fi.Name(), err)
				}
			}(fi)
		} else {
			log.Printf("ImgOptions.Exec unhandled file type %v\n", fi.Name())
		}
	}

	wg.Wait()
	return nil
}

type ImgConvertSvg struct {
	src     string
	dst     string
	density int
	size    string
}

func (img ImgConvertSvg) Exec() error {
	dstpref := filepath.Join(img.dst, strings.TrimSuffix(filepath.Base(img.src), ".svg"))
	args := []string{"-background", "none", "-density", strconv.Itoa(img.density), "-resize", img.size, img.src, "-write", dstpref + ".webp", dstpref + ".png"}

	cmd := exec.Command("convert", args...)
	_, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("ImgConvertSvg.Exec: %w", err)
	}
	return nil
}

type ImgConvertIcon struct {
	src     string // source file name
	dst     string // destination dir
	quality int
	sizes   []int
}

func (img ImgConvertIcon) Exec() error {
	args := []string{img.src, "-flatten"}
	for i, s := range img.sizes {
		is := strconv.Itoa(s)
		a := []string{"-resize", is + "x" + is, "-quality", strconv.Itoa(img.quality), "-write", filepath.Join(img.dst, "icon-"+is+".png")}
		if i < len(img.sizes)-1 {
			a = append([]string{"(", "+clone"}, a...)
			a = append(a, []string{"+delete", ")"}...)
		}
		args = append(args, a...)
	}

	cmd := exec.Command("convert", args...)
	_, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("ImgConvertIcon.Exec: %w", err)
	}
	return nil
}
