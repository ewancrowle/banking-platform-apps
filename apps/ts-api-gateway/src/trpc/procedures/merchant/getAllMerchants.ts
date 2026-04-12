import merchantClient from "../../../clients/merchant";
import { protectedProcedure } from "../..";

const getAllMerchants = protectedProcedure.query(async () => {
	const response = await merchantClient.getAllMerchants({});
	return response;
});

export default getAllMerchants;
