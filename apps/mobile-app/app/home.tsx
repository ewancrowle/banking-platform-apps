import { Ionicons } from "@expo/vector-icons";
import * as Haptics from "expo-haptics";
import { router } from "expo-router";
import { useEffect, useState } from "react";
import {
	Alert,
	FlatList,
	Pressable,
	StyleSheet,
	Text,
	useColorScheme,
	View,
} from "react-native";
import { SafeAreaView } from "react-native-safe-area-context";
import { getAccount } from "@/api/auth";
import trpc from "@/api/trpc";
import { ThemedButton } from "@/components/themed-button";
import { ThemedText } from "@/components/themed-text";
import {
	type Transaction,
	TransactionItem,
} from "@/components/transaction-item";
import { useAuthStore } from "@/store/auth";

export default function HomeScreen() {
	const colorScheme = useColorScheme();
	const { account, setAccount } = useAuthStore();
	const [balance, setBalance] = useState<string | null>(null);
	const [transactions, setTransactions] = useState<Transaction[]>([]);

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
			<FlatList
				data={transactions}
				keyExtractor={(item) => item.id}
				contentContainerStyle={styles.listContent}
				showsVerticalScrollIndicator={false}
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
								variant="ghost"
								style={{
									width: "auto",
								}}
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
									color={colorScheme === "dark" ? "#fff" : "#000"}
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
									fontSize: 32,
									fontWeight: "700",
								}}
							>
								{balance || "£0.00"}
							</ThemedText>
						</View>
						<View style={styles.buttonContainer}>
							<ThemedButton icon="add" style={styles.flexButton}>
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
				renderItem={({ item }) => (
					<View style={styles.itemWrapper}>
						<TransactionItem {...item} />
					</View>
				)}
			/>
		</SafeAreaView>
	);
}

const styles = StyleSheet.create({
	container: {
		flex: 1,
	},
	listContent: {
		flexGrow: 1,
		paddingBottom: 40,
	},
	headerContainer: {
		paddingHorizontal: 24,
		gap: 24,
		marginBottom: 12,
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
	itemWrapper: {
		paddingHorizontal: 24,
		marginVertical: 6,
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
