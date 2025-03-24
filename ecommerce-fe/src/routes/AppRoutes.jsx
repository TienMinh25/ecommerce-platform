import { Routes, Route, Navigate } from 'react-router-dom';
import PrivateRoute from './PrivateRoute';
import PublicRoute from './PublicRoute';

import Login from '../pages/Auth/Login';
import Register from '../pages/Auth/Register';
import Home from '../pages/Home';
import ProductListing from '../pages/ProductListing';
import ProductDetail from '../pages/ProductDetail';
import NotFound from '../pages/NotFound';

import MainLayout from '../components/layout/MainLayout';
import DashboardLayout from '../components/layout/DashboardLayout';

import Dashboard from '../pages/Dashboard';
import UserManagementComponent from '../pages/Module/Dashboard/UserManagementComponent';
import RoleManagementComponent from '../pages/Module/Dashboard/RoleManagementComponent';

const AppRoutes = () => {
  return (
    <Routes>
      {/* Public Routes */}
      <Route element={<PublicRoute />}>
        <Route path='/login' element={<Login />} />
        <Route path='/register' element={<Register />} />
      </Route>

      {/* Private Routes */}
      <Route element={<PrivateRoute />}>
        <Route element={<MainLayout />}>
          <Route path='/products' element={<ProductListing />} />
          <Route path='/products/:id' element={<ProductDetail />} />
        </Route>
      </Route>

      {/* Dashboard Routes */}
      <Route element={<PrivateRoute />}>
        <Route path='/dashboard' element={<DashboardLayout />}>
          {/* Main dashboard */}
          <Route index element={<Dashboard />} />

          {/* Routes - render when clicks into the sidebar button */}
          <Route path='users' element={<UserManagementComponent />} />
          <Route path='roles' element={<RoleManagementComponent />} />
          {/* <Route path='permissions' element={<DashboardComponent />} />
          <Route path='resources' element={<DashboardComponent />} />
          <Route path='suppliers' element={<DashboardComponent />} />
          <Route path='deliverers' element={<DashboardComponent />} /> */}
        </Route>
      </Route>

      {/* Mixed Access Routes */}
      <Route element={<MainLayout />}>
        <Route path='/' element={<Home />} />
      </Route>

      {/* 404 Route */}
      <Route path='/404' element={<NotFound />} />
      <Route path='*' element={<Navigate to='/404' replace />} />
    </Routes>
  );
};

export default AppRoutes;