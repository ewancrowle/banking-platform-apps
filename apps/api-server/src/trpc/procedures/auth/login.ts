import { z } from "zod";
import { publicProcedure } from "../../index";
import oauthService from "../../../clients/oauth";
import serverError from "../../../utils/server-error";

const login = publicProcedure
  .input(
    z.object({
      email: z.email(),
      password: z.string().regex(/^(?=.*?[A-Z])(?=.*?[a-z])(?=.*?[0-9])(?=.*?[#?!@$ %^&*-]).{8,}$/),
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
      return await oauthService.token(opts.input);
    } catch (err) {
      throw serverError(err);
    }
  });

export default login;
