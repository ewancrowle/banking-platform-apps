import type { CreateHTTPContextOptions } from "@trpc/server/adapters/standalone";
import oauthService from "../rpc-clients/oauth-service.ts";
import accountService from "../rpc-clients/account-service.ts";

export async function createContext({ req }: CreateHTTPContextOptions) {
	if (!req.headers.authorization) {
		return {
			account: null,
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
		};
	}
}
export type Context = Awaited<ReturnType<typeof createContext>>;
