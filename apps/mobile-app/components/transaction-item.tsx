import { Ionicons } from "@expo/vector-icons";
import {StyleSheet, Text, useColorScheme, View} from "react-native";
import {ThemedText} from "@/components/themed-text";

type TransactionItemProps = {
	merchant: string;
	location: string;
	amount: string;
	time: string;
	icon: keyof typeof Ionicons.glyphMap;
};

export function TransactionItem({
	merchant,
	location,
	amount,
	time,
	icon,
}: TransactionItemProps) {
	const colorScheme = useColorScheme();

	return (
		<View style={styles.container}>
			<View style={styles.leftSection}>
				<View style={styles.iconContainer}>
					<Ionicons name={icon} size={24} color={colorScheme === "dark" ? "#fff" : "#000"} />
				</View>
				<View style={styles.textDetails}>
					<ThemedText style={styles.title}>{merchant}</ThemedText>
					<ThemedText style={styles.subtitle}>{location}</ThemedText>
				</View>
			</View>
			<View style={styles.rightSection}>
				<ThemedText style={styles.title}>{amount}</ThemedText>
				<ThemedText style={styles.subtitle}>{time}</ThemedText>
			</View>
		</View>
	);
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
		width: 48,
		height: 48,
		borderRadius: 24,
		justifyContent: "center",
		alignItems: "center",
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
		fontSize: 16,
		fontWeight: "600",
	},
});
