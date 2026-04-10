import { router } from "../index";
import getTotalSpending from "../procedures/spending/get-total-spending";

export const spendingRouter = router({
	getTotalSpending,
});

export type SpendingRouter = typeof spendingRouter;
