import { createClient } from "@connectrpc/connect";
import { createConnectTransport } from "@connectrpc/connect-node";
import { LedgerService } from "protos/ledger";
import { env } from "../utils/env";

const transport = createConnectTransport({
	httpVersion: "2",
	baseUrl: env.AUTH_SERVICE_URL,
});

const client = createClient(LedgerService, transport);

export default client;
