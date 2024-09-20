package connect

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/netip"
	"strings"

	authConnect "buf.build/gen/go/mreg/protobuf/connectrpc/go/mreg/auth/v1alpha1/authv1alpha1connect"
	auth "buf.build/gen/go/mreg/protobuf/protocolbuffers/go/mreg/auth/v1alpha1"
	"connectrpc.com/connect"

	"gitlab.mreg.io/my-registry/auth/service/registration"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type registrationHandler struct {
	registrationService registration.Service
}

func NewRegistrationHandler(registrationService registration.Service) authConnect.RegistrationServiceHandler {
	return &registrationHandler{registrationService}
}

func (r *registrationHandler) CreateRegistrationFlow(ctx context.Context, req *connect.Request[auth.CreateRegistrationFlowRequest]) (*connect.Response[auth.CreateRegistrationFlowResponse], error) {
	headers := req.Header()
	userAgent := headers.Get("User-Agent")
	xForwardedFor := headers.Get("X-Forwarded-For")
	var clientIP netip.Addr
	var err error
	if userAgent == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("missing User-Agent header"))
	}
	if xForwardedFor != "" {
		clientIPString, _, _ := strings.Cut(xForwardedFor, ",")
		clientIP, err = netip.ParseAddr(clientIPString)
		if err != nil {
			return nil, connect.NewError(connect.CodeInvalidArgument, err)
		}
	} else {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("cannot determine client IP address"))
	}

	flow, session, err := r.registrationService.CreateRegistrationFlow(ctx, clientIP, userAgent)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, errors.New("unable to create registration flow"))
	}
	eTag, err := flow.ETag()
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, errors.New("unable to generate etag in registration flow"))
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
		Name:     "__Host-session_id",
		Value:    session.ID,
		Expires:  session.ExpiresAt,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}
	csrfToken, err := session.GetCSRFToken()
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, errors.New("unable to generate CSRF token"))
	}
	csrfCookie := &http.Cookie{
		Name:     "csrf_token",
		Value:    csrfToken,
		Expires:  session.ExpiresAt,
		Path:     "/",
		Secure:   true,
		HttpOnly: false,
		SameSite: http.SameSiteStrictMode,
	}
	response := connect.NewResponse[auth.CreateRegistrationFlowResponse](res)
	response.Header().Add("Set-Cookie", cookie.String())
	response.Header().Add("Set-Cookie", csrfCookie.String())
	return response, nil
}

func (r *registrationHandler) CompleteRegistrationFlow(context.Context, *connect.Request[auth.CompleteRegistrationFlowRequest]) (*connect.Response[auth.CompleteRegistrationFlowResponse], error) {
	// TODO implement me
	panic("implement me")
}
