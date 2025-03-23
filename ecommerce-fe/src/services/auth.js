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
      const response = await api.post('/auth/register', userData);
      return response.data;
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

  logout: async () => {
    try {
      await api.post('/auth/logout');
      return true;
    } catch (error) {
      console.error('Logout error:', error);
      return false;
    }
  },
};
