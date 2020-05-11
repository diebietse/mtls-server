package server

import (
	"net/http"
)

type certificateChain []certificate

const dateFormat = "2006-01-02"

type certificate struct {
	IssuerCommonName  string `json:"issuer_common_name"`
	SubjectCommonName string `json:"subject_common_name"`
	NotBefore         string `json:"not_before"`
	NotAfter          string `json:"not_after"`
	IsCA              bool   `json:"is_ca"`
}

type response struct {
	MTLSValid             bool               `json:"mtls_valid"`
	ClientCertificatePath string             `json:"client_certificate_path"`
	PresentedCertificates []certificate      `json:"presented_certificates"`
	VerifiedChains        []certificateChain `json:"verified_certificate_chains"`
}

func generateResponse(r *http.Request, certificatePath string) response {
	if r.TLS == nil {
		return response{}
	}

	resp := response{
		ClientCertificatePath: certificatePath,
	}

	for _, cert := range r.TLS.PeerCertificates {
		certInfo := certificate{
			IssuerCommonName:  cert.Issuer.CommonName,
			SubjectCommonName: cert.Subject.CommonName,
			NotAfter:          cert.NotAfter.Format(dateFormat),
			NotBefore:         cert.NotBefore.Format(dateFormat),
			IsCA:              cert.IsCA,
		}
		resp.PresentedCertificates = append(resp.PresentedCertificates, certInfo)
	}

	if len(r.TLS.VerifiedChains) > 0 {
		resp.MTLSValid = true
		for _, chain := range r.TLS.VerifiedChains {
			chainInfo := certificateChain{}
			for _, cert := range chain {
				certInfo := certificate{
					IssuerCommonName:  cert.Issuer.CommonName,
					SubjectCommonName: cert.Subject.CommonName,
					NotAfter:          cert.NotAfter.Format(dateFormat),
					NotBefore:         cert.NotBefore.Format(dateFormat),
					IsCA:              cert.IsCA,
				}
				chainInfo = append(chainInfo, certInfo)
			}
			resp.VerifiedChains = append(resp.VerifiedChains, chainInfo)
		}
	}
	return resp
}
