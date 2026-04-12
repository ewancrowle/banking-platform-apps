import { deleteItemAsync, setItemAsync } from "expo-secure-store";
import { create } from "zustand";
import type { Account, Tokens } from "@/types/auth";

type AuthState = {
	accessToken: string | null;
	refreshToken: string | null;
	account: Account | null;
	setTokens: (tokens: Tokens) => Promise<void>;
	setAccount: (account: Account) => void;
	reset: () => Promise<void>;
};

export const useAuthStore = create<AuthState>((set) => ({
	accessToken: null,
	refreshToken: null,
	account: null,
	setTokens: async (tokens) => {
		await setItemAsync("accessToken", tokens.accessToken);
		await setItemAsync("refreshToken", tokens.refreshToken);
		set({ accessToken: tokens.accessToken, refreshToken: tokens.refreshToken });
	},
	setAccount: (account) => set({ account }),
	reset: async () => {
		await deleteItemAsync("accessToken");
		await deleteItemAsync("refreshToken");
		set({ accessToken: null, refreshToken: null, account: null });
	},
}));
