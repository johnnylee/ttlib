package ttlib

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"io/ioutil"
	"math/big"
	"net"
	"os"
	"time"
)

// Function that generates certificates. Mostly coppied from
// http://golang.org/src/pkg/crypto/tls/generate_cert.go.
func GenerateKeyAndCert(baseDir, host string) error {

	rsaBits := 2048

	// Generate the private key.
	priv, err := rsa.GenerateKey(rand.Reader, rsaBits)
	if err != nil {
		return err
	}

	// Generate a random serial number.
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return err
	}

	// Create a certificate template.
	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject:      pkix.Name{Organization: []string{"None"}},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().Add(24 * 365 * 128 * time.Hour),
		KeyUsage: (x509.KeyUsageKeyEncipherment |
			x509.KeyUsageDigitalSignature),
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	// Set the host name (or IP).
	if ip := net.ParseIP(host); ip != nil {
		template.IPAddresses = append(template.IPAddresses, ip)
	} else {
		template.DNSNames = append(template.DNSNames, host)
	}

	// Set certificate authority flags.
	template.IsCA = true
	template.KeyUsage |= x509.KeyUsageCertSign

	// Create the certificate.
	derBytes, err := x509.CreateCertificate(
		rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		return err
	}

	// Encode and save the key.
	keyPath := ServerKeyPath(baseDir)
	keyOut, err := os.OpenFile(keyPath, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer func() {
		_ = keyOut.Close()
	}()

	err = pem.Encode(
		keyOut,
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(priv),
		})
	if err != nil {
		return err
	}

	// Encode and save the certificate.
	certPath := ServerCertPath(baseDir)
	certOut, err := os.OpenFile(certPath, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer func() {
		_ = certOut.Close()
	}()

	err = pem.Encode(
		certOut,
		&pem.Block{
			Type:  "CERTIFICATE",
			Bytes: derBytes,
		})
	if err != nil {
		return err
	}

	return nil
}

// LoadKeyAndCert reads the key and certificate files into memory and
// returns them in that order.
func LoadKeyAndCert(baseDir string) ([]byte, []byte, error) {
	key, err := ioutil.ReadFile(ServerKeyPath(baseDir))
	if err != nil {
		return nil, nil, err
	}

	cert, err := ioutil.ReadFile(ServerCertPath(baseDir))
	if err != nil {
		return nil, nil, err
	}

	return key, cert, nil
}
