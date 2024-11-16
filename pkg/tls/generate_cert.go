package tls

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"log"
	"math/big"
	"os"
	"time"
)

// GeneratePrivateKey - функция, которая генерирует приватный ключ.
func GeneratePrivateKey() *ecdsa.PrivateKey {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		log.Fatalf("Failed generate private key")
	}
	return privateKey
}

// CreateAcCertificate - шаблон создания CA сертификата, который может подписывать другие сертификаты
func CreateAcCertificate() (*ecdsa.PrivateKey, *x509.Certificate) {
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		log.Fatalf("Failed generate serial number")
	}

	caTemplate := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"MyCA"},
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(365 * 24 * time.Hour),

		KeyUsage:              x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
		IsCA:                  true, // Yes, I'm CA
	}

	privateKey := GeneratePrivateKey()

	derBytes, err := x509.CreateCertificate(rand.Reader, &caTemplate, &caTemplate, &privateKey.PublicKey, privateKey)
	if err != nil {
		log.Fatalf("Failed created certificate")
	}
	pemCertificate := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	if pemCertificate == nil {
		log.Fatalf("Failed to encode certificate to PEM")
	}
	if err := os.WriteFile("/caCertificate.pem", pemCertificate, 0644); err != nil {
		log.Fatal(err)
	}
	log.Println("caCertificate.pem was successfully created")

	privBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		log.Fatalf("Unable to marshal private key: %v", err)
	}
	pemKey := pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: privBytes})
	if pemKey == nil {
		log.Fatal("Failed to encode key to PEM")
	}
	if err := os.WriteFile("/caPrivateKey.pem", pemKey, 0600); err != nil {
		log.Fatal(err)
	}
	log.Println("caPrivateKey.pem was successfully created")

	return privateKey, &caTemplate
}

// CreateServerCertificate - шаблон создания сертификата server, в нашем случае самостоятельно-подписывающийся сертификат.
func CreateServerCertificate(caKey *ecdsa.PrivateKey, caCertificate *x509.Certificate) {
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		log.Fatalf("Failed generate serial number")
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"BusinezzHack"},
		},
		DNSNames:  []string{"localhost"},
		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(3 * time.Hour),

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	privateKey := GeneratePrivateKey()

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, caCertificate, &privateKey.PublicKey, caKey)
	if err != nil {
		log.Fatalf("Failed created certificate")
	}
	pemCertificate := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	if pemCertificate == nil {
		log.Fatalf("Failed to encode certificate to PEM")
	}
	if err := os.WriteFile("/serverCertificate.pem", pemCertificate, 0644); err != nil {
		log.Fatal(err)
	}
	log.Println("serverCertificate.pem was successfully created")

	privBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		log.Fatalf("Unable to marshal private key: %v", err)
	}
	pemKey := pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: privBytes})
	if pemKey == nil {
		log.Fatal("Failed to encode key to PEM")
	}
	if err := os.WriteFile("/serverPrivateKey.pem", pemKey, 0600); err != nil {
		log.Fatal(err)
	}
	log.Println("serverPrivateKey.pem was successfully created")
}

// CreateClientCertificate - шаблон создания сертификата client, в нашем случае самостоятельно-подписывающийся сертификат.
func CreateClientCertificate(caKey *ecdsa.PrivateKey, caCertificate *x509.Certificate) {
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		log.Fatalf("Failed generate serial number")
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"BusinezzHack"},
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(3 * time.Hour),

		KeyUsage:              x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		BasicConstraintsValid: true,
	}

	privateKey := GeneratePrivateKey()

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, caCertificate, &privateKey.PublicKey, caKey)
	if err != nil {
		log.Fatalf("Failed created certificate")
	}
	pemCertificate := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	if pemCertificate == nil {
		log.Fatalf("Failed to encode certificate to PEM")
	}
	if err := os.WriteFile("/clientCertificate.pem", pemCertificate, 0644); err != nil {
		log.Fatal(err)
	}
	log.Println("clientCertificate.pem was successfully created")

	privBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		log.Fatalf("Unable to marshal private key: %v", err)
	}
	pemKey := pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: privBytes})
	if pemKey == nil {
		log.Fatal("Failed to encode key to PEM")
	}
	if err := os.WriteFile("/clientPrivateKey.pem", pemKey, 0600); err != nil {
		log.Fatal(err)
	}
	log.Println("clientPrivateKey.pem was successfully created")
}
