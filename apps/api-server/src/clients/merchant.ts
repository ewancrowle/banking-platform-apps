import { createClient } from "@connectrpc/connect";
import { createConnectTransport } from "@connectrpc/connect-node";
import { MerchantService } from "protos/merchant";
import { env } from "../utils/env";

const transport = createConnectTransport({
	httpVersion: "2",
	baseUrl: env.MERCHANT_SERVICE_URL,
});

const client = createClient(MerchantService, transport);

export default client;
