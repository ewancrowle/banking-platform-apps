import { TRPCError } from "@trpc/server";
import { refreshTokens } from "@/api/auth";
import { useAuthStore } from "@/store/auth";

export default async function <T>(fn: () => Promise<T>) {
	try {
		return await fn();
	} catch (err) {
		if (err instanceof TRPCError) {
			if (err.code === "UNAUTHORIZED") {
				const newTokens = await refreshTokens();
				if (!newTokens) return null;

				await useAuthStore.getState().setTokens(newTokens);
				return await fn();
			}
		}
		throw err;
	}
}
