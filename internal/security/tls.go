package security

import (
	"crypto/tls"
	"crypto/x509"
	"os"

	"github.com/codejago/polypully-openai-checker/internal/config"
)

type Tls struct {
}

func NewTls() *Tls {
	return &Tls{}
}

func (s *Tls) SetupTLS(c config.Tls) (*tls.Config, error) {
	certPEMBlock, err := os.ReadFile(c.ServerCert)
	if err != nil {
		return nil, err
	}
	keyPEMBlock, err := os.ReadFile(c.ServerKey)
	if err != nil {
		return nil, err
	}
	cert, err := tls.X509KeyPair(certPEMBlock, keyPEMBlock)
	if err != nil {
		return nil, err
	}
	ca, err := os.ReadFile(c.ServerCaCert)
	if err != nil {
		return nil, err
	}
	caCertPool := x509.NewCertPool()
	if ok := caCertPool.AppendCertsFromPEM(ca); !ok {
		return nil, err
	}
	tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{cert},
		ClientCAs:          caCertPool,
		ClientAuth:         tls.RequireAndVerifyClientCert,
		InsecureSkipVerify: c.SelfSigned,
	}
	return tlsConfig, nil
}
