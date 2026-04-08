import { Ionicons } from "@expo/vector-icons";
import { useTheme } from "@react-navigation/native";
import dayjs from "dayjs";
import { Stack, useLocalSearchParams } from "expo-router";
import { ScrollView, StyleSheet, View } from "react-native";
import { Section } from "@/components/section";
import { ThemedText } from "@/components/themed-text";

export default function TransactionDetailsScreen() {
	const params = useLocalSearchParams();
	const theme = useTheme();

	const { id, amount, currencyCode, description, type, status, createdAt } =
		params;

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

	const getIconForType = (t: string): keyof typeof Ionicons.glyphMap => {
		switch (t) {
			case "deposit":
				return "arrow-down";
			case "withdrawal":
				return "arrow-up";
			case "card":
				return "card-outline";
			case "account_to_account":
				return "swap-horizontal-outline";
			default:
				return "cash-outline";
		}
	};

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
			maxWidth: "60%",
			textAlign: "right",
		},
	});

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
					<Ionicons
						name={getIconForType(type as string)}
						size={32}
						color={theme.colors.text}
					/>
				</View>
				<ThemedText style={styles.amount}>
					{type === "deposit" ? "+" : "-"}
					{formattedAmount}
				</ThemedText>
				<ThemedText style={styles.description}>{description}</ThemedText>
			</View>

			<Section title="Transaction information">
				<View style={styles.row}>
					<ThemedText style={styles.label}>Status</ThemedText>
					<ThemedText style={[styles.value, { color: statusColor }]}>
						{statusLabel}
					</ThemedText>
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
