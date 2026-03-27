import {
	Text,
	TextInput,
	type TextInputProps,
	useColorScheme,
} from "react-native";

export function ThemedInput({ style, ...props }: TextInputProps) {
	const colorScheme = useColorScheme();

	return (
		<TextInput
			style={[
				{
					color: colorScheme === "dark" ? "#fff" : "#000",
					backgroundColor: colorScheme === "dark" ? "#222" : "#fff",
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
