import { createClient } from "@connectrpc/connect";
import { createConnectTransport } from "@connectrpc/connect-node";
import { ConfirmationOfPayeeService } from "protos/confirmation_of_payee";
import { env } from "../utils/env";

const transport = createConnectTransport({
	httpVersion: "2",
	baseUrl: env.COP_SERVICE_URL,
});

const client = createClient(ConfirmationOfPayeeService, transport);

export default client;
