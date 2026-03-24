import { createClient } from "@connectrpc/connect";
import { createConnectTransport } from "@connectrpc/connect-node";
import { AccountService } from "protos/account";
import { env } from "../utils/env";

const transport = createConnectTransport({
  httpVersion: "2",
  baseUrl: env.ACCOUNT_SERVICE_URL,
});

const client = createClient(AccountService, transport);

export default client;
