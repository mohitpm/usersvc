package saml

import (
	"context"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"net/http"
	"net/url"
	"os"

	"github.com/crewjam/saml/samlsp"
	cfg "github.com/mohitpm/usersvc/config"
)

type ServiceProvider struct {
	config cfg.SAML
}

func NewSP(samlConfig cfg.SAML) (*samlsp.Middleware, error) {
	cert, err := os.ReadFile(samlConfig.SP.CertFile)
	if err != nil {
		return nil, err
	}

	key, err := os.ReadFile(samlConfig.SP.KeyFile)
	if err != nil {
		return nil, err
	}

	keyPair, err := tls.X509KeyPair(cert, key)
	if err != nil {
		return nil, err
	}

	keyPair.Leaf, err = x509.ParseCertificate(keyPair.Certificate[0])
	if err != nil {
		return nil, err
	}

	idpMetadataURL, err := url.Parse(samlConfig.IDP.MetadataURL)
	if err != nil {
		return nil, err
	}

	idpMetadata, err := samlsp.FetchMetadata(context.Background(), http.DefaultClient, *idpMetadataURL)
	if err != nil {
		return nil, err
	}

	rootURL, err := url.Parse(samlConfig.SP.RootURL)
	if err != nil {
		return nil, err
	}

	samlSP, err := samlsp.New(samlsp.Options{
		URL:               *rootURL,
		Key:               keyPair.PrivateKey.(*rsa.PrivateKey),
		Certificate:       keyPair.Leaf,
		AllowIDPInitiated: true,
		IDPMetadata:       idpMetadata,
	})
	if err != nil {
		return nil, err
	}

	return samlSP, nil
}
