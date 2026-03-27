export type Account = {
	id: bigint;
	firstName: string;
	middleNames: string;
	lastName: string;
	email: string;
	line1: string;
	line2: string;
	town: string;
	postcode: string;
	createdAt: string;
};

export type Tokens = {
	accessToken: string;
	refreshToken: string;
};
