package usvc

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

type TLSOpts struct {
	KeyFile string
	CrtFile string
	CAFile  string
}

func (o *TLSOpts) Flags(fs *flag.FlagSet) {
	fs.StringVar(&o.CrtFile, "tls.crt", "/var/secret/tls/tls.crt", "tls cert file")
	fs.StringVar(&o.KeyFile, "tls.key", "/var/secret/tls/tls.key", "tls key file")
	fs.StringVar(&o.CAFile, "ca.crt", "/var/secret/tls/ca.crt", "ca crt for client connections")
}

func (o *TLSOpts) Config() (*tls.Config, error) {
	c := &tls.Config{}

	cert, err := tls.LoadX509KeyPair(o.CrtFile, o.KeyFile)
	if errors.Is(err, os.ErrNotExist) {
		// skip no certs
	} else if err != nil {
		return nil, fmt.Errorf("usvc.tls crt=%s key=%s: %w", o.CrtFile, o.KeyFile, err)
	} else {
		c.Certificates = []tls.Certificate{cert}
	}

	b, err := ioutil.ReadFile(o.CAFile)
	if errors.Is(err, os.ErrNotExist) {
		// skip no certs
	} else if err != nil {
		return nil, fmt.Errorf("usvc.tls ca=%s: %w", o.CAFile, err)
	} else {
		cp := x509.NewCertPool()
		if !cp.AppendCertsFromPEM(b) {
			return nil, fmt.Errorf("usvc.tls append ca cert: %w", err)
		}
		c.RootCAs = cp
	}
	return c, nil
}
