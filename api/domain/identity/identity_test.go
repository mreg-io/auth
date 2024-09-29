package identity

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type IdentityTestSuite struct {
	suite.Suite
}

func (i *IdentityTestSuite) SetupTest() {
}

func (f *EmailTestSuite) TestIdentityETag() {
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

	identity := Identity{
		ID:              "c935b23d-6cb4-448a-814e-b42aec9ef6cf",
		State:           IDState(StateSuspended),
		FullName:        "John Doe",
		DisplayName:     "John",
		AvatarURL:       "https://example.com/avatar.jpg",
		Emails:          []Email{email},
		Timezone:        "HELL",
		CreateTime:      createTime,
		UpdateTime:      updatedAt,
		StateUpdateTime: verifiedAt,
		PasswordHash:    "hashedpassword",
	}

	// Act: calculate the ETag
	etag, err := identity.ETag()
	f.Require().NoError(err)
	// Assert: make sure the ETag is not empty and starts with 'W/"'
	f.Require().NotEmpty(etag, "ETag should not be empty")
	f.Require().Contains(etag, "W/\"", "ETag should be in the weak format")
}

func (f *EmailTestSuite) TestIdentityETag_NoID() {
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

	identity := Identity{
		State:           IDState(StateSuspended),
		FullName:        "John Doe",
		DisplayName:     "John",
		AvatarURL:       "https://example.com/avatar.jpg",
		Emails:          []Email{email},
		Timezone:        "HELL",
		CreateTime:      createTime,
		UpdateTime:      updatedAt,
		StateUpdateTime: verifiedAt,
		PasswordHash:    "hashedpassword",
	}

	// Act: calculate the ETag
	_, err = identity.ETag()
	f.Require().Error(err)
}

func (f *EmailTestSuite) TestIdentityETag_NoEmail() {
	createTime, err := time.Parse(time.UnixDate, "Wed Feb 28 11:06:39 UTC 1069")
	f.Require().NoError(err)

	verifiedAt, err := time.Parse(time.UnixDate, "Wed Feb 28 11:06:39 UTC 2069")
	f.Require().NoError(err)

	updatedAt, err := time.Parse(time.UnixDate, "Wed Feb 28 11:06:39 UTC 3069")
	f.Require().NoError(err)

	identity := Identity{
		ID:              "c935b23d-6cb4-448a-814e-b42aec9ef6cf",
		State:           IDState(StateSuspended),
		FullName:        "John Doe",
		DisplayName:     "John",
		AvatarURL:       "https://example.com/avatar.jpg",
		Timezone:        "HELL",
		CreateTime:      createTime,
		UpdateTime:      updatedAt,
		StateUpdateTime: verifiedAt,
		PasswordHash:    "hashedpassword",
	}

	// Act: calculate the ETag
	_, err = identity.ETag()
	f.Require().Equal(fmt.Errorf("identity must have at least one email").Error(), err.Error())
	// Assert: make sure the ETag is not empty and starts with 'W/"'
}

func (f *EmailTestSuite) TestIdentityETag_AdditionOfEmail() {
	createTime, err := time.Parse(time.UnixDate, "Wed Feb 28 11:06:39 UTC 1069")
	f.Require().NoError(err)

	verifiedAt, err := time.Parse(time.UnixDate, "Wed Feb 28 11:06:39 UTC 2069")
	f.Require().NoError(err)

	updatedAt, err := time.Parse(time.UnixDate, "Wed Feb 28 11:06:39 UTC 3069")
	f.Require().NoError(err)

	email1 := Email{
		Value:      "c935b23d-6cb4-448a-814e-b42aec9ef6cf",
		Verified:   true,
		VerifiedAt: verifiedAt,
		CreateTime: createTime,
		UpdateTime: updatedAt,
	}
	email2 := Email{
		Value:      "c8763c8763-6cb4-448a-814e-b42aec9ef6cf",
		Verified:   false,
		VerifiedAt: verifiedAt,
		CreateTime: createTime,
		UpdateTime: updatedAt,
	}

	identity1Email := Identity{
		ID:              "c935b23d-6cb4-448a-814e-b42aec9ef6cf",
		State:           IDState(StateSuspended),
		FullName:        "John Doe",
		DisplayName:     "John",
		AvatarURL:       "https://example.com/avatar.jpg",
		Emails:          []Email{email1},
		Timezone:        "HELL",
		CreateTime:      createTime,
		UpdateTime:      updatedAt,
		StateUpdateTime: verifiedAt,
		PasswordHash:    "hashedpassword",
	}

	identity2Email := Identity{
		ID:              "c935b23d-6cb4-448a-814e-b42aec9ef6cf",
		State:           IDState(StateSuspended),
		FullName:        "John Doe",
		DisplayName:     "John",
		AvatarURL:       "https://example.com/avatar.jpg",
		Emails:          []Email{email1, email2},
		Timezone:        "HELL",
		CreateTime:      createTime,
		UpdateTime:      updatedAt,
		StateUpdateTime: verifiedAt,
		PasswordHash:    "hashedpassword",
	}

	// Act: calculate the ETag
	etag1, err := identity1Email.ETag()
	f.Require().NoError(err)
	etag2, err := identity2Email.ETag()
	f.Require().NoError(err)
	f.Require().NotEqual(etag1, etag2, "ETags should be different for identities with different emails")
	// Assert: make sure the ETag is not empty and starts with 'W/"'
	f.Require().NotEmpty(etag1, "ETag should not be empty")
	f.Require().Contains(etag1, "W/\"", "ETag should be in the weak format")
	f.Require().NotEmpty(etag2, "ETag should not be empty")
	f.Require().Contains(etag2, "W/\"", "ETag should be in the weak format")
}

func TestIdentityEtag(t *testing.T) {
	suite.Run(t, new(IdentityTestSuite))
}
