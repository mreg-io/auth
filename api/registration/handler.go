package registration

import (
	"context"
	"log"

	authconnect "buf.build/gen/go/mreg/protobuf/connectrpc/go/mreg/auth/v1alpha1/authv1alpha1connect"
	auth "buf.build/gen/go/mreg/protobuf/protocolbuffers/go/mreg/auth/v1alpha1"
	"connectrpc.com/connect"
)

type handler struct {
	service Service
}

func NewHandler(service Service) authconnect.RegistrationServiceHandler {
	return &handler{
		service: service,
	}
}

func (h *handler) CreateRegistrationFlow(ctx context.Context, req *connect.Request[auth.CreateRegistrationFlowRequest]) (*connect.Response[auth.CreateRegistrationFlowResponse], error) {
	//TODO implement me

	_, err := h.service.CreateRegistrationFlow()
	if err != nil {
		log.Fatal("Failed to handle request:", err)
	}
	return nil, err
}

func (h *handler) CompleteRegistrationFlow(ctx context.Context, req *connect.Request[auth.CompleteRegistrationFlowRequest]) (*connect.Response[auth.CompleteRegistrationFlowResponse], error) {
	//TODO implement me
	panic("implement me")
}
