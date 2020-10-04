package usvc

import (
	"flag"
	"fmt"

	"go.seankhliao.com/apis/saver/v1"
	"google.golang.org/grpc"
)

type SaverOpts struct {
	Addr string
}

func (o *SaverOpts) Flag(fs *flag.FlagSet) {
	fs.StringVar(&o.Addr, "saver.addr", "saver:443", "host:port of saver")
}

func (o SaverOpts) Saver(opts ...grpc.DialOption) (client saver.SaverClient, shutdown func() error, err error) {
	cc, err := grpc.Dial(o.Addr, opts...)
	if err != nil {
		return nil, nil, fmt.Errorf("usvc.saver dial addr=%v: %w", o.Addr, err)
	}
	client = saver.NewSaverClient(cc)
	shutdown = func() error {
		return cc.Close()
	}
	return client, shutdown, nil

}
