import { z } from "zod";
import accountService from "../rpc-clients/account-service.ts";
import oauthService from "../rpc-clients/oauth-service.ts";
import { publicProcedure } from "../server.ts";

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
	.output(
		z.object({
			accessToken: z.string(),
			expiresIn: z.number(),
			refreshToken: z.string(),
		}),
	)
	.mutation(async (opts) => {
		await accountService.createAccount(opts.input);
		return oauthService.token({
			email: opts.input.email,
			ipAddress: opts.ctx.ipAddress,
			password: opts.input.password,
			userAgent: opts.ctx.userAgent,
		});
	});

export default signUp;
