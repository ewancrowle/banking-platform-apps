import DateTimePicker from "@react-native-community/datetimepicker";
import { formOptions, useForm, useStore } from "@tanstack/react-form";
import { useState } from "react";
import { Alert, Keyboard, Pressable, ScrollView } from "react-native";
import * as z from "zod";
import trpc from "@/api/trpc";
import { Section } from "@/components/section";
import { ThemedButton } from "@/components/themed-button";
import { ThemedInput } from "@/components/themed-input";
import { ThemedText } from "@/components/themed-text";
import { useAuthStore } from "@/store/auth";

const formSchema = z
	.object({
		firstName: z.string().min(1, {
			message: "Please enter your first name.",
		}),
		middleNames: z.string().optional(),
		lastName: z.string().min(1, {
			message: "Please enter your last name.",
		}),
		email: z.email({
			message: "Please enter a valid email address.",
		}),
		phoneNumber: z.string().min(1, {
			message: "Please enter your phone number.",
		}),
		password: z
			.string()
			.regex(
				/^(?=.*?[A-Z])(?=.*?[a-z])(?=.*?[0-9])(?=.*?[#?!@$ %^&*-]).{8,}$/,
				{
					message:
						"Your password must be at least 8 characters and include uppercase, lowercase, a number, and a special character.",
				},
			),
		confirmPassword: z.string().min(1, {
			message: "Please confirm your password.",
		}),
		line1: z.string().min(1, {
			message: "Please enter your address line 1.",
		}),
		line2: z.string().optional(),
		town: z.string().min(1, {
			message: "Please enter your town.",
		}),
		postcode: z.string().min(1, {
			message: "Please enter your postcode.",
		}),
	})
	.refine((arg) => arg.password === arg.confirmPassword, {
		message: "Passwords do not match.",
		path: ["confirmPassword"],
	});

const formOpts = formOptions({
	defaultValues: {
		firstName: "",
		middleNames: "",
		lastName: "",
		email: "",
		phoneNumber: "",
		password: "",
		confirmPassword: "",
		line1: "",
		line2: "",
		town: "",
		postcode: "",
	} as z.infer<typeof formSchema>,
	validators: {
		onChange: formSchema,
	},
});

export default function NewAccount() {
	const form = useForm({
		...formOpts,
		onSubmit: async ({ value }) => {
			try {
				const tokens = await trpc.auth.signUp.mutate(value);
				const store = useAuthStore.getState();
				await store.setTokens(tokens);
				const account = await trpc.account.getAccount.query();
				store.setAccount(account);
			} catch (err) {
				Alert.alert("An error occurred. Please try again later.");
				console.log(err);
			}
		},
	});

	const [date, setDate] = useState(new Date());

	const formErrorMap = useStore(form.store, (state) => state.errorMap);

	return (
		<Pressable style={{ flex: 1, padding: 16 }} onPress={Keyboard.dismiss}>
			<ScrollView showsVerticalScrollIndicator={false}>
				<ThemedText>Let&apos;s get you set up with a new account.</ThemedText>

				<Section title="Enter your name">
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

					<form.Field name="middleNames">
						{(field) => (
							<>
								<ThemedInput
									placeholder="Middle Names"
									value={field.state.value}
									onChangeText={field.handleChange}
									textContentType="middleName"
									autoCapitalize="words"
									autoComplete="name-middle"
									autoCorrect={false}
									keyboardType="default"
								/>
								{formErrorMap.onChange?.middleNames && (
									<ThemedText style={{ color: "red", marginTop: 4 }}>
										{formErrorMap.onChange.middleNames
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

				<Section title="Enter your date of birth">
					<ThemedText>
						You must be at least 18 years old to open an account.
					</ThemedText>
					<DateTimePicker
						value={date}
						mode="date"
						display="spinner"
						onChange={(_, selectedDate) => setDate(selectedDate ?? new Date())}
						maximumDate={
							new Date(
								new Date().getFullYear() - 18,
								new Date().getMonth(),
								new Date().getDate(),
							)
						}
					/>
				</Section>

				<Section title="Enter your contact info">
					<form.Field name="email">
						{(field) => (
							<>
								<ThemedInput
									placeholder="Email Address"
									value={field.state.value}
									onChangeText={field.handleChange}
									textContentType="emailAddress"
									autoCapitalize="none"
									autoComplete="email"
									autoCorrect={false}
									keyboardType="email-address"
								/>
								{formErrorMap.onChange?.email && (
									<ThemedText style={{ color: "red", marginTop: 4 }}>
										{formErrorMap.onChange.email
											.map((issue) => issue.message)
											.join(", ")}
									</ThemedText>
								)}
							</>
						)}
					</form.Field>

					<form.Field name="phoneNumber">
						{(field) => (
							<>
								<ThemedInput
									placeholder="Phone Number"
									value={field.state.value}
									onChangeText={field.handleChange}
									textContentType="telephoneNumber"
									autoComplete="tel"
									autoCorrect={false}
									keyboardType="phone-pad"
								/>
								{formErrorMap.onChange?.phoneNumber && (
									<ThemedText style={{ color: "red", marginTop: 4 }}>
										{formErrorMap.onChange.phoneNumber
											.map((issue) => issue.message)
											.join(", ")}
									</ThemedText>
								)}
							</>
						)}
					</form.Field>
				</Section>

				<Section title="Enter your address">
					<form.Field name="line1">
						{(field) => (
							<>
								<ThemedInput
									placeholder="Line 1"
									value={field.state.value}
									onChangeText={field.handleChange}
									textContentType="streetAddressLine1"
									autoComplete="address-line1"
									autoCorrect={false}
									keyboardType="default"
								/>
								{formErrorMap.onChange?.line1 && (
									<ThemedText style={{ color: "red", marginTop: 4 }}>
										{formErrorMap.onChange.line1
											.map((issue) => issue.message)
											.join(", ")}
									</ThemedText>
								)}
							</>
						)}
					</form.Field>

					<form.Field name="line2">
						{(field) => (
							<>
								<ThemedInput
									placeholder="Line 2"
									value={field.state.value}
									onChangeText={field.handleChange}
									textContentType="streetAddressLine2"
									autoComplete="address-line2"
									autoCorrect={false}
									keyboardType="default"
								/>
								{formErrorMap.onChange?.line2 && (
									<ThemedText style={{ color: "red", marginTop: 4 }}>
										{formErrorMap.onChange.line2
											.map((issue) => issue.message)
											.join(", ")}
									</ThemedText>
								)}
							</>
						)}
					</form.Field>

					<form.Field name="town">
						{(field) => (
							<>
								<ThemedInput
									placeholder="Town"
									value={field.state.value}
									onChangeText={field.handleChange}
									textContentType="addressCity"
									autoComplete="postal-address-locality"
									autoCorrect={false}
									keyboardType="default"
								/>
								{formErrorMap.onChange?.town && (
									<ThemedText style={{ color: "red", marginTop: 4 }}>
										{formErrorMap.onChange.town
											.map((issue) => issue.message)
											.join(", ")}
									</ThemedText>
								)}
							</>
						)}
					</form.Field>

					<form.Field name="postcode">
						{(field) => (
							<>
								<ThemedInput
									placeholder="Postcode"
									value={field.state.value}
									onChangeText={field.handleChange}
									textContentType="postalCode"
									autoComplete="postal-code"
									autoCorrect={false}
									keyboardType="default"
								/>
								{formErrorMap.onChange?.postcode && (
									<ThemedText style={{ color: "red", marginTop: 4 }}>
										{formErrorMap.onChange.postcode
											.map((issue) => issue.message)
											.join(", ")}
									</ThemedText>
								)}
							</>
						)}
					</form.Field>
				</Section>

				<Section title="Pick a strong password">
					<form.Field name="password">
						{(field) => (
							<>
								<ThemedInput
									placeholder="Password"
									value={field.state.value}
									onChangeText={field.handleChange}
									textContentType="password"
									secureTextEntry={true}
									autoCapitalize="none"
									autoComplete="password"
									autoCorrect={false}
									keyboardType="visible-password"
								/>
								{formErrorMap.onChange?.password && (
									<ThemedText style={{ color: "red", marginTop: 4 }}>
										{formErrorMap.onChange.password
											.map((issue) => issue.message)
											.join(", ")}
									</ThemedText>
								)}
							</>
						)}
					</form.Field>

					<form.Field
						name="confirmPassword"
						validators={{
							onSubmit: ({ value }) => {
								if (value !== form.state.values.password)
									return "Passwords do not match";
							},
						}}
					>
						{(field) => (
							<>
								<ThemedInput
									placeholder="Confirm Password"
									value={field.state.value}
									onChangeText={field.handleChange}
									textContentType="password"
									secureTextEntry={true}
									autoCapitalize="none"
									autoComplete="password"
									autoCorrect={false}
									keyboardType="visible-password"
								/>
								{formErrorMap.onChange?.confirmPassword && (
									<ThemedText style={{ color: "red", marginTop: 4 }}>
										{formErrorMap.onChange.confirmPassword
											.map((issue) => issue.message)
											.join(", ")}
									</ThemedText>
								)}
							</>
						)}
					</form.Field>
				</Section>

				<Section>
					<ThemedButton onPress={() => form.handleSubmit()}>
						Submit
					</ThemedButton>
				</Section>
			</ScrollView>
		</Pressable>
	);
}
