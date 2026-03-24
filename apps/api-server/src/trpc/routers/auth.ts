import { router } from "../index";
import login from "../procedures/auth/login";
import signUp from "../procedures/auth/sign-up";
import refresh from "../procedures/auth/refresh";

export const authRouter = router({
  signUp,
  login,
  refresh,
});

export type AuthRouter = typeof authRouter;
