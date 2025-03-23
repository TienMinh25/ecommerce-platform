import { Navigate, Outlet, useLocation } from 'react-router-dom';
import { useAuth } from '../hooks/useAuth';

const PublicRoute = () => {
  const { isAuthenticated } = useAuth();
  const location = useLocation();

  // If user is already authenticated and tries to access login/register page,
  // redirect them to the page they came from or to the home page
  if (isAuthenticated) {
    const from = location.state?.from || '/';
    return <Navigate to={from} replace />;
  }

  return <Outlet />;
};

export default PublicRoute;
