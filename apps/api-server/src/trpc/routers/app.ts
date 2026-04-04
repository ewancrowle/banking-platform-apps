import { router } from "../index";
import { accountRouter } from "./account";
import { authRouter } from "./auth";
import { paymentRouter } from "./payment";

export const appRouter = router({
	auth: authRouter,
	account: accountRouter,
	payment: paymentRouter,
});

export type AppRouter = typeof appRouter;
