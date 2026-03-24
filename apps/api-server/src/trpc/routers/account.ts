import { router } from "../index";
import getAccount from "../procedures/account/get-account";

export const accountRouter = router({
  getAccount,
});

export type AccountRouter = typeof accountRouter;
