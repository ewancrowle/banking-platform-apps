import { useActionSheet } from "@expo/react-native-action-sheet";
import { Ionicons } from "@expo/vector-icons";
import { useTheme } from "@react-navigation/native";
import * as Haptics from "expo-haptics";
import { router } from "expo-router";
import { useEffect, useState } from "react";
import { Alert, Pressable, StyleSheet, Text, View } from "react-native";
import { SafeAreaView } from "react-native-safe-area-context";
import { getAccount } from "@/api/auth";
import trpc from "@/api/trpc";
import { ThemedButton } from "@/components/themed-button";
import { ThemedText } from "@/components/themed-text";
import type { Transaction } from "@/components/transaction-item";
import { TransactionList } from "@/components/transaction-list";
import { useAuthStore } from "@/store/auth";

export default function HomeScreen() {
	const theme = useTheme();
	const { showActionSheetWithOptions } = useActionSheet();
	const { account, setAccount, reset } = useAuthStore();
	const [balance, setBalance] = useState<string | null>(null);
	const [transactions, setTransactions] = useState<Transaction[]>([]);

	const onPressHelp = () => {
		showActionSheetWithOptions(
			{
				options: ["Log out"],
				destructiveButtonIndex: 0,
			},
			async () => {
				await reset();
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
				.then((res) => {
					const amount = Number(res.availableBalance) / 100;
					setBalance(
						new Intl.NumberFormat("en-GB", {
							style: "currency",
							currency: res.currencyCode || "GBP",
						}).format(amount),
					);
				})
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
					const mapped = res.payments
						.filter((p) => validTypes.includes(p.type))
						.map((p) => ({
							id: p.id.toString(),
							amount: Number(p.amount),
							currencyCode: p.currencyCode,
							description: p.description,
							type: p.type as Transaction["type"],
							status: p.status,
							createdAt: p.createdAt,
						}));
					setTransactions(mapped);
				})
				.catch((err) => {
					console.error("Failed to load transactions", err);
				});
		}
	}, [account]);

	const initials = account
		? `${account.firstName[0]}${account.lastName[0]}`.toUpperCase()
		: "??";

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
								<ThemedText>{account?.id.toString().slice(0, 8)}</ThemedText>
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
								Add money
							</ThemedButton>
							<ThemedButton
								icon="arrow-forward"
								style={styles.flexButton}
								onPress={() => router.push("/new-transfer")}
							>
								Pay someone
							</ThemedButton>
						</View>
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
});
