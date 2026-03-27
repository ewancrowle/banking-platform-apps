import { StyleSheet, Text, type TextProps, useColorScheme } from "react-native";

export function ThemedText({ style, ...props }: TextProps) {
	const colorScheme = useColorScheme();

	return (
		<Text
			style={[
				{
					color: colorScheme === "dark" ? "#fff" : "#000",
					fontSize: 16,
				},
				style,
			]}
			{...props}
		/>
	);
}
