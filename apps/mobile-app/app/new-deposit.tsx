import { useTheme } from "@react-navigation/native";
import { formOptions, useForm, useStore } from "@tanstack/react-form";
import { router } from "expo-router";
import { Decision } from "protos/payment";
import { Alert, Keyboard, Pressable, ScrollView } from "react-native";
import CurrencyInput from "react-native-currency-input";
import * as z from "zod";
import trpc from "@/api/trpc";
import { Section } from "@/components/section";
import { ThemedButton } from "@/components/themed-button";
import { ThemedText } from "@/components/themed-text";
import { useBalanceStore } from "@/store/balance";
import { usePaymentsStore } from "@/store/payments";
import { useSpendingStore } from "@/store/spending";

const formSchema = z.object({
	amount: z.number().min(0.01, {
		message: "Please enter a valid amount.",
	}),
});

const formOpts = formOptions({
	defaultValues: {
		amount: 0,
	} as z.infer<typeof formSchema>,
	validators: {
		onChange: formSchema,
	},
});

export default function NewDeposit() {
	const theme = useTheme();

	const { refresh: refreshPayments } = usePaymentsStore();
	const { refresh: refreshBalance } = useBalanceStore();
	const { refresh: refreshSpending } = useSpendingStore();

	const form = useForm({
		...formOpts,
		onSubmit: async ({ value }) => {
			try {
				const payment = await trpc.payment.newDeposit.mutate({
					amount: value.amount * 100,
				});

				if (payment.decision === Decision.DECLINED) {
					Alert.alert("Deposit declined. Please try again later.");
					return;
				}
				Alert.alert("Deposit successful.");

				await refreshPayments();
				await refreshBalance();
				refreshSpending();

				form.reset();
				router.back();
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
			<ThemedText>Add money to your account.</ThemedText>

			<Section title="Enter the deposit amount">
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
				<ThemedButton onPress={() => form.handleSubmit()}>Deposit</ThemedButton>
			</Section>
		</Pressable>
	);
}
