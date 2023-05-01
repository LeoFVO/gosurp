package dkim

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"time"
)

// Signer represents a DKIM signer.
type Signer struct {
	PrivateKey *rsa.PrivateKey
	Domain     string
	Selector   string
}

// Sign signs the provided email message and returns the DKIM signature header.
func (s *Signer) Sign(msg string) (string, error) {
	h := sha256.New()
	h.Write([]byte(msg))
	hash := h.Sum(nil)

	header := fmt.Sprintf("v=1; a=rsa-sha256; d=%s; s=%s; t=%d; c=relaxed/relaxed; bh=%s; b=", s.Domain, s.Selector, time.Now().Unix(), base64.StdEncoding.EncodeToString(hash))

	sig, err := rsa.SignPKCS1v15(rand.Reader, s.PrivateKey, crypto.SHA256, hash)
	if err != nil {
		return "", err
	}

	header += base64.StdEncoding.EncodeToString(sig)

	return header, nil
}