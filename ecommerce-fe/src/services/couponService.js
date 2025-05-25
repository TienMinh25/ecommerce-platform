import api from './api.js';

const couponService = {
    // Get list of coupons with pagination and filters
    getCoupons: async (params = {}) => {
        try {
            const response = await api.get('/coupons', { params });
            return response.data;
        } catch (error) {
            console.error('Error fetching coupons:', error);
            throw error;
        }
    },

    // Get coupon by ID
    getCouponById: async (couponId) => {
        try {
            const response = await api.get(`/coupons/${couponId}`);
            return response.data;
        } catch (error) {
            console.error('Error fetching coupon by ID:', error);
            throw error;
        }
    },

    getCouponsForClient: async (params = {}) => {
        try {
            // Add current date as required parameter
            const clientParams = {
                current_date: new Date().toISOString(),
                ...params
            };
            const response = await api.get('/coupons/client', { params: clientParams });
            return response.data;
        } catch (error) {
            console.error('Error fetching client coupons:', error);
            throw error;
        }
    },

    // Create new coupon
    createCoupon: async (couponData) => {
        try {
            const response = await api.post('/coupons', couponData);
            return response.data;
        } catch (error) {
            console.error('Error creating coupon:', error);
            throw error;
        }
    },

    // Update coupon
    updateCoupon: async (couponId, couponData) => {
        try {
            const response = await api.patch(`/coupons/${couponId}`, couponData);
            return response.data;
        } catch (error) {
            console.error('Error updating coupon:', error);
            throw error;
        }
    },

    // Delete coupon
    deleteCoupon: async (couponId) => {
        try {
            const response = await api.delete(`/coupons/${couponId}`);
            return response.data;
        } catch (error) {
            console.error('Error deleting coupon:', error);
            throw error;
        }
    }
};

export default couponService;