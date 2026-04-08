import { useRouter } from "expo-router";
import { StyleSheet, Text, View } from "react-native";
import { SafeAreaView } from "react-native-safe-area-context";
import Background from "@/components/background";
import { ThemedButton } from "@/components/themed-button";

export default function AuthScreen() {
	const router = useRouter();

	return (
		<View style={styles.container}>
			<Background />
			<SafeAreaView style={styles.content}>
				<View style={styles.section}>
					<Text style={styles.title}>Welcome</Text>
					<Text style={styles.subtitle}>
						The bank you&apos;ll actually like. Experience banking that stays
						out of your way.
					</Text>
				</View>
				<View style={styles.section}>
					<ThemedButton
						backgroundColor="#fff"
						textColor="#000"
						onPress={() => router.push("/new-account")}
					>
						Open a New Account
					</ThemedButton>
					<ThemedButton
						onPress={() => router.push("/login")}
						backgroundColor="transparent"
						textColor="#fff"
					>
						Log In as a Returning Customer
					</ThemedButton>
				</View>
			</SafeAreaView>
		</View>
	);
}

const styles = StyleSheet.create({
	container: {
		flex: 1,
		backgroundColor: "#000",
	},
	content: {
		flex: 1,
		paddingHorizontal: 16,
		justifyContent: "flex-end",
		gap: 32,
	},
	section: {
		gap: 8,
	},
	title: {
		fontSize: 18,
		fontWeight: "bold",
		color: "#fff",
	},
	subtitle: {
		fontSize: 16,
		color: "#fff",
	},
});
