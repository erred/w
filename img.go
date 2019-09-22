package main

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
	return fmt.Errorf("ErrNotImplemented")
}

