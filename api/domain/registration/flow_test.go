package registration

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type FlowTestSuite struct {
	suite.Suite
}

func (f *FlowTestSuite) SetupTest() {
}

func (f *FlowTestSuite) TestFlow_ETag() {
	// Arrange: create a Flow instance with predefined values
	IssuedAt, err := time.Parse(time.UnixDate, "Wed Feb 28 11:06:39 UTC 2069")
	f.Require().NoError(err)
	ExpiresAt, err := time.Parse(time.UnixDate, "Wed Feb 28 11:06:39 UTC 2069")
	f.Require().NoError(err)
	flow := Flow{
		FlowID:    "c935b23d-6cb4-448a-814e-b42aec9ef6cf",
		IssuedAt:  IssuedAt,
		ExpiresAt: ExpiresAt,
		SessionID: "269e873b-38ee-4904-bb4f-4207a33137df",
		Interval:  time.Hour,
	}

	// Act: calculate the ETag
	etag, err := flow.ETag()
	f.Require().NoError(err)
	fmt.Println("ETag:", etag)
	// Assert: make sure the ETag is not empty and starts with 'W/"'
	f.Require().NotEmpty(etag, "ETag should not be empty")
	f.Require().Contains(etag, "W/\"", "ETag should be in the weak format")
}

func (f *FlowTestSuite) TestFlow_ETag_NoExpireAt() {
	// Arrange: create a Flow instance with predefined values
	IssuedAt, err := time.Parse(time.UnixDate, "Wed Feb 28 11:06:39 UTC 2069")
	f.Require().NoError(err)
	flow := Flow{
		FlowID:    "c935b23d-6cb4-448a-814e-b42aec9ef6cf",
		IssuedAt:  IssuedAt,
		SessionID: "269e873b-38ee-4904-bb4f-4207a33137df",
		Interval:  time.Hour,
	}

	// Act: calculate the ETag
	etag, err := flow.ETag()
	fmt.Println("ETag:", etag)
	f.Require().Error(err)
}

func (f *FlowTestSuite) TestFlow_IsExpired() {
	// Test case where the flow is not expired
	notExpiredFlow := Flow{
		ExpiresAt: time.Now().Add(1 * time.Hour), // Expires 1 hour from now
	}

	f.False(notExpiredFlow.IsExpired(), "Expected flow to not be expired, but it is marked as expired")

	// Test case where the flow is expired
	expiredFlow := Flow{
		ExpiresAt: time.Now().Add(-1 * time.Hour), // Expired 1 hour ago
	}

	f.True(expiredFlow.IsExpired(), "Expected flow to be expired, but it is marked as not expired")

	// Test case where the expiration time is exactly now (boundary test)
	exactNowFlow := Flow{
		ExpiresAt: time.Now(),
	}

	f.True(exactNowFlow.IsExpired(), "Expected flow to be expired when expiration time is exactly now, but it is marked as not expired")
}

func TestSessionTestSuite(t *testing.T) {
	suite.Run(t, new(FlowTestSuite))
}
