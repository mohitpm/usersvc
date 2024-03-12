/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/xml"
	"net/http"
	"net/url"
	"os"

	"github.com/crewjam/saml/samlsp"
	cfg "github.com/mohitpm/usersvc/config"
	"github.com/spf13/cobra"
)

// samlMetadataCmd represents the samlMetadata command
var samlMetadataCmd = &cobra.Command{
	Use:   "samlMetadata",
	Short: "Generate saml metadata of service provider",
	Long:  `Generate saml metadata of service provider`,
	Run: func(cmd *cobra.Command, args []string) {
		generateSAMLMetadata()
	},
}

func init() {
	generateCmd.AddCommand(samlMetadataCmd)
}

func generateSAMLMetadata() {
	config := cfg.LoadConfig()

	cert, err := os.ReadFile(config.SAML.SP.CertFile)
	if err != nil {
		panic(err)
	}

	key, err := os.ReadFile(config.SAML.SP.KeyFile)
	if err != nil {
		panic(err)
	}

	keyPair, err := tls.X509KeyPair(cert, key)
	if err != nil {
		panic(err)
	}

	keyPair.Leaf, err = x509.ParseCertificate(keyPair.Certificate[0])
	if err != nil {
		panic(err)
	}

	idpMetadataURL, err := url.Parse(config.SAML.IDP.MetadataURL)
	if err != nil {
		panic(err)
	}

	idpMetadata, err := samlsp.FetchMetadata(context.Background(), http.DefaultClient, *idpMetadataURL)
	if err != nil {
		panic(err)
	}

	rootURL, err := url.Parse(config.SAML.SP.RootURL)
	if err != nil {
		panic(err)
	}

	samlSP, err := samlsp.New(samlsp.Options{
		URL:               *rootURL,
		Key:               keyPair.PrivateKey.(*rsa.PrivateKey),
		Certificate:       keyPair.Leaf,
		AllowIDPInitiated: true,
		IDPMetadata:       idpMetadata,
	})
	if err != nil {
		panic(err)
	}

	// register with the service provider
	spMetadataBuf, _ := xml.MarshalIndent(samlSP.ServiceProvider.Metadata(), "", "  ")
	file, err := os.OpenFile("sp.xml", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	_, err = file.Write(spMetadataBuf)
	if err != nil {
		panic(err)
	}
}
