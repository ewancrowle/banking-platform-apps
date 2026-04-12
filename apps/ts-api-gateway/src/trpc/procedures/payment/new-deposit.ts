import { Decision } from "protos/payment";
import * as z from "zod";
import paymentService from "../../../clients/payment";
import serverError from "../../../utils/server-error";
import { protectedProcedure } from "../..";

const newDeposit = protectedProcedure
	.input(
		z.object({
			amount: z.number(),
		}),
	)
	.output(
		z.object({
			paymentId: z.bigint(),
			decision: z.enum(Decision),
			decisionId: z.bigint(),
		}),
	)
	.mutation(async (opts) => {
		try {
			return await paymentService.authorisePayment({
				accountId: opts.ctx.account.id,
				amount: BigInt(opts.input.amount),
				currencyCode: "GBP",
				description: "Deposit",
				type: "deposit",
			});
		} catch (err) {
			throw serverError(err);
		}
	});

export default newDeposit;
