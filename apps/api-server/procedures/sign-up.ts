import { publicProcedure } from "../server.ts";
import { z } from "zod";
import accountService from "../rpc-clients/account-service.ts";

const signUp = publicProcedure
	.input(
		z.object({
			firstName: z.string(),
			middleNames: z.string(),
			lastName: z.string(),
			email: z.email(),
			password: z
				.string()
				.regex(
					/^(?=.*?[A-Z])(?=.*?[a-z])(?=.*?[0-9])(?=.*?[#?!@$ %^&*-]).{8,}$/,
				),
			line1: z.string(),
			line2: z.string(),
			town: z.string(),
			postcode: z.string(),
		}),
	)
	.output(z.object({ id: z.bigint() }))
	.mutation((opts) => accountService.createAccount(opts.input));

export default signUp;
