import { deleteItemAsync, getItemAsync, setItemAsync } from "expo-secure-store";
import { getTRPCErrorCode } from "@/utils/get-trpc-error-code";
import trpc from "./trpc";

export const refreshTokens = async () => {
	const refreshToken = await getItemAsync("refreshToken");
	if (!refreshToken) return null;

	try {
		const newTokens = await trpc.auth.refresh.mutate({ refreshToken });
		await setItemAsync("accessToken", newTokens.accessToken);
		await setItemAsync("refreshToken", newTokens.refreshToken);
		return newTokens;
	} catch (err) {
		const code = getTRPCErrorCode(err);
		if (code === "UNAUTHORIZED" || code === "FORBIDDEN") {
			await deleteItemAsync("accessToken");
			await deleteItemAsync("refreshToken");
		}
		return null;
	}
};

export const getAccount = async () => {
	try {
		return await trpc.account.getAccount.query();
	} catch (error) {
		console.log(error);
		return null;
	}
};
