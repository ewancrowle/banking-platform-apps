import type { CreateHTTPContextOptions } from "@trpc/server/adapters/standalone";
import accountService from "../clients/account";
import oauthService from "../clients/oauth";

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
  } catch {
    return {
      account: null,
      ipAddress: req.socket.remoteAddress,
      userAgent: req.headers["user-agent"],
    };
  }
}
export type Context = Awaited<ReturnType<typeof createContext>>;
