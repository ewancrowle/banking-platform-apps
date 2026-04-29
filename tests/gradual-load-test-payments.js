import http from "k6/http";
import faker from "k6/x/faker";

const PAYMENTS_URL = "http://localhost:8080";

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
				{ target: 50, duration: "5s" }, // Ramp to 50
				{ target: 50, duration: "10s" }, // Hold 50
				{ target: 60, duration: "5s" }, // Ramp to 60
				{ target: 60, duration: "10s" }, // Hold 60
				{ target: 70, duration: "5s" }, // Ramp to 70
				{ target: 70, duration: "10s" }, // Hold 70
				{ target: 80, duration: "5s" }, // Ramp to 80
				{ target: 80, duration: "10s" }, // Hold 80
				{ target: 90, duration: "5s" }, // Ramp to 90
				{ target: 90, duration: "10s" }, // Hold 90
				{ target: 100, duration: "5s" }, // Ramp to 100
				{ target: 100, duration: "10s" }, // Hold 100
			],
			preAllocatedVUs: 100,
			maxVUs: 1000,
		},
	},
};

export default function () {
	http.post(
		`${PAYMENTS_URL}/payment.v1.PaymentService/AuthorisePayment`,
		JSON.stringify({
			accountId: "70114039698554880",
			merchantId: "70109715689897984",
			amount: 100,
			currencyCode: "GBP",
			description: "Test Payment",
			type: "card",
		}),
		{
			headers: {
				"Content-Type": "application/json",
				"X-Forwarded-For": faker.internet.ipv4Address(),
			},
		},
	);
}
