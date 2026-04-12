import type { CreateHTTPContextOptions } from "@trpc/server/adapters/standalone";
import type { IncomingMessage } from "http";
import accountService from "../clients/account";
import oauthService from "../clients/oauth";

export const ipAddress = (req: IncomingMessage) => {
	const forwarded = req.headers["x-forwarded-for"];
	const ip = forwarded
		? (typeof forwarded === "string" ? forwarded : forwarded[0])?.split(/, /)[0]
		: req.socket.remoteAddress;
	return ip || "127.0.0.1";
};

export async function createContext({ req }: CreateHTTPContextOptions) {
	if (!req.headers.authorization) {
		return {
			account: null,
			ipAddress: ipAddress(req),
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
			ipAddress: ipAddress(req),
			userAgent: req.headers["user-agent"],
		};
	} catch {
		return {
			account: null,
			ipAddress: ipAddress(req),
			userAgent: req.headers["user-agent"],
		};
	}
}
export type Context = Awaited<ReturnType<typeof createContext>>;
