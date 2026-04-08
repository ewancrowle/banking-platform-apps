import { router } from "../index";
import confirmPayee from "../procedures/payment/confirm-payee";
import getPayments from "../procedures/payment/get-payments";
import newCardPayment from "../procedures/payment/new-card-payment";
import newDeposit from "../procedures/payment/new-deposit";
import newTransfer from "../procedures/payment/new-transfer";

export const paymentRouter = router({
	confirmPayee,
	newDeposit,
	newCardPayment,
	newTransfer,
	getPayments,
});

export type PaymentRouter = typeof paymentRouter;
