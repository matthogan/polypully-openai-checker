package security

import (
	"crypto/tls"
	"testing"

	"github.com/codejago/polypully-openai-checker/internal/config"
	"github.com/codejago/polypully-openai-checker/internal/security"
)

func TestSetupTLS(t *testing.T) {
	t.Run("setup tls", func(t *testing.T) {
		tlsConfig := config.Tls{
			ServerCert:   "testdata/server.crt",
			ServerKey:    "testdata/server.key",
			ServerCaCert: "testdata/ca.crt",
			SelfSigned:   true,
		}
		securityTls := security.NewTls()
		tlsCfg, err := securityTls.SetupTLS(tlsConfig)

		if err != nil {
			t.Fatalf("failed: %v", err)
		}
		if tlsCfg == nil {
			t.Fatal("expected non-nil tls.Config")
		}
		if len(tlsCfg.Certificates) != 1 {
			t.Fatalf("expected 1 certificate, got %d", len(tlsCfg.Certificates))
		}
		if tlsCfg.ClientAuth != tls.RequireAndVerifyClientCert {
			t.Fatalf("expected ClientAuth to be RequireAndVerifyClientCert, got %v", tlsCfg.ClientAuth)
		}
		if !tlsCfg.InsecureSkipVerify {
			t.Fatal("expected InsecureSkipVerify to be true")
		}
	})
}
