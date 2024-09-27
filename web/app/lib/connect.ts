import {
  Message,
  MethodInfo,
  MethodInfoBiDiStreaming,
  MethodInfoClientStreaming,
  MethodInfoServerStreaming,
  MethodInfoUnary,
  MethodKind,
  PartialMessage,
  ServiceType,
} from "@bufbuild/protobuf";
import {
  CallOptions,
  Code,
  ConnectError,
  makeAnyClient,
  StreamResponse,
  Transport,
} from "@connectrpc/connect";
import { createAsyncIterable } from "@connectrpc/connect/protocol";

export type PromiseCallOptions = Omit<CallOptions, "onHeader" | "onTrailer">;

export interface UnaryPromiseResponse<O extends Message<O>> {
  response: O;
  headers: Headers;
  trailers: Headers;
}

export type PromiseClient<T extends ServiceType> = {
  [P in keyof T["methods"]]: T["methods"][P] extends MethodInfoUnary<
    infer I,
    infer O
  >
    ? UnaryFn<I, O>
    : T["methods"][P] extends MethodInfoServerStreaming<infer I, infer O>
      ? (request: PartialMessage<I>, options?: CallOptions) => AsyncIterable<O>
      : T["methods"][P] extends MethodInfoClientStreaming<infer I, infer O>
        ? ClientStreamingFn<I, O>
        : T["methods"][P] extends MethodInfoBiDiStreaming<infer I, infer O>
          ? (
              request: PartialMessage<I>,
              options?: CallOptions,
            ) => AsyncIterable<O>
          : never;
};

export const createPromiseClient = <T extends ServiceType>(
  service: T,
  transport: Transport,
) =>
  makeAnyClient(service, (method) => {
    switch (method.kind) {
      case MethodKind.Unary:
        return createUnaryFn(transport, service, method);
      case MethodKind.ServerStreaming:
        return createServerStreamingFn(transport, service, method);
      case MethodKind.ClientStreaming:
        return createClientStreamingFn(transport, service, method);
      case MethodKind.BiDiStreaming:
        return createBiDiStreamingFn(transport, service, method);
    }
  }) as PromiseClient<T>;

type UnaryFn<I extends Message<I>, O extends Message<O>> = (
  request: PartialMessage<I>,
  options?: PromiseCallOptions,
) => Promise<UnaryPromiseResponse<O>>;

const createUnaryFn =
  <I extends Message<I>, O extends Message<O>>(
    transport: Transport,
    service: ServiceType,
    method: MethodInfo<I, O>,
  ): UnaryFn<I, O> =>
  async (input, options) => {
    const response = await transport.unary(
      service,
      method,
      options?.signal,
      options?.timeoutMs,
      options?.headers,
      input,
      options?.contextValues,
    );
    return {
      response: response.message,
      headers: response.header,
      trailers: response.trailer,
    };
  };

type ClientStreamingFn<I extends Message<I>, O extends Message<O>> = (
  request: AsyncIterable<PartialMessage<I>>,
  options?: PromiseCallOptions,
) => Promise<UnaryPromiseResponse<O>>;

const createClientStreamingFn =
  <I extends Message<I>, O extends Message<O>>(
    transport: Transport,
    service: ServiceType,
    method: MethodInfo<I, O>,
  ): ClientStreamingFn<I, O> =>
  async (input, options) => {
    const response = await transport.stream<I, O>(
      service,
      method,
      options?.signal,
      options?.timeoutMs,
      options?.headers,
      input,
      options?.contextValues,
    );

    let singleMessage: O | undefined;
    let count = 0;
    for await (const message of response.message) {
      singleMessage = message;
      count++;
    }
    if (!singleMessage) {
      throw new ConnectError(
        "protocol error: missing response message",
        Code.Unimplemented,
      );
    }
    if (count > 1) {
      throw new ConnectError(
        "protocol error: received extra messages for client streaming method",
        Code.Unimplemented,
      );
    }

    return {
      headers: response.header,
      trailers: response.trailer,
      response: singleMessage,
    };
  };

// Fallback to ConnectRPC default implementation
type ServerStreamingFn<I extends Message<I>, O extends Message<O>> = (
  request: PartialMessage<I>,
  options?: CallOptions,
) => AsyncIterable<O>;

function createServerStreamingFn<I extends Message<I>, O extends Message<O>>(
  transport: Transport,
  service: ServiceType,
  method: MethodInfo<I, O>,
): ServerStreamingFn<I, O> {
  return function (input, options): AsyncIterable<O> {
    return handleStreamResponse(
      transport.stream<I, O>(
        service,
        method,
        options?.signal,
        options?.timeoutMs,
        options?.headers,
        createAsyncIterable([input]),
        options?.contextValues,
      ),
      options,
    );
  };
}

type BiDiStreamingFn<I extends Message<I>, O extends Message<O>> = (
  request: AsyncIterable<PartialMessage<I>>,
  options?: CallOptions,
) => AsyncIterable<O>;

function createBiDiStreamingFn<I extends Message<I>, O extends Message<O>>(
  transport: Transport,
  service: ServiceType,
  method: MethodInfo<I, O>,
): BiDiStreamingFn<I, O> {
  return function (
    request: AsyncIterable<PartialMessage<I>>,
    options?: CallOptions,
  ): AsyncIterable<O> {
    return handleStreamResponse(
      transport.stream<I, O>(
        service,
        method,
        options?.signal,
        options?.timeoutMs,
        options?.headers,
        request,
        options?.contextValues,
      ),
      options,
    );
  };
}

function handleStreamResponse<I extends Message<I>, O extends Message<O>>(
  stream: Promise<StreamResponse<I, O>>,
  options?: CallOptions,
): AsyncIterable<O> {
  const it = (async function* () {
    const response = await stream;
    options?.onHeader?.(response.header);
    yield* response.message;
    options?.onTrailer?.(response.trailer);
  })()[Symbol.asyncIterator]();
  // Create a new iterable to omit throw/return.
  return {
    [Symbol.asyncIterator]: () => ({
      next: () => it.next(),
    }),
  };
}
