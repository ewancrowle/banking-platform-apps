import { Decision } from "protos/payment";
import * as z from "zod";
import paymentService from "../../../clients/payment";
import serverError from "../../../utils/server-error";
import { protectedProcedure } from "../..";

const newTransfer = protectedProcedure
	.input(
		z.object({
			confirmationOfPayeeToken: z.string(),
			amount: z.number(),
			reference: z.string(),
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
				confirmationOfPayeeToken: opts.input.confirmationOfPayeeToken,
				amount: BigInt(opts.input.amount),
				currencyCode: "GBP",
				description: opts.input.reference,
				type: "account_to_account",
			});
		} catch (err) {
			throw serverError(err);
		}
	});

export default newTransfer;
