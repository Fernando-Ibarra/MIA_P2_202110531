import axios from 'axios';
import { getEnvVariables } from '../helpers';

const { VITE_API_URL } = getEnvVariables();

const appApi = axios.create({
    baseURL: VITE_API_URL,
});

appApi.interceptors.request.use( config => {
    
    config.headers = {
        ...config.headers,
    }

    return config;
})

export default appApi;