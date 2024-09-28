import { describe, expect, it } from "vitest";
import {
  generateCSRFToken,
  generateCSRFTokenFromHeaders,
  isCSRFTokenValid,
} from "~/lib/csrf.server";

describe("generateCSRFToken", () => {
  it("should generate valid CSRF token", () => {
    const sessionID = crypto.randomUUID();
    const csrfToken = generateCSRFToken(sessionID);
    expect(isCSRFTokenValid(sessionID, csrfToken)).toBeTruthy();
  });
});

describe("generateCSRFTokenFromHeaders", () => {
  it("should generate valid CSRF token from session_id Set-Cookie header", () => {
    const sessionID = crypto.randomUUID();
    const setCookie = `session_id=${sessionID}; Path=/; Expires=Fri, 27 Sep 2024 18:00:46 GMT; HttpOnly; Secure; SameSite=Strict`;
    const headers = new Headers({ "Set-Cookie": setCookie });
    const csrfToken = generateCSRFTokenFromHeaders(headers);
    expect(isCSRFTokenValid(sessionID, csrfToken)).toBeTruthy();
  });
  it("should return empty string if session_id is missing", () => {
    expect(generateCSRFTokenFromHeaders(new Headers())).toBe("");
  });
});

describe("isCSRFTokenValid", () => {
  it("should return true if CSRF token valid", () => {
    const sessionID = crypto.randomUUID();
    const csrfToken = generateCSRFToken(sessionID);
    expect(isCSRFTokenValid(sessionID, csrfToken)).toBeTruthy();
  });
  it("should return false if invalid CSRF token", () => {
    const sessionID1 = crypto.randomUUID();
    const sessionID2 = crypto.randomUUID();
    const csrfToken = generateCSRFToken(sessionID1);
    expect(isCSRFTokenValid(sessionID2, csrfToken)).toBeFalsy();
  });
});
