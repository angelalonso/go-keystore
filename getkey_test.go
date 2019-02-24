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
	//"net/http"
	"os"
	"testing"
)

// test getting a key
func TestGetKey(t *testing.T) {
	// do we get a key?
	// given a testing pair, can it decrypt an encrypted message?

	// Let's test decryption
	/*
		secretMessage := "Testing ~n &/() Keys is not enough§"
		testPublicKey := loadPublicPemKey("./testpub.pem")
		encryptedMessage := EncryptOAEP(secretMessage, *testPublicKey)
		testRSAPrivateKey := *loadRSAPrivatePemKey("./testpriv.pem")
		decryptedMessage := DecryptOAEP(encryptedMessage, testRSAPrivateKey)
		fmt.Println(decryptedMessage)
		//TODO: test here that the private key can decrypt what was encrypted with the public one

		getkey_req, _ := http.NewRequest("GET", "/getkey/?user=test", nil)
		getkey_rsp := executeRequest(getkey_req)

		checkResponseCode(t, http.StatusOK, getkey_rsp.Code)

		if body := getkey_rsp.Body.String(); body != `"Public key for user test saved"` {
			t.Errorf("Expected an empty message. Got %s", body)
		} else {
			fmt.Println("- Test OK: upload key without keys folder")
		}
	*/
}

// test key registry

func TestKeyStore(t *testing.T) {
	// list key files vs keys on memory
	// remove a key
	// clean up all keys
}

func loadRSAPrivatePemKey(fileName string) *rsa.PrivateKey {
	privateKeyFile, err := os.Open(fileName)

	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		fmt.Println(" - You can generate the key pair like this:")
		fmt.Println("     openssl genrsa -out testpriv.pem 4096")
		fmt.Println("     ssh-keygen -f testpriv.pem -e -m pem > testpub.pem")
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

	return string(plaintext)
}
