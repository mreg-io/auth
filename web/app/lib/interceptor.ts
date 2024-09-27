import { Interceptor } from "@connectrpc/connect";

export const stripHeaders =
  (...headerNames: string[]): Interceptor =>
  (next) =>
  async (req) => {
    const response = await next(req);
    for (const name of headerNames) {
      response.header.delete(name);
    }
    return response;
  };
