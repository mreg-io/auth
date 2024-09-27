import {
  AnyMessage,
  Message,
  MessageType,
  protoBase64,
} from "@bufbuild/protobuf";
import { useLoaderData } from "@remix-run/react";

export function protobuf<T extends Message<T> = AnyMessage>(
  message: Message<T>,
  init?: ResponseInit,
) {
  const headers = new Headers(init?.headers);
  headers.set("Content-Type", "application/octet-stream");
  return new Response(protoBase64.enc(message.toBinary()), {
    ...init,
    headers,
  });
}

export function useLoaderProtobuf<T extends Message<T> = AnyMessage>(
  ProtobufMessage: MessageType<T>,
) {
  return new ProtobufMessage().fromBinary(protoBase64.dec(useLoaderData()));
}
