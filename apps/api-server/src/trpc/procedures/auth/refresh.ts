import * as z from "zod";
import oauthService from "../../../clients/oauth";
import serverError from "../../../utils/server-error";
import { publicProcedure } from "../..";

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
	.mutation(async (opts) => {
		try {
			return await oauthService.refresh(opts.input);
		} catch (err) {
			throw serverError(err);
		}
	});

export default refresh;
