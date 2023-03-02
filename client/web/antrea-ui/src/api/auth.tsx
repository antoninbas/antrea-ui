import api from './axios'
import { handleError } from './common'
import { encode } from 'base-64';

interface Token {
    tokenType: string
    accessToken: string
    expiresIn: number
}

export const authAPI = {
    login: async (username: string, password: string): Promise<Token> => {
        return api.get(`auth/login`, {
            headers: {
                "Authorization": "Basic " + encode(username + ":" + password),
            },
        }).then((response) => response.data as Token).catch(error => handleError(error, "Error when trying to log in"))
    },

    logout: async (): Promise<void> => {
        return api.get(`auth/logout`).then(_ => {}).catch((error) => handleError(error, "Error when trying to log out"))
    },

    refreshToken: async (): Promise<void> => {
        return api.get(`auth/refresh_token`).then(_ => {}).catch((error) => handleError(error))
    },
}
