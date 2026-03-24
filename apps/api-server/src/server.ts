import {createHTTPServer} from "@trpc/server/adapters/standalone";
import {appRouter} from "./trpc/routers/app";
import {createContext} from "./utils/context";
import serverError from "./utils/server-error";

const server = createHTTPServer({
  router: appRouter,
  createContext,
  onError: (opts) => {
    const {cause} = opts.error;
    if (cause) console.error(cause);
    serverError(opts.error)
  },
});

server.listen(3000);
