import { getLocales } from "expo-localization";

export default function (
	amount: number | bigint,
	currencyCode: string = "GBP",
) {
	return new Intl.NumberFormat(getLocales()[0].languageTag, {
		style: "currency",
		currency: currencyCode,
	}).format(Number(amount) / 100);
}
