package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"io/ioutil"
	"log"
	"math/big"
	"os"
	"path/filepath"
	"time"
)

func main() {
	// Создаем CA-сертификат
	caKey, caCert := CreateCACertificate()

	// Создаем сертификат сервера, подписанный CA
	CreateServerCertificate(caKey, caCert)

	// Создаем сертификат клиента, подписанный CA
	CreateClientCertificate(caKey, caCert)
}

// GeneratePrivateKey генерирует приватный ключ ECDSA
func GeneratePrivateKey() *ecdsa.PrivateKey {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		log.Fatalf("Не удалось сгенерировать приватный ECDSA ключ: %v", err)
	}
	return privateKey
}

func GeneratePrivateKeyRSA() *rsa.PrivateKey {
	privateKey, err := rsa.GenerateKey(rand.Reader, 256)
	if err != nil {
		log.Fatalf("Не удалось сгенерировать приватный RSA ключ: %v", err)
	}
	return privateKey
}

// CreateCACertificate создает CA-сертификат и приватный ключ
func CreateCACertificate() (*ecdsa.PrivateKey, *x509.Certificate) {
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		log.Fatalf("Не удалось сгенерировать серийный номер: %v", err)
	}

	caTemplate := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"MyCA"},
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(365 * 24 * time.Hour),

		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageDigitalSignature,
		BasicConstraintsValid: true,
		IsCA:                  true, // Это CA-сертификат
	}

	privateKey := GeneratePrivateKey()

	derBytes, err := x509.CreateCertificate(rand.Reader, &caTemplate, &caTemplate, &privateKey.PublicKey, privateKey)
	if err != nil {
		log.Fatalf("Не удалось создать CA-сертификат: %v", err)
	}

	// Убедимся, что директория certificates существует
	outputDir := "certificates"
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		log.Fatalf("Не удалось создать директорию %s: %v", outputDir, err)
	}

	// Записываем CA-сертификат в PEM формате
	caCertPEMPath := filepath.Join(outputDir, "caCertificate.pem")
	pemCertificate := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	if err := ioutil.WriteFile(caCertPEMPath, pemCertificate, 0644); err != nil {
		log.Fatalf("Не удалось записать CA-сертификат PEM: %v", err)
	}
	log.Println("caCertificate.pem был успешно создан")

	// Записываем CA-сертификат в DER формате (.crt)
	caCertCRTPath := filepath.Join(outputDir, "caCertificate.crt")
	if err := ioutil.WriteFile(caCertCRTPath, derBytes, 0644); err != nil {
		log.Fatalf("Не удалось записать CA-сертификат CRT: %v", err)
	}
	log.Println("caCertificate.crt был успешно создан")

	// Записываем приватный ключ CA
	privBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		log.Fatalf("Не удалось маршализовать приватный ключ: %v", err)
	}
	pemKey := pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: privBytes})
	caKeyPath := filepath.Join(outputDir, "caPrivateKey.pem")
	if err := ioutil.WriteFile(caKeyPath, pemKey, 0600); err != nil {
		log.Fatalf("Не удалось записать приватный ключ CA: %v", err)
	}
	log.Println("caPrivateKey.pem был успешно создан")

	return privateKey, &caTemplate
}

// CreateServerCertificate создает сертификат сервера, подписанный CA
func CreateServerCertificate(caKey *ecdsa.PrivateKey, caCertificate *x509.Certificate) {
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		log.Fatalf("Не удалось сгенерировать серийный номер: %v", err)
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"BusinezzHack"},
		},
		DNSNames:  []string{"localhost", "auth_service"},
		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(365 * 24 * time.Hour),

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	privateKey := GeneratePrivateKey()

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, caCertificate, &privateKey.PublicKey, caKey)
	if err != nil {
		log.Fatalf("Не удалось создать сертификат сервера: %v", err)
	}

	// Убедимся, что директория certificates существует
	outputDir := "certificates"
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		log.Fatalf("Не удалось создать директорию %s: %v", outputDir, err)
	}

	// Записываем сертификат сервера в PEM формате
	serverCertPEMPath := filepath.Join(outputDir, "serverCertificate.pem")
	pemCertificate := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	if err := ioutil.WriteFile(serverCertPEMPath, pemCertificate, 0644); err != nil {
		log.Fatalf("Не удалось записать сертификат сервера PEM: %v", err)
	}
	log.Println("serverCertificate.pem был успешно создан")

	// Записываем сертификат сервера в DER формате (.crt)
	serverCertCRTPath := filepath.Join(outputDir, "serverCertificate.crt")
	if err := ioutil.WriteFile(serverCertCRTPath, derBytes, 0644); err != nil {
		log.Fatalf("Не удалось записать сертификат сервера CRT: %v", err)
	}
	log.Println("serverCertificate.crt был успешно создан")

	// Записываем приватный ключ сервера
	privBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		log.Fatalf("Не удалось маршализовать приватный ключ сервера: %v", err)
	}
	pemKey := pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: privBytes})
	serverKeyPath := filepath.Join(outputDir, "serverPrivateKey.pem")
	if err := ioutil.WriteFile(serverKeyPath, pemKey, 0600); err != nil {
		log.Fatalf("Не удалось записать приватный ключ сервера: %v", err)
	}
	log.Println("serverPrivateKey.pem был успешно создан")
}

// CreateClientCertificate создает сертификат клиента, подписанный CA
func CreateClientCertificate(caKey *ecdsa.PrivateKey, caCertificate *x509.Certificate) {
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		log.Fatalf("Не удалось сгенерировать серийный номер: %v", err)
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"BusinezzHack"},
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(365 * 24 * time.Hour),

		KeyUsage:              x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		BasicConstraintsValid: true,
	}

	privateKey := GeneratePrivateKey()

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, caCertificate, &privateKey.PublicKey, caKey)
	if err != nil {
		log.Fatalf("Не удалось создать сертификат клиента: %v", err)
	}

	// Убедимся, что директория certificates существует
	outputDir := "certificates"
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		log.Fatalf("Не удалось создать директорию %s: %v", outputDir, err)
	}

	// Записываем сертификат клиента в PEM формате
	clientCertPEMPath := filepath.Join(outputDir, "clientCertificate.pem")
	pemCertificate := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	if err := ioutil.WriteFile(clientCertPEMPath, pemCertificate, 0644); err != nil {
		log.Fatalf("Не удалось записать сертификат клиента PEM: %v", err)
	}
	log.Println("clientCertificate.pem был успешно создан")

	// Записываем сертификат клиента в DER формате (.crt)
	clientCertCRTPath := filepath.Join(outputDir, "clientCertificate.crt")
	if err := ioutil.WriteFile(clientCertCRTPath, derBytes, 0644); err != nil {
		log.Fatalf("Не удалось записать сертификат клиента CRT: %v", err)
	}
	log.Println("clientCertificate.crt был успешно создан")

	// Записываем приватный ключ клиента
	privBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		log.Fatalf("Не удалось маршализовать приватный ключ клиента: %v", err)
	}
	pemKey := pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: privBytes})
	clientKeyPath := filepath.Join(outputDir, "clientPrivateKey.pem")
	if err := ioutil.WriteFile(clientKeyPath, pemKey, 0600); err != nil {
		log.Fatalf("Не удалось записать приватный ключ клиента: %v", err)
	}
	log.Println("clientPrivateKey.pem был успешно создан")
}
