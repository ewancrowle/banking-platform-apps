import { createConnectTransport } from "@connectrpc/connect-web";
import { createClient } from "@connectrpc/connect";
import { OAuthService } from "protos/oauth/oauth/v1/oauth_pb.ts";

const transport = createConnectTransport({
	baseUrl: "https://demo.connectrpc.com",
});

const client = createClient(OAuthService, transport);

export default client;
