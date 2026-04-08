import { Ionicons } from "@expo/vector-icons";
import { useTheme } from "@react-navigation/native";
import * as Haptics from "expo-haptics";
import type { PropsWithChildren } from "react";
import { Pressable, type PressableProps, StyleSheet, Text } from "react-native";

type ThemedButtonProps = PropsWithChildren<PressableProps> & {
	icon?: keyof typeof Ionicons.glyphMap;
	textColor?: string;
	backgroundColor?: string;
};

export function ThemedButton({
	icon,
	style,
	children,
	onPressIn,
	textColor,
	backgroundColor,
	...rest
}: ThemedButtonProps) {
	const theme = useTheme();

	const styles = StyleSheet.create({
		container: {
			backgroundColor: backgroundColor ?? theme.colors.text,
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
			color: textColor ?? theme.colors.background,
			textAlign: "center",
			fontSize: 16,
			fontWeight: "600",
		},
	});

	return (
		<Pressable
			style={[styles.container, style as any]}
			onPressIn={(event) => {
				if (process.env.EXPO_OS === "ios") {
					Haptics.impactAsync(Haptics.ImpactFeedbackStyle.Light);
				}
				onPressIn?.(event);
			}}
			{...rest}
		>
			{icon && (
				<Ionicons
					name={icon}
					size={16}
					color={textColor ?? theme.colors.background}
					style={styles.icon}
				/>
			)}
			<Text style={styles.text}>{children}</Text>
		</Pressable>
	);
}
