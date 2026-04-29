import superjson from "https://cdn.jsdelivr.net/npm/superjson@2.2.6/+esm";
import { expect } from "https://jslib.k6.io/k6chaijs/4.5.0.1/index.js";
import "https://jslib.k6.io/url/1.0.0/index.js";
import { sleep } from "k6";
import http from "k6/http";
import faker from "k6/x/faker";

const API_URL = "http://localhost:3000";

const ipAddress = faker.internet.ipv4Address();

const customerDetails = {
	firstName: faker.person.firstName(),
	lastName: faker.person.lastName(),
	email: faker.person.email(),
	password: "@oK9t@b$",
	line1: faker.address.street(),
	town: faker.address.city(),
	postcode: faker.address.zip(),
};

export default function () {
	const signUpRes = http.post(
		`${API_URL}/auth.signUp`,
		JSON.stringify({ json: customerDetails }),
		{
			headers: {
				"Content-Type": "application/json",
				"X-Forwarded-For": ipAddress,
			},
		},
	);

	const accessToken = signUpRes.json("result.data.json.accessToken");

	expect(signUpRes.status, "response status").to.equal(200);
	expect(accessToken, "access token").to.exist;
	expect(signUpRes.json("result.data.json.refreshToken"), "refresh token").to
		.exist;

	const getAccountRes = http.get(`${API_URL}/account.getAccount`, {
		headers: {
			"Content-Type": "application/json",
			"X-Forwarded-For": ipAddress,
			Authorization: `Bearer ${accessToken}`,
		},
	});

	const accountId = getAccountRes.json("result.data.json.id");

	expect(getAccountRes.status, "response status").to.equal(200);
	expect(accountId, "account id").to.exist;

	const newDepositRes = http.post(
		`${API_URL}/payment.newDeposit`,
		superjson.stringify({
			amount: 500,
		}),
		{
			headers: {
				"Content-Type": "application/json",
				"X-Forwarded-For": ipAddress,
				Authorization: `Bearer ${accessToken}`,
			},
		},
	);

	expect(newDepositRes.status, "response status").to.equal(200);
	expect(newDepositRes.json("result.data.json.paymentId"), "payment id").to
		.exist;

	sleep(5);

	const newCardPaymentRes = http.post(
		`${API_URL}/payment.newCardPayment`,
		superjson.stringify({
			merchantId: BigInt("70109715689897984"),
			amount: 100,
		}),
		{
			headers: {
				"Content-Type": "application/json",
				"X-Forwarded-For": ipAddress,
				Authorization: `Bearer ${accessToken}`,
			},
		},
	);

	expect(newCardPaymentRes.status, "response status").to.equal(200);
	expect(newCardPaymentRes.json("result.data.json.paymentId"), "payment id").to
		.exist;

	const getPaymentsRes = http.get(`${API_URL}/payment.getPayments`, {
		headers: {
			"Content-Type": "application/json",
			"X-Forwarded-For": ipAddress,
			Authorization: `Bearer ${accessToken}`,
		},
	});

	expect(getPaymentsRes.json("result.data.json.payments"), "payments").to.not.be
		.empty;

	sleep(10);

	const getBalancesRes = http.get(`${API_URL}/balance.getBalances`, {
		headers: {
			"Content-Type": "application/json",
			"X-Forwarded-For": ipAddress,
			Authorization: `Bearer ${accessToken}`,
		},
	});

	expect(getBalancesRes.status, "response status").to.equal(200);
	expect(
		getBalancesRes.json("result.data.json.availableBalance"),
		"available balance",
	).to.equal("400");
}
