import { protectedProcedure } from "../../index";

const getAccount = protectedProcedure.query((opts) => opts.ctx.account);

export default getAccount;
