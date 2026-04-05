import { router } from "../index";
import { accountRouter } from "./account";
import { authRouter } from "./auth";
import { balanceRouter } from "./balance";
import { paymentRouter } from "./payment";

export const appRouter = router({
	auth: authRouter,
	account: accountRouter,
	payment: paymentRouter,
	balance: balanceRouter,
});

export type AppRouter = typeof appRouter;
