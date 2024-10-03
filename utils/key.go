package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
)

func GenKey() (string, string) {
	// Tạo cặp khóa RSA
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}

	// Định dạng khóa riêng PKCS#1
	privKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privKeyBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY", // Thay đổi ở đây
		Bytes: privKeyBytes,
	}
	privKeyString := string(pem.EncodeToMemory(privKeyBlock))

	// Định dạng khóa công khai
	pubKeyBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		panic(err)
	}
	pubKeyBlock := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubKeyBytes,
	}
	pubKeyString := string(pem.EncodeToMemory(pubKeyBlock))

	return pubKeyString, privKeyString
}
