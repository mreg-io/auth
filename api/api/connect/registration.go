package connect

import (
	"context"

	authConnect "buf.build/gen/go/mreg/protobuf/connectrpc/go/mreg/auth/v1alpha1/authv1alpha1connect"
	auth "buf.build/gen/go/mreg/protobuf/protocolbuffers/go/mreg/auth/v1alpha1"
	"connectrpc.com/connect"

	"gitlab.mreg.io/my-registry/auth/service/registration"
)

type registrationHandler struct {
	registrationService registration.Service
}

func NewRegistrationHandler(registrationService registration.Service) authConnect.RegistrationServiceHandler {
	return &registrationHandler{registrationService}
}

func (r *registrationHandler) CreateRegistrationFlow(context.Context, *connect.Request[auth.CreateRegistrationFlowRequest]) (*connect.Response[auth.CreateRegistrationFlowResponse], error) {
	// TODO implement me
	panic("implement me")
}

func (r *registrationHandler) CompleteRegistrationFlow(context.Context, *connect.Request[auth.CompleteRegistrationFlowRequest]) (*connect.Response[auth.CompleteRegistrationFlowResponse], error) {
	// TODO implement me
	panic("implement me")
}
