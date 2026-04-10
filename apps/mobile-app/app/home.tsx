import { useActionSheet } from "@expo/react-native-action-sheet";
import { Ionicons } from "@expo/vector-icons";
import { useTheme } from "@react-navigation/native";
import * as Haptics from "expo-haptics";
import { router } from "expo-router";
import type { GetTotalSpendingResponse } from "protos/ledger";
import { useEffect, useState } from "react";
import { Alert, Pressable, StyleSheet, Text, View } from "react-native";
import { SafeAreaView } from "react-native-safe-area-context";
import { getAccount } from "@/api/auth";
import trpc from "@/api/trpc";
import { Section } from "@/components/section";
import { ThemedButton } from "@/components/themed-button";
import { ThemedText } from "@/components/themed-text";
import { TransactionList } from "@/components/transaction-list";
import { useAuthStore } from "@/store/auth";
import { usePaymentsStore } from "@/store/payments";

const formatCurrency = (amount: bigint, currencyCode: string): string => {
	return new Intl.NumberFormat("en-GB", {
		style: "currency",
		currency: currencyCode || "GBP",
	}).format(Number(amount) / 100);
};

export default function Home() {
	const theme = useTheme();
	const { showActionSheetWithOptions } = useActionSheet();
	const { account, setAccount, reset } = useAuthStore();
	const { payments: transactions, setPayments } = usePaymentsStore();
	const [balance, setBalance] = useState<string | null>(null);
	const [totalSpent, setTotalSpent] = useState<{
		today: string;
		thisWeek: string;
		thisMonth: string;
	} | null>(null);

	const onPressHelp = () => {
		showActionSheetWithOptions(
			{
				options: ["Log out", "Cancel"],
				cancelButtonIndex: 2,
				destructiveButtonIndex: 0,
			},
			async (selectedIndex) => {
				if (selectedIndex === 0) {
					await reset();
				}
			},
		);
	};

	const onPressSpendMoney = () => {
		showActionSheetWithOptions(
			{
				options: ["New transfer", "New card payment", "Cancel"],
				cancelButtonIndex: 3,
			},
			(selectedIndex) => {
				if (selectedIndex === 0) {
					router.push("/new-transfer");
				} else if (selectedIndex === 1) {
					router.push("/new-card-payment");
				}
			},
		);
	};

	useEffect(() => {
		if (!account) {
			getAccount().then((acc) => {
				if (acc) setAccount(acc);
			});
		}
	}, [account, setAccount]);

	useEffect(() => {
		if (account) {
			trpc.balance.getBalances
				.query()
				.then((res) =>
					setBalance(formatCurrency(res.availableBalance, res.currencyCode)),
				)
				.catch((err) => {
					console.error(err);
					Alert.alert(
						"Your balance could not be loaded at this time. Please try again later.",
					);
				});

			trpc.payment.getPayments
				.query()
				.then((res) => {
					const validTypes = [
						"deposit",
						"withdrawal",
						"card",
						"account_to_account",
					];
					setPayments(res.payments.filter((p) => validTypes.includes(p.type)));
				})
				.catch((err) => {
					console.error("Failed to load transactions", err);
				});
		}
	}, [account, setPayments]);

	useEffect(() => {
		trpc.spending.getTotalSpending
			.query()
			.then((res) => {
				setTotalSpent({
					today: formatCurrency(res.totalSpentToday, res.currencyCode),
					thisWeek: formatCurrency(res.totalSpentThisWeek, res.currencyCode),
					thisMonth: formatCurrency(res.totalSpentThisMonth, res.currencyCode),
				});
			})
			.catch((err) => {
				console.error("Failed to load total spending", err);
			});
	}, []);

	const initials = account
		? `${account.firstName[0]}${account.lastName[0]}`.toUpperCase()
		: "??";

	const styles = StyleSheet.create({
		container: {
			flex: 1,
		},
		headerContainer: {
			paddingHorizontal: 16,
			gap: 24,
		},
		accountInfoContainer: {
			alignItems: "flex-start",
			gap: 2,
		},
		accountDetails: {
			flexDirection: "row",
			alignItems: "center",
			gap: 6,
		},
		balanceContainer: {
			alignItems: "center",
			gap: 8,
		},
		buttonContainer: {
			flexDirection: "row",
			gap: 8,
		},
		flexButton: {
			flex: 1,
			width: "auto",
		},
		topBar: {
			flexDirection: "row",
			justifyContent: "space-between",
			alignItems: "center",
		},
		profileCircle: {
			width: 48,
			height: 48,
			borderRadius: 24,
			backgroundColor: "#00F",
			justifyContent: "center",
			alignItems: "center",
		},
		initialsText: {
			fontSize: 18,
			color: "#fff",
		},
		helpButton: {
			flexDirection: "row",
			alignItems: "center",
			gap: 8,
		},
		buttonPressed: {
			opacity: 0.7,
			transform: [{ scale: 0.98 }],
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

	return (
		<SafeAreaView style={styles.container}>
			<TransactionList
				transactions={transactions}
				ListHeaderComponent={
					<View style={styles.headerContainer}>
						<View style={styles.topBar}>
							<Pressable
								onPressIn={() =>
									Haptics.impactAsync(Haptics.ImpactFeedbackStyle.Light)
								}
								style={({ pressed }) => [
									styles.profileCircle,
									pressed && styles.buttonPressed,
								]}
							>
								<Text style={styles.initialsText}>{initials}</Text>
							</Pressable>
							<ThemedButton
								icon="help-buoy-outline"
								textColor={theme.colors.text}
								backgroundColor="transparent"
								style={{
									width: "auto",
								}}
								onPress={onPressHelp}
							>
								Help
							</ThemedButton>
						</View>
						<View style={styles.accountInfoContainer}>
							<ThemedText
								style={{
									fontSize: 18,
									fontWeight: "600",
								}}
							>
								{account?.firstName}&apos;s Current Account
							</ThemedText>
							<View style={styles.accountDetails}>
								<Ionicons
									name="card-outline"
									size={16}
									color={theme.colors.text}
								/>
								<ThemedText>{account?.accountNum}</ThemedText>
							</View>
						</View>
						<View style={styles.balanceContainer}>
							<ThemedText
								style={{
									fontSize: 18,
									fontWeight: "600",
								}}
							>
								Available balance
							</ThemedText>
							<ThemedText
								style={{
									fontSize: 36,
									fontWeight: "700",
								}}
							>
								{balance || "£0.00"}
							</ThemedText>
						</View>
						<View style={styles.buttonContainer}>
							<ThemedButton
								icon="add"
								style={styles.flexButton}
								onPress={() => router.push("/new-deposit")}
							>
								Deposit money
							</ThemedButton>
							<ThemedButton
								icon="arrow-forward"
								style={styles.flexButton}
								onPress={onPressSpendMoney}
							>
								Spend money
							</ThemedButton>
						</View>

						{totalSpent && (
							<Section title="Your spending">
								<View style={styles.row}>
									<ThemedText style={styles.label}>Spent Today</ThemedText>
									<ThemedText style={styles.value}>
										{totalSpent.today}
									</ThemedText>
								</View>
								<View style={styles.row}>
									<ThemedText style={styles.label}>Spent This Week</ThemedText>
									<ThemedText style={styles.value}>
										{totalSpent.thisWeek}
									</ThemedText>
								</View>
								<View style={[styles.row, styles.noBorder]}>
									<ThemedText style={styles.label}>Spent This Month</ThemedText>
									<ThemedText style={styles.value}>
										{totalSpent.thisMonth}
									</ThemedText>
								</View>
							</Section>
						)}

						<ThemedText
							style={{
								fontSize: 18,
								fontWeight: "600",
							}}
						>
							Recent transactions
						</ThemedText>
					</View>
				}
			/>
		</SafeAreaView>
	);
}
