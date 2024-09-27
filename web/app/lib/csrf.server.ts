import { randomUUID, createHmac } from "node:crypto";
import { parseString } from "set-cookie-parser";

const CSRF_TOKEN_SECRET = process.env.CSRF_TOKEN_SECRET;

export function generateCSRFToken(session: string | Headers) {
  switch (typeof session) {
    case "string":
      return generateTokenFromID(session);
    case "object":
      return generateTokenFromHeaders(session);
    default:
      return "";
  }
}

function generateTokenFromID(id: string) {
  const message = `${id}!${randomUUID()}`;
  const hmac = createHmac("sha256", CSRF_TOKEN_SECRET);
  hmac.update(message);
  return `${hmac.digest("base64")}.${Buffer.from(message).toString("base64")}`;
}

function generateTokenFromHeaders(headers: Headers) {
  let csrfToken = "";
  for (const setCookie of headers.getSetCookie()) {
    const cookie = parseString(setCookie);
    if (cookie.name === "session_id") {
      csrfToken = generateTokenFromID(cookie.value);
    }
  }
  return csrfToken;
}

export function isCSRFTokenValid(sessionID: string, csrfToken: string) {
  const [digest, base64Message] = csrfToken.split(".");
  if (!digest || !base64Message) return false;
  const message = Buffer.from(base64Message, "base64").toString("ascii");
  const [tokenSessionID] = message.split("!");
  if (tokenSessionID !== sessionID) return false;

  const hmac = createHmac("sha256", CSRF_TOKEN_SECRET);
  hmac.update(message);
  return hmac.digest("base64") === digest;
}
