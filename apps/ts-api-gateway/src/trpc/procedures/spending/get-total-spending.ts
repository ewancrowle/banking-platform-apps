import ledgerService from "../../../clients/ledger";
import serverError from "../../../utils/server-error";
import { protectedProcedure } from "../..";

const getTotalSpending = protectedProcedure.query(async (opts) => {
	try {
		return await ledgerService.getTotalSpending({
			accountId: opts.ctx.account.id,
		});
	} catch (err) {
		throw serverError(err);
	}
});

export default getTotalSpending;
