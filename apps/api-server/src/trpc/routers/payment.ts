import { router } from "../index";
import confirmPayee from "../procedures/payment/confirm-payee";
import getPayments from "../procedures/payment/get-payments";
import newTransfer from "../procedures/payment/new-transfer";

export const paymentRouter = router({
	confirmPayee,
	newTransfer,
	getPayments,
});

export type PaymentRouter = typeof paymentRouter;
