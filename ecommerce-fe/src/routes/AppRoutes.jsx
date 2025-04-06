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
import EmailVerificationPage from "../pages/Auth/EmailVerificationPage.jsx";
import ForgotPasswordPage from "../pages/Auth/ForgotPasswordPage.jsx";
import ResetPasswordPage from "../pages/Auth/ResetPasswordPage.jsx";
import OAuthCallbackPage from "../pages/Auth/OAuthCallbackPage.jsx";
import PermissionManagementComponent from "../pages/Module/Dashboard/PermissionManagementComponent.jsx";
import ModuleManagementComponent from "../pages/Module/Dashboard/ModuleManagementComponent.jsx";
import AddressTypesManagementComponent from "../pages/Module/Dashboard/AddressTypesManagementComponent.jsx";

const AppRoutes = () => {
    return (
        <Routes>
            {/* Public Routes */}
            <Route element={<PublicRoute />}>
                <Route path="/login" element={<Login />} />
                <Route path="/register" element={<Register />} />

                <Route path="/verify-email" element={<EmailVerificationPage />} />
                <Route path="/forgot-password" element={<ForgotPasswordPage />} />
                <Route path="/reset-password" element={<ResetPasswordPage />} />
                <Route path="/oauth" element={<OAuthCallbackPage />} />
            </Route>

            {/* Dashboard Routes */}
            <Route element={<PrivateRoute />}>
                <Route element={<MainLayout />}>
                    <Route path='/' element={<Home />} />
                </Route>
                <Route path='/dashboard' element={<DashboardLayout />}>
                    {/* Main dashboard */}
                    <Route index element={<Dashboard />} />

                    {/* Routes - render when clicks into the sidebar button */}
                    <Route path='users' element={<UserManagementComponent />} />
                    <Route path='roles' element={<RoleManagementComponent />} />
                    <Route path='permissions' element={<PermissionManagementComponent />} />
                    <Route path='modules' element={<ModuleManagementComponent />} />
                    <Route path='address-types' element={<AddressTypesManagementComponent />} />
                    {/*
                      <Route path='onboarding/suppliers' element={<DashboardComponent />} />
                      <Route path='onboarding/deliverers' element={<DashboardComponent />} /> */}
                </Route>
            </Route>



            {/* Private Routes */}
            <Route element={<PrivateRoute />}>
                <Route element={<MainLayout />}>
                    <Route path='/products' element={<ProductListing />} />
                    <Route path='/products/:id' element={<ProductDetail />} />
                </Route>
            </Route>

            {/* 404 Route */}
            <Route path='/404' element={<NotFound />} />
            <Route path='*' element={<Navigate to='/404' replace />} />
        </Routes>
    );
};

export default AppRoutes;