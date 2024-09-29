package connect

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/netip"
	"testing"
	"time"
	_ "time/tzdata"

	"google.golang.org/genproto/googleapis/type/datetime"

	"gitlab.mreg.io/my-registry/auth/domain/identity"

	authConnect "buf.build/gen/go/mreg/protobuf/connectrpc/go/mreg/auth/v1alpha1/authv1alpha1connect"
	auth "buf.build/gen/go/mreg/protobuf/protocolbuffers/go/mreg/auth/v1alpha1"
	"connectrpc.com/connect"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"gitlab.mreg.io/my-registry/auth/domain/registration"
	"gitlab.mreg.io/my-registry/auth/domain/session"
	registrationService "gitlab.mreg.io/my-registry/auth/service/registration"
)

type mockRegistrationService struct {
	mock.Mock
}

func (m *mockRegistrationService) CreateRegistrationFlow(ctx context.Context, ipAddress netip.Addr, userAgent string) (*registration.Flow, *session.Session, error) {
	// Call the mock method and get the arguments
	args := m.Called(ctx, ipAddress, userAgent)

	// Extract the returned values from the mock call
	flow, _ := args.Get(0).(*registration.Flow)
	sessionModel, _ := args.Get(1).(*session.Session)
	err := args.Error(2)

	// Return the values
	return flow, sessionModel, err
}

func (m *mockRegistrationService) CompleteRegistrationFlow(ctx context.Context, f *registration.Flow, addr netip.Addr, s string) (*session.Session, error) {
	args := m.Called(ctx, f, addr, s)

	// Extract the returned values from the mock call
	sessionModel, _ := args.Get(0).(*session.Session)
	err := args.Error(1)
	return sessionModel, err
}

type handlerTestSuite struct {
	suite.Suite
	mockService *mockRegistrationService
	handler     authConnect.RegistrationServiceHandler
}

func (h *handlerTestSuite) SetupSuite() {
	h.mockService = new(mockRegistrationService)
	h.handler = NewRegistrationHandler(h.mockService)
}

var (
	filledEmail    = "test@example.com"
	filledPassword = "strong_password_123"
	UA             = "pro-n-hub"
	IP             = "192.0.2.43"
	preSessionID   = "a02bdf6d-a87b-439b-9140-87287b8d0a96"
	timezone       = "America/New_York"
)

func (h *handlerTestSuite) TestCreateRegistrationFlow() {
	req := connect.NewRequest[auth.CreateRegistrationFlowRequest](&auth.CreateRegistrationFlowRequest{})
	req.Header().Set("User-Agent", UA)
	req.Header().Set("X-Forwarded-For", IP)
	ctx := context.Background()

	clientIP, err := netip.ParseAddr(IP)
	h.Require().NoError(err)
	issuedTime, err := time.Parse(time.UnixDate, "Wed Feb 25 11:06:39 UTC 1069")
	h.Require().NoError(err)
	flowExpiresTime, err := time.Parse(time.UnixDate, "Wed Feb 28 11:06:39 UTC 2069")
	h.Require().NoError(err)
	interval, err := time.ParseDuration("24h")
	h.Require().NoError(err)
	sessionExpiresTime, err := time.Parse(time.UnixDate, "Wed Feb 28 11:06:39 UTC 3069")
	h.Require().NoError(err)
	call := h.mockService.
		On("CreateRegistrationFlow", ctx, clientIP, UA).
		Run(func(_ mock.Arguments) {
		}).
		Return(&registration.Flow{
			FlowID:    "0dc909cb-8c0f-4dc8-98b9-e82d77eb9d79", // Replace with actual return value
			IssuedAt:  issuedTime,
			ExpiresAt: flowExpiresTime,
			SessionID: "c2e577de-2fbc-4fa4-8dcd-321a960ebb36",
			Interval:  interval,
		}, &session.Session{ID: "c2e577de-2fbc-4fa4-8dcd-321a960ebb36", ExpiresAt: sessionExpiresTime}, nil).
		Once()

	res, err := h.handler.CreateRegistrationFlow(ctx, req)
	h.Require().NoError(err)
	h.mockService.AssertExpectations(h.T())
	// Access the underlying message from the response
	message := res.Msg

	h.Require().Equal("registrationFlows/0dc909cb-8c0f-4dc8-98b9-e82d77eb9d79", message.GetRegistrationFlow().GetName())
	h.Require().Equal("0dc909cb-8c0f-4dc8-98b9-e82d77eb9d79", message.GetRegistrationFlow().GetFlowId())
	h.Require().Equal(message.GetRegistrationFlow().GetIssuedAt().AsTime(), issuedTime)
	h.Require().Equal(message.GetRegistrationFlow().GetExpiresAt().AsTime(), flowExpiresTime)
	h.Require().NotEmpty(message.GetRegistrationFlow().GetEtag())
	realCookie := &http.Cookie{
		Name:     "session_id",
		Value:    "c2e577de-2fbc-4fa4-8dcd-321a960ebb36",
		Expires:  sessionExpiresTime,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}

	cookie := res.Header()["Set-Cookie"][0]
	h.Require().NoError(err)
	h.Require().Equal(cookie, realCookie.String())
	call.Unset()
}

func (h *handlerTestSuite) TestCreateRegistrationFlow_WithoutIP() {
	req := connect.NewRequest[auth.CreateRegistrationFlowRequest](&auth.CreateRegistrationFlowRequest{})
	req.Header().Set("User-Agent", UA)
	ctx := context.Background()

	_, err := h.handler.CreateRegistrationFlow(ctx, req)
	// Assert proper call to repository
	h.Require().Error(err) // no ip address provided
}

func (h *handlerTestSuite) TestCreateRegistrationFlow_WithoutUA() {
	req := connect.NewRequest[auth.CreateRegistrationFlowRequest](&auth.CreateRegistrationFlowRequest{})
	req.Header().Set("X-Forwarded-For", IP)
	ctx := context.Background()

	_, err := h.handler.CreateRegistrationFlow(ctx, req)
	// Assert proper call to repository
	h.Require().Error(err) // no ip address provided
}

func (h *handlerTestSuite) TestCompleteRegistrationFlow() {
	cTimezone := &datetime.TimeZone{
		Id:      timezone, // You can set the ID to the timezone string
		Version: "1.0",    // Set a default version or your desired value
	}
	req := connect.NewRequest[auth.CompleteRegistrationFlowRequest](&auth.CompleteRegistrationFlowRequest{
		RegistrationFlow: &auth.RegistrationFlow{
			Traits: &auth.IdentityTraits{
				Email:    &filledEmail,
				Timezone: cTimezone,
			},
			Credential: &auth.RegistrationFlow_Password{
				Password: &auth.Password{
					Password: filledPassword,
				},
			},
		},
	})
	cookie := &http.Cookie{
		Name:  "session_id",
		Value: preSessionID, // Your session ID value

	}
	req.Header().Add("Cookie", cookie.String())
	req.Header().Set("User-Agent", UA)
	req.Header().Set("X-Forwarded-For", IP)

	clientIP, err := netip.ParseAddr(IP)
	h.Require().NoError(err)

	ctx := context.Background()
	flow := &registration.Flow{
		SessionID: preSessionID, // Use the extracted session ID
		Password:  req.Msg.GetRegistrationFlow().GetPassword().GetPassword(),
		Identity: &identity.Identity{
			Emails: []identity.Email{
				{
					Value: req.Msg.GetRegistrationFlow().GetTraits().GetEmail(),
				},
			},
			Timezone: timezone,
		},
	}
	// Mock CSRF verification

	// Mock CompleteRegistrationFlow
	CreateTime, err := time.Parse(time.UnixDate, "Wed Feb 25 11:06:39 UTC 1069")
	h.Require().NoError(err)
	sessionExpiresTime, err := time.Parse(time.UnixDate, "Wed Feb 28 11:06:39 UTC 2069")
	h.Require().NoError(err)

	identityID := "IamBatMan"
	newSessionID := "c2e577de-2fbc-4fa4-8dcd-321a960ebb36"
	call1 := h.mockService.
		On("CompleteRegistrationFlow", ctx, flow, clientIP, UA).
		Run(func(args mock.Arguments) {
			flow := args.Get(1).(*registration.Flow)
			identityData := flow.Identity
			identityData.ID = identityID
			identityData.State = identity.StateActive
			identityData.Emails[0].Value = filledEmail
			identityData.Emails[0].Verified = false
			identityData.Emails[0].CreateTime = CreateTime
			identityData.Emails[0].UpdateTime = CreateTime
			identityData.CreateTime = CreateTime
			identityData.UpdateTime = CreateTime
			identityData.StateUpdateTime = CreateTime
		}).
		Return(&session.Session{
			ID:        newSessionID,
			ExpiresAt: sessionExpiresTime,
		}, nil).Once()

	res, err := h.handler.CompleteRegistrationFlow(ctx, req)
	h.Require().NoError(err)
	h.mockService.AssertExpectations(h.T())

	// Validate response
	message := res.Msg
	h.Require().Equal(fmt.Sprintf("identities/%s", identityID), message.GetIdentity().GetName())
	h.Require().Equal(identityID, message.GetIdentity().GetIdentityId())
	h.Require().Equal(fmt.Sprintf("identities/%s/addresses/%s", identityID, filledEmail), message.GetIdentity().GetAddresses()[0].GetName())
	h.Require().Equal(filledEmail, message.GetIdentity().GetAddresses()[0].GetValue())
	h.Require().False(message.GetIdentity().GetAddresses()[0].GetVerified())
	h.Require().Equal(CreateTime, message.GetIdentity().GetCreateTime().AsTime())
	h.Require().Equal(CreateTime, message.GetIdentity().GetUpdateTime().AsTime())
	h.Require().Equal(CreateTime, message.GetIdentity().GetStateUpdateTime().AsTime())
	// Validate cookies
	setCookies := res.Header()["Set-Cookie"]

	sessionCookie := setCookies[0]
	// Validate session cookie
	sessionC, err := http.ParseSetCookie(sessionCookie)
	h.Require().NoError(err)
	h.Require().Equal("session_id", sessionC.Name)
	h.Require().Equal(newSessionID, sessionC.Value)
	h.Require().Equal(sessionExpiresTime, sessionC.Expires)
	h.Require().Equal("/", sessionC.Path)
	h.Require().True(sessionC.Secure)
	h.Require().True(sessionC.HttpOnly)
	h.Require().Equal(http.SameSiteStrictMode, sessionC.SameSite)

	call1.Unset()
}

func (h *handlerTestSuite) TestCompleteRegistrationFlow_NoCookie() {
	cTimezone := &datetime.TimeZone{
		Id:      timezone, // You can set the ID to the timezone string
		Version: "1.0",    // Set a default version or your desired value
	}
	req := connect.NewRequest[auth.CompleteRegistrationFlowRequest](&auth.CompleteRegistrationFlowRequest{
		RegistrationFlow: &auth.RegistrationFlow{
			Traits: &auth.IdentityTraits{
				Email:    &filledEmail,
				Timezone: cTimezone,
			},
			Credential: &auth.RegistrationFlow_Password{
				Password: &auth.Password{
					Password: filledPassword,
				},
			},
		},
	})
	req.Header().Set("User-Agent", UA)
	req.Header().Set("X-Forwarded-For", IP)

	clientIP, err := netip.ParseAddr(IP)
	h.Require().NoError(err)

	ctx := context.Background()
	flow := &registration.Flow{
		SessionID: preSessionID, // Use the extracted session ID
		Password:  req.Msg.GetRegistrationFlow().GetPassword().GetPassword(),
		Identity: &identity.Identity{
			Emails: []identity.Email{
				{
					Value: req.Msg.GetRegistrationFlow().GetTraits().GetEmail(),
				},
			},
			Timezone: timezone,
		},
	}
	// Mock CSRF verification

	call1 := h.mockService.
		On("CompleteRegistrationFlow", ctx, flow, clientIP, UA).
		Run(func(_ mock.Arguments) {
		}).
		Return(&session.Session{}, nil).Once()
	_, err = h.handler.CompleteRegistrationFlow(ctx, req)
	h.Require().Equal(connect.NewError(connect.CodeUnauthenticated, errors.New("unauthenticated")).Error(), err.Error())
	call1.Unset()
}

func (h *handlerTestSuite) TestCompleteRegistrationFlow_NoCookieWithNameSessionID() {
	cTimezone := &datetime.TimeZone{
		Id:      timezone, // You can set the ID to the timezone string
		Version: "1.0",    // Set a default version or your desired value
	}
	req := connect.NewRequest[auth.CompleteRegistrationFlowRequest](&auth.CompleteRegistrationFlowRequest{
		RegistrationFlow: &auth.RegistrationFlow{
			Traits: &auth.IdentityTraits{
				Email:    &filledEmail,
				Timezone: cTimezone,
			},
			Credential: &auth.RegistrationFlow_Password{
				Password: &auth.Password{
					Password: filledPassword,
				},
			},
		},
	})
	cookie := &http.Cookie{
		Name:  "session_id123",
		Value: preSessionID, // Your session ID value

	}
	req.Header().Add("Cookie", cookie.String())
	req.Header().Set("User-Agent", UA)
	req.Header().Set("X-Forwarded-For", IP)

	clientIP, err := netip.ParseAddr(IP)
	h.Require().NoError(err)

	ctx := context.Background()
	flow := &registration.Flow{
		SessionID: preSessionID, // Use the extracted session ID
		Password:  req.Msg.GetRegistrationFlow().GetPassword().GetPassword(),
		Identity: &identity.Identity{
			Emails: []identity.Email{
				{
					Value: req.Msg.GetRegistrationFlow().GetTraits().GetEmail(),
				},
			},
			Timezone: timezone,
		},
	}
	// Mock CSRF verification

	call1 := h.mockService.
		On("CompleteRegistrationFlow", ctx, flow, clientIP, UA).
		Run(func(_ mock.Arguments) {
		}).
		Return(&session.Session{}, nil).Once()

	_, err = h.handler.CompleteRegistrationFlow(ctx, req)
	h.Require().Equal(connect.NewError(connect.CodeUnauthenticated, errors.New("unauthenticated")).Error(), err.Error())
	call1.Unset()
}

func (h *handlerTestSuite) TestCompleteRegistrationFlow_ServiceError() {
	cTimezone := &datetime.TimeZone{
		Id:      timezone, // You can set the ID to the timezone string
		Version: "1.0",    // Set a default version or your desired value
	}
	req := connect.NewRequest[auth.CompleteRegistrationFlowRequest](&auth.CompleteRegistrationFlowRequest{
		RegistrationFlow: &auth.RegistrationFlow{
			Traits: &auth.IdentityTraits{
				Email:    &filledEmail,
				Timezone: cTimezone,
			},
			Credential: &auth.RegistrationFlow_Password{
				Password: &auth.Password{
					Password: filledPassword,
				},
			},
		},
	})
	cookie := &http.Cookie{
		Name:  "session_id",
		Value: preSessionID, // Your session ID value

	}
	req.Header().Add("Cookie", cookie.String())
	req.Header().Set("User-Agent", UA)
	req.Header().Set("X-Forwarded-For", IP)

	clientIP, err := netip.ParseAddr(IP)
	h.Require().NoError(err)

	ctx := context.Background()
	flow := &registration.Flow{
		SessionID: preSessionID, // Use the extracted session ID
		Password:  req.Msg.GetRegistrationFlow().GetPassword().GetPassword(),
		Identity: &identity.Identity{
			Emails: []identity.Email{
				{
					Value: req.Msg.GetRegistrationFlow().GetTraits().GetEmail(),
				},
			},
			Timezone: timezone,
		},
	}
	// Mock CSRF verification

	// Mock CompleteRegistrationFlow
	call1 := h.mockService.
		On("CompleteRegistrationFlow", ctx, flow, clientIP, UA).
		Run(func(_ mock.Arguments) {
		}).
		Return(nil, errors.New("internal")).Once()

	_, err = h.handler.CompleteRegistrationFlow(ctx, req)
	h.Require().Equal(err.Error(), internalError().Error())
	h.mockService.AssertExpectations(h.T())
	call1.Unset()

	call2 := h.mockService.
		On("CompleteRegistrationFlow", ctx, flow, clientIP, UA).
		Run(func(_ mock.Arguments) {
		}).
		Return(nil, registrationService.ErrEmailExists).Once()

	_, err = h.handler.CompleteRegistrationFlow(ctx, req)
	h.Require().Equal(err.Error(), errorEmailExist().Error())
	h.mockService.AssertExpectations(h.T())
	call2.Unset()
}

func TestHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(handlerTestSuite))
}
