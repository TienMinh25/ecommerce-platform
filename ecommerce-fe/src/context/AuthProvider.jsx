import { createContext, useEffect, useState } from 'react';

// eslint-disable-next-line react-refresh/only-export-components
export const AuthContext = createContext(null);

// Provider component
export const AuthProvider = ({ children }) => {
  const [user, setUser] = useState(null);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    // Check if user is already logged in
    const checkAuthStatus = async () => {
      try {
        // Try to get the current user from localStorage or cookies
        const token = localStorage.getItem('token');
        if (token) {
          // TODO: Validate token with backend
          setUser({ id: '1', name: 'John Doe', email: 'john@example.com' }); // Placeholder user
        }
      } catch (error) {
        console.error('Authentication error:', error);
        // Clear any invalid auth data
        localStorage.removeItem('token');
      } finally {
        setIsLoading(false);
      }
    };

    checkAuthStatus();
  }, []);

  // Login function
  const login = async (credentials) => {
    setIsLoading(true);
    try {
      // TODO: Implement actual login logic with API
      // const response = await authService.login(credentials);
      const mockResponse = {
        user: { id: '1', name: 'John Doe', email: credentials.email },
        token: 'mock-jwt-token',
      };

      // Save token
      localStorage.setItem('token', mockResponse.token);

      // Set user state
      setUser(mockResponse.user);
      return { success: true };
    } catch (error) {
      console.error('Login error:', error);
      return { success: false, error: error.message || 'Login failed' };
    } finally {
      setIsLoading(false);
    }
  };

  // Register function
  const register = async (userData) => {
    setIsLoading(true);
    try {
      // TODO: Implement actual registration logic with API
      // const response = await authService.register(userData);
      const mockResponse = {
        user: { id: '1', name: userData.name, email: userData.email },
        token: 'mock-jwt-token',
      };

      // Save token
      localStorage.setItem('token', mockResponse.token);

      // Set user state
      setUser(mockResponse.user);
      return { success: true };
    } catch (error) {
      console.error('Registration error:', error);
      return { success: false, error: error.message || 'Registration failed' };
    } finally {
      setIsLoading(false);
    }
  };

  // Social login function
  const socialLogin = async (provider) => {
    setIsLoading(true);
    try {
      // TODO: Implement social login logic
      // const response = await authService.socialLogin(provider);
      const mockResponse = {
        user: { id: '1', name: 'Social User', email: 'social@example.com' },
        token: 'mock-jwt-token',
      };

      localStorage.setItem('token', mockResponse.token);
      setUser(mockResponse.user);
      return { success: true };
    } catch (error) {
      console.error('Social login error:', error);
      return { success: false, error: error.message || 'Social login failed' };
    } finally {
      setIsLoading(false);
    }
  };

  // Logout function
  const logout = () => {
    // Clear local storage
    localStorage.removeItem('token');
    // Clear user state
    setUser(null);
  };

  const value = {
    user,
    isAuthenticated: !!user,
    isLoading,
    login,
    register,
    socialLogin,
    logout,
  };

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
};
