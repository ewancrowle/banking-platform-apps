import { Alert } from "react-native";
import { create } from "zustand";
import trpc from "@/api/trpc";
import formatCurrency from "@/utils/format-currency";

type SpendingState = {
	spentToday: bigint | null;
	spentThisWeek: bigint | null;
	spentThisMonth: bigint | null;
	formattedSpentToday: string | null;
	formattedSpentThisWeek: string | null;
	formattedSpentThisMonth: string | null;
	refresh: () => Promise<void>;
};

export const useSpendingStore = create<SpendingState>((set) => ({
	spentToday: null,
	spentThisWeek: null,
	spentThisMonth: null,
	formattedSpentToday: null,
	formattedSpentThisWeek: null,
	formattedSpentThisMonth: null,
	refresh: async () => {
		try {
			const {
				totalSpentToday,
				totalSpentThisWeek,
				totalSpentThisMonth,
				currencyCode,
			} = await trpc.spending.getTotalSpending.query();
			set({
				spentToday: totalSpentToday,
				spentThisWeek: totalSpentThisWeek,
				spentThisMonth: totalSpentThisMonth,
				formattedSpentToday: formatCurrency(totalSpentToday, currencyCode),
				formattedSpentThisWeek: formatCurrency(
					totalSpentThisWeek,
					currencyCode,
				),
				formattedSpentThisMonth: formatCurrency(
					totalSpentThisMonth,
					currencyCode,
				),
			});
		} catch {
			Alert.alert("Sorry, we couldn't load your spending at this time.");
		}
	},
}));
