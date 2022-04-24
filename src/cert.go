package src

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"net"
	"os"
	"time"
)

// genCertFiles creates a self-signed certificate and it's private key in config/certificates directory.
//
//  Inspired by https://gist.github.com/shaneutt/5e1995295cff6721c89a71d13a71c251
func genCertFiles() (err error) {
	showInfo("Web server:" + "Generating certificate")

	subject := pkix.Name{
		CommonName:    "xTeVe",
		Country:       []string{"US"},
		Locality:      []string{"San Francisco"},
		Organization:  []string{"xTeVe, Inc."},
		PostalCode:    []string{"94016"},
		Province:      []string{""},
		StreetAddress: []string{"Golden Gate Bridge"},
	}

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return
	}

	certPrivKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return
	}

	certPrivKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(certPrivKey),
	})

	cert := &x509.Certificate{
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		IPAddresses:  append(System.IPAddressesV4Raw, net.IPv6loopback),
		KeyUsage:     x509.KeyUsageDigitalSignature,
		NotAfter:     time.Now().AddDate(10, 0, 0),
		NotBefore:    time.Now(),
		SerialNumber: serialNumber,
		Subject:      subject,
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, cert, cert, &certPrivKey.PublicKey, certPrivKey)
	if err != nil {
		return
	}

	certPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	})

	err = os.WriteFile(System.File.ServerCertPrivKey, certPrivKeyPEM, 0644)
	if err != nil {
		return
	}

	err = os.WriteFile(System.File.ServerCert, certPEM, 0644)
	if err != nil {
		return
	}

	return
}
