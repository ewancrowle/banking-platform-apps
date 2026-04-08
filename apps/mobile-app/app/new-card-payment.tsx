import { useActionSheet } from "@expo/react-native-action-sheet";
import { useTheme } from "@react-navigation/native";
import { formOptions, useForm, useStore } from "@tanstack/react-form";
import { router } from "expo-router";
import type { Merchant } from "protos/merchant";
import { Decision } from "protos/payment";
import { useEffect, useState } from "react";
import { Alert, Keyboard, Pressable } from "react-native";
import CurrencyInput from "react-native-currency-input";
import * as z from "zod";
import trpc from "@/api/trpc";
import { Section } from "@/components/section";
import { ThemedButton } from "@/components/themed-button";
import { ThemedText } from "@/components/themed-text";

const formSchema = z.object({
	merchantId: z.bigint().min(BigInt(1), {
		message: "Please select a merchant.",
	}),
	amount: z.number().min(0.01, {
		message: "Please enter a valid amount.",
	}),
});

const formOpts = formOptions({
	defaultValues: {
		merchantId: BigInt(0),
		amount: 0,
	} as z.infer<typeof formSchema>,
	validators: {
		onChange: formSchema,
	},
});

export default function NewCardPayment() {
	const theme = useTheme();
	const { showActionSheetWithOptions } = useActionSheet();
	const [merchants, setMerchants] = useState<Merchant[]>([]);

	useEffect(() => {
		const fetchMerchants = async () => {
			const { merchants } = await trpc.merchant.getAllMerchants.query();
			setMerchants(merchants);
		};
		fetchMerchants();
	}, []);

	const form = useForm({
		...formOpts,
		onSubmit: async ({ value }) => {
			try {
				const payment = await trpc.payment.newCardPayment.mutate({
					merchantId: value.merchantId,
					amount: value.amount * 100,
				});
				if (payment.decision === Decision.DECLINED) {
					Alert.alert("Payment declined. Please try again later.");
					return;
				}
				Alert.alert("Payment successful.");
				form.reset();
				router.replace("/home");
			} catch (err) {
				console.log(err);
				Alert.alert("An error occurred. Please try again later.");
				return;
			}
		},
	});

	const formErrorMap = useStore(form.store, (state) => state.errorMap);

	return (
		<Pressable style={{ flex: 1, padding: 16 }} onPress={Keyboard.dismiss}>
			<ThemedText>Spend your money right now.</ThemedText>

			<Section title="Select a merchant">
				<form.Field name="merchantId">
					{(field) => {
						const selectedMerchant = merchants.find(
							(m) => m.id === field.state.value,
						);
						return (
							<>
								<Pressable
									style={{
										backgroundColor: theme.colors.card,
										padding: 16,
										borderRadius: 8,
									}}
									onPress={() => {
										const options = [
											...merchants.map((m) => m.shortDescriptor),
											"Cancel",
										];
										const cancelButtonIndex = merchants.length;

										showActionSheetWithOptions(
											{
												options,
												cancelButtonIndex,
											},
											(selectedIndex) => {
												if (
													selectedIndex !== undefined &&
													selectedIndex !== cancelButtonIndex
												) {
													field.handleChange(merchants[selectedIndex].id);
												}
											},
										);
									}}
								>
									<ThemedText>
										{selectedMerchant
											? selectedMerchant.shortDescriptor
											: "Tap to select merchant..."}
									</ThemedText>
								</Pressable>
								{formErrorMap.onChange?.merchantId && (
									<ThemedText style={{ color: "red", marginTop: 4 }}>
										{formErrorMap.onChange.merchantId
											.map((issue) => issue.message)
											.join(", ")}
									</ThemedText>
								)}
							</>
						);
					}}
				</form.Field>
			</Section>

			<Section title="Enter the payment amount">
				<form.Field name="amount">
					{(field) => (
						<>
							<CurrencyInput
								value={field.state.value}
								onChangeValue={(value) => field.handleChange(value || 0)}
								prefix="£"
								separator="."
								precision={2}
								placeholder="Amount"
								style={{
									color: theme.colors.text,
									backgroundColor: theme.colors.card,
									borderRadius: 8,
									fontSize: 16,
									padding: 12,
								}}
							/>
							{formErrorMap.onChange?.amount && (
								<ThemedText style={{ color: "red", marginTop: 4 }}>
									{formErrorMap.onChange.amount
										.map((issue) => issue.message)
										.join(", ")}
								</ThemedText>
							)}
						</>
					)}
				</form.Field>
			</Section>

			<Section>
				<ThemedButton onPress={() => form.handleSubmit()}>Pay</ThemedButton>
			</Section>
		</Pressable>
	);
}
