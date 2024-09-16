package session

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"
)

type SessionTestSuite struct {
	suite.Suite
}

func (s *SessionTestSuite) SetupTest() {
	s.T().Setenv("CSRF_SECRET", "supersecretkey") // Set environment variable for the test
}

func (s *SessionTestSuite) TestGenerateCSRFToken() {
	// Simulate a session
	session := &Session{ID: "123456789"}

	// Generate a CSRF token
	csrfContent, err := session.GetCSRFToken()
	s.Require().NoError(err)

	// Split the token to retrieve the message and MAC
	parts := strings.Split(csrfContent, ".")
	s.Require().Len(parts, 2, "Invalid CSRF token format")
	messageMAC := []byte(parts[0])
	message := []byte(parts[1])

	// Verify the CSRF token
	s.True(VerifyCSRFToken(message, messageMAC), "CSRF token verification failed")
}

func TestSessionTestSuite(t *testing.T) {
	suite.Run(t, new(SessionTestSuite))
}
