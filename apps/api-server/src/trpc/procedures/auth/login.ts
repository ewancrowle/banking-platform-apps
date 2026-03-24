import { z } from "zod";
import { publicProcedure } from "../../index";
import oauthService from "../../../clients/oauth";

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
  .mutation((opts) => oauthService.token(opts.input));

export default login;
