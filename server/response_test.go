package server

import (
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"net/http"
	"reflect"
	"testing"
	"time"
)

var (
	testDate   = time.Date(2000, time.January, 1, 0, 0, 0, 0, time.Local)
	clientCert = x509.Certificate{
		Subject: pkix.Name{
			CommonName: "testClient",
		},
		Issuer: pkix.Name{
			CommonName: "testCA",
		},
		NotBefore: testDate,
		NotAfter:  testDate.Add(time.Hour * 24),
		IsCA:      false,
	}
	caCert = x509.Certificate{
		Subject: pkix.Name{
			CommonName: "testCA",
		},
		Issuer: pkix.Name{
			CommonName: "testCA",
		},
		NotBefore: testDate,
		NotAfter:  testDate.Add(time.Hour * 24),
		IsCA:      true,
	}
	responseClientCert = certificate{
		IssuerCommonName:  "testCA",
		SubjectCommonName: "testClient",
		NotAfter:          "2000-01-02",
		NotBefore:         "2000-01-01",
		IsCA:              false,
	}
	responseCACert = certificate{
		IssuerCommonName:  "testCA",
		SubjectCommonName: "testCA",
		NotAfter:          "2000-01-02",
		NotBefore:         "2000-01-01",
		IsCA:              true,
	}
)

func Test_generateResponse(t *testing.T) {
	type args struct {
		r               *http.Request
		certificatePath string
	}
	tests := []struct {
		name string
		args args
		want response
	}{
		{
			name: "Nil TLS test",
			args: args{
				certificatePath: "path",
				r:               &http.Request{},
			},
			want: response{},
		},
		{
			name: "Valid mTLS",
			args: args{
				certificatePath: "path",
				r:               &http.Request{TLS: verifiedTLS()},
			},
			want: response{
				MTLSValid:             true,
				ClientCertificatePath: "path",
				PresentedCertificates: []certificate{responseClientCert},
				VerifiedChains:        []certificateChain{{responseClientCert, responseCACert}},
			},
		},
		{
			name: "No verified chains",
			args: args{
				certificatePath: "path",
				r:               &http.Request{TLS: noVerifiedTLS()},
			},
			want: response{
				MTLSValid:             false,
				ClientCertificatePath: "path",
				PresentedCertificates: []certificate{responseClientCert},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := generateResponse(tt.args.r, tt.args.certificatePath); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("generateResponse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func verifiedTLS() *tls.ConnectionState {
	state := &tls.ConnectionState{
		PeerCertificates: []*x509.Certificate{&clientCert},
		VerifiedChains:   [][]*x509.Certificate{{&clientCert, &caCert}},
	}
	return state
}

func noVerifiedTLS() *tls.ConnectionState {
	state := &tls.ConnectionState{
		PeerCertificates: []*x509.Certificate{&clientCert},
	}
	return state
}
