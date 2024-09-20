package connect

import (
	"context"
	"net/http"
	"net/netip"
	"testing"
	"time"

	authConnect "buf.build/gen/go/mreg/protobuf/connectrpc/go/mreg/auth/v1alpha1/authv1alpha1connect"
	auth "buf.build/gen/go/mreg/protobuf/protocolbuffers/go/mreg/auth/v1alpha1"
	"connectrpc.com/connect"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"gitlab.mreg.io/my-registry/auth/domain/registration"
	"gitlab.mreg.io/my-registry/auth/domain/session"
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

type handlerTestSuite struct {
	suite.Suite
	mockService *mockRegistrationService
	handler     authConnect.RegistrationServiceHandler
}

func (h *handlerTestSuite) SetupSuite() {
	h.mockService = new(mockRegistrationService)
	h.handler = NewRegistrationHandler(h.mockService)
}

func (h *handlerTestSuite) TestCreateRegistrationFlow() {
	req := connect.NewRequest[auth.CreateRegistrationFlowRequest](&auth.CreateRegistrationFlowRequest{})
	UA := "pro-n-hub"
	IP := "192.0.2.43"
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
		Name:     "__Host-session_id",
		Value:    "c2e577de-2fbc-4fa4-8dcd-321a960ebb36",
		Expires:  sessionExpiresTime,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}

	cookie := res.Header()["Set-Cookie"][0]
	csrfCookie := res.Header()["Set-Cookie"][1]
	c, err := http.ParseSetCookie(csrfCookie)
	h.Require().NoError(err)
	h.Require().Equal(cookie, realCookie.String())
	h.Require().Equal(c.Expires, sessionExpiresTime)
	h.Require().NotEmpty(c.Value)
	call.Unset()
}

func (h *handlerTestSuite) TestCreateRegistrationFlow_WithoutIP() {
	req := connect.NewRequest[auth.CreateRegistrationFlowRequest](&auth.CreateRegistrationFlowRequest{})
	UA := "pro-n-hub"
	req.Header().Set("User-Agent", UA)
	ctx := context.Background()

	_, err := h.handler.CreateRegistrationFlow(ctx, req)
	// Assert proper call to repository
	h.Require().Error(err) // no ip address provided
}

func (h *handlerTestSuite) TestCreateRegistrationFlow_WithoutUA() {
	req := connect.NewRequest[auth.CreateRegistrationFlowRequest](&auth.CreateRegistrationFlowRequest{})
	IP := "192.0.2.43"
	req.Header().Set("X-Forwarded-For", IP)
	ctx := context.Background()

	_, err := h.handler.CreateRegistrationFlow(ctx, req)
	// Assert proper call to repository
	h.Require().Error(err) // no ip address provided
}

func TestHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(handlerTestSuite))
}
