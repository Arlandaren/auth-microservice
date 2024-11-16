package utils

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"google.golang.org/grpc/credentials"
	"io/ioutil"
)

func LoadServerTLS() (*tls.Config, error) {
	caCertificate := "./certificates/caCertificate.pem"
	serverCertificate := "./certificates/serverCertificate.pem"
	serverPrivateKey := "./certificates/serverPrivateKey.pem"

	// Загрузка сертификата сервера
	serverCert, err := tls.LoadX509KeyPair(serverCertificate, serverPrivateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to load server key pair: %w", err)
	}

	// Загрузка CA сертификата
	caCert, err := ioutil.ReadFile(caCertificate)
	if err != nil {
		return nil, fmt.Errorf("failed to read CA certificate: %w", err)
	}
	caCertPool := x509.NewCertPool()
	if !caCertPool.AppendCertsFromPEM(caCert) {
		return nil, fmt.Errorf("failed to append CA certificate to pool")
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{serverCert},  // Сертификат сервера
		ClientCAs:    caCertPool,                     // CA, которому доверяем
		ClientAuth:   tls.RequireAndVerifyClientCert, // Требовать сертификаты от клиентов
	}
	return tlsConfig, nil
}

func LoadServerTLSCredentials() (credentials.TransportCredentials, error) {
	tlsConfig, err := LoadServerTLS()
	if err != nil {
		return nil, err
	}
	return credentials.NewTLS(tlsConfig), nil
}

func LoadClientTLSCredentials() (credentials.TransportCredentials, error) {
	caCertificate := "./certificates/caCertificate.pem"
	clientCertificate := "./certificates/clientCertificate.pem"
	clientPrivateKey := "./certificates/clientPrivateKey.pem"

	// Загрузка сертификата клиента
	clientCert, err := tls.LoadX509KeyPair(clientCertificate, clientPrivateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to load client key pair: %w", err)
	}

	// Загрузка CA сертификата
	caCert, err := ioutil.ReadFile(caCertificate)
	if err != nil {
		return nil, fmt.Errorf("failed to read CA certificate: %w", err)
	}
	caCertPool := x509.NewCertPool()
	if !caCertPool.AppendCertsFromPEM(caCert) {
		return nil, fmt.Errorf("failed to append CA certificate to pool")
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{clientCert}, // Сертификат клиента
		RootCAs:      caCertPool,                    // CA, которому доверяем
	}

	return credentials.NewTLS(tlsConfig), nil
}
