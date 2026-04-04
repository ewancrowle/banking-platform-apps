import { router } from "../index";
import confirmPayee from "../procedures/payment/confirm-payee";

export const paymentRouter = router({
	confirmPayee,
});

export type PaymentRouter = typeof paymentRouter;
