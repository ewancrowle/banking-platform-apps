import { protectedProcedure } from "../server.ts";

const getAccount = protectedProcedure.query((opts) => opts.ctx.account);

export default getAccount;
