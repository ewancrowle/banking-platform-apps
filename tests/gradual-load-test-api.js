import superjson from "https://cdn.jsdelivr.net/npm/superjson@2.2.6/+esm";
import "https://jslib.k6.io/url/1.0.0/index.js";
import http from "k6/http";
import faker from "k6/x/faker";

const API_URL = "http://localhost:3000";

export const options = {
	scenarios: {
		staircase_load_test: {
			executor: "ramping-arrival-rate",
			startRate: 0,
			timeUnit: "1s",
			stages: [
				{ target: 10, duration: "5s" }, // Ramp to 10
				{ target: 10, duration: "10s" }, // Hold 10
				{ target: 20, duration: "5s" }, // Ramp to 20
				{ target: 20, duration: "10s" }, // Hold 20
				{ target: 30, duration: "5s" }, // Ramp to 30
				{ target: 30, duration: "10s" }, // Hold 30
				{ target: 40, duration: "5s" }, // Ramp to 40
				{ target: 40, duration: "10s" }, // Hold 40
			],
			preAllocatedVUs: 50,
			maxVUs: 400,
		},
	},
};

export function setup() {
	const res = http.post(
		`${API_URL}/auth.signUp`,
		JSON.stringify({
			json: {
				firstName: faker.person.firstName(),
				lastName: faker.person.lastName(),
				email: faker.person.email(),
				password: "@oK9t@b$",
				line1: faker.address.street(),
				town: faker.address.city(),
				postcode: faker.address.zip(),
			},
		}),
		{
			headers: {
				"Content-Type": "application/json",
				"X-Forwarded-For": faker.internet.ipv4Address(),
			},
		},
	);

	return {
		accessToken: res.json("result.data.json.accessToken"),
	};
}

export default function (data) {
	http.post(
		`${API_URL}/payment.newCardPayment`,
		superjson.stringify({
			merchantId: BigInt("70109715689897984"),
			amount: 500,
		}),
		{
			headers: {
				"Content-Type": "application/json",
				"X-Forwarded-For": faker.internet.ipv4Address(),
				Authorization: `Bearer ${data.accessToken}`,
			},
		},
	);
}
