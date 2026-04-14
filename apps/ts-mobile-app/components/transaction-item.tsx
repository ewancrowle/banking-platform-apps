import { Ionicons } from "@expo/vector-icons";
import { useTheme } from "@react-navigation/native";
import dayjs from "dayjs";
import relativeTime from "dayjs/plugin/relativeTime";
import type { Payment } from "protos/payment";
import { StyleSheet, View } from "react-native";
import { ThemedText } from "@/components/themed-text";
import { getPaymentIcon } from "@/utils/get-payment-icon";

type TransactionItemProps = Payment;

const formatAmount = (amount: number, currencyCode: string) => {
	return new Intl.NumberFormat("en-GB", {
		style: "currency",
		currency: currencyCode || "GBP",
	}).format(amount / 100);
};

dayjs.extend(relativeTime);

export function TransactionItem(payment: TransactionItemProps) {
	const theme = useTheme();
	const icon = getPaymentIcon(payment);
	const {
		amount,
		currencyCode,
		description,
		merchant,
		type,
		status,
		createdAt,
	} = payment;
	let time = dayjs(`${createdAt.split(".")[0]}Z`).fromNow();
	let typeLabel =
		{
			deposit: "Deposit",
			withdrawal: "Withdrawal",
			card: "Card Payment",
			outbound_transfer: "Transfer",
			inbound_transfer: "Transfer",
		}[type] || type;

	if (status === "declined") {
		time = `Declined ${time}`;
	} else if (status === "authorised") {
		typeLabel = `Pending ${typeLabel}`;
	}

	const styles = StyleSheet.create({
		container: {
			flexDirection: "row",
			justifyContent: "space-between",
			alignItems: "center",
		},
		leftSection: {
			flexDirection: "row",
			alignItems: "center",
			gap: 12,
		},
		iconContainer: {
			width: 40,
			height: 40,
			borderRadius: 20,
			justifyContent: "center",
			alignItems: "center",
			backgroundColor: theme.colors.card,
		},
		textDetails: {
			gap: 2,
		},
		subtitle: {
			opacity: 0.5,
		},
		rightSection: {
			alignItems: "flex-end",
			gap: 2,
		},
		title: {
			fontWeight: "600",
		},
	});

	return (
		<View style={styles.container}>
			<View style={styles.leftSection}>
				<View style={styles.iconContainer}>
					<Ionicons name={icon} size={20} color={theme.colors.text} />
				</View>
				<View style={styles.textDetails}>
					<ThemedText style={styles.title}>
						{merchant?.shortDescriptor ?? description}
					</ThemedText>
					<ThemedText style={styles.subtitle}>{typeLabel}</ThemedText>
				</View>
			</View>
			<View style={styles.rightSection}>
				<ThemedText style={styles.title}>
					{type === "deposit" && "+"}
					{formatAmount(Number(amount), currencyCode)}
				</ThemedText>
				<ThemedText
					style={[
						styles.subtitle,
						status === "declined" && { color: "#f00", opacity: 1 },
					]}
				>
					{time}
				</ThemedText>
			</View>
		</View>
	);
}
