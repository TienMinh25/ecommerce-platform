import {createContext, useEffect, useLayoutEffect, useState} from 'react';
import {authService} from "../services/auth.js";
import {User} from "./type.js"

// eslint-disable-next-line react-refresh/only-export-components
export const AuthContext = createContext(null);

// Provider component
export const AuthProvider = ({ children }) => {
  const [user, setUser] = useState(null);
  const [isLoading, setIsLoading] = useState(true);

  useLayoutEffect(() => {
    // Check if user is already logged in
    const checkAuthStatus = async () => {
      try {
        let authToken = localStorage.getItem("access_token")

        if (authToken) {
          let data = authService.validateToken().data

          if (data) {
            let user = new User(data["full_name"], data["avatar_url"], data["roles"])
            localStorage.setItem("user", JSON.stringify(user));
            setUser(user);
          }
        }
      } catch (error) {
        console.error('Authentication error:', error);
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
      const response = await authService.login(credentials);

      let data = response.data;
      // Save token
      localStorage.setItem('access_token', data["access_token"]);
      localStorage.setItem("refresh_token", data["refresh_token"]);

      let user = new User(data["full_name"], data["avatar_url"], data["roles"])

      localStorage.setItem("user", JSON.stringify(user));
      // Set user state
      setUser(user);

      return { success: true };
    } catch (error) {
      return { success: false, error: error.response.data.error.message || 'Login failed' };
    } finally {
      setIsLoading(false);
    }
  };

  // Register function
  const register = async (userData) => {
    setIsLoading(true);
    try {
      await authService.register(userData);

      return { success: true };
    } catch (error) {
      return { success: false, error: error.response.data.error.message || 'Registration failed' };
    } finally {
      setIsLoading(false);
    }
  };

  // todo: lam sau
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
    try {
      authService.logout(localStorage.getItem('refresh_token'))
      // Clear local storage
      localStorage.removeItem('access_token');
      localStorage.removeItem('refresh_token');
      localStorage.removeItem('user')
      // Clear user state
      setUser(null);

      return {success: true}
    } catch (error) {
      return {success: false, error: error.response.data.error.message || "Something was wrong!"}
    }
  };

  const value = {
    user,
    isLoading,
    login,
    register,
    socialLogin,
    logout,
  };

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
};