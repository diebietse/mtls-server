package server

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"

	"golang.org/x/crypto/acme"
	"golang.org/x/crypto/acme/autocert"
)

type MTLSServer struct {
	httpsServer *http.Server
	httpServer  *http.Server
}

const (
	stagingACMEDirectoryURL = "https://acme-staging-v02.api.letsencrypt.org/directory"
	httpAddress             = ":80"
	httpsAddress            = ":443"
)

// Config contains the configuration forht the mTLS server instance.
type Config struct {
	SiteFQDN       string
	TemplateFile   string
	ClientCertName string
	ClientCAPool   *x509.CertPool
	UseStaging     bool
}

// New create a mTLS server with an ACME certificate manager.
func New(config Config) (*MTLSServer, error) {
	// Setup ACME client
	certManager := &autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(config.SiteFQDN),
	}

	if config.UseStaging {
		certManager.Client = &acme.Client{DirectoryURL: stagingACMEDirectoryURL}
	}
	tlsConfig := certManager.TLSConfig()

	// Setup client certificate verification
	tlsConfig.ClientCAs = config.ClientCAPool
	tlsConfig.ClientAuth = tls.VerifyClientCertIfGiven

	// Load index web template
	webTemplate, err := loadWebTemplate(config.TemplateFile)
	if err != nil {
		return nil, err
	}

	httpsServer := &http.Server{
		Addr:      httpsAddress,
		TLSConfig: tlsConfig,
	}

	httpServer := &http.Server{
		Addr:    httpAddress,
		Handler: certManager.HTTPHandler(nil),
	}

	mTLSServer := &MTLSServer{
		httpServer:  httpServer,
		httpsServer: httpsServer,
	}

	// Setup handlers for the web endpoints
	setupHandlers(config.ClientCertName, webTemplate)

	return mTLSServer, nil
}

// ListenAndServe sets both the HTTP and HTTPS server to listen and serve
func (m *MTLSServer) ListenAndServe() {
	go func() {
		// serve HTTP, which will redirect automatically to HTTPS
		if err := m.httpServer.ListenAndServe(); err != nil {
			log.Fatalf("HTTP server failed: %v", err)
		}
	}()

	if err := m.httpsServer.ListenAndServeTLS("", ""); err != nil {
		log.Fatalf("HTTPS server failed: %v", err)
	}
}

func setupHandlers(clientCert string, webTemplate *template.Template) {
	clientCertPath := fmt.Sprintf("/%s", clientCert)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		requestHandlerHTML(w, r, webTemplate, clientCertPath)
	})
	http.HandleFunc("/json", func(w http.ResponseWriter, r *http.Request) {
		requestHandlerJSON(w, r, clientCertPath)
	})
	http.HandleFunc("/images/mtls-on.svg", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "images/mtls-on.svg")
	})
	http.HandleFunc("/images/mtls-off.svg", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "images/mtls-off.svg")
	})
	http.HandleFunc(clientCertPath, func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, clientCert)
	})
}

func requestHandlerHTML(w http.ResponseWriter, r *http.Request, webTemplate *template.Template, downloadLink string) {
	defer r.Body.Close()
	resp := generateResponse(r, downloadLink)

	err := webTemplate.Execute(w, resp)
	if err != nil {
		log.Printf("Could not generate response page: %v", err)
	}
}

func requestHandlerJSON(w http.ResponseWriter, r *http.Request, downloadLink string) {
	defer r.Body.Close()

	resp := generateResponse(r, downloadLink)
	info, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !resp.MTLSValid {
		w.WriteHeader(http.StatusUnauthorized)
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(info)
	if err != nil {
		log.Printf("Write error: %v", err)
	}
}

func loadWebTemplate(filename string) (*template.Template, error) {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return template.New("index").Parse(string(file))
}
