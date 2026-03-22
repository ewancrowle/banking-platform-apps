import { z } from "zod";
import oauthService from "../rpc-clients/oauth-service.ts";
import { publicProcedure } from "../server.ts";

const refresh = publicProcedure
	.input(
		z.object({
			refreshToken: z.string(),
		}),
	)
	.output(
		z.object({
			accessToken: z.string(),
			expiresIn: z.number(),
			refreshToken: z.string(),
		}),
	)
	.mutation((opts) => oauthService.refresh(opts.input));

export default refresh;
