import { createGrpcTransport } from "@connectrpc/connect-node";
import { RegistrationService } from "@buf/mreg_protobuf.connectrpc_es/mreg/auth/v1alpha1/registration_service_connect";
import { createPromiseClient } from "~/lib/connect";
import { stripHeaders } from "~/lib/interceptor";

const transport = createGrpcTransport({
  baseUrl: process.env.AUTH_API_URL,
  httpVersion: "2",
  interceptors: [
    stripHeaders(
      "grpc-accept-encoding",
      "grpc-encoding",
      "Content-Type",
      "Date",
    ),
  ],
});

export const registrationService = createPromiseClient(
  RegistrationService,
  transport,
);
