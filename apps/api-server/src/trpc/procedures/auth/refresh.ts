import { z } from "zod";
import oauthService from "../../../clients/oauth";
import { publicProcedure } from "../../index";
import serverError from "../../../utils/server-error";

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
  .mutation((opts) => {
    try {
      return oauthService.refresh(opts.input);
    } catch (err) {
      throw serverError(err);
    }
  });

export default refresh;
