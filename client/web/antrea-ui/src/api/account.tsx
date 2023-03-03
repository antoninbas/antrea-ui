import api from './axios'
import { handleError } from './common'
import { getToken } from './token'
import { encode } from 'base-64';

export const accountAPI = {
    updatePassword: async (newPassword: string): Promise<void> => {
        return api.put(`account/password`, JSON.stringify({password: encode(newPassword)}), {
            headers: {
                "Content-Type": "application/json",
                "Authorization": `Bearer ${getToken}`,
            },
        }).then((response) => {}).catch((error) => handleError(error, "Error when updating password"))
    },
}
