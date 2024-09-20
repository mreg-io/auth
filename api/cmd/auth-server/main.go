package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	authConnect "buf.build/gen/go/mreg/protobuf/connectrpc/go/mreg/auth/v1alpha1/authv1alpha1connect"
	"connectrpc.com/grpcreflect"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"gitlab.mreg.io/my-registry/auth/api/connect"
	"gitlab.mreg.io/my-registry/auth/infrastructure/cockroachdb"
	"gitlab.mreg.io/my-registry/auth/service/registration"
)

const DatabaseURLEnvName string = "DATABASE_URL"

func main() {
	// Application-wide context
	ctx := context.Background()

	// Create CockroachDB pgx pool
	connString, ok := os.LookupEnv(DatabaseURLEnvName)
	if !ok {
		log.Fatalf("Cannot find %s environement variable\n", DatabaseURLEnvName)
	}
	pool := cockroachdb.NewPgxPool(ctx, connString)
	defer pool.Close()

	// Initialize repositories
	sessionRepository := cockroachdb.NewSessionRepository(pool)
	registrationFlowRepository := cockroachdb.NewRegistrationRepository(pool)

	// Initialize services
	registrationService := registration.NewService(sessionRepository, registrationFlowRepository)

	// Initialize handlers
	registrationHandler := connect.NewRegistrationHandler(registrationService)

	// Create ConnectRPC server
	mux := http.NewServeMux()
	reflector := grpcreflect.NewStaticReflector(
		"mreg.auth.v1alpha1.RegistrationService",
	)
	mux.Handle(grpcreflect.NewHandlerV1(reflector))
	mux.Handle(grpcreflect.NewHandlerV1Alpha(reflector))
	mux.Handle(authConnect.NewRegistrationServiceHandler(registrationHandler))

	if err := http.ListenAndServe(
		"0.0.0.0:8080",
		h2c.NewHandler(mux, &http2.Server{}),
	); err != nil {
		panic(fmt.Sprintf("Server failed: %v", err))
	}
}
