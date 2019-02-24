package main

import (
	"bufio"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"os"
)

func loadRSAPrivatePemKey(fileName string) *rsa.PrivateKey {
	privateKeyFile, err := os.Open(fileName)

	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}

	pemfileinfo, _ := privateKeyFile.Stat()
	var size int64 = pemfileinfo.Size()
	pembytes := make([]byte, size)
	buffer := bufio.NewReader(privateKeyFile)
	_, err = buffer.Read(pembytes)

	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
	data, _ := pem.Decode([]byte(pembytes))
	privateKeyFile.Close()
	privateKeyImported, err := x509.ParsePKCS1PrivateKey(data.Bytes)
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
	return privateKeyImported
}

func loadPublicPemKey(fileName string) *rsa.PublicKey {

	publicKeyFile, err := os.Open(fileName)
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}

	pemfileinfo, _ := publicKeyFile.Stat()

	size := pemfileinfo.Size()
	pembytes := make([]byte, size)
	buffer := bufio.NewReader(publicKeyFile)
	_, err = buffer.Read(pembytes)
	data, _ := pem.Decode([]byte(pembytes))
	publicKeyFile.Close()
	publicKeyFileImported, err := x509.ParsePKCS1PublicKey(data.Bytes)
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
	return publicKeyFileImported
}

func EncryptOAEP(secretMessage string, pubkey rsa.PublicKey) string {
	label := []byte("OAEP Encrypted")
	// crypto/rand.Reader is a good source of entropy for randomizing the
	// encryption function.
	rng := rand.Reader
	ciphertext, err := rsa.EncryptOAEP(sha256.New(), rng, &pubkey, []byte(secretMessage), label)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error from encryption: %s\n", err)
		return "Error from encryption"
	}
	return base64.StdEncoding.EncodeToString(ciphertext)
}

func DecryptOAEP(cipherText string, privKey rsa.PrivateKey) string {
	ct, _ := base64.StdEncoding.DecodeString(cipherText)
	label := []byte("OAEP Encrypted")

	// crypto/rand.Reader is a good source of entropy for blinding the RSA
	// operation.
	rng := rand.Reader
	plaintext, err := rsa.DecryptOAEP(sha256.New(), rng, &privKey, ct, label)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error from decryption: %s\n", err)
		return "Error from Decryption"
	}
	fmt.Printf("Plaintext: %s\n", string(plaintext))

	return string(plaintext)
}

func main() {
	// https://8gwifi.org/docs/go-rsa.jsp
	secretMessage := "Hello 8gwifi.org"
	// to generate: openssl genrsa -out manualpriv.pem 4096
	manualRSAPrivateKey := *loadRSAPrivatePemKey("./manualpriv.pem")
	// to generate: ssh-keygen -f manualpriv.pem -e -m pem > manualpub.pem
	manualPublicKey := loadPublicPemKey("./manualpub.pem")
	encryptedNewMessage := EncryptOAEP(secretMessage, *manualPublicKey)
	DecryptOAEP(encryptedNewMessage, manualRSAPrivateKey)
}
