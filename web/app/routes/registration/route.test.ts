import { beforeEach, describe, expect, it, vi } from "vitest";
import { loader } from "~/routes/registration/route";
import * as services from "~/lib/connect.server";
import { createPromiseClient } from "~/lib/connect";
import { RegistrationService } from "@buf/mreg_protobuf.connectrpc_es/mreg/auth/v1alpha1/registration_service_connect";
import {
  Code,
  ConnectError,
  createRouterTransport,
  ServiceImpl,
} from "@connectrpc/connect";
import { CreateRegistrationFlowResponse } from "@buf/mreg_protobuf.bufbuild_es/mreg/auth/v1alpha1/registration_service_pb";
import { isCSRFTokenValid } from "~/lib/csrf.server";

const mockRegistrationService = (
  implementation: Partial<ServiceImpl<typeof RegistrationService>>,
) => {
  vi.spyOn(services, "registrationService", "get").mockReturnValue(
    createPromiseClient(
      RegistrationService,
      createRouterTransport((router) => {
        router.service(RegistrationService, implementation);
      }),
    ),
  );
};

const mockResponse = new CreateRegistrationFlowResponse({
  registrationFlow: {
    name: "registrationFlows/01923634-d98c-9563-8c9e-3a676d49ac00",
    flowId: "01923634-d98c-9563-8c9e-3a676d49ac00",
    etag: "UWu5u//dU1PuukmRaEmO1RUNSN5NkUgXV/3gpUhMHow=.MDE5MjM2MzQtZDk4My05MjM1LTU2OWMtN2E1ZTc2NjcxOWUzIWJmMDFiODA3LWIwZTQtNDkxOS1hZGFiLTM5NTQzZDU5OTgwYg==",
  },
});

describe("loader", () => {
  interface TestContext {
    sessionID: string;
    setCookie: string;
  }

  beforeEach<TestContext>((context) => {
    context.sessionID = crypto.randomUUID();
    context.setCookie = `session_id=${context.sessionID}; Path=/; Expires=Fri, 27 Sep 2024 18:00:46 GMT; HttpOnly; Secure; SameSite=Strict`;
  });

  it<TestContext>("should return registration flow with Set-Cookie header and CSRF token", async ({
    sessionID,
    setCookie,
  }) => {
    mockRegistrationService({
      createRegistrationFlow: (request, context) => {
        expect(request).toEqual({});
        expect(context.requestHeader.has("X-Forwarded-For"));
        expect(context.requestHeader.has("User-Agent"));

        context.responseHeader.set("Set-Cookie", setCookie);
        return mockResponse;
      },
    });

    const response = await loader({
      request: new Request("http://localhost:3000/registration"),
      params: {},
      context: {},
    });

    expect(response.data.response).toEqual(mockResponse);
    expect(isCSRFTokenValid(sessionID, response.data.csrfToken)).toBeTruthy();
    expect(new Headers(response.init?.headers).getSetCookie()).toContain(
      setCookie,
    );
  });

  it("should throw error if upstream error", async () => {
    const error = new ConnectError(
      "missing header",
      Code.InvalidArgument,
      undefined,
    );
    mockRegistrationService({
      createRegistrationFlow: () => {
        throw error;
      },
    });

    const response = loader({
      request: new Request("http://localhost:3000/registration"),
      params: {},
      context: {},
    });
    await expect(response).rejects.toThrowError(error);
  });
});
