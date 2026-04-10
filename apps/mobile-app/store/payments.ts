import type { Payment } from "protos/payment";
import { Alert } from "react-native";
import { create } from "zustand";
import trpc from "@/api/trpc";

type PaymentsState = {
	payments: Payment[];
	clear: () => void;
	refresh: () => Promise<void>;
};

const VALID_TYPES = ["deposit", "withdrawal", "card", "account_to_account"];

export const usePaymentsStore = create<PaymentsState>((set) => ({
	payments: [],
	clear: () => set({ payments: [] }),
	refresh: async () => {
		try {
			const { payments } = await trpc.payment.getPayments.query();
			set({ payments: payments.filter((p) => VALID_TYPES.includes(p.type)) });
		} catch {
			Alert.alert("Sorry, we couldn't load your transactions at this time.");
		}
	},
}));
