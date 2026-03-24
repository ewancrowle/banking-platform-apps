import {createClient} from "@connectrpc/connect";
import {createConnectTransport} from "@connectrpc/connect-node";
import {env} from "../utils/env";
import {OAuthService} from "protos/oauth";

const transport = createConnectTransport({
  httpVersion: "2",
  baseUrl: env.ACCOUNT_SERVICE_URL,
});

const client = createClient(OAuthService, transport);

export default client;
