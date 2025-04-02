import api from './api';

export const authService = {
  login: async (credentials) => {
    // eslint-disable-next-line no-useless-catch
    try {
      const response = await api.post('/auth/login', credentials);
      return response.data;
    } catch (error) {
      throw error;
    }
  },

  register: async (userData) => {
    // eslint-disable-next-line no-useless-catch
    try {
      await api.post('/auth/register', userData);
    } catch (error) {
      throw error;
    }
  },

  socialLogin: async (provider, token) => {
    // eslint-disable-next-line no-useless-catch
    try {
      const response = await api.post(`/auth/${provider}`, { token });
      return response.data;
    } catch (error) {
      throw error;
    }
  },

  logout: async (refreshToken) => {
    try {
      await api.post('/auth/logout', {"refresh-token": refreshToken});
    } catch (error) {
      throw error
    }
  },

  // Thêm function để validate token
  validateToken: async () => {
    try {
      const response = await api.get('/auth/check-token');
      return response.data;
    } catch (error) {
      console.error('Token validation error:', error);
      return null;
    }
  },
};