package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"os"
)

func main() {
	if err := run(); err != nil {
		log.Fatalln(err)
	}
}

func run() error {

	// Generate a new private key.
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}

	// // Create a file for the private key in PEM format.
	// privateFile, err := os.Create("private.pem")
	// if err != nil {
	// 	return fmt.Errorf("error creating private key file: %w", err)
	// }
	// defer privateFile.Close()

	// // Construct a PEM block for the private key.
	// privateBlock := pem.Block{
	// 	Type:  "PRIVATE KEY",
	// 	Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	// }

	// // Write the private key to the file.
	// if err := pem.Encode(privateFile, &privateBlock); err != nil {
	// 	return fmt.Errorf("error encoding private key: %w", err)
	// }

	// ############################################################

	// publicFile, err := os.Create("public.pem")
	// if err != nil {
	// 	return fmt.Errorf("error creating public key file: %w", err)
	// }
	// defer publicFile.Close()

	// Mashal the public key from the private key.
	asn1Bytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return fmt.Errorf("error marshaling public key: %w", err)
	}

	// Construct a PEM block for the public key.
	publicBlock := pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: asn1Bytes,
	}

	// // Write the public key to the public key file.
	// if err := pem.Encode(publicFile, &publicBlock); err != nil {
	// 	return fmt.Errorf("error encoding public key: %w", err)
	// }
	// Write the public key to the public key file.
	if err := pem.Encode(os.Stdout, &publicBlock); err != nil {
		return fmt.Errorf("error encoding public key: %w", err)
	}

	fmt.Println("Private and Public Keys generated successfully.")

	return nil
}
