import { Code, ConnectError } from "@connectrpc/connect";
import { TRPCError } from "@trpc/server";

export default function (error: unknown): TRPCError {
  const connectErr = ConnectError.from(error);
  switch (connectErr.code) {
    case Code.AlreadyExists:
      return new TRPCError({
        code: "CONFLICT",
        message: connectErr.message,
        cause: error,
      });
    case Code.InvalidArgument:
      return new TRPCError({
        code: "BAD_REQUEST",
        message: connectErr.message,
        cause: error,
      });
    case Code.NotFound:
      return new TRPCError({
        code: "NOT_FOUND",
        message: connectErr.message,
        cause: error,
      });
    case Code.Unauthenticated:
      return new TRPCError({
        code: "UNAUTHORIZED",
        message: connectErr.message,
        cause: error,
      });
    case Code.PermissionDenied:
      return new TRPCError({
        code: "FORBIDDEN",
        message: connectErr.message,
        cause: error,
      });
    default:
      return new TRPCError({
        code: "INTERNAL_SERVER_ERROR",
        message: connectErr.message,
        cause: error,
      });
  }
}
