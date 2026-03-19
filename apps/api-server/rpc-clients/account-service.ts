import { createConnectTransport } from "@connectrpc/connect-web";
import { createClient } from "@connectrpc/connect";
import { AccountService } from "protos/account/account/v1/account_pb.ts";

const transport = createConnectTransport({
	baseUrl: "https://demo.connectrpc.com",
});

const client = createClient(AccountService, transport);

export default client;
