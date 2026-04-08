import { useTheme } from "@react-navigation/native";
import { Text, type TextProps } from "react-native";

export function ThemedText({ style, ...props }: TextProps) {
	const theme = useTheme();

	return (
		<Text
			style={[
				{
					color: theme.colors.text,
					fontSize: 16,
				},
				style,
			]}
			{...props}
		/>
	);
}
