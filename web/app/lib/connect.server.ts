import { createGrpcTransport } from "@connectrpc/connect-node";
import { RegistrationService } from "@buf/mreg_protobuf.connectrpc_es/mreg/auth/v1alpha1/registration_service_connect";
import { createPromiseClient } from "~/lib/connect";

const transport = createGrpcTransport({
  baseUrl: process.env.AUTH_API_URL,
  httpVersion: "2",
});

export const registrationService = createPromiseClient(
  RegistrationService,
  transport
);
