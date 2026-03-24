import { initTRPC, TRPCError } from "@trpc/server";
import { Context } from "../utils/context";
import { authRouter } from "./routers/auth";

const t = initTRPC.context<Context>().create();

export const router = t.router;

export const publicProcedure = t.procedure;

export const protectedProcedure = t.procedure.use(async (opts) => {
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

export const appRouter = router({
  auth: authRouter,
});

export type AppRouter = typeof appRouter;
