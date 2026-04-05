import { router } from "../index";
import getBalances from "../procedures/balance/get-balances";

export const balanceRouter = router({
	getBalances,
});

export type BalanceRouter = typeof balanceRouter;
