import axios from 'axios';
import config from '../config';

const { apiUri } = config;

export default axios.create({
    baseURL: apiUri,
});
