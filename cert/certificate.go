package cert

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"net"
	"strings"
	"time"
)

func Generate(host string) (_key, _cert *bytes.Buffer) {
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		logger.Error("failed to generate private key: " + err.Error())
		return
	}

	var notBefore = time.Now()
	notAfter := notBefore.Add(365 * 24 * time.Hour)

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		logger.Error("failed to generate serial number: " + err.Error())
		return
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"hyper-scheduler"},
		},
		NotBefore: notBefore,
		NotAfter:  notAfter,

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	hosts := strings.Split(host, ",")
	for _, h := range hosts {
		if ip := net.ParseIP(h); ip != nil {
			template.IPAddresses = append(template.IPAddresses, ip)
		} else {
			template.DNSNames = append(template.DNSNames, h)
		}
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		logger.Error("Failed to create certificate: " + err.Error())
		return
	}
	cert := new(bytes.Buffer)
	key := new(bytes.Buffer)

	if err := pem.Encode(cert, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes}); err != nil {
		logger.Error("failed to write data to cert.pem: " + err.Error())
		return
	}

	b, err := x509.MarshalECPrivateKey(priv)
	if err != nil {
		logger.Error("Unable to marshal ECDSA private key: " + err.Error())
		return
	}

	if err := pem.Encode(key, &pem.Block{Type: "EC PRIVATE KEY", Bytes: b}); err != nil {
		logger.Error("failed to write data to key.pem: " + err.Error())
		return
	}
	return key, cert
}
