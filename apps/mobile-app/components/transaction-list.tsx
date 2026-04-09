import { useTheme } from "@react-navigation/native";
import { router } from "expo-router";
import type { Payment } from "protos/payment";
import type React from "react";
import { FlatList, Pressable, StyleSheet, View } from "react-native";
import { TransactionItem } from "./transaction-item";

interface TransactionListProps {
	transactions: Payment[];
	ListHeaderComponent?: React.ReactElement;
}

export function TransactionList({
	transactions,
	ListHeaderComponent,
}: TransactionListProps) {
	const theme = useTheme();

	return (
		<FlatList
			data={transactions}
			keyExtractor={(item) => item.id.toString()}
			contentContainerStyle={styles.listContent}
			showsVerticalScrollIndicator={false}
			ListHeaderComponent={ListHeaderComponent}
			ItemSeparatorComponent={() => (
				<View
					style={[styles.separator, { backgroundColor: theme.colors.border }]}
				/>
			)}
			renderItem={({ item }) => (
				<Pressable
					style={styles.itemWrapper}
					onPress={() =>
						router.push({
							pathname: "/transaction/[id]",
							params: {
								id: item.id.toString(),
							},
						})
					}
				>
					<TransactionItem {...item} />
				</Pressable>
			)}
		/>
	);
}

const styles = StyleSheet.create({
	listContent: {
		flexGrow: 1,
		paddingBottom: 40,
	},
	itemWrapper: {
		paddingHorizontal: 16,
		paddingVertical: 12,
	},
	separator: {
		height: 1,
		marginLeft: 16,
	},
});
