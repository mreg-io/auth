import { describe, expect, it } from "vitest";
import { createRouterTransport } from "@connectrpc/connect";
import { RegistrationService } from "@buf/mreg_protobuf.connectrpc_es/mreg/auth/v1alpha1/registration_service_connect";
import { CreateRegistrationFlowResponse } from "@buf/mreg_protobuf.bufbuild_es/mreg/auth/v1alpha1/registration_service_pb";
import { createPromiseClient } from "~/lib/connect";

describe("createPromiseClient", () => {
  it("should able to return Unary response", async () => {
    const sessionID = crypto.randomUUID();
    const setCookie = `session_id=${sessionID}; Path=/; Expires=Fri, 27 Sep 2024 18:00:46 GMT; HttpOnly; Secure; SameSite=Strict`;
    const mockResponse = new CreateRegistrationFlowResponse({
      registrationFlow: {
        name: "registrationFlows/01923634-d98c-9563-8c9e-3a676d49ac00",
        flowId: "01923634-d98c-9563-8c9e-3a676d49ac00",
        etag: "UWu5u//dU1PuukmRaEmO1RUNSN5NkUgXV/3gpUhMHow=.MDE5MjM2MzQtZDk4My05MjM1LTU2OWMtN2E1ZTc2NjcxOWUzIWJmMDFiODA3LWIwZTQtNDkxOS1hZGFiLTM5NTQzZDU5OTgwYg==",
      },
    });

    const mockTransport = createRouterTransport((router) => {
      router.service(RegistrationService, {
        createRegistrationFlow: (request, context) => {
          expect(request).toEqual({});

          context.responseHeader.set("Set-Cookie", setCookie);
          return mockResponse;
        },
      });
    });

    const registrationService = createPromiseClient(
      RegistrationService,
      mockTransport,
    );
    const { response, headers } =
      await registrationService.createRegistrationFlow({});
    expect(response).toEqual(mockResponse);
    expect(headers.getSetCookie()).toContain(setCookie);
  });
});
