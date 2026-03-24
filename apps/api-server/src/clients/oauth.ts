import { createClient, Interceptor } from "@connectrpc/connect";
import { createConnectTransport } from "@connectrpc/connect-web";
import { env } from "../utils/env";
import { OAuthService } from "protos/oauth";

const logger: Interceptor = (next) => async (req) => {
  console.log(`sending message to ${req.url}`);
  return await next(req);
};

const transport = createConnectTransport({
  baseUrl: env.AUTH_SERVICE_URL,
  interceptors: [logger],
});

const client = createClient(OAuthService, transport);

export default client;
