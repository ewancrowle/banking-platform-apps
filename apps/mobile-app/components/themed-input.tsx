import { useTheme } from "@react-navigation/native";
import { TextInput, type TextInputProps } from "react-native";

export function ThemedInput({ style, ...props }: TextInputProps) {
	const theme = useTheme();

	return (
		<TextInput
			style={[
				{
					color: theme.colors.text,
					backgroundColor: theme.colors.card,
					borderRadius: 8,
					fontSize: 16,
					padding: 12,
				},
				style,
			]}
			{...props}
		/>
	);
}
