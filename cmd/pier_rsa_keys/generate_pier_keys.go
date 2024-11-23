package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func main() {
	GenerateKeysRSA()
}

func GenerateKeysRSA() {
	outputDir := "pierkeys"
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		log.Fatalf("Не удалось создать директорию %s: %v", outputDir, err)
	}

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatalf("Не удалось сгенерировать приватный RSA ключ: %v", err)
	}

	publicBytes := x509.MarshalPKCS1PublicKey(&privateKey.PublicKey)

	pemKeyPBL := pem.EncodeToMemory(&pem.Block{Type: "RSA PUBLIC KEY", Bytes: publicBytes})
	rsaKeyPublicPath := filepath.Join(outputDir, "rsaPublicKey.pem")
	if err := ioutil.WriteFile(rsaKeyPublicPath, pemKeyPBL, 0600); err != nil {
		log.Fatalf("Не удалось записать приватный ключ CA: %v", err)
	}
	log.Println("rsaPublicKey.pem был успешно создан")

	privBytes := x509.MarshalPKCS1PrivateKey(privateKey)

	pemKeyPR := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: privBytes})
	rsaKeyPrivatePath := filepath.Join(outputDir, "rsaPrivateKey.pem")
	if err := ioutil.WriteFile(rsaKeyPrivatePath, pemKeyPR, 0600); err != nil {
		log.Fatalf("Не удалось записать приватный ключ CA: %v", err)
	}
	log.Println("rsaPrivateKey.pem был успешно создан")

}
