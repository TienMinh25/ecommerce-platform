import axios from 'axios';

// Create an axios instance with default config
const api = axios.create({
    baseURL: 'http://server.local:3000/api/v1',
    headers: {
        'Content-Type': 'application/json',
    },
});

// Flag to prevent multiple refresh token requests
let isRefreshing = false;
// Store pending requests that should be retried after token refresh
let failedQueue = [];

const processQueue = (error, token = null) => {
    failedQueue.forEach(prom => {
        if (error) {
            prom.reject(error);
        } else {
            prom.resolve(token);
        }
    });

    failedQueue = [];
};

// Add request interceptor for auth token
api.interceptors.request.use(
    (config) => {
        const token = localStorage.getItem('access_token');
        if (token) {
            config.headers.Authorization = `Bearer ${token}`;
        }
        return config;
    },
    (error) => Promise.reject(error)
);

// Add response interceptor for error handling
api.interceptors.response.use(
    (response) => response,
    async (error) => {
        const originalRequest = error.config;

        // Prevent infinite loops
        if (originalRequest._retry) {
            return Promise.reject(error);
        }

        // Handle 401 Unauthorized errors
        if (error.response &&
            error.response.status === 401 &&
            error.response.data &&
            error.response.data.code === 4001) {

            if (isRefreshing) {
                // If a refresh is already in progress, add the request to queue
                return new Promise((resolve, reject) => {
                    failedQueue.push({ resolve, reject });
                })
                    .then(token => {
                        // Create a NEW request with the same config but new token
                        const newRequest = { ...originalRequest };
                        newRequest.headers.Authorization = `Bearer ${token}`;
                        newRequest._retry = true; // Mark as retried to prevent loops
                        return axios(newRequest);
                    })
                    .catch(err => Promise.reject(err));
            }

            originalRequest._retry = true;
            isRefreshing = true;

            try {
                // Try refreshing the token
                const refreshToken = localStorage.getItem('refresh_token');

                if (!refreshToken) {
                    // If no refresh token, logout
                    throw new Error('No refresh token available');
                }

                // Use a direct axios call (NOT through the api instance) to avoid interceptors
                const response = await axios({
                    method: 'post',
                    url: 'http://server.local:3000/api/v1/auth/refresh',
                    headers: {
                        'X-Authorization': `Bearer ${refreshToken}`
                    }
                });

                if (response.data && response.data.access_token) {
                    const { access_token, refresh_token } = response.data;

                    // Store new tokens
                    localStorage.setItem('access_token', access_token);
                    if (refresh_token) {
                        localStorage.setItem('refresh_token', refresh_token);
                    }

                    // Update authorization header
                    api.defaults.headers.common['Authorization'] = `Bearer ${access_token}`;

                    // Create a new request with the same config
                    const newRequest = { ...originalRequest };
                    newRequest.headers['Authorization'] = `Bearer ${access_token}`;
                    newRequest._retry = true; // Mark as retried to prevent loops

                    // Process all queued requests
                    processQueue(null, access_token);

                    // Retry original request with a new axios call
                    return axios(newRequest);
                } else {
                    throw new Error('Refresh token request successful but no token returned');
                }
            } catch (refreshError) {
                // Handle refresh token errors
                processQueue(refreshError, null);

                // If refresh token also returns 401, logout
                if (refreshError.response &&
                    refreshError.response.status === 401) {
                    // Clear storage and redirect to login
                    localStorage.removeItem('access_token');
                    localStorage.removeItem('refresh_token');
                    localStorage.removeItem('user_info');
                    window.location.href = '/login';
                }

                return Promise.reject(refreshError);
            } finally {
                isRefreshing = false;
            }
        }

        // For login 401 (code 4008) - just reject, don't try to refresh
        if (error.response &&
            error.response.status === 401 &&
            error.response.data &&
            error.response.data.code === 4008) {
            return Promise.reject(error);
        }

        // For other errors, just reject
        return Promise.reject(error);
    }
);

export default api;