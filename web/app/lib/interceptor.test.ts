import { describe, expect, it } from "vitest";
import { CreateRegistrationFlowResponse } from "@buf/mreg_protobuf.bufbuild_es/mreg/auth/v1alpha1/registration_service_pb";
import { createRouterTransport } from "@connectrpc/connect";
import { RegistrationService } from "@buf/mreg_protobuf.connectrpc_es/mreg/auth/v1alpha1/registration_service_connect";
import { createPromiseClient } from "~/lib/connect";
import { stripHeaders } from "~/lib/interceptor";

describe("stripHeaders", () => {
  const mockHeaders = new Headers({
    "set-cookie":
      "session_id=01923644-5726-3de2-f6d9-c90722cefdde; Path=/; Expires=Sat, 28 Sep 2024 03:35:01 GMT; HttpOnly; Secure; SameSite=Strict",
    date: "Sat, 28 Sep 2024 01:35:01 GMT",
  });
  const mockResponse = new CreateRegistrationFlowResponse({
    registrationFlow: {
      name: "registrationFlows/01923634-d98c-9563-8c9e-3a676d49ac00",
      flowId: "01923634-d98c-9563-8c9e-3a676d49ac00",
      etag: "UWu5u//dU1PuukmRaEmO1RUNSN5NkUgXV/3gpUhMHow=.MDE5MjM2MzQtZDk4My05MjM1LTU2OWMtN2E1ZTc2NjcxOWUzIWJmMDFiODA3LWIwZTQtNDkxOS1hZGFiLTM5NTQzZDU5OTgwYg==",
    },
  });

  it("should correctly strip specified headers", async () => {
    const mockTransport = createRouterTransport(
      (router) => {
        router.service(RegistrationService, {
          createRegistrationFlow: (request, context) => {
            expect(request).toEqual({});

            for (const [name, value] of mockHeaders.entries()) {
              context.responseHeader.set(name, value);
            }
            return mockResponse;
          },
        });
      },
      { transport: { interceptors: [stripHeaders("Set-Cookie")] } },
    );
    const registrationService = createPromiseClient(
      RegistrationService,
      mockTransport,
    );

    const { response, headers } =
      await registrationService.createRegistrationFlow({});
    expect(response).toEqual(mockResponse);
    expect(headers.get("Date")).toBe(mockHeaders.get("Date"));
    expect(headers.has("Set-Cookie")).toBeFalsy();
  });

  it("should ignore unknown headers", async () => {
    const mockTransport = createRouterTransport(
      (router) => {
        router.service(RegistrationService, {
          createRegistrationFlow: (request, context) => {
            expect(request).toEqual({});

            for (const [name, value] of mockHeaders.entries()) {
              context.responseHeader.set(name, value);
            }
            return mockResponse;
          },
        });
      },
      { transport: { interceptors: [stripHeaders("Content-Type")] } },
    );
    const registrationService = createPromiseClient(
      RegistrationService,
      mockTransport,
    );
    const { response, headers } =
      await registrationService.createRegistrationFlow({});
    expect(response).toEqual(mockResponse);
    expect(headers.get("Date")).toBe(mockHeaders.get("Date"));
    expect(headers.getSetCookie()).toStrictEqual(mockHeaders.getSetCookie());
  });
});
