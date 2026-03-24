import { z } from "zod";
import oauthService from "../../../clients/oauth";
import { publicProcedure } from "../../index";

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
