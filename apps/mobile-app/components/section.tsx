import type { PropsWithChildren } from "react";
import { View } from "react-native";
import { ThemedText } from "./themed-text";

type SectionProps = PropsWithChildren<{
	title?: string;
}>;

export function Section(props: SectionProps) {
	return (
		<View style={{ marginTop: 16, gap: 12 }}>
			{props.title && (
				<ThemedText style={{ fontWeight: "bold", fontSize: 18 }}>
					{props.title}
				</ThemedText>
			)}
			{props.children}
		</View>
	);
}
