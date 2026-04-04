import { Ionicons } from "@expo/vector-icons";
import * as Haptics from "expo-haptics";
import { useEffect } from "react";
import {
	FlatList,
	Pressable,
	StyleSheet,
	Text,
	useColorScheme,
	View,
} from "react-native";
import { SafeAreaView } from "react-native-safe-area-context";
import { getAccount } from "@/api/auth";
import { ThemedButton } from "@/components/themed-button";
import { ThemedText } from "@/components/themed-text";
import { TransactionItem } from "@/components/transaction-item";

const transactions = [
	{
		id: "1",
		merchant: "Salary",
		location: "Workplace Ltd.",
		amount: "+£2,500.00",
		time: "Yesterday",
		icon: "cash-outline" as const,
	},
	{
		id: "2",
		merchant: "Starbucks",
		location: "The Hayes, Cardiff",
		amount: "-£5.50",
		time: "3 days ago",
		icon: "cafe-outline" as const,
	},
	{
		id: "3",
		merchant: "Transport for Wales",
		location: "Central Square, Cardiff",
		amount: "-£3.50",
		time: "4 days ago",
		icon: "train-outline" as const,
	},
];

import { router } from "expo-router";
import { useAuthStore } from "@/store/auth";

export default function HomeScreen() {
	const colorScheme = useColorScheme();
	const { account, setAccount } = useAuthStore();

	useEffect(() => {
		if (!account) {
			getAccount().then((acc) => {
				if (acc) setAccount(acc);
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
								£50.00
							</ThemedText>
						</View>
						<View style={styles.buttonContainer}>
							<ThemedButton icon="add" style={styles.flexButton}>
								Add money
							</ThemedButton>
							<ThemedButton
								icon="arrow-forward"
								style={styles.flexButton}
								onPress={() => router.push("/new-payment")}
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
