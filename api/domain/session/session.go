package session

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"log"
	"math"
	"math/big"
	"os"
	"strconv"
	"time"
)

var csrfKey = []byte(os.Getenv("CSRF_SECRET"))

type Session struct {
	ID                          string
	Active                      bool
	AuthenticatorAssuranceLevel uint8
	IssuedAt                    time.Time
	ExpiresAt                   time.Time
	AuthenticatedAt             time.Time
	Devices                     []Device

	ExpiryInterval time.Duration
	csrfToken      string
}

func (s *Session) GetCSRFToken() (string, error) {
	if s.csrfToken != "" {
		return s.csrfToken, nil
	}
	if err := s.generateCSRFToken(); err != nil {
		return "", err
	}
	return s.csrfToken, nil
}

func (s *Session) generateCSRFToken() error {
	// Gather the values
	nBig, err := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))
	if err != nil || nBig == nil {
		log.Printf("Error generating random number: %v (nBig: %v)\n", err, nBig)
		return err
	}
	randomValue := nBig.Int64()
	randomValueString := strconv.FormatInt(randomValue, 16)

	// Create the CSRF Token
	message := append([]byte(s.ID), "!"+randomValueString...) // HMAC message payload
	mac := hmac.New(sha256.New, csrfKey)
	mac.Write(message) // Generate the HMAC hash
	macBytes := mac.Sum(nil)
	csrfToken := string(macBytes) + "." + string(message) // Combine HMAC hash with message to generate the token. The plain message is required to later authenticate it against its HMAC hash
	s.csrfToken = csrfToken
	return nil
}

func VerifyCSRFToken(message, messageMAC []byte) bool {
	mac := hmac.New(sha256.New, csrfKey)
	mac.Write(message)
	expectedMAC := mac.Sum(nil)
	return hmac.Equal(messageMAC, expectedMAC)
}
