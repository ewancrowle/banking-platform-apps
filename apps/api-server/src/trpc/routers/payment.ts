import { router } from "../index";
import confirmPayee from "../procedures/payment/confirm-payee";
import newTransfer from "../procedures/payment/new-transfer";

export const paymentRouter = router({
	confirmPayee,
	newTransfer,
});

export type PaymentRouter = typeof paymentRouter;
