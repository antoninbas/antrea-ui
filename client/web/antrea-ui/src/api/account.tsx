import { APIError, handleErrorResponse } from './common'
import config from '../config';
import { encode } from 'base-64';
const { apiServer, apiUri } = config;

export const accountAPI = {
    updatePassword: async (newPassword: string, token: string) => {
        try {
            let url = `${apiUri}/account/password`
            let response = await fetch(url, {
                method: "PUT",
                mode: "cors",
                headers: {
                    "Content-Type": "application/json",
                    "Authorization": `Bearer ${token}`,
                },
                body: JSON.stringify({password: encode(newPassword)}),
            });

            if (response.status !== 200) {
                throw new APIError(response.status, response.statusText, "Error when updating password");
            }
        } catch (err) {
            console.error("Failed to update password");
            throw err;
        }
    }
}
