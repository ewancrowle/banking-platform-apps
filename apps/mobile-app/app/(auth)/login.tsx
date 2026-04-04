import { formOptions, useForm, useStore } from "@tanstack/react-form";
import { Alert, Keyboard, Pressable, ScrollView } from "react-native";
import * as z from "zod";
import trpc from "@/api/trpc";
import { getAccount } from "@/api/auth";
import { Section } from "@/components/section";
import { ThemedButton } from "@/components/themed-button";
import { ThemedInput } from "@/components/themed-input";
import { ThemedText } from "@/components/themed-text";
import { useAuthStore } from "@/store/auth";
import { isTRPCClientError } from "@trpc/client";
import { TRPC_ERROR_CODE_KEY } from "@trpc/server";
import { getTRPCErrorCode } from "@/utils/get-trpc-error-code";

const formSchema = z.object({
    email: z.email({
        message: "Enter a valid email address.",
    }),
    password: z
        .string()
        .regex(/^(?=.*?[A-Z])(?=.*?[a-z])(?=.*?[0-9])(?=.*?[#?!@$ %^&*-]).{8,}$/, {
            message:
                "Password must be at least 8 characters and include uppercase, lowercase, a number, and a special character.",
        }),
});

const formOpts = formOptions({
    defaultValues: {
        email: "",
        password: "",
    } as z.infer<typeof formSchema>,
    validators: {
        onChange: formSchema,
    }
});

export default function LoginScreen() {
    const form = useForm({
        ...formOpts,
        onSubmit: async ({ value }) => {
            try {
                const tokens = await trpc.auth.login.mutate(value);
                const store = useAuthStore.getState();
                await store.setTokens(tokens);
                const account = await getAccount();
                if (account) {
                    store.setAccount(account);
                }
            } catch (err) {
                const code = getTRPCErrorCode(err);
                if (code === "UNAUTHORIZED") {
                    Alert.alert("Incorrect email or password.");
                    form.reset();
                    return;
                }
                console.log(err);
                Alert.alert("An error occurred. Please try again later.");
            }
        },
    });

    const formErrorMap = useStore(form.store, (state) => state.errorMap)

    return (
        <Pressable style={{ flex: 1, padding: 16 }} onPress={Keyboard.dismiss}>
            <ScrollView showsVerticalScrollIndicator={false}>
                <ThemedText>Let&apos;s get you logged in.</ThemedText>

                <Section title="Enter your login details">
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
                                        {formErrorMap.onChange.email.map((issue) => issue.message).join(", ")}
                                    </ThemedText>
                                )}
                            </>
                        )}
                    </form.Field>

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
                                        {formErrorMap.onChange.password.map((issue) => issue.message).join(", ")}
                                    </ThemedText>
                                )}
                            </>
                        )}
                    </form.Field>
                </Section>

                <Section>
                    <ThemedButton onPress={() => form.handleSubmit()}>Login</ThemedButton>
                </Section>
            </ScrollView>
        </Pressable>
    );
}
