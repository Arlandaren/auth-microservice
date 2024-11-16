package utils

import (
	"crypto/tls"
	"crypto/x509"
	"google.golang.org/grpc/credentials"
	"io/ioutil"
)

func LoadServerTLS() (*tls.Config, error) {
	caCertificate := "tls_certificates/ca.crt"
	serverCertificate := "tls_certificates/server.crt"
	serverPrivateKey := "tls_certificates/server.key"

	serverCert, err := tls.LoadX509KeyPair(serverCertificate, serverPrivateKey)
	if err != nil {
		return nil, err
	}

	caCert, err := ioutil.ReadFile(caCertificate)
	if err != nil {
		return nil, err
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	TlsConfig := &tls.Config{
		Certificates: []tls.Certificate{serverCert},  // Отправка сертификата сервера
		ClientCAs:    caCertPool,                     // CA, которому доверяем
		ClientAuth:   tls.RequireAndVerifyClientCert, // Требовать сертификаты от клиентов
	}
	return TlsConfig, nil
}

func LoadServerTLSCredentials() (credentials.TransportCredentials, error) {
	caCertificate := "tls_certificates/ca.crt"
	serverCertificate := "tls_certificates/server.crt"
	serverPrivateKey := "tls_certificates/server.key"

	serverCert, err := tls.LoadX509KeyPair(serverCertificate, serverPrivateKey)
	if err != nil {
		return nil, err
	}

	caCert, err := ioutil.ReadFile(caCertificate)
	if err != nil {
		return nil, err
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	TlsConfig := &tls.Config{
		Certificates: []tls.Certificate{serverCert},  // Отправка сертификата сервера
		ClientCAs:    caCertPool,                     // CA, которому доверяем
		ClientAuth:   tls.RequireAndVerifyClientCert, // Требовать сертификаты от клиентов
	}

	return credentials.NewTLS(TlsConfig), nil
}

func LoadClientTLSCredentials() (credentials.TransportCredentials, error) {
	caCertificate := "tls_certificates/ca.crt"
	clientCertificate := "tls_certificates/client.crt"
	clientPrivateKey := "tls_certificates/client.key"

	clientCert, err := tls.LoadX509KeyPair(clientCertificate, clientPrivateKey)
	if err != nil {
		return nil, err
	}

	caCert, err := ioutil.ReadFile(caCertificate)
	if err != nil {
		return nil, err
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{clientCert}, // Отправка сертификата клиента
		RootCAs:      caCertPool,                    // CA, которому доверяем
	}

	return credentials.NewTLS(tlsConfig), nil
}
