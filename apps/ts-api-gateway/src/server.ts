import { createHTTPServer } from "@trpc/server/adapters/standalone";
import { appRouter } from "./trpc/routers/app";
import { createContext } from "./utils/context";

const server = createHTTPServer({
  router: appRouter,
  createContext,
});

server.listen(3000);
