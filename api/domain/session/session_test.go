package session

import (
	"net/netip"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"gitlab.mreg.io/my-registry/auth/domain/identity"
)

type SessionTestSuite struct {
	suite.Suite
}

func (s *SessionTestSuite) SetupTest() {
}

func (s *SessionTestSuite) TestSessionETag() {
	createTime, err := time.Parse(time.UnixDate, "Wed Feb 28 11:06:39 UTC 1069")
	s.Require().NoError(err)
	verifiedAt, err := time.Parse(time.UnixDate, "Wed Feb 28 11:06:39 UTC 2069")
	s.Require().NoError(err)

	device := Device{
		ID:          "device-12345",
		IPAddress:   netip.MustParseAddr("192.0.2.1"),
		GeoLocation: "New York, USA",
		UserAgent:   "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
		SessionID:   "c935b23d-6cb4-448a-814e-b42aec9ef6cf",
	}
	UserIdentity := &identity.Identity{}
	session := Session{
		ID:                          "c935b23d-6cb4-448a-814e-b42aec9ef6cf",
		Active:                      true,
		AuthenticatorAssuranceLevel: 2,
		IssuedAt:                    createTime,
		ExpiresAt:                   verifiedAt.Add(24 * time.Hour),
		AuthenticatedAt:             verifiedAt,
		Devices:                     []Device{device},
		Identity:                    UserIdentity,
		ExpiryInterval:              2 * time.Hour,
	}

	// Act: calculate the ETag
	etag, err := session.ETag()
	s.Require().NoError(err)
	// Assert: make sure the ETag is not empty and starts with 'W/"'
	s.Require().NotEmpty(etag, "ETag should not be empty")
	s.Require().Contains(etag, "W/\"", "ETag should be in the weak format")
}

func (s *SessionTestSuite) TestSessionETag_NoID() {
	createTime, err := time.Parse(time.UnixDate, "Wed Feb 28 11:06:39 UTC 1069")
	s.Require().NoError(err)
	verifiedAt, err := time.Parse(time.UnixDate, "Wed Feb 28 11:06:39 UTC 2069")
	s.Require().NoError(err)

	device := Device{
		ID:          "device-12345",
		IPAddress:   netip.MustParseAddr("192.0.2.1"),
		GeoLocation: "New York, USA",
		UserAgent:   "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
		SessionID:   "c935b23d-6cb4-448a-814e-b42aec9ef6cf",
	}

	UserIdentity := &identity.Identity{}
	session := Session{
		Active:                      true,
		AuthenticatorAssuranceLevel: 2,
		IssuedAt:                    createTime,
		ExpiresAt:                   verifiedAt.Add(24 * time.Hour),
		AuthenticatedAt:             verifiedAt,
		Devices:                     []Device{device},
		Identity:                    UserIdentity,
		ExpiryInterval:              2 * time.Hour,
	}

	// Act: calculate the ETag
	_, err = session.ETag()
	s.Require().Error(err)
}

func (s *SessionTestSuite) TestSessionETag_NoVerified_NoDevice() {
	createTime, err := time.Parse(time.UnixDate, "Wed Feb 28 11:06:39 UTC 1069")
	s.Require().NoError(err)
	verifiedAt, err := time.Parse(time.UnixDate, "Wed Feb 28 11:06:39 UTC 2069")
	s.Require().NoError(err)
	UserIdentity := &identity.Identity{}
	session := Session{
		ID:                          "c935b23d-6cb4-448a-814e-b42aec9ef6cf",
		Active:                      true,
		AuthenticatorAssuranceLevel: 2,
		IssuedAt:                    createTime,
		ExpiresAt:                   verifiedAt.Add(24 * time.Hour),
		AuthenticatedAt:             verifiedAt,
		Identity:                    UserIdentity,
		ExpiryInterval:              2 * time.Hour,
	}

	// Act: calculate the ETag
	_, err = session.ETag()
	s.Require().Error(err)
	// Assert: make sure the ETag is not empty and starts with 'W/"'
}

func (s *SessionTestSuite) TestSessionETag_AddDevice() {
	createTime, err := time.Parse(time.UnixDate, "Wed Feb 28 11:06:39 UTC 1069")
	s.Require().NoError(err)
	verifiedAt, err := time.Parse(time.UnixDate, "Wed Feb 28 11:06:39 UTC 2069")
	s.Require().NoError(err)

	device1 := Device{
		ID:          "device-12345",
		IPAddress:   netip.MustParseAddr("192.0.2.1"),
		GeoLocation: "New York, USA",
		UserAgent:   "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
		SessionID:   "c935b23d-6cb4-448a-814e-b42aec9ef6cf",
	}
	device2 := Device{
		ID:          "device-Gy",
		IPAddress:   netip.MustParseAddr("192.0.34.1"),
		GeoLocation: "Taiwan, Thailand",
		UserAgent:   "LEEPOKI",
		SessionID:   "c935b23d-6cb4-448a-814e-b42aec9ef6cf",
	}

	session1 := Session{
		ID:                          "c935b23d-6cb4-448a-814e-b42aec9ef6cf",
		Active:                      true,
		AuthenticatorAssuranceLevel: 2,
		IssuedAt:                    createTime,
		ExpiresAt:                   verifiedAt.Add(24 * time.Hour),
		AuthenticatedAt:             verifiedAt,
		Devices:                     []Device{device1},
		Identity:                    nil,
		ExpiryInterval:              2 * time.Hour,
	}
	session2 := Session{
		ID:                          "c935b23d-6cb4-448a-814e-b42aec9ef6cf",
		Active:                      true,
		AuthenticatorAssuranceLevel: 2,
		IssuedAt:                    createTime,
		ExpiresAt:                   verifiedAt.Add(24 * time.Hour),
		AuthenticatedAt:             verifiedAt,
		Devices:                     []Device{device1, device2},
		Identity:                    nil,
		ExpiryInterval:              2 * time.Hour,
	}

	// Act: calculate the ETag
	etag1, err := session1.ETag()
	s.Require().NoError(err)
	etag2, err := session2.ETag()
	s.Require().NoError(err)
	// Assert: make sure the ETag is not empty and starts with 'W/"'
	s.Require().NotEqual(etag1, etag2, "ETags should be different")
	s.Require().NotEmpty(etag1, "ETag should not be empty")
	s.Require().Contains(etag1, "W/\"", "ETag should be in the weak format")
	s.Require().NotEmpty(etag2, "ETag should not be empty")
	s.Require().Contains(etag2, "W/\"", "ETag should be in the weak format")
}

func TestEmailEtag(t *testing.T) {
	suite.Run(t, new(SessionTestSuite))
}
