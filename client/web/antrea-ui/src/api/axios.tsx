import axios from 'axios';
import config from '../config';
import { getToken, setToken } from './token';

const { apiUri } = config;

const api = axios.create({
    baseURL: apiUri,
});

api.interceptors.request.use((request) => {
    request.headers['Authorization'] = `Bearer ${getToken()}`;
    return request;
});

export default api;
