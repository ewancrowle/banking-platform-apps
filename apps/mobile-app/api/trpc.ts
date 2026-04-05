import { createTRPCClient, httpBatchLink } from "@trpc/client";
import type { AppRouter } from "api-server";
import superjson from "superjson";
import { useAuthStore } from "@/store/auth";

export default createTRPCClient<AppRouter>({
	links: [
		httpBatchLink({
			url: "http://localhost:3000",
			transformer: superjson,
			async headers() {
				const accessToken = useAuthStore.getState().accessToken;
				return {
					Authorization: accessToken ? `Bearer ${accessToken}` : undefined,
				};
			},
		}),
	],
});
