import { APIError, handleErrorResponse } from './common'
import config from '../config';
import { encode } from 'base-64';
const { apiServer, apiUri } = config;

interface Token {
    tokenType: string
    accessToken: string
    expiresIn: number
}

export const authAPI = {
    login: async (username: string, password: string): Promise<Token | undefined > => {
        try {
            let url = `${apiUri}/auth/login`
            let response = await fetch(url, {
                method: "GET",
                mode: "cors",
                headers: {
                    "Authorization": "Basic " + encode(username + ":" + password),
                },
            });

            if (response.status !== 200) {
                throw new APIError(response.status, response.statusText, "Error when trying to log in");
            }

            return response.json().then((data) => data as Token);
        } catch (err) {
            console.error("Login error");
            throw err;
        }
    },

    logout: async () => {
        try {
            let url = `${apiUri}/auth/logout`
            let response = await fetch(url, {
                method: "GET",
                mode: "cors",
            });

            if (response.status !== 200) {
                throw new APIError(response.status, response.statusText, "Error when trying to log out");
            }
        } catch (err) {
            console.error("Logout error");
            throw err;
        }
    },

    refreshToken: async () => {
        try {
            let url = `${apiUri}/auth/refresh_token`
            let response = await fetch(url, {
                method: "GET",
                mode: "cors",
            });

            if (response.status !== 200) {
                throw new APIError(response.status, response.statusText, "Error when trying to log out");
            }
        } catch (err) {
            console.error("Logout error");
            throw err;
        }
    },
}
