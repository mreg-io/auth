package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	authConnect "buf.build/gen/go/mreg/protobuf/connectrpc/go/mreg/auth/v1alpha1/authv1alpha1connect"
	"connectrpc.com/connect"
	"connectrpc.com/grpcreflect"
	"connectrpc.com/validate"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	apiConnect "gitlab.mreg.io/my-registry/auth/api/connect"
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
	identityRepository := cockroachdb.NewIdentityRepository(pool)

	// Initialize services
	registrationService := registration.NewService(sessionRepository, registrationFlowRepository, identityRepository)

	// Initialize handlers
	registrationHandler := apiConnect.NewRegistrationHandler(registrationService)

	// Create ConnectRPC server
	mux := http.NewServeMux()
	reflector := grpcreflect.NewStaticReflector(
		"mreg.auth.v1alpha1.RegistrationService",
	)
	mux.Handle(grpcreflect.NewHandlerV1(reflector))
	mux.Handle(grpcreflect.NewHandlerV1Alpha(reflector))

	interceptor, err := validate.NewInterceptor()
	if err != nil {
		panic(fmt.Sprintf("NewInterceptor failed: %v", err))
	}

	mux.Handle(authConnect.NewRegistrationServiceHandler(registrationHandler, connect.WithInterceptors(interceptor)))
	server := &http.Server{
		Addr:           "0.0.0.0:8080",
		Handler:        h2c.NewHandler(mux, &http2.Server{}), // Enable HTTP/2 over cleartext (H2C)
		MaxHeaderBytes: 8192,                                 // Limit header size to 8 KB
		ReadTimeout:    2 * time.Second,
		WriteTimeout:   5 * time.Second,
	}

	// Start the server
	if err := server.ListenAndServe(); err != nil {
		panic(fmt.Sprintf("Server failed: %v", err))
	}
}
