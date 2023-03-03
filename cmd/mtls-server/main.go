package main

import (
	"crypto/x509"
	"errors"
	"flag"
	"log"
	"os"

	"github.com/diebietse/mtls-server/server"
)

const (
	fqdnEnv = "DEMO_FQDN"
)

func main() {
	staging := flag.Bool("staging", false, "Use letsencrypt staging ACME server")
	indexTemplate := flag.String("index-template", "index.html", "Set the filename for the index template to load")
	rootCA := flag.String("root-ca", "root.pem", "root CA")
	clientCert := flag.String("client-cert", "mtls-example-client.p12", "Relative link to the client certificate")
	flag.Parse()

	caCertPool, err := loadClientCA(*rootCA)
	if err != nil {
		log.Fatalf("Could not load client CA: %+v", err)
	}

	fqdn := os.Getenv(fqdnEnv)
	if len(fqdn) == 0 {
		log.Fatalf("Please provide an FQDN by setting the environment variable '%s'", fqdnEnv)
	}

	s, err := server.New(server.Config{
		SiteFQDN:       fqdn,
		TemplateFile:   *indexTemplate,
		ClientCertName: *clientCert,
		ClientCAPool:   caCertPool,
		UseStaging:     *staging,
	})
	if err != nil {
		log.Fatal(err)
	}

	s.ListenAndServe()
}

func loadClientCA(rootCA string) (*x509.CertPool, error) {
	pem, err := os.ReadFile(rootCA)
	if err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pem) {
		return nil, errors.New("could not decode CA pem")
	}

	return certPool, nil
}
