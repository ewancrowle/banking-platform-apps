import { router } from "../index";
import { accountRouter } from "./account";
import { authRouter } from "./auth";
import { balanceRouter } from "./balance";
import { merchantRouter } from "./merchant";
import { paymentRouter } from "./payment";

export const appRouter = router({
	auth: authRouter,
	account: accountRouter,
	payment: paymentRouter,
	balance: balanceRouter,
	merchant: merchantRouter,
});

export type AppRouter = typeof appRouter;
