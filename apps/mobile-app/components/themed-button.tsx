import { Ionicons } from "@expo/vector-icons";
import * as Haptics from "expo-haptics";
import type { PropsWithChildren } from "react";
import {
	Pressable,
	type PressableProps,
	StyleSheet,
	Text,
	useColorScheme,
	View,
} from "react-native";

type ThemedButtonProps = PropsWithChildren<PressableProps> & {
	icon?: keyof typeof Ionicons.glyphMap;
	variant?: "light" | "dark" | "ghost";
	textColor?: string;
	backgroundColor?: string;
};

export function ThemedButton({
	icon,
	style,
	children,
	onPressIn,
	variant,
	textColor: color,
	backgroundColor: background,
	...rest
}: ThemedButtonProps) {
	const colorScheme = useColorScheme();

	const backgroundColor = (() => {
		if (background) return background;
		switch (variant) {
			case "light":
				return "#000";
			case "dark":
				return "#fff";
			case "ghost":
				return "transparent";
			default:
				return colorScheme === "dark" ? "#fff" : "#000";
		}
	})();

	const textColor = (() => {
		if (color) return color;
		switch (variant) {
			case "light":
				return "#fff";
			case "dark":
				return "#000";
			case "ghost":
				return colorScheme === "dark" ? "#fff" : "#000";
			default:
				return colorScheme === "dark" ? "#000" : "#fff";
		}
	})();

	return (
		<Pressable
			style={[styles.container, { backgroundColor }, style as any]}
			onPressIn={(event) => {
				if (process.env.EXPO_OS === "ios") {
					Haptics.impactAsync(Haptics.ImpactFeedbackStyle.Light);
				}
				onPressIn?.(event);
			}}
			{...rest}>
			{icon && (
				<Ionicons
					name={icon}
					size={16}
					color={textColor}
					style={styles.icon}
				/>
			)}
			<Text style={[{ color: textColor }, styles.text]}>{children}</Text>
		</Pressable>
	);
}

const styles = StyleSheet.create({
	container: {
		borderRadius: 999,
		paddingVertical: 14,
		alignItems: "center",
		flexDirection: "row",
		justifyContent: "center",
		width: "100%",
	},
	icon: {
		marginRight: 8,
	},
	text: {
		textAlign: "center",
		fontSize: 16,
		fontWeight: "600",
	},
});
