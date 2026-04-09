import type { Payment } from "protos/payment";
import { create } from "zustand";

type PaymentsState = {
	payments: Payment[];
	setPayments: (payments: Payment[]) => void;
	clear: () => void;
};

export const usePaymentsStore = create<PaymentsState>((set) => ({
	payments: [],
	setPayments: (payments) => set({ payments }),
	clear: () => set({ payments: [] }),
}));
