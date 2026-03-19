import { protectedProcedure } from "../server.ts";

const getAccount = protectedProcedure.mutation((opts) => opts.ctx.account);

export default getAccount;
