import type { Payment } from "protos/payment";
import { Alert } from "react-native";
import { create } from "zustand";
import trpc from "@/api/trpc";
import formatCurrency from "@/utils/format-currency";

type BalanceState = {
	balance: bigint | null;
	formattedBalance: string | null;
	refresh: () => Promise<void>;
};

export const useBalanceStore = create<BalanceState>((set) => ({
	balance: null,
	formattedBalance: null,
	refresh: async () => {
		try {
			const { availableBalance, currencyCode } =
				await trpc.balance.getBalances.query();
			set({
				balance: availableBalance,
				formattedBalance: formatCurrency(availableBalance, currencyCode),
			});
		} catch {
			Alert.alert("Sorry, we couldn't load your balance at this time.");
		}
	},
}));
