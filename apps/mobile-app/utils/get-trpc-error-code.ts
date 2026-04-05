import { isTRPCClientError } from "@trpc/client";
import type { TRPC_ERROR_CODE_KEY } from "@trpc/server";
import { TRPC_ERROR_CODES_BY_KEY } from "@trpc/server/rpc";

export function getTRPCErrorCode(error: unknown) {
	if (isTRPCClientError(error)) {
		const code = error.data?.code;
		if (typeof code === "string" && code in TRPC_ERROR_CODES_BY_KEY) {
			return code as TRPC_ERROR_CODE_KEY;
		}
	}
}
