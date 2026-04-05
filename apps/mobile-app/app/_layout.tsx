import { Stack, useRouter, useSegments } from "expo-router";
import { getItemAsync } from "expo-secure-store";
import * as SplashScreen from "expo-splash-screen";
import { StatusBar } from "expo-status-bar";
import { useEffect, useState } from "react";
import { getAccount, refreshTokens } from "@/api/auth";
import { useAuthStore } from "@/store/auth";
import "react-native-reanimated";
import {
	DarkTheme,
	DefaultTheme,
	ThemeProvider,
} from "@react-navigation/native";
import { Alert, useColorScheme } from "react-native";

// Keep the splash screen visible while we fetch resources
SplashScreen.preventAutoHideAsync();

export const unstable_settings = {
	initialRouteName: "(auth)/index",
};

export default function RootLayout() {
	const colorScheme = useColorScheme();
	const [isReady, setIsReady] = useState(false);
	const { accessToken, setTokens, setAccount, reset } = useAuthStore();
	const segments = useSegments();
	const router = useRouter();

	useEffect(() => {
		const initAuth = async () => {
			try {
				const storedAccessToken = await getItemAsync("accessToken");
				const storedRefreshToken = await getItemAsync("refreshToken");

				if (storedAccessToken && storedRefreshToken) {
					useAuthStore.setState({
						accessToken: storedAccessToken,
						refreshToken: storedRefreshToken,
					});

					let account = await getAccount();

					if (!account) {
						const tokens = await refreshTokens();
						if (tokens) {
							await setTokens(tokens);
							account = await getAccount();
						}

						await reset();
					}

					if (account) {
						setAccount(account);
					}
				}
			} catch (err) {
				console.error("Auth init error", err);
				Alert.alert("An error occurred. Please try again later.");
			}

			setIsReady(true);
			await SplashScreen.hideAsync();
		};

		initAuth();
	}, []);

	useEffect(() => {
		if (!isReady) return;

		const inAuthGroup = segments[0] === "(auth)";

		if (accessToken && inAuthGroup) {
			router.replace("/home");
		} else if (!accessToken && !inAuthGroup) {
			router.replace("/(auth)");
		}
	}, [accessToken, isReady, segments]);

	if (!isReady) return null;

	return (
		<ThemeProvider value={colorScheme === "dark" ? DarkTheme : DefaultTheme}>
			<Stack>
				<Stack.Screen
					name="(auth)/index"
					options={{ title: "", headerShown: false }}
				/>
				<Stack.Screen
					name="(auth)/new-account"
					options={{
						title: "Open a New Account",
					}}
				/>
				<Stack.Screen
					name="(auth)/login"
					options={{
						title: "Log In",
					}}
				/>
				<Stack.Screen name="home" options={{ headerShown: false }} />
				<Stack.Screen
					name="new-transfer"
					options={{
						title: "New Transfer",
					}}
				/>
			</Stack>
			<StatusBar style="auto" />
		</ThemeProvider>
	);
}
