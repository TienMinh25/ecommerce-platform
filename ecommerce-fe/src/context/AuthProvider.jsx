import { createContext, useEffect, useState } from 'react';
import { authService } from "../services/auth.js";
import { User } from "./type.js";

// eslint-disable-next-line react-refresh/only-export-components
export const AuthContext = createContext(null);

// Provider component
export const AuthProvider = ({ children }) => {
  const [user, setUser] = useState(null);
  const [isLoading, setIsLoading] = useState(true);
  const [authCheckComplete, setAuthCheckComplete] = useState(false);

  // Sử dụng useEffect thay vì useLayoutEffect để đảm bảo không chặn render UI
  useEffect(() => {
    // Check if user is already logged in
    const checkAuthStatus = async () => {
      setIsLoading(true);
      try {
        // Kiểm tra xem có access_token trong localStorage không
        let authToken = localStorage.getItem("access_token");

        if (authToken) {
          try {
            // Gọi API để validate token
            const response = await authService.validateToken();
            const data = response.data;

            if (data) {
              // Tạo đối tượng user từ data nhận được
              let user = new User(data["full_name"], data["avatar_url"], data["role"]);
              localStorage.setItem("user", JSON.stringify(user));
              setUser(user);
            }
          } catch (validationError) {
            console.error('Token validation failed:', validationError);
            // Xóa token không hợp lệ
            localStorage.removeItem('access_token');
            localStorage.removeItem('refresh_token');
            localStorage.removeItem('user');
          }
        } else {
          // Không có token
          console.log("No authentication token found");
        }
      } catch (error) {
        console.error('Authentication check error:', error);
      } finally {
        setIsLoading(false);
        setAuthCheckComplete(true);
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

      let user = new User(data["full_name"], data["avatar_url"], data["role"])

      localStorage.setItem("user", JSON.stringify(user));
      // Set user state
      setUser(user);

      return { success: true, error: null, needVerification: false };
    } catch (error) {
      if (error.response?.data?.error?.["error_code"] === "4002") {
        return {success: false, error: error.response.data.error.message, needVerification: true}
      }

      return {
        success: false,
        error: error.response?.data?.error?.message || 'Login failed',
        needVerification: false
      };
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
      return {
        success: false,
        error: error.response?.data?.error?.message || 'Registration failed'
      };
    } finally {
      setIsLoading(false);
    }
  };

  // Social login function
  const socialLogin = async (code, state, provider, getUrlOnly = false) => {
    // Nếu chỉ cần lấy URL xác thực
    if (getUrlOnly) {
      try {
        const response = await authService.getAuthorizationURL(provider);
        return { url: response.data["authorization_url"] };
      } catch (error) {
        console.error('Failed to get authorization URL:', error);
        return {
          success: false,
          error: error.response?.data?.error?.message || `Không thể kết nối với ${provider}`
        };
      }
    }

    // Nếu đã có code, tiến hành xác thực
    if (code && state) {
      setIsLoading(true);
      try {
        const response = await authService.exchangeOAuthCode(code, state, provider);

        let data = response.data;

        // Lưu token
        localStorage.setItem('access_token', data.access_token);
        localStorage.setItem("refresh_token", data.refresh_token);

        // Tạo đối tượng user từ response
        let user = new User(
            data.full_name,
            data.avatar_url,
            data.role
        );

        localStorage.setItem("user", JSON.stringify(user));
        setUser(user);

        return { success: true };
      } catch (error) {
        console.error('Social login error:', error);
        return {
          success: false,
          error: error.response?.data?.error?.message || 'Social login failed'
        };
      } finally {
        setIsLoading(false);
      }
    }

    return { success: false, error: 'Invalid request' };
  };

  // Logout function
  const logout = async () => {
    try {
      setIsLoading(true);
      await authService.logout(localStorage.getItem('refresh_token'))
      // Clear local storage
      localStorage.removeItem('access_token');
      localStorage.removeItem('refresh_token');
      localStorage.removeItem('user')
      // Clear user state
      setUser(null);

      return {success: true}
    } catch (error) {
      return {
        success: false,
        error: error.response?.data?.error?.message || "Something was wrong!"
      }
    } finally {
      setIsLoading(false);
    }
  };

  const verifyOTP = async (dataVerify) => {
    try {
      await authService.verifyOTP(dataVerify)

      return {success: true}
    } catch (error) {
      return {
        success: false,
        error: error.response?.data?.error?.message || "Verify email failed, please try again!"
      }
    }
  }

  const resendVerifyEmailOTP = async (dataResend) => {
    try {
      await authService.resendEmailOTP(dataResend)

      return {success: true}
    } catch (error) {
      return {
        success: false,
        error: error.response?.data?.error?.message || "Resend otp to verify email failed, please try again!"
      }
    }
  }

  const sendPasswordResetOTP = async (data) => {
    try {
      await authService.sendPasswordResetOTP(data)

      return {success: true, error: null}
    } catch (error) {
      return {
        success: false,
        error: error.response?.data?.error?.message || "Send otp forgot password failed, please try again!"
      }
    }
  }

  const resetPassword = async (data) => {
    try {
      await authService.resetPassword(data)

      return {success: true, error: null}
    } catch (error) {
      return {
        success: false,
        error: error.response?.data?.error?.message || "Reset password failed, please try again!"
      }
    }
  }

  const value = {
    user,
    isLoading,
    authCheckComplete, // Biến mới để đảm bảo kiểm tra xác thực đã hoàn tất
    login,
    register,
    socialLogin,
    logout,
    verifyOTP,
    resendVerifyEmailOTP,
    sendPasswordResetOTP,
    resetPassword,
  };

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
};