import { createClient } from "@connectrpc/connect";
import { createConnectTransport } from "@connectrpc/connect-node";
import { PaymentService } from "protos/payment";
import { env } from "../utils/env";

const transport = createConnectTransport({
	httpVersion: "2",
	baseUrl: env.PAYMENT_SERVICE_URL,
});

const client = createClient(PaymentService, transport);

export default client;
