import { Code, ConnectError } from "@connectrpc/connect";
import { TRPCError } from "@trpc/server";

export default function (error: TRPCError) {
  const connectErr = ConnectError.from(error.cause);
  switch (connectErr.code) {
    case Code.AlreadyExists:
      throw new TRPCError({
        code: "CONFLICT",
        message: connectErr.message,
        cause: error.cause,
      });
    case Code.InvalidArgument:
      throw new TRPCError({
        code: "BAD_REQUEST",
        message: connectErr.message,
        cause: error.cause,
      });
    case Code.NotFound:
      throw new TRPCError({
        code: "NOT_FOUND",
        message: connectErr.message,
        cause: error.cause,
      });
    case Code.Unauthenticated:
      throw new TRPCError({
        code: "UNAUTHORIZED",
        message: connectErr.message,
        cause: error.cause,
      });
    case Code.PermissionDenied:
      throw new TRPCError({
        code: "FORBIDDEN",
        message: connectErr.message,
        cause: error.cause,
      });
    default:
      throw new TRPCError({
        code: "INTERNAL_SERVER_ERROR",
        message: connectErr.message,
        cause: error.cause,
      });
  }
}
