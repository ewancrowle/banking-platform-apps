import {
	Canvas,
	FractalNoise,
	Group,
	LinearGradient,
	Rect,
	vec,
} from "@shopify/react-native-skia";
import React from "react";
import { StyleSheet, useWindowDimensions, View } from "react-native";

export default function Background() {
	const { width, height } = useWindowDimensions();

	return (
		<View style={StyleSheet.absoluteFill}>
			<Canvas style={{ flex: 1 }}>
				<Rect x={0} y={0} width={width} height={height}>
					<LinearGradient
						start={vec(width / 2, 0)}
						end={vec(width / 2, height)}
						colors={["#00F", "#000"]}
						positions={[0.25, 0.75]}
					/>

					<Group blendMode="multiply" opacity={0.5}>
						<Rect x={0} y={0} width={width} height={height}>
							<FractalNoise freqX={0.9} freqY={0.9} octaves={2} seed={1} />
						</Rect>
					</Group>
				</Rect>
			</Canvas>
		</View>
	);
}
