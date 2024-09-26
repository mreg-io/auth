package connect

import (
	"errors"
	"fmt"

	"connectrpc.com/connect"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/protobuf/proto"
)

func wrapErrorAsConnectResponse(err *connect.Error, msg proto.Message) *connect.Error {
	detail, detailErr := connect.NewErrorDetail(msg)
	if detailErr == nil {
		err.AddDetail(detail) // return error if creating the detail fails
	}
	return err
}

func internalError() error {
	err := connect.NewError(connect.CodeInternal, errors.New("an internal server error occurred"))

	// Optionally, you can add additional error info using DebugInfo or ErrorInfo
	debugInfo := &errdetails.DebugInfo{
		Detail: "An unexpected error occurred",
	}
	return wrapErrorAsConnectResponse(err, debugInfo)
}

func errorMissingHeader(header string) error {
	err := connect.NewError(connect.CodeInvalidArgument, errors.New("missing header"))
	violation := &errdetails.BadRequest_FieldViolation{
		Field:       header,
		Description: fmt.Sprintf("The %s header is required but was not provided.", header),
	}

	// Create a BadRequest error detail message
	badRequest := &errdetails.BadRequest{
		FieldViolations: []*errdetails.BadRequest_FieldViolation{violation},
	}

	return wrapErrorAsConnectResponse(err, badRequest)
}

func errorEmailExist() error {
	err := connect.NewError(connect.CodeAlreadyExists, errors.New("email exist"))

	// Create a BadRequest error detail message
	badRequest := &errdetails.ResourceInfo{
		ResourceType: "User",
		ResourceName: "email",
		Description:  "The email address is already registered.",
	}
	return wrapErrorAsConnectResponse(err, badRequest)
}

func errorInsecurePassword() error {
	err := connect.NewError(connect.CodeInvalidArgument, errors.New("insecure password"))

	// Create a BadRequest error detail message
	info := &errdetails.ResourceInfo{
		ResourceType: "Credential",
		ResourceName: "password",
		Description:  "The password is not strong enough.",
	}
	return wrapErrorAsConnectResponse(err, info)
}

func errorSessionExpired() error {
	err := connect.NewError(connect.CodeUnauthenticated, errors.New("session expired"))

	// Create a ResourceInfo error detail message
	info := &errdetails.ResourceInfo{
		ResourceType: "Authentication",
		ResourceName: "session",
		Description:  "The user session has expired. Please log in again.",
	}
	return wrapErrorAsConnectResponse(err, info)
}

func errorFlowExpired() error {
	err := connect.NewError(connect.CodeUnauthenticated, errors.New("flow expired"))

	// Create a ResourceInfo error detail message
	info := &errdetails.ResourceInfo{
		ResourceType: "Authentication",
		ResourceName: "flow",
		Description:  "The registration flow has expired. Please refresh the page.",
	}
	return wrapErrorAsConnectResponse(err, info)
}
