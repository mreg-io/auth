package main

import (
	"net/http"

	"buf.build/gen/go/mreg/protobuf/connectrpc/go/mreg/auth/v1alpha1/authv1alpha1connect"
	"connectrpc.com/grpcreflect"

	"gitlab.mreg.io/my-registry/auth/registration"
	"gitlab.mreg.io/my-registry/auth/session"
)

var (
	sessionRepository session.Repository
	sessionService    session.Service
)

var (
	registrationRepository registration.Repository
	registrationService    registration.Service
	registrationHandler    authv1alpha1connect.RegistrationServiceHandler
)

func bootstrap() {
	sessionRepository = session.NewRepository(conn)
	sessionService = session.NewService(sessionRepository)

	registrationRepository = registration.NewRepository(conn)
	registrationService = registration.NewService(registrationRepository, sessionService)
	registrationHandler = registration.NewHandler(registrationService)
}

func registerHandlers(mux *http.ServeMux) {
	// gRPC reflection
	reflector := grpcreflect.NewStaticReflector(
		"mreg.auth.v1alpha1.RegistrationService",
	)
	mux.Handle(grpcreflect.NewHandlerV1(reflector))
	mux.Handle(grpcreflect.NewHandlerV1Alpha(reflector))

	// Register service handlers
	mux.Handle(authv1alpha1connect.NewRegistrationServiceHandler(registrationHandler))
}
