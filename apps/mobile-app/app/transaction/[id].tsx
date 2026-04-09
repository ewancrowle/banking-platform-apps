import { Ionicons } from "@expo/vector-icons";
import { useTheme } from "@react-navigation/native";
import dayjs from "dayjs";
import { Stack, useLocalSearchParams } from "expo-router";
import { DeclineReason } from "protos/payment";
import { ScrollView, StyleSheet, View } from "react-native";
import { Section } from "@/components/section";
import { ThemedText } from "@/components/themed-text";
import { usePaymentsStore } from "@/store/payments";
import { getPaymentIcon } from "@/utils/get-payment-icon";

export default function TransactionInfo() {
	const theme = useTheme();

	const { id } = useLocalSearchParams();
	const { payments } = usePaymentsStore();
	const transaction = payments.find((p) => p.id.toString() === id);

	if (!transaction) {
		return (
			<View style={{ flex: 1, justifyContent: "center", alignItems: "center" }}>
				<Stack.Screen options={{ title: "Transaction Details" }} />
				<ThemedText>Transaction not found</ThemedText>
			</View>
		);
	}

	const {
		amount,
		currencyCode,
		createdAt,
		type,
		status,
		description,
		merchant,
		declineReason,
	} = transaction;

	const formattedAmount = new Intl.NumberFormat("en-GB", {
		style: "currency",
		currency: (currencyCode as string) || "GBP",
	}).format(Number(amount) / 100);

	const time = dayjs(`${(createdAt as string).split(".")[0]}Z`).format(
		"D MMMM YYYY, HH:mm",
	);

	const typeLabel =
		{
			deposit: "Deposit",
			withdrawal: "Withdrawal",
			card: "Card Payment",
			account_to_account: "Transfer",
		}[type as string] || type;

	const statusLabel =
		{
			declined: "Declined",
			authorised: "Pending",
			captured: "Complete",
		}[status as string] ||
		(status
			? (status as string).charAt(0).toUpperCase() + (status as string).slice(1)
			: "");

	const icon = getPaymentIcon(transaction);

	let statusColor = theme.colors.text;
	switch (status) {
		case "declined":
			statusColor = "#ff0000";
			break;
		case "authorised":
			statusColor = "#ffa500";
			break;
		case "captured":
			statusColor = "#008000";
			break;
	}

	const styles = StyleSheet.create({
		header: {
			alignItems: "center",
			marginTop: 32,
			marginBottom: 16,
			gap: 16,
			paddingRight: 16,
		},
		iconContainer: {
			width: 64,
			height: 64,
			borderRadius: 32,
			justifyContent: "center",
			alignItems: "center",
		},
		amount: {
			fontSize: 36,
			fontWeight: "700",
		},
		description: {
			fontSize: 18,
			opacity: 0.7,
		},
		row: {
			flexDirection: "row",
			justifyContent: "space-between",
			alignItems: "center",
			paddingBottom: 12,
			paddingRight: 16,
			borderBottomWidth: 1,
			borderBottomColor: theme.colors.border,
		},
		noBorder: {
			borderBottomWidth: 0,
		},
		label: {
			fontSize: 16,
			fontWeight: "500",
			opacity: 0.7,
		},
		value: {
			fontSize: 16,
			textAlign: "right",
		},
	});

	const getMerchantAddress = () => {
		if (!merchant) return undefined;
		return [merchant.line1, merchant.line2, merchant.town, merchant.postcode]
			.filter(Boolean)
			.join(", ");
	};

	return (
		<ScrollView
			style={{ flex: 1, paddingLeft: 16 }}
			showsVerticalScrollIndicator={false}
			contentContainerStyle={{ paddingBottom: 40 }}
		>
			<Stack.Screen options={{ title: "Transaction Details" }} />

			<View style={styles.header}>
				<View
					style={[styles.iconContainer, { backgroundColor: theme.colors.card }]}
				>
					<Ionicons name={icon} size={32} color={theme.colors.text} />
				</View>
				<ThemedText style={styles.amount}>
					{type === "deposit" && "+"}
					{formattedAmount}
				</ThemedText>
				<ThemedText style={styles.description}>{description}</ThemedText>
			</View>

			{merchant && (
				<Section title="Merchant">
					<View style={styles.row}>
						<ThemedText style={styles.label}>Merchant</ThemedText>
						<ThemedText style={styles.value}>{merchant.descriptor}</ThemedText>
					</View>
					<View style={styles.row}>
						<ThemedText style={styles.label}>Trading Name</ThemedText>
						<ThemedText style={styles.value}>
							{merchant.shortDescriptor}
						</ThemedText>
					</View>
					<View style={styles.row}>
						<ThemedText style={styles.label}>Address</ThemedText>
						<ThemedText style={styles.value}>{getMerchantAddress()}</ThemedText>
					</View>
				</Section>
			)}

			<Section title="Transaction information">
				<View style={styles.row}>
					<ThemedText style={styles.label}>Status</ThemedText>
					{declineReason === DeclineReason.INSUFFICIENT_FUNDS ? (
						<ThemedText style={[styles.value, { color: statusColor }]}>
							{statusLabel}
						</ThemedText>
					) : (
						<ThemedText style={[styles.value, { color: statusColor }]}>
							{statusLabel} Due to Insufficient Funds
						</ThemedText>
					)}
				</View>
				<View style={styles.row}>
					<ThemedText style={styles.label}>Type</ThemedText>
					<ThemedText style={styles.value}>{typeLabel}</ThemedText>
				</View>
				<View style={styles.row}>
					<ThemedText style={styles.label}>Date</ThemedText>
					<ThemedText style={styles.value}>{time}</ThemedText>
				</View>
				<View style={[styles.row, styles.noBorder]}>
					<ThemedText style={styles.label}>Transaction ID</ThemedText>
					<ThemedText
						style={styles.value}
						numberOfLines={1}
						ellipsizeMode="middle"
					>
						{id}
					</ThemedText>
				</View>
			</Section>
		</ScrollView>
	);
}
