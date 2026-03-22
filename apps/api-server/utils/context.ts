import type { CreateHTTPContextOptions } from "@trpc/server/adapters/standalone";
import accountService from "../rpc-clients/account-service.ts";
import oauthService from "../rpc-clients/oauth-service.ts";

export async function createContext({ req }: CreateHTTPContextOptions) {
	if (!req.headers.authorization) {
		return {
			account: null,
			ipAddress: req.socket.remoteAddress,
			userAgent: req.headers["user-agent"],
		};
	}

	try {
		const { accountId } = await oauthService.introspect({
			accessToken: req.headers.authorization.split(" ")[1],
		});

		const account = await accountService.getAccount({
			id: accountId,
		});

		return {
			account,
		};
	} catch (err) {
		return {
			account: null,
			ipAddress: req.socket.remoteAddress,
			userAgent: req.headers["user-agent"],
		};
	}
}
export type Context = Awaited<ReturnType<typeof createContext>>;
