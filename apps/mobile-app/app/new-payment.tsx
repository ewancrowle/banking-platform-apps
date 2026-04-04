import { formOptions, useForm, useStore } from "@tanstack/react-form";
import { Keyboard, Pressable, ScrollView, useColorScheme } from "react-native";
import CurrencyInput from "react-native-currency-input";
import * as z from "zod";
import { Section } from "@/components/section";
import { ThemedInput } from "@/components/themed-input";
import { ThemedText } from "@/components/themed-text";

const formSchema = z.object({
	firstName: z.string().min(1, {
		message: "Please enter your first name.",
	}),
	lastName: z.string().min(1, {
		message: "Please enter your last name.",
	}),
	accountNumber: z
		.string()
		.min(1, {
			message: "Please enter the payee's account number.",
		})
		.max(8, {
			message: "Please enter a valid account number.",
		})
		.refine(
			(val) => {
				if (val.length !== 8 || !/^\d+$/.test(val)) return false;

				const digits = val.split("").map(Number);
				const checkDigit = digits[7];
				const payload = digits.slice(0, 7);

				let sum = 0;
				let double = true;

				for (let i = payload.length - 1; i >= 0; i--) {
					let d = payload[i];

					if (double) {
						d *= 2;
						if (d > 9) {
							d -= 9;
						}
					}

					sum += d;
					double = !double;
				}

				const expected = (10 - (sum % 10)) % 10;
				return checkDigit === expected;
			},
			{ message: "Please enter a valid account number." },
		),
	amount: z.number().min(1, {
		message: "Please enter a valid amount.",
	}),
	reference: z.string().optional(),
});

const formOpts = formOptions({
	defaultValues: {
		firstName: "",
		lastName: "",
		accountNumber: "",
		amount: 0,
		reference: "",
	} as z.infer<typeof formSchema>,
	validators: {
		onChange: formSchema,
	},
});

export default function NewPaymentScreen() {
	const colorScheme = useColorScheme();

	const form = useForm({
		...formOpts,
	});

	const formErrorMap = useStore(form.store, (state) => state.errorMap);

	return (
		<Pressable style={{ flex: 1, padding: 16 }} onPress={Keyboard.dismiss}>
			<ScrollView showsVerticalScrollIndicator={false}>
				<ThemedText>Transfer money to someone you know.</ThemedText>

				<Section title="Enter the payee's name">
					<form.Field name="firstName">
						{(field) => (
							<>
								<ThemedInput
									placeholder="First Name"
									value={field.state.value}
									onChangeText={field.handleChange}
									textContentType="givenName"
									autoCapitalize="words"
									autoComplete="name-given"
									autoCorrect={false}
									keyboardType="default"
								/>
								{formErrorMap.onChange?.firstName && (
									<ThemedText style={{ color: "red", marginTop: 4 }}>
										{formErrorMap.onChange.firstName
											.map((issue) => issue.message)
											.join(", ")}
									</ThemedText>
								)}
							</>
						)}
					</form.Field>

					<form.Field name="lastName">
						{(field) => (
							<>
								<ThemedInput
									placeholder="Last Name"
									value={field.state.value}
									onChangeText={field.handleChange}
									textContentType="familyName"
									autoCapitalize="words"
									autoComplete="name-family"
									autoCorrect={false}
									keyboardType="default"
								/>
								{formErrorMap.onChange?.lastName && (
									<ThemedText style={{ color: "red", marginTop: 4 }}>
										{formErrorMap.onChange.lastName
											.map((issue) => issue.message)
											.join(", ")}
									</ThemedText>
								)}
							</>
						)}
					</form.Field>
				</Section>

				<Section title="Enter the payee's account info">
					<form.Field name="accountNumber" asyncDebounceMs={300}>
						{(field) => (
							<>
								<ThemedInput
									placeholder="Account Number"
									value={field.state.value}
									onChangeText={field.handleChange}
									autoCorrect={false}
									keyboardType="numeric"
								/>
								{formErrorMap.onChange?.accountNumber && (
									<ThemedText style={{ color: "red", marginTop: 4 }}>
										{formErrorMap.onChange.accountNumber
											.map((issue) => issue.message)
											.join(", ")}
									</ThemedText>
								)}
							</>
						)}
					</form.Field>
				</Section>

				<Section title="Enter the payment info">
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
										color: colorScheme === "dark" ? "#fff" : "#000",
										backgroundColor: colorScheme === "dark" ? "#222" : "#fff",
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
					<form.Field name="reference">
						{(field) => (
							<>
								<ThemedInput
									placeholder="What&apos;s this payment for?"
									value={field.state.value}
									onChangeText={field.handleChange}
									autoCorrect={false}
									keyboardType="default"
								/>
								{formErrorMap.onChange?.reference && (
									<ThemedText style={{ color: "red", marginTop: 4 }}>
										{formErrorMap.onChange.reference
											.map((issue) => issue.message)
											.join(", ")}
									</ThemedText>
								)}
							</>
						)}
					</form.Field>
				</Section>
			</ScrollView>
		</Pressable>
	);
}
