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

  // Thêm các phương thức OAuth
  getAuthorizationURL: async (provider) => {
    try {
      const response = await api.get(`/auth/oauth/url?oauth_provider=${provider}`);
      return response.data;
    } catch (error) {
      throw error;
    }
  },

  exchangeOAuthCode: async (code, state, provider) => {
    try {
      const response = await api.get(`/auth/oauth/exchange?code=${code}&state=${state}&oauth_provider=${provider}`);
      return response.data;
    } catch (error) {
      throw error;
    }
  },

  logout: async (refreshToken) => {
    try {
      await api.post('/auth/logout', {"refresh_token": refreshToken});
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

  verifyOTP: async (dataVerify) => {
    try {
      await api.post('/auth/verify-email', dataVerify)
    } catch (error) {
      throw error
    }
  },

  resendEmailOTP: async (dataResend) => {
    try {
      await api.post('/auth/resend-verify-email', dataResend)
    } catch (error) {
      throw error
    }
  },

  sendPasswordResetOTP: async (data) => {
    try {
      await api.post('/auth/forgot-password', data)
    } catch (error) {
      throw error
    }
  },

  resetPassword: async (data) => {
    try {
      await api.post('/auth/reset-password', data)
    } catch (error) {
      throw error
    }
  }
};