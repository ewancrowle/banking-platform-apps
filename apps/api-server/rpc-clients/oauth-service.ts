import { createClient } from "@connectrpc/connect";
import { createConnectTransport } from "@connectrpc/connect-web";
import { OAuthService } from "protos/oauth/oauth/v1/oauth_pb.ts";
import { env } from "../utils/env.ts";

const transport = createConnectTransport({
	baseUrl: env.AUTH_SERVICE_URL,
});

const client = createClient(OAuthService, transport);

export default client;
