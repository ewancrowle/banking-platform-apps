import * as z from "zod";
import copService from "../../../clients/confirmation-of-payee";
import serverError from "../../../utils/server-error";
import { publicProcedure } from "../..";

const confirmPayee = publicProcedure
	.input(
		z.object({
			firstName: z.string(),
			lastName: z.string(),
			accountNumber: z.string(),
		}),
	)
	.output(
		z.object({
			confirmationOfPayeeToken: z.bigint(),
		}),
	)
	.mutation(async (opts) => {
		try {
			return await copService.confirmPayee({
				firstName: opts.input.firstName,
				lastName: opts.input.lastName,
				accountNum: opts.input.accountNumber,
			});
		} catch (err) {
			throw serverError(err);
		}
	});

export default confirmPayee;
