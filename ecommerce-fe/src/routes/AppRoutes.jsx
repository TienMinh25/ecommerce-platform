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
import UserManagementComponent from '../components/dashboard/admin/user-management/UserManagementComponent.jsx';
import RoleManagementComponent from '../components/dashboard/admin/role-management/RoleManagementComponent.jsx';
import EmailVerificationPage from "../pages/Auth/EmailVerificationPage.jsx";
import ForgotPasswordPage from "../pages/Auth/ForgotPasswordPage.jsx";
import ResetPasswordPage from "../pages/Auth/ResetPasswordPage.jsx";
import OAuthCallbackPage from "../pages/Auth/OAuthCallbackPage.jsx";
import PermissionManagementComponent from "../components/dashboard/admin/permission-management/PermissionManagementComponent.jsx";
import ModuleManagementComponent from "../components/dashboard/admin/module-management/ModuleManagementComponent.jsx";
import AddressTypesManagementComponent from "../components/dashboard/admin/address-types-management/AddressTypesManagementComponent.jsx";
import UserAccountLayout from "../components/user/UserAccountLayout.jsx";
import UserProfile from "../components/user/UserProfile.jsx";
import ChangePassword from "../components/user/ChangePassword.jsx";
import UserAddresses from "../components/user/UserAddresses.jsx";
import NotificationSettings from "../components/user/NotificationSettings.jsx";
import UserOrders from "../components/user/UserOrders.jsx";
import UserNotifications from "../components/user/UserNotifications.jsx";
import CartPage from "../pages/Cart/CartPage.jsx";
import CouponManagementComponent from "../components/dashboard/admin/coupon-management/CouponManagementComponent.jsx";
import CheckoutPage from "../pages/Checkout/CheckoutPage.jsx";
import SupplierRegistration from "../pages/SupplierRegistration.jsx";
import DelivererRegistration from "../pages/DelivererRegistration.jsx";

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
                    {/* User Account Routes */}
                    <Route path='/user/account' element={<UserAccountLayout />}>
                        <Route index element={<Navigate to='/user/account/profile' replace />} />
                        <Route path='profile' element={<UserProfile />} />
                        <Route path='password' element={<ChangePassword />} />
                        <Route path='addresses' element={<UserAddresses />} />
                        <Route path='notifications/settings' element={<NotificationSettings />} />
                        <Route path='orders' element={<UserOrders />} />
                        <Route path='notifications/see' element={<UserNotifications />} />
                    </Route>

                    {/* Registration Routes */}
                    <Route path='/register/supplier' element={<SupplierRegistration />} />
                    <Route path='/register/deliverer' element={<DelivererRegistration />} />

                    <Route path='/products' element={<ProductListing />} />
                    <Route path='/products/:id' element={<ProductDetail />} />
                    <Route path='carts' element={<CartPage />} />
                    <Route path='/checkout' element={<CheckoutPage />} />
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
                    <Route path='coupons' element={<CouponManagementComponent />} />
                    {/*
                      <Route path='onboarding/suppliers' element={<DashboardComponent />} />
                      <Route path='onboarding/deliverers' element={<DashboardComponent />} /> */}
                </Route>
            </Route>

            {/* 404 Route */}
            <Route path='/404' element={<NotFound />} />
            <Route path='*' element={<Navigate to='/404' replace />} />
        </Routes>
    );
};

export default AppRoutes;