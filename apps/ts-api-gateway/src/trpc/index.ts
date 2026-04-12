import { initTRPC, TRPCError } from "@trpc/server";
import { createTRPCStoreLimiter } from "@trpc-limiter/memory";
import superjson from "superjson";
import type { Context } from "../utils/context";

const t = initTRPC.context<Context>().create({
	transformer: superjson,
});

export const router = t.router;

const rateLimiter = createTRPCStoreLimiter<typeof t>({
	fingerprint: (ctx) => ctx.ipAddress,
	message: (retryAfter) =>
		`Too many requests. Slow down. Retry in ${retryAfter}s.`,
	max: 15,
	windowMs: 10000,
});

export const publicProcedure = t.procedure.use(rateLimiter);

export const protectedProcedure = t.procedure
	.use(rateLimiter)
	.use(async (opts) => {
		const { ctx } = opts;
		if (!ctx.account) {
			throw new TRPCError({ code: "UNAUTHORIZED" });
		}

		return opts.next({
			ctx: {
				account: ctx.account,
			},
		});
	});
