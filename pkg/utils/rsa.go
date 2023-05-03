package utils

import (
	"crypto/rsa"
	"os"

	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
)

func mustNot(err error) {
    if err != nil {
        panic(err)
    }
}

func LoadRSAPrivateKey(path string) *rsa.PrivateKey {
	log.Tracef("Loading RSA private key from %s", path)
	bytes, err := os.ReadFile(path)
	mustNot(err)
	key, err := ssh.ParseRawPrivateKey(bytes)
	mustNot(err)
	
	return key.(*rsa.PrivateKey)
}