import { protectedProcedure } from "../..";

const getAccount = protectedProcedure.query((opts) => opts.ctx.account);

export default getAccount;
