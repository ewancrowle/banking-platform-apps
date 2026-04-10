import type { Ionicons } from "@expo/vector-icons";
import type { Payment } from "protos/payment";

export function getPaymentIcon(
	payment: Payment,
): keyof typeof Ionicons.glyphMap {
	switch (payment.type) {
		case "deposit":
			return "arrow-down";
		case "withdrawal":
			return "arrow-up";
		case "card":
			if (payment.merchant) {
				switch (payment.merchant.mcc) {
					case "5411":
						return "cart-outline";
					case "5814":
						return "fast-food-outline";
					default:
						return "card-outline";
				}
			}
			return "card-outline";
		case "outbound_transfer":
			return "swap-horizontal-outline";
		case "inbound_transfer":
			return "swap-horizontal-outline";
		default:
			return "cash-outline";
	}
}
