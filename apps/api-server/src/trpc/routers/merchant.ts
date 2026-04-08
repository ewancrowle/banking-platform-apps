import { router } from "../index";
import getAllMerchants from "../procedures/merchant/getAllMerchants";

export const merchantRouter = router({
	getAllMerchants,
});

export type MerchantRouter = typeof merchantRouter;
