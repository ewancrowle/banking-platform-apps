import { createClient } from "@connectrpc/connect";
import { createConnectTransport } from "@connectrpc/connect-web";
import { AccountService } from "protos/account";
import { env } from "../utils/env";

const transport = createConnectTransport({
  baseUrl: env.ACCOUNT_SERVICE_URL,
});

const client = createClient(AccountService, transport);

export default client;
