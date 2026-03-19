import { initTRPC, TRPCError } from "@trpc/server";
import { createHTTPServer } from "@trpc/server/adapters/standalone";
import signUp from "./procedures/sign-up.ts";
import serverError from "./utils/server-error.ts";
import login from "./procedures/login.ts";
import refresh from "./procedures/refresh.ts";
import { type Context, createContext } from "./utils/context.ts";
import getAccount from "./procedures/get-account.ts";

/**
 * Initialization of tRPC backend
 * Should be done only once per backend!
 */
const t = initTRPC.context<Context>().create();

/**
 * Export reusable router and procedure helpers
 * that can be used throughout the router
 */
export const router = t.router;
export const publicProcedure = t.procedure;
export const protectedProcedure = t.procedure.use(async (opts) => {
	if (!opts.ctx.account) {
		throw new TRPCError({ code: "UNAUTHORIZED" });
	}
	return opts.next({
		ctx: opts.ctx,
	});
});

const appRouter = router({
	signUp,
	login,
	refresh,
	getAccount,
});

// Export type router type signature,
// NOT the router itself.
export type AppRouter = typeof appRouter;

const server = createHTTPServer({
	router: appRouter,
	createContext,
	onError: (opts) => serverError(opts.error),
});

server.listen(3000);
