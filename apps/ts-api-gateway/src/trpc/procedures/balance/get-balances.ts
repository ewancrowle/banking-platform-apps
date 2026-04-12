import ledgerService from "../../../clients/ledger";
import serverError from "../../../utils/server-error";
import { protectedProcedure } from "../..";

const getBalances = protectedProcedure.query(async (opts) => {
	try {
		return await ledgerService.getBalances({
			accountId: opts.ctx.account.id,
		});
	} catch (err) {
		throw serverError(err);
	}
});

export default getBalances;
