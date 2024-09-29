package identity

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type EmailTestSuite struct {
	suite.Suite
}

func (f *EmailTestSuite) SetupTest() {
}

func (f *EmailTestSuite) TestEmailETag() {
	createTime, err := time.Parse(time.UnixDate, "Wed Feb 28 11:06:39 UTC 1069")
	f.Require().NoError(err)
	verifiedAt, err := time.Parse(time.UnixDate, "Wed Feb 28 11:06:39 UTC 2069")
	f.Require().NoError(err)
	updatedAt, err := time.Parse(time.UnixDate, "Wed Feb 28 11:06:39 UTC 3069")
	f.Require().NoError(err)
	email := Email{
		Value:      "c935b23d-6cb4-448a-814e-b42aec9ef6cf",
		Verified:   true,
		VerifiedAt: verifiedAt,
		CreateTime: createTime,
		UpdateTime: updatedAt,
	}

	// Act: calculate the ETag
	etag, err := email.ETag()
	f.Require().NoError(err)
	// Assert: make sure the ETag is not empty and starts with 'W/"'
	f.Require().NotEmpty(etag, "ETag should not be empty")
	f.Require().Contains(etag, "W/\"", "ETag should be in the weak format")
}

func (f *EmailTestSuite) TestEmailETag_NoValue() {
	createTime, err := time.Parse(time.UnixDate, "Wed Feb 28 11:06:39 UTC 1069")
	f.Require().NoError(err)
	verifiedAt, err := time.Parse(time.UnixDate, "Wed Feb 28 11:06:39 UTC 2069")
	f.Require().NoError(err)
	updatedAt, err := time.Parse(time.UnixDate, "Wed Feb 28 11:06:39 UTC 3069")
	f.Require().NoError(err)
	email := Email{
		Verified:   true,
		VerifiedAt: verifiedAt,
		CreateTime: createTime,
		UpdateTime: updatedAt,
	}

	// Act: calculate the ETag
	_, err = email.ETag()
	f.Require().Error(err)
}

func (f *EmailTestSuite) TestEmailETag_NoVerified_NoErr() {
	createTime, err := time.Parse(time.UnixDate, "Wed Feb 28 11:06:39 UTC 1069")
	f.Require().NoError(err)
	updatedAt, err := time.Parse(time.UnixDate, "Wed Feb 28 11:06:39 UTC 3069")
	f.Require().NoError(err)
	email := Email{
		Value:      "c935b23d-6cb4-448a-814e-b42aec9ef6cf",
		Verified:   true,
		CreateTime: createTime,
		UpdateTime: updatedAt,
	}

	// Act: calculate the ETag
	etag, err := email.ETag()
	f.Require().NoError(err)
	// Assert: make sure the ETag is not empty and starts with 'W/"'
	f.Require().NotEmpty(etag, "ETag should not be empty")
	f.Require().Contains(etag, "W/\"", "ETag should be in the weak format")
}

func (f *EmailTestSuite) TestEmailETag_ChangeOfField() {
	createTime, err := time.Parse(time.UnixDate, "Wed Feb 28 11:06:39 UTC 1069")
	f.Require().NoError(err)
	updatedAt, err := time.Parse(time.UnixDate, "Wed Feb 28 11:06:39 UTC 3069")
	f.Require().NoError(err)
	email1 := Email{
		Value:      "c935b23d-6cb4-448a-814e-b42aec9ef6cf",
		Verified:   true,
		CreateTime: createTime,
		UpdateTime: updatedAt,
	}
	email2 := Email{
		Value:      "c935b23d-6cb4-448a-814e-b42aec9ef6cf",
		Verified:   false,
		CreateTime: createTime,
		UpdateTime: updatedAt,
	}
	// Act: calculate the ETag
	etag1, err := email1.ETag()
	f.Require().NoError(err)
	etag2, err := email2.ETag()
	f.Require().NoError(err)
	f.Require().NotEqual(etag1, etag2, "ETag should change when verified status changes")
}

func TestEmailEtag(t *testing.T) {
	suite.Run(t, new(EmailTestSuite))
}
