import { Code, ConnectError } from "@connectrpc/connect";
import { TRPCError } from "@trpc/server";

export default function (error: TRPCError) {
	const connectErr = ConnectError.from(error.cause);
	switch (connectErr.code) {
		case Code.AlreadyExists:
			throw new TRPCError({
				code: "CONFLICT",
				message: connectErr.message,
			});
		case Code.InvalidArgument:
			throw new TRPCError({
				code: "BAD_REQUEST",
				message: connectErr.message,
			});
		case Code.NotFound:
			throw new TRPCError({
				code: "NOT_FOUND",
				message: connectErr.message,
			});
		default:
			throw new TRPCError({
				code: "INTERNAL_SERVER_ERROR",
				message: connectErr.message,
			});
	}
}
