import { router } from "../index";
import { accountRouter } from "./account";
import { authRouter } from "./auth";
import { balanceRouter } from "./balance";
import { merchantRouter } from "./merchant";
import { paymentRouter } from "./payment";
import { spendingRouter } from "./spending";

export const appRouter = router({
	auth: authRouter,
	account: accountRouter,
	payment: paymentRouter,
	balance: balanceRouter,
	merchant: merchantRouter,
	spending: spendingRouter,
});

export type AppRouter = typeof appRouter;
