import { useContext } from 'react';
import { AuthContext } from '../context/AuthProvider';

// Custom hook for using the auth context
export const useAuth = () => {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
};

// Default export for importing the hook
export default useAuth;
