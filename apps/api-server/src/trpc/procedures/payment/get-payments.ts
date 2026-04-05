import paymentService from "../../../clients/payment";
import serverError from "../../../utils/server-error";
import { protectedProcedure } from "../..";

const getPayments = protectedProcedure.query(async (opts) => {
	try {
		return await paymentService.getPayments({
			accountId: opts.ctx.account.id,
		});
	} catch (err) {
		throw serverError(err);
	}
});

export default getPayments;
