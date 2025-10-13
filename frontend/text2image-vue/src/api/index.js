// src/api/index.js
import axios from 'axios';

// 从环境变量读取 API 地址，如果没设置则 fallback 到本地
const API_BASE_URL = process.env.VUE_APP_API_BASE_URL || 'http://localhost:8080';

const apiClient = axios.create({
    baseURL: API_BASE_URL,
    timeout: 10000,
});

// 请求拦截器
apiClient.interceptors.request.use(
    config => {
        const token = localStorage.getItem('token');
        if (token) {
            // 将 token 添加到请求头
            config.headers.Authorization = `${token}`;
        }
        return config;
    },
    error => {
        return Promise.reject(error);
    }
);

// 响应拦截器
apiClient.interceptors.response.use(
    response => {
        return response;
    },
    error => {
        return Promise.reject(error);
    }
);

export default apiClient;