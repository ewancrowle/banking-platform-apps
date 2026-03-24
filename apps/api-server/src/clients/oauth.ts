import { createClient } from "@connectrpc/connect";
import { createConnectTransport } from "@connectrpc/connect-web";
import { env } from "../utils/env";
import { OAuthService } from "protos/oauth";

const transport = createConnectTransport({
  baseUrl: env.AUTH_SERVICE_URL,
});

const client = createClient(OAuthService, transport);

export default client;
