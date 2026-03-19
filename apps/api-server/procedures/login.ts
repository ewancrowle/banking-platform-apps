import { publicProcedure } from "../server.ts";
import { z } from "zod";
import oauthService from "../rpc-clients/oauth-service.ts";

const login = publicProcedure
	.input(
		z.object({
			email: z.email(),
			password: z
				.string()
				.regex(
					/^(?=.*?[A-Z])(?=.*?[a-z])(?=.*?[0-9])(?=.*?[#?!@$ %^&*-]).{8,}$/,
				),
		}),
	)
	.output(
		z.object({
			accessToken: z.string(),
			expiresIn: z.number(),
			refreshToken: z.string(),
		}),
	)
	.mutation((opts) => oauthService.token(opts.input));

export default login;
