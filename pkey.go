package main

import (
	"crypto/ed25519"
	"crypto/pbkdf2"
	"crypto/sha3"
	"os"

	"golang.org/x/crypto/ssh"
)

var saltValue = [16]byte{ 0x94, 0xd6, 0x65, 0x0c, 0x62, 0xb7, 0x86, 0x2f, 0x62, 0xe7, 0x9b, 0x0f, 0x4c, 0xed, 0xe2, 0x21 }

// Generate private key based on password
func GeneratePrivateKey(passwd string) (ssh.Signer, error) {
	// 1. Generate random key:
	seed, err := pbkdf2.Key(sha3.New256, passwd, saltValue[:], 16384, ed25519.SeedSize)
	if err != nil {
		return nil, err
	}

	// 2. And now we can create the key:
	priv := ed25519.NewKeyFromSeed(seed)

	// 3. And finally signer:
	return ssh.NewSignerFromKey(priv)
}

// Load private key from file or generate based on command line parameter
func LoadPrivateKey(key string) (ssh.Signer, error) {
	if len(key) > 1 && key[0] == '@' {
		return GeneratePrivateKey(key)
	}

	// Try to read and parse the key:
	keyData, err := os.ReadFile(key)
	if err != nil {
		return nil, err
	}

	return ssh.ParsePrivateKey(keyData)
}
