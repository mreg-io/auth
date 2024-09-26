package connect

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/netip"
	"strings"

	"gitlab.mreg.io/my-registry/auth/infrastructure/cockroachdb"

	authConnect "buf.build/gen/go/mreg/protobuf/connectrpc/go/mreg/auth/v1alpha1/authv1alpha1connect"
	auth "buf.build/gen/go/mreg/protobuf/protocolbuffers/go/mreg/auth/v1alpha1"
	"connectrpc.com/connect"
	"gitlab.mreg.io/my-registry/auth/domain/identity"
	"gitlab.mreg.io/my-registry/auth/domain/registration"

	serviceRegistration "gitlab.mreg.io/my-registry/auth/service/registration"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type registrationHandler struct {
	registrationService serviceRegistration.Service
}

func NewRegistrationHandler(registrationService serviceRegistration.Service) authConnect.RegistrationServiceHandler {
	return &registrationHandler{registrationService}
}

func (r *registrationHandler) CreateRegistrationFlow(ctx context.Context, req *connect.Request[auth.CreateRegistrationFlowRequest]) (*connect.Response[auth.CreateRegistrationFlowResponse], error) {
	headers := req.Header()
	userAgent := headers.Get("User-Agent")
	xForwardedFor := headers.Get("X-Forwarded-For")
	var clientIP netip.Addr
	var err error
	if userAgent == "" {
		return nil, errorMissingHeader("User-Agent")
	}
	if xForwardedFor != "" {
		clientIPString, _, _ := strings.Cut(xForwardedFor, ",")
		clientIP, err = netip.ParseAddr(clientIPString)
		if err != nil {
			return nil, connect.NewError(connect.CodeInvalidArgument, err)
		}
	} else {
		return nil, errorMissingHeader("X-Forwarded-For")
	}

	flow, sessionData, err := r.registrationService.CreateRegistrationFlow(ctx, clientIP, userAgent)
	if err != nil {
		fmt.Printf("error creating registration flow: %v\n", err)
		return nil, connect.NewError(connect.CodeInternal, errors.New(""))
	}
	eTag, err := flow.ETag()
	if err != nil {
		fmt.Printf("error generating etag in registration flow: %v\n", err)
		return nil, connect.NewError(connect.CodeInternal, errors.New(""))
	}
	res := &auth.CreateRegistrationFlowResponse{
		RegistrationFlow: &auth.RegistrationFlow{
			Name:      fmt.Sprintf("registrationFlows/%s", flow.FlowID),
			FlowId:    flow.FlowID,
			IssuedAt:  timestamppb.New(flow.IssuedAt),
			ExpiresAt: timestamppb.New(flow.ExpiresAt),
			Etag:      eTag,
		},
	}
	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    sessionData.ID,
		Expires:  sessionData.ExpiresAt,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}
	response := connect.NewResponse[auth.CreateRegistrationFlowResponse](res)
	response.Header().Add("Set-Cookie", cookie.String())
	return response, nil
}

func (r *registrationHandler) CompleteRegistrationFlow(ctx context.Context, req *connect.Request[auth.CompleteRegistrationFlowRequest]) (*connect.Response[auth.CompleteRegistrationFlowResponse], error) {
	// Extract and verify CSRF token from the request
	cookie := req.Header().Get("Cookie")
	parsedCookies, err := http.ParseCookie(cookie)
	if err != nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("unauthenticated"))
	}
	var sessionID string
	for _, cookie := range parsedCookies {
		if cookie.Name == "session_id" {
			sessionID = cookie.Value
			break
		}
	}
	if sessionID == "" {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("unauthenticated"))
	}
	flow := &registration.Flow{
		SessionID: sessionID, // Use the extracted session ID
		Password:  req.Msg.GetRegistrationFlow().GetPassword().GetPassword(),
		Identity: &identity.Identity{
			Emails: []identity.Email{
				{
					Value: req.Msg.GetRegistrationFlow().GetTraits().GetEmail(),
				},
			},
		},
	}

	// Complete the registration flow, handling potential errors
	sessionData, err := r.registrationService.CompleteRegistrationFlow(ctx, flow)
	if err != nil {
		switch {
		case errors.Is(err, cockroachdb.ErrConstraint):
			return nil, internalError()
		case errors.Is(err, serviceRegistration.ErrEmailExists):
			return nil, errorEmailExist()
		case errors.Is(err, serviceRegistration.ErrInsecurePassword):
			return nil, errorInsecurePassword()
		case errors.Is(err, serviceRegistration.ErrSessionExpired):
			return nil, errorSessionExpired()
		case errors.Is(err, serviceRegistration.ErrFlowExpired):
			return nil, errorFlowExpired()
		default:
			return nil, internalError()
		}
	}
	identityData := flow.Identity
	// Generate ETag for the identity and address
	identityEtag, err := identityData.ETag()
	if err != nil {
		fmt.Printf("error generating identityEtag in registration flow: %v\n", err)
		return nil, internalError()
	}
	email := identityData.Emails[0]
	addressEtag, err := email.ETag()
	if err != nil {
		fmt.Printf("error generating addressEtag in registration flow: %v\n", err)
		return nil, internalError()
	}

	// Prepare the response message with identity data
	res := &auth.CompleteRegistrationFlowResponse{
		Identity: &auth.Identity{
			Name:       fmt.Sprintf("identities/%s", identityData.ID),
			IdentityId: identityData.ID,
			State:      auth.Identity_State(identityData.State),
			Addresses: []*auth.Address{
				{
					Name:       fmt.Sprintf("identities/%s/addresses/%s", identityData.ID, email.Value),
					Identity:   identityData.ID,
					Value:      email.Value,
					Via:        auth.Address_DeliveryMethod(1),
					Verified:   email.Verified,
					VerifiedAt: timestamppb.New(email.VerifiedAt),
					Etag:       addressEtag,
					CreateTime: timestamppb.New(email.CreateTime),
					UpdateTime: timestamppb.New(email.UpdateTime),
				},
			},
			Etag:            identityEtag,
			CreateTime:      timestamppb.New(identityData.CreateTime),
			UpdateTime:      timestamppb.New(identityData.UpdateTime),
			StateUpdateTime: timestamppb.New(identityData.StateUpdateTime),
		},
	}

	// Create response and set cookies for session and CSRF token
	response := connect.NewResponse[auth.CompleteRegistrationFlowResponse](res)
	setCookie := &http.Cookie{
		Name:     "session_id",
		Value:    sessionData.ID,
		Expires:  sessionData.ExpiresAt,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}
	response.Header().Add("Set-Cookie", setCookie.String())

	return response, nil
}
