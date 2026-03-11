import { initTRPC } from '@trpc/server';
import { createHTTPServer } from '@trpc/server/adapters/standalone';

/**
 * Initialization of tRPC backend
 * Should be done only once per backend!
 */
const t = initTRPC.create();

/**
 * Export reusable router and procedure helpers
 * that can be used throughout the router
 */
export const router = t.router;
export const publicProcedure = t.procedure;

const appRouter = router({
    // ...
});

// Export type router type signature,
// NOT the router itself.
export type AppRouter = typeof appRouter;

const server = createHTTPServer({
    router: appRouter,
});

server.listen(3000);
